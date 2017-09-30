package schema

type Indication struct {
	Value      float32 `json:"value"`
	Pin        string  `json:"pin,omitempty"`
	CreateDate int64   `json:"create_date,omitempty"`
}
