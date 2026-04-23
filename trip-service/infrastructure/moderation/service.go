package moderation

import (
	"context"
	"math"
	"strings"
	"trip-service/internal/domain"
	"unicode"
)

const similarityThreshold = 0.80

var referenceTexts = []string{
	"seni öldüreceğim",
	"orospu çocuğu",
	"gerizekalı birisin",
	"seni tehdit ediyorum",
	"bok gibi",
	"defol git",
	"rezil insan",
}

var blacklist = []string{
	"sik", "orospu", "göt", "amk", "bok",
	"yarrak", "piç", "ibne", "kahpe",
	"gerizekalı", "öldürürüm", "gebertrim",
}

type Service struct {
	aiSvc           domain.AIService
	referenceEmbeds [][]float32
}

func NewModerationService(ctx context.Context, aiSvc domain.AIService) (domain.ModerationService, error) {
	svc := &Service{aiSvc: aiSvc}

	// Uygulama ayağa kalkarken referansları embed et
	for _, text := range referenceTexts {
		emb, err := aiSvc.GetVector(ctx, text)
		if err != nil {
			return nil, err
		}
		svc.referenceEmbeds = append(svc.referenceEmbeds, emb)
	}

	return svc, nil
}

func (s *Service) Moderate(ctx context.Context, content string) (*domain.ModerationResult, error) {
	// 1. Blacklist
	normalized := normalize(content)
	for _, word := range blacklist {
		if strings.Contains(normalized, word) {
			return &domain.ModerationResult{
				IsAppropriate: false,
				Reason:        "inappropriate language detected",
			}, nil
		}
	}

	// 2. Semantic similarity
	contentEmb, err := s.aiSvc.GetVector(ctx, content)
	if err != nil {
		// Ollama erişilemezse sadece blacklist sonucu yeterli
		return &domain.ModerationResult{IsAppropriate: true}, nil
	}

	for _, refEmb := range s.referenceEmbeds {
		if cosineSimilarity(contentEmb, refEmb) >= similarityThreshold {
			return &domain.ModerationResult{
				IsAppropriate: false,
				Reason:        "content flagged as inappropriate",
			}, nil
		}
	}

	return &domain.ModerationResult{IsAppropriate: true}, nil
}

func normalize(text string) string {
	replacements := map[rune]rune{
		'0': 'o', '1': 'i', '3': 'e',
		'4': 'a', '5': 's', '@': 'a',
		'!': 'i', '$': 's',
	}
	text = strings.ToLower(text)
	var sb strings.Builder
	for _, r := range text {
		if mapped, ok := replacements[r]; ok {
			sb.WriteRune(mapped)
		} else if unicode.IsLetter(r) || unicode.IsSpace(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func cosineSimilarity(a, b []float32) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
