package models

type Model struct {
	Features      int       `json:"J"`
	Fields        int       `json:"F"`
	LatentFactors int       `json:"K"`
	W             []float64 `json:"W"`
	Normalization bool      `json:"N"`
}
