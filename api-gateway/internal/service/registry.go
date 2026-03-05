package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int32

const (
	StateClosed   CircuitState = iota // Normal operation, requests pass through
	StateOpen                         // Service is down, requests are blocked
	StateHalfOpen                     // Recovery mode, one probe request allowed
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

const (
	defaultFailThreshold  = int32(3)
	defaultRecoveryWindow = 30 * time.Second
	defaultHTTPTimeout    = 5 * time.Second
	healthCheckPath       = "/health"
)

// BackendURL tracks the health of a single upstream instance.
type BackendURL struct {
	URL       string
	failCount atomic.Int32
	healthy   atomic.Bool
}

func newBackendURL(url string) *BackendURL {
	b := &BackendURL{URL: url}
	b.healthy.Store(true)
	return b
}

func (b *BackendURL) IsHealthy() bool      { return b.healthy.Load() }
func (b *BackendURL) FailCount() int32     { return b.failCount.Load() }
func (b *BackendURL) markHealthy()         { b.healthy.Store(true); b.failCount.Store(0) }
func (b *BackendURL) markUnhealthy()       { b.healthy.Store(false) }
func (b *BackendURL) incrementFail() int32 { return b.failCount.Add(1) }

// ServiceHealth holds circuit-breaker state for a service.
type ServiceHealth struct {
	mu           sync.RWMutex
	state        CircuitState
	lastFailTime time.Time
	LastCheck    time.Time
}

func newServiceHealth() *ServiceHealth {
	return &ServiceHealth{
		state:     StateClosed,
		LastCheck: time.Now(),
	}
}

func (h *ServiceHealth) State() CircuitState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.state
}

// CanRequest checks whether a request is allowed given the current circuit state.
func (h *ServiceHealth) CanRequest() bool {
	h.mu.RLock()
	state := h.state
	lastFail := h.lastFailTime
	h.mu.RUnlock()

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(lastFail) >= defaultRecoveryWindow {
			h.mu.Lock()
			// Re-check after acquiring write lock (double-checked locking).
			if h.state == StateOpen && time.Since(h.lastFailTime) >= defaultRecoveryWindow {
				h.state = StateHalfOpen
				log.Printf("🟡 Circuit entering HALF-OPEN state")
			}
			h.mu.Unlock()
			return h.state == StateHalfOpen
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (h *ServiceHealth) recordSuccess(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.state != StateClosed {
		log.Printf("🟢 Service UP (Circuit CLOSED): %s", name)
	}
	h.state = StateClosed
	h.LastCheck = time.Now()
}

func (h *ServiceHealth) recordFailure(name string, threshold int32, currentFails int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.LastCheck = time.Now()
	if currentFails >= threshold && h.state == StateClosed {
		log.Printf("🔴 Service DOWN (Circuit OPEN): %s (fails: %d)", name, currentFails)
		h.state = StateOpen
		h.lastFailTime = time.Now()
	}
}

// Service represents a registered upstream service with load-balancing support.
type Service struct {
	Name       string
	Backends   []*BackendURL
	PathPrefix string
	Health     *ServiceHealth
	Timeout    time.Duration
	FailThresh int32
	nextIndex  atomic.Uint64
}

// GetNextBackend returns the next healthy backend using round-robin.
// Returns ("", false) when no healthy backend is available or the circuit is open.
func (s *Service) GetNextBackend() (string, bool) {
	if !s.Health.CanRequest() {
		return "", false
	}

	n := uint64(len(s.Backends))
	// Try every backend at most once to find a healthy one.
	start := s.nextIndex.Add(1) - 1
	for i := uint64(0); i < n; i++ {
		b := s.Backends[(start+i)%n]
		if b.IsHealthy() {
			return b.URL, true
		}
	}
	return "", false
}

// BaseURLs returns the raw URL strings (for backwards-compatibility / logging).
func (s *Service) BaseURLs() []string {
	urls := make([]string, len(s.Backends))
	for i, b := range s.Backends {
		urls[i] = b.URL
	}
	return urls
}

// HealthSummary returns a human-readable summary for status endpoints.
type HealthSummary struct {
	ServiceName string
	Circuit     string
	LastCheck   time.Time
	Backends    []BackendSummary
}

type BackendSummary struct {
	URL       string
	Healthy   bool
	FailCount int32
}

func (s *Service) HealthSummary() HealthSummary {
	backends := make([]BackendSummary, len(s.Backends))
	for i, b := range s.Backends {
		backends[i] = BackendSummary{
			URL:       b.URL,
			Healthy:   b.IsHealthy(),
			FailCount: b.FailCount(),
		}
	}
	return HealthSummary{
		ServiceName: s.Name,
		Circuit:     s.Health.State().String(),
		LastCheck:   s.Health.LastCheck,
		Backends:    backends,
	}
}

// ServiceRegistry holds all registered services and manages their lifecycle.
type ServiceRegistry struct {
	services   sync.Map
	httpClient *http.Client
}

// NewServiceRegistry creates an empty registry.
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Global üst limit
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				MaxIdleConnsPerHost: 10,
			},
		},
	}
}

// Regi}ster adds (or replaces) a service in the registry.
func (sr *ServiceRegistry) Register(name string, baseURLs []string, pathPrefix string) error {
	if len(baseURLs) == 0 {
		return fmt.Errorf("service %q: at least one base URL is required", name)
	}
	if pathPrefix == "" {
		return fmt.Errorf("service %q: pathPrefix must not be empty", name)
	}

	backends := make([]*BackendURL, len(baseURLs))
	for i, u := range baseURLs {
		backends[i] = newBackendURL(strings.TrimRight(u, "/"))
	}

	svc := &Service{
		Name:       name,
		Backends:   backends,
		PathPrefix: pathPrefix,
		Health:     newServiceHealth(),
		Timeout:    defaultHTTPTimeout,
		FailThresh: defaultFailThreshold,
	}

	sr.services.Store(name, svc)
	log.Printf("✅ Service registered: %s -> %v (prefix: %s)", name, baseURLs, pathPrefix)
	return nil
}

// Deregister removes a service from the registry.
func (sr *ServiceRegistry) Deregister(name string) {
	sr.services.Delete(name)
	log.Printf("🗑️  Service deregistered: %s", name)
}

// GetByPath returns the service whose PathPrefix is the longest match for path.
func (sr *ServiceRegistry) GetByPath(path string) (*Service, bool) {
	var best *Service
	bestLen := 0

	sr.services.Range(func(_, value interface{}) bool {
		svc := value.(*Service)
		if strings.HasPrefix(path, svc.PathPrefix) && len(svc.PathPrefix) > bestLen {
			best = svc
			bestLen = len(svc.PathPrefix)
		}
		return true
	})

	return best, best != nil
}

// GetByName looks up a service by its registered name.
func (sr *ServiceRegistry) GetByName(name string) (*Service, bool) {
	v, ok := sr.services.Load(name)
	if !ok {
		return nil, false
	}
	return v.(*Service), true
}

// List returns all registered services (order is non-deterministic).
func (sr *ServiceRegistry) List() []*Service {
	var list []*Service
	sr.services.Range(func(_, value interface{}) bool {
		list = append(list, value.(*Service))
		return true
	})
	return list
}

// IsHealthy reports whether the service circuit is closed (i.e. traffic is allowed).
func (sr *ServiceRegistry) IsHealthy(svc *Service) bool {
	return svc.Health.CanRequest()
}

// HealthSummaries returns a snapshot of every service's health for status pages.
func (sr *ServiceRegistry) HealthSummaries() []HealthSummary {
	var out []HealthSummary
	sr.services.Range(func(_, value interface{}) bool {
		out = append(out, value.(*Service).HealthSummary())
		return true
	})
	return out
}

// StartHealthChecks launches a background goroutine that periodically probes
// every registered service.  It stops cleanly when ctx is cancelled.
func (sr *ServiceRegistry) StartHealthChecks(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 Health checks stopped")
				return
			case <-ticker.C:
				sr.services.Range(func(_, value interface{}) bool {
					go sr.checkHealth(value.(*Service))
					return true
				})
			}
		}
	}()
	log.Printf("🏥 Health checks started (interval: %v)", interval)
}

// checkHealth probes every backend of svc and updates circuit-breaker state.
// HTTP calls are intentionally performed OUTSIDE any mutex to prevent lock contention.
func (sr *ServiceRegistry) checkHealth(svc *Service) {
	//client := &http.Client{Timeout: svc.Timeout}

	var wg sync.WaitGroup
	var totalFails int32

	for _, b := range svc.Backends {
		wg.Add(1)
		go func(backend *BackendURL) {
			defer wg.Done()

			url := backend.URL + healthCheckPath

			resp, err := sr.httpClient.Get(url)

			if err != nil || resp.StatusCode != http.StatusOK {
				fails := backend.incrementFail()
				if fails >= svc.FailThresh {
					backend.markUnhealthy()
					log.Printf("⚠️  Backend unhealthy: %s (fails: %d)", backend.URL, fails)
				}
				atomic.AddInt32(&totalFails, 1)
			} else {
				resp.Body.Close()
				if !backend.IsHealthy() {
					log.Printf("✅ Backend recovered: %s", backend.URL)
				}
				backend.markHealthy()
			}
		}(b)
	}

	wg.Wait() // Tüm backend'lerin kontrolü bitene kadar bekle

	// Servis seviyesindeki Circuit Breaker'ı güncelle
	if totalFails == 0 {
		svc.Health.recordSuccess(svc.Name)
	} else {
		svc.Health.recordFailure(svc.Name, svc.FailThresh, totalFails)
	}
}
