package api

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

type Publisher struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	PublishedOn string    `json:"published_on"`
	Author      Author    `json:"author"`
	Publisher   Publisher `json:"publisher"`
}
