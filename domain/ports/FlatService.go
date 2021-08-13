package ports

import "idealista/domain"

type FlatRepository interface {
	Add([]domain.Flat, string) bool
	Get(string, bool, bool) []domain.Flat
}

type FlatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string, bool, bool) []domain.Flat
}
