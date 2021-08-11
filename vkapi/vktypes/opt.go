package vktypes

// CropPhoto struct.
type CropPhoto struct {
	Photo Photo `json:"photo"`
	Crop  struct {
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
		X2 float64 `json:"x2"`
		Y2 float64 `json:"y2"`
	} `json:"crop"`
}
