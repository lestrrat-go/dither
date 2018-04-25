package dither

func NewFilter(name string, m *Matrix) *Filter {
	return &Filter{
		name:   name,
		matrix: m,
	}
}

func (f *Filter) Name() string {
	return f.name
}

func (f *Filter) Matrix() *Matrix {
	return f.matrix
}

var FloydSteinberg = NewFilter(
	"FloydSteinberg",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0}).
		AddRow([]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0}).
		AddRow([]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0}).
		Build(),
)

var Stucki = NewFilter(
	"Stucki",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0}).
		AddRow([]float32{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0}).
		AddRow([]float32{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0}).
		Build(),
)

var Atkinson = NewFilter(
	"Atkinson",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 1.0 / 8.0, 1.0 / 8.0}).
		AddRow([]float32{1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0.0}).
		AddRow([]float32{0.0, 1.0 / 8.0, 0.0, 0.0}).
		Build(),
)

var Burkes = NewFilter(
	"Burkes",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0}).
		AddRow([]float32{2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}).
		AddRow([]float32{0.0, 0.0, 0.0, 0.0, 0.0}).
		Build(),
)

var Sierra3 = NewFilter(
	"Sierra-3",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0}).
		AddRow([]float32{2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}).
		AddRow([]float32{0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0}).
		Build(),
)

var Sierra2 = NewFilter(
	"Sierra-2",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 0.0, 4.0 / 16.0, 3.0 / 16.0}).
		AddRow([]float32{1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0}).
		AddRow([]float32{0.0, 0.0, 0.0, 0.0, 0.0}).
		Build(),
)

var SierraLite = NewFilter(
	"Sierra-Lite",
	NewMatrixBuilder(5, 3).
		AddRow([]float32{0.0, 0.0, 2.0 / 4.0}).
		AddRow([]float32{1.0 / 4.0, 1.0 / 4.0, 0.0}).
		AddRow([]float32{0.0, 0.0, 0.0}).
		Build(),
)
