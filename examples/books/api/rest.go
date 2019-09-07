package api

type CreateBookRequest struct {
	Title       string `json:"title"`
	AuthorID    string `json:"author_id"`
	PublisherID string `json:"publisher_id"`
	PublishedOn string `json:"published_on"`
	Pages       int    `json:"pages"`
}

type ListBooksResponse struct {
	Books []*Book `json:"books"`
}
