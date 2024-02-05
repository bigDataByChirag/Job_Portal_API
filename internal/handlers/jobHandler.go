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

// Job struct represents the handler for job-related operations
type Job struct {
	jobService *services.JobService
	a          *auth.Auth
}

// NewJob creates a new Job handler with the provided services and authentication
func NewJob(js *services.JobService, a *auth.Auth) (*Job, error) {
	if js == nil || a == nil {
		return nil, errors.New("please provide all the values")
	}
	return &Job{
		jobService: js,
		a:          a,
	}, nil
}

// CreateJob handles the creation of a new job
func (j Job) CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract company ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	companyID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Decode the request body into a new job model
	var newJob models.NewJob
	err = json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the new job model
	validate := validator.New()
	if err := validate.Struct(newJob); err != nil {
		log.Error().Err(err).Send()
		sendErrorResp(w, "send valid values", http.StatusBadRequest)
		return
	}

	// Create the job using the job service
	jobRole := newJob.JobRole
	salary := newJob.Salary

	_, err = j.jobService.CreateJob(jobRole, salary, companyID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "something went wrong in creating job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Job created successfully")
}

// GetJobByCompanyID handles the retrieval of jobs by company ID
func (j Job) GetJobByCompanyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract company ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	companyID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get jobs by company ID using the job service
	jobs, err := j.jobService.GetJobsByCompaniesID(companyID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not get jobs by company id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
}

// GetAllJob handles the retrieval of all jobs
func (j Job) GetAllJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all jobs using the job service
	jobs, err := j.jobService.GetAllJobs()
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
}

// GetJobByID handles the retrieval of a job by ID
func (j Job) GetJobByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract job ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get the job by ID using the job service
	job, err := j.jobService.GetJobsByID(jobID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not get jobs by this id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
}

// sendErrorResp is a helper function to send error responses with a custom message and status code
func sendErrorResp(w http.ResponseWriter, msg string, statusCode int) {
	errorMsg := struct {
		Msg string `json:"msg"`
	}{
		Msg: msg,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorMsg)
}

// DeleteJobByUserID handles the deletion of a job by user ID
func (j Job) DeleteJobByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract job ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(idStr)
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

	// Delete the job by user ID using the job service
	err = j.jobService.DeleteJobsByUserID(userID, jobID)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not delete job by user id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("job deleted successfully")
}

// UpdateJobByUserID handles the updating of a job by user ID
func (j Job) UpdateJobByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract job ID from the URL parameter
	idStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "invalid job id", http.StatusBadRequest)
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

	// Perform the update using the job service
	err = j.jobService.UpdateJobByUserID(userID, jobID, updates)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "could not update job by user id and job id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("job updated successfully")
}
