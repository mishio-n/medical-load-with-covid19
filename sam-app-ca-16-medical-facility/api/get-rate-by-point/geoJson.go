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
	// スタイル設定
	Stroke        string  `json:"stroke"`
	StrokeOpacity float64 `json:"stroke-opacity"`
	StrokeWidth   int     `json:"stroke-width"`
	Fill          string  `json:"fill"`
	FillOpacity   float64 `json:"fill-opacity"`
	MarkerColor   string  `json:"marker-color"`
	// メタデータ
	Name string  `json:"病院"`
	Type string  `json:"種別"`
	Rate float64 `json:"rate"`
}

type GeoJsonFeatureGeometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}
