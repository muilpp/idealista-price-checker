package ports

import "idealista/domain"

type FlatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string, bool) [][]domain.Flat
}
