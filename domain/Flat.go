package domain

type FlatList struct {
	Flats []Flat `json:"elementList"`
}

type Flat struct {
	Location  string   `json:"-"`
	AreaId    int      `json:"-"`
	Price     float64  `json:"price"`
	AreaPrice float64  `json:"priceByArea"`
	Date      string   `json:"date"`
	Size      FlatSize `json:"-"`
}

func NewFlat(areaId int, totalPrice float64, areaPrice float64, size FlatSize) *Flat {
	f := new(Flat)
	f.AreaId = areaId
	f.Price = totalPrice
	f.AreaPrice = areaPrice
	f.Size = size
	return f
}

func NewFlatWithDate(location string, totalPrice float64, areaPrice float64, date string, size FlatSize) *Flat {
	f := new(Flat)
	f.Location = location
	f.Price = totalPrice
	f.AreaPrice = areaPrice
	f.Date = date
	f.Size = size
	return f
}
