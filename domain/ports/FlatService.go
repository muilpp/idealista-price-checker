package ports

import "idealista/domain"

type FlatRepository interface {
	Add([]domain.Flat, string, int) bool
	Get(string, bool, bool, int) []domain.Flat
}

type FlatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string, bool, bool) ([]domain.Flat, []domain.Flat)
}
