package schema

type Indication struct {
	Value     float32 `json:"value"`
	Pin       string  `json:"pin,omitempty"`
	CreatedAt int64   `json:"create_date,omitempty"`
}
