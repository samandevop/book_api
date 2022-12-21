package models

type BookPrimarKey struct {
	Id string `json:"book_id"`
}

type CreateBook struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
type Book struct {
	Id        string  `json:"book_id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UpdateBookSwagger struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

type UpdateBook struct {
	Id     string  `json:"book_id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

type GetListBookRequest struct {
	Limit  int32
	Offset int32
}

type GetListBookResponse struct {
	Count int32   `json:"count"`
	Books []*Book `json:"books"`
}
