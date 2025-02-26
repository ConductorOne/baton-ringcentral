package client

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Generic structures -->

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

// <-- Generic structures

// Extension Response Structures -->

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

// <-- Extension Response Structures

// Role Response Structures -->

type RoleResponse struct {
	Uri        string     `json:"uri,omitempty"`
	Records    []Role     `json:"records,omitempty"`
	Paging     Paging     `json:"paging,omitempty"`
	Navigation Navigation `json:"navigation,omitempty"`
}

type Role struct {
	Uri            string `json:"uri,omitempty"`
	Id             string `json:"id,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	Description    string `json:"description,omitempty"`
	Custom         bool   `json:"custom,omitempty"`
	Scope          string `json:"scope,omitempty"`
	Hidden         bool   `json:"hidden,omitempty"`
	SiteCompatible bool   `json:"siteCompatible,omitempty"`
}

// <-- Role Response Structures

// Role Per User Response Structures -->

type UserRoleResponse struct {
	Records []UserRole `json:"records,omitempty"`
}

type UserRole struct {
	Id             string `json:"id,omitempty"`
	AutoAssigned   bool   `json:"autoAssigned,omitempty"`
	SiteRestricted bool   `json:"siteRestricted,omitempty"`
	SiteCompatible bool   `json:"siteCompatible,omitempty"`
}

// <-- Role Per User Response Structures
