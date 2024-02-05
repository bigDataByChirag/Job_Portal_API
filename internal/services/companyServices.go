package services

import (
	"database/sql"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"log"
	"strconv"
	"strings"
)

// CompanyService handles business logic related to company operations.
type CompanyService struct {
	db *sql.DB
}

// NewCompanyService creates a new CompanyService instance.
func NewCompanyService(db *sql.DB) (*CompanyService, error) {
	if db == nil {
		return nil, errors.New("db connection cannot be nil")
	}
	return &CompanyService{db: db}, nil
}

// CreateCompany creates a new company record in the database.
func (cs *CompanyService) CreateCompany(userId int, name, address string) (*models.Company, error) {
	// Convert name and address to lowercase
	name = strings.ToLower(name)
	address = strings.ToLower(address)

	company := models.Company{
		Name:    name,
		Address: address,
		UserId:  userId,
	}

	// Execute the SQL query to insert a new company and retrieve the generated ID
	row := cs.db.QueryRow(`
		INSERT INTO companies (name, address, userId)
		VALUES ($1, $2, $3) RETURNING id`, name, address, userId)

	err := row.Scan(&company.ID)
	if err != nil {
		return nil, fmt.Errorf("create company: %w", err)
	}
	return &company, nil
}

// GetAllCompanies retrieves all companies from the database.
func (cs *CompanyService) GetAllCompanies() ([]*models.Company, error) {
	fmt.Println("**************")
	var companies []*models.Company

	// Execute the SQL query to select all companies
	rows, err := cs.db.Query("SELECT id, name, address, userid FROM companies")
	if err != nil {
		return nil, fmt.Errorf("get all companies: %w", err)
	}
	defer rows.Close()

	// Iterate over the result rows and populate the companies slice
	for rows.Next() {
		log.Println("**************")
		var company models.Company
		if err := rows.Scan(&company.ID, &company.Name, &company.Address, &company.UserId); err != nil {
			return nil, fmt.Errorf("get all companies: %w", err)
		}
		log.Println("**************", company)
		companies = append(companies, &company)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all companies: %w", err)
	}

	return companies, nil
}

// GetCompanyByID retrieves a company by its ID from the database.
func (cs *CompanyService) GetCompanyByID(id int) (*models.Company, error) {
	var company models.Company

	// Execute the SQL query to select a company by ID
	err := cs.db.QueryRow("SELECT * FROM companies WHERE id= $1", id).Scan(&company.ID, &company.Name, &company.Address, &company.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("get company by ID: %w", err)
	}
	return &company, nil
}

// GetCompaniesByUserID retrieves all companies associated with a user from the database.
func (cs *CompanyService) GetCompaniesByUserID(userID int) ([]*models.Company, error) {
	// Execute the SQL query to select companies by user ID
	rows, err := cs.db.Query("SELECT * FROM companies WHERE userId = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("query companies by user ID: %w", err)
	}
	defer rows.Close()

	var companies []*models.Company

	// Iterate over the result rows and populate the companies slice
	for rows.Next() {
		var company models.Company
		err := rows.Scan(&company.ID, &company.Name, &company.Address, &company.UserId)
		if err != nil {
			return nil, fmt.Errorf("scan company row: %w", err)
		}
		companies = append(companies, &company)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate over rows: %w", err)
	}

	return companies, nil
}

// DeleteCompaniesByUserID deletes a company associated with a user from the database.
func (cs *CompanyService) DeleteCompaniesByUserID(userID, companyID int) error {
	// Check if the company exists for the given user
	var count int
	err := cs.db.QueryRow("SELECT COUNT(*) FROM companies WHERE userId = $1 AND id = $2", userID, companyID).Scan(&count)
	if err != nil {
		return fmt.Errorf("query company existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("company not found for user with ID %d and company ID %d", userID, companyID)
	}

	// Delete the company
	_, err = cs.db.Exec("DELETE FROM companies WHERE userId = $1 AND id = $2", userID, companyID)
	if err != nil {
		return fmt.Errorf("delete company with ID %d: %w", companyID, err)
	}

	return nil
}

// UpdateCompaniesByUserID updates a company associated with a user in the database.
func (cs *CompanyService) UpdateCompaniesByUserID(userID, companyID int, updates map[string]interface{}) error {
	// Check if the company exists for the given user
	var count int
	err := cs.db.QueryRow("SELECT COUNT(*) FROM companies WHERE userId = $1 AND id = $2", userID, companyID).Scan(&count)
	if err != nil {
		return fmt.Errorf("query company existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("company not found for user with ID %d and company ID %d", userID, companyID)
	}

	// Build the UPDATE query dynamically based on the fields provided in the updates map
	query := "UPDATE companies SET "
	values := []interface{}{}

	i := 1
	for key, value := range updates {
		// Append each field to the query
		query += key + " = $" + strconv.Itoa(i) + ", "
		values = append(values, value)
		i++
	}

	// Remove the trailing comma and space from the query
	query = query[:len(query)-2]

	// Add the WHERE conditions to update the specific company
	query += " WHERE userId = $" + strconv.Itoa(i) + " AND id = $" + strconv.Itoa(i+1)
	values = append(values, userID, companyID)

	// Execute the dynamic UPDATE query
	_, err = cs.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("patch company with ID %d: %w", companyID, err)
	}

	return nil
}
