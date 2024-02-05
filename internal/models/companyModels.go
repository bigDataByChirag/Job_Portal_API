package models

// NewCompany represents the structure for creating a new company. It includes fields for the company's name and address.
type NewCompany struct {
	Name    string `json:"name" validate:"required"`    // Name is the name of the company and is required.
	Address string `json:"address" validate:"required"` // Address is the physical address of the company and is required.
}

// Company represents the structure for a company entity. It includes fields such as ID, name, address, and UserId.
type Company struct {
	ID      int    `json:"id"`      // ID is a unique identifier for the company.
	Name    string `json:"name"`    // Name is the name of the company.
	Address string `json:"address"` // Address is the physical address of the company.
	UserId  int    `json:"userId"`  // UserId is the identifier of the user associated with the company.
}
