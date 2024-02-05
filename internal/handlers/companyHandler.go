package handlers

import (
	"encoding/json"
	"errors"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Company struct represents the handler for company-related operations
type Company struct {
	companyService *services.CompanyService
	a              *auth.Auth
}

// NewCompany creates a new Company handler with the provided services and authentication
func NewCompany(cs *services.CompanyService, a *auth.Auth) (*Company, error) {
	if cs == nil || a == nil {
		return nil, errors.New("please provide all the values")
	}
	return &Company{
		companyService: cs,
		a:              a,
	}, nil
}

// CreateCompany handles the creation of a new company
func (c Company) CreateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a new company model
	var newCompany models.NewCompany
	err := json.NewDecoder(r.Body).Decode(&newCompany)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the new company model
	validate := validator.New()
	if err := validate.Struct(newCompany); err != nil {
		log.Error().Err(err).Send()
		appErr := HandlerError{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	// Extract user ID from the request context
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		log.Error().Err(err).Send()
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid user id in context", http.StatusUnauthorized)
		return
	}

	// Create the company using the company service
	name := newCompany.Name
	address := newCompany.Address
	_, err = c.companyService.CreateCompany(userID, name, address)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Company Created Successfully")
}

// GetAllCompanies handles the retrieval of all companies
func (c Company) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all companies using the company service
	companies, err := c.companyService.GetAllCompanies()
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

// GetCompanyByID handles the retrieval of a company by ID
func (c Company) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract company ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get the company by ID using the company service
	company, err := c.companyService.GetCompanyByID(id)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not get companies by id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

// GetCompanyByUserID handles the retrieval of companies by user ID
func (c Company) GetCompanyByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from the request context
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid user id in context", http.StatusUnauthorized)
		return
	}

	// Get companies by user ID using the company service
	company, err := c.companyService.GetCompaniesByUserID(userID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not get companies by user id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

// DeleteCompanyByUserID handles the deletion of a company by user ID
func (c Company) DeleteCompanyByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract company ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	compID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Extract user ID from the request context
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid user id in context", http.StatusUnauthorized)
		return
	}

	// Delete the company by user ID using the company service
	err = c.companyService.DeleteCompaniesByUserID(userID, compID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not delete company by user id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("company deleted successfully")
}

// UpdateCompanyByUserID handles the updating of a company by user ID
func (c Company) UpdateCompanyByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract company ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	compID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Extract user ID from the request context
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid user id in context", http.StatusUnauthorized)
		return
	}

	// Parse the request body to get the updates
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Perform the update
	err = c.companyService.UpdateCompaniesByUserID(userID, compID, updates)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not update company by user id and company id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("company updated successfully")
}
