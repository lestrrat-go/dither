package dither

type Matrix struct {
	cols   int // horizontal
	rows   int // vertical
	values []float32
}

func NewMatrix(cols, rows int) *Matrix {
	return &Matrix{
		cols:   cols,
		rows:   rows,
		values: make([]float32, cols*rows, cols*rows),
	}
}

func (m *Matrix) Get(x, y int) float32 {
	return m.values[x+m.cols*y]
}

func (m *Matrix) Set(x, y int, v float32) {
	m.values[x+m.cols*y] = v
}

func (m *Matrix) SetRow(y int, v []float32) *Matrix {
	copy(m.values[m.cols*y:], v)
	return m
}

func (m *Matrix) Cols() int {
	return m.cols
}

func (m *Matrix) Rows() int {
	return m.rows
}

type MatrixBuilder struct {
	currow int
	matrix *Matrix
}

func NewMatrixBuilder(cols, rows int) *MatrixBuilder {
	return &MatrixBuilder{
		matrix: NewMatrix(cols, rows),
	}
}

func (b *MatrixBuilder) AddRow(v []float32) *MatrixBuilder {
	b.matrix.SetRow(b.currow, v)
	b.currow++
	return b
}

func (b *MatrixBuilder) Build() *Matrix {
	return b.matrix
}
