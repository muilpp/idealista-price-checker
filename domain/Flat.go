package domain

type FlatList struct {
	Flats []Flat `json:"elementList"`
}

type Flat struct {
	Price     float64 `json:"price"`
	AreaPrice float64 `json:"priceByArea"`
	Date      string  `json:"date"`
}

func NewFlat(totalPrice, areaPrice float64) *Flat {
	f := new(Flat)
	f.Price = totalPrice
	f.AreaPrice = areaPrice
	return f
}

func NewFlatWithDate(totalPrice, areaPrice float64, date string) *Flat {
	f := new(Flat)
	f.Price = totalPrice
	f.AreaPrice = areaPrice
	f.Date = date
	return f
}
