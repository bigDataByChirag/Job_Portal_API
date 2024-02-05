package services

import (
	"database/sql"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"strconv"
	"strings"
)

// JobService handles business logic related to job operations.
type JobService struct {
	db *sql.DB
}

// NewJobService creates a new JobService instance.
func NewJobService(db *sql.DB) (*JobService, error) {
	if db == nil {
		return nil, errors.New("db connection cannot be nil")
	}
	return &JobService{db: db}, nil
}

// CreateJob creates a new job record in the database.
func (js *JobService) CreateJob(jobRole string, salary int, companyId int) (*models.Job, error) {
	// Convert jobRole to lowercase
	jobRole = strings.ToLower(jobRole)

	job := models.Job{
		JobRole:   jobRole,
		Salary:    salary,
		CompanyId: companyId,
	}

	// Execute the SQL query to insert a new job and retrieve the generated ID
	row := js.db.QueryRow(`
		INSERT INTO jobs (jobRole, salary, companyId)
		VALUES ($1, $2, $3) RETURNING id`, jobRole, salary, companyId)

	err := row.Scan(&job.ID)
	if err != nil {
		fmt.Println("No Company Associated with that Id")
		return nil, fmt.Errorf("create job: %w", err)
	}
	return &job, nil
}

// GetJobsByCompaniesID retrieves all jobs associated with a company from the database.
func (js *JobService) GetJobsByCompaniesID(id int) ([]*models.Job, error) {
	var jobs []*models.Job

	// Execute the SQL query to select jobs by company ID
	rows, err := js.db.Query("SELECT * FROM Jobs WHERE companyId = $1", id)
	if err != nil {
		return nil, fmt.Errorf("get jobs by company ID: %w", err)
	}
	defer rows.Close()

	// Iterate over the result rows and populate the jobs slice
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(&job.ID, &job.JobRole, &job.Salary, &job.CompanyId); err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return jobs, nil
}

// GetAllJobs retrieves all jobs from the database.
func (js *JobService) GetAllJobs() ([]*models.Job, error) {
	var jobs []*models.Job

	// Execute the SQL query to select all jobs
	rows, err := js.db.Query("SELECT id, jobRole, salary, companyID FROM jobs")
	if err != nil {
		return nil, fmt.Errorf("get all jobs: %w", err)
	}
	defer rows.Close()

	// Iterate over the result rows and populate the jobs slice
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(&job.ID, &job.JobRole, &job.Salary, &job.CompanyId); err != nil {
			return nil, fmt.Errorf("get all jobs: %w", err)
		}
		jobs = append(jobs, &job)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all jobs: %w", err)
	}

	return jobs, nil
}

// GetJobsByID retrieves a job by its ID from the database.
func (js *JobService) GetJobsByID(id int) (*models.Job, error) {
	var job models.Job

	// Execute the SQL query to select a job by ID
	err := js.db.QueryRow("SELECT id, jobRole, salary, companyId FROM jobs WHERE id= $1", id).Scan(&job.ID, &job.JobRole, &job.Salary, &job.CompanyId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("get job by ID: %w", err)
	}
	return &job, nil
}

// DeleteJobsByUserID deletes a job associated with a user from the database.
func (js *JobService) DeleteJobsByUserID(userID, jobID int) error {
	// Check if the job exists for the given user
	var count int
	err := js.db.QueryRow("SELECT COUNT(*) FROM jobs j INNER JOIN companies c ON j.companyId = c.id WHERE c.userId = $1 AND j.id = $2", userID, jobID).Scan(&count)
	if err != nil {
		return fmt.Errorf("query job existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("job not found for user with ID %d and job ID %d", userID, jobID)
	}

	// Delete the job
	_, err = js.db.Exec("DELETE FROM jobs WHERE id = $1", jobID)
	if err != nil {
		return fmt.Errorf("delete job with ID %d: %w", jobID, err)
	}

	return nil
}

// UpdateJobByUserID updates a job associated with a user in the database.
func (js *JobService) UpdateJobByUserID(userID, jobID int, updates map[string]interface{}) error {
	// Check if the job exists for the given user
	var count int
	err := js.db.QueryRow("SELECT COUNT(*) FROM jobs j INNER JOIN companies c ON j.companyId = c.id WHERE c.userId = $1 AND j.id = $2", userID, jobID).Scan(&count)
	if err != nil {
		return fmt.Errorf("query job existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("job not found for user with ID %d and job ID %d", userID, jobID)
	}

	// Build the UPDATE query dynamically based on the fields provided in the updates map
	query := "UPDATE jobs SET "
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

	// Add the WHERE conditions to update the specific job
	query += " WHERE id = $" + strconv.Itoa(i)
	values = append(values, jobID)

	// Execute the dynamic UPDATE query
	_, err = js.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("update job with ID %d: %w", jobID, err)
	}

	return nil
}
