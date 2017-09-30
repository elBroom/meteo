package schema

type MeteoData struct {
	Color string          `json:"color,omitempty"`
	Name  string          `json:"name"`
	Data  [][]interface{} `json:"data"`
	Pin   string          `json:"pin"`
	Unit  string          `json:"unit,omitempty"`
}
