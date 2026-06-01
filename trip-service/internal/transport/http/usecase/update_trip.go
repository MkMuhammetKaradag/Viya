package usecase

import (
	"context"
	"errors"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

// Kontrolcüde id beklediğin için dönüş tipini (uuid.UUID, error) yaptık usta
type UpdateTripUseCase interface {
	Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error)
}

type updateTripUseCase struct {
	repo domain.TripRepository
}

func NewUpdateTripUseCase(repo domain.TripRepository) UpdateTripUseCase {
	return &updateTripUseCase{
		repo: repo,
	}
}

func (u *updateTripUseCase) Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	// 1. Veritabanından mevcut seyahat kaydını çekiyoruz
	existingTrip, err := u.repo.GetTripByID(ctx, trip.ID)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. 🛡️ GÜVENLİK DUVARI: Bu seyahat gerçekten bu kullanıcıya mı ait?
	if existingTrip.UserID != trip.UserID {
		return uuid.Nil, errors.New("unauthorized: you do not own this trip")
	}

	// 3. 🔄 HİBRİT MERGE OPERASYONU
	// Sadece kontrolcüden (dolayısıyla ön yüzden) içi dolu gelen alanları eziyoruz.
	// Gelmeyen alanlarda veritabanındaki eski (existingTrip) değerler korunuyor.

	if trip.Title != "" {
		existingTrip.Title = trip.Title
	}
	if trip.Description != nil {
		existingTrip.Description = trip.Description
	}
	if trip.CoverImageURL != nil {
		existingTrip.CoverImageURL = trip.CoverImageURL
	}

	// Boolean alanlar için kontrolcüde pointer mantığı kurmuştuk.
	// Eğer ilgili alanlar domain modeline taşınırken kontrolcüde set edildiyse
	// (yani sıfır değer sorunu bypass edildiyse) doğrudan ezebiliriz.
	existingTrip.IsPublic = trip.IsPublic
	existingTrip.IsActive = trip.IsActive

	// Kategori ilişkilerini güncelleme
	if len(trip.CategoryIDs) > 0 {
		existingTrip.CategoryIDs = trip.CategoryIDs
	}

	// PublishedAt mantığı: Kontrolcüde sıfırlanmasını engellemiştik,
	// gelen yeni bir tarih varsa güncelliyoruz.
	if !trip.PublishedAt.IsZero() {
		existingTrip.PublishedAt = trip.PublishedAt
	}

	// 4. Güncellenmiş nesneyi veritabanına yazması için repository'e paslıyoruz
	err = u.repo.UpdateTrip(ctx, existingTrip)
	if err != nil {
		return uuid.Nil, err
	}

	// 5. Başarılıysa güncellenen seyahatin ID'sini dönüyoruz
	return existingTrip.ID, nil
}
