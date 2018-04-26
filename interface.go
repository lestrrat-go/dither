package dither

type Filter struct {
	name   string
	matrix *Matrix
}

type Matrixer interface {
	Matrix() *Matrix
}
