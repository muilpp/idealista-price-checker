package domain

type FlatList struct {
	Flats []Flat `json:"elementList"`
}

type Flat struct {
	Price     float64 `json:"price"`
	AreaPrice float64 `json:"priceByArea"`
}

func NewFlat(totalPrice, areaPrice float64) *Flat {
	f := new(Flat)
	f.Price = totalPrice
	f.AreaPrice = areaPrice
	return f
}
