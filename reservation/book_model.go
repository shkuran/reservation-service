package reservation

type Book struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	PublicationYear int64  `json:"publication_year"`
	AvailableCopies int64  `json:"available_copies"`
}