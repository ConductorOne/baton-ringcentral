package client

type ExtensionResponse struct {
	Uri        string      `json:"uri,omitempty"`
	Records    []Extension `json:"records,omitempty"`
	Paging     Paging      `json:"paging,omitempty"`
	Navigation Navigation  `json:"navigation,omitempty"`
}

type Extension struct {
	ID          int64            `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Type        string           `json:"type,omitempty"`
	Status      string           `json:"status,omitempty"`
	ContactInfo ExtensionContact `json:"contact,omitempty"`
}

type ExtensionContact struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}

type Paging struct {
	Page       int `json:"page,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
	PageEnd    int `json:"pageEnd,omitempty"`
}

type Navigation struct {
	FirstPage NavPage `json:"firstPage,omitempty"`
	LastPage  NavPage `json:"lastPage,omitempty"`
}

type NavPage struct {
	Uri string `json:"uri,omitempty"`
}
