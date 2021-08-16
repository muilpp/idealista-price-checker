package domain

type Place struct {
	id        int
	name      string
	center    string
	distance  int
	min_size  int
	max_size  int
	bedrooms  string
	terrace   bool
	operation string
}

func NewPlace(id int, name string, center string, distance int, min_size int, max_size int, bedrooms string, terrace bool, operation string) *Place {
	return &Place{id, name, center, distance, min_size, max_size, bedrooms, terrace, operation}
}

func (p Place) GetId() int {
	return p.id
}

func (p Place) GetName() string {
	return p.name
}

func (p Place) GetCenter() string {
	return p.center
}

func (p Place) GetDistance() int {
	return p.distance
}

func (p Place) GetMinSize() int {
	return p.min_size
}

func (p Place) GetMaxSize() int {
	return p.max_size
}

func (p Place) GetBedrooms() string {
	return p.bedrooms
}

func (p Place) HasTerrace() bool {
	return p.terrace
}

func (p Place) GetOperation() string {
	return p.operation
}
