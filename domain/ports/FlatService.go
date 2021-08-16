package ports

import "idealista/domain"

type FlatRepository interface {
	Add([][]domain.Flat) bool
	Get(string, bool) [][]domain.Flat
	GetPlacesToSearch() []domain.Place
}

type FlatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string, bool) [][]domain.Flat
}
