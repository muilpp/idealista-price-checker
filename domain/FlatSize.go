package domain

type FlatSize struct {
	min int
	max int
}

func NewFlatSize(minSize int, maxSize int) *FlatSize {
	return &FlatSize{minSize, maxSize}
}

func (f FlatSize) GetFlatSize() (int, int) {
	return f.min, f.max
}

func (f FlatSize) GetMinSize() int {
	return f.min
}

func (f FlatSize) GetMaxSize() int {
	return f.max
}
