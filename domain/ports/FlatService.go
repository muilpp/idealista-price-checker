package ports

import "idealista/domain"

type FlatRepository interface {
	Add([][]domain.Flat) bool
	Get(string, bool, bool, int) []domain.Flat
	GetPlacesToSearch() []domain.Place
}

type FlatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string, bool, bool) ([]domain.Flat, []domain.Flat)
}
