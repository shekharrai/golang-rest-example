package book

type Book struct {
	Id     int     `json:"id"`
	Name   string  `json:"name, omitempty"`
	Author string  `json:"author, omitempty"`
	Price  float32 `json:"price, omitempty"`
}
