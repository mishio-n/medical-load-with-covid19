package main

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type GeoJson struct {
	Type     string           `json:"type"`
	Features []GeoJsonFeature `json:"features"`
}

type GeoJsonFeature struct {
	Type       string                 `json:"type"`
	Properties GeoJsonFeatureProperty `json:"properties"`
	Geometry   GeoJsonFeatureGeometry `json:"geometry"`
}

type GeoJsonFeatureProperty struct {
	Stroke        string  `json:"stroke"`
	StrokeOpacity float64 `json:"stroke-opacity"`
	StrokeWidth   int     `json:"stroke-width"`
	Fill          string  `json:"fill"`
	FillOpacity   float64 `json:"fill-opacity"`
	MarkerColor   string  `json:"marker-color"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
}

type GeoJsonFeatureGeometry struct {
	Type        string          `json:"type"`
	Coordinates [1][][2]float64 `json:"coordinates"`
}
