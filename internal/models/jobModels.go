package models

// NewJob represents the structure for creating a new job. It includes fields for the job role and salary.
type NewJob struct {
	JobRole string `json:"jobRole" validate:"required"` // JobRole is the role or title of the job and is required.
	Salary  int    `json:"salary" validate:"required"`  // Salary is the salary associated with the job and is required.
}

// Job represents the structure for a job entity. It includes fields such as ID, job role, salary, and CompanyId.
type Job struct {
	ID        int    `json:"id"`        // ID is a unique identifier for the job.
	JobRole   string `json:"jobRole"`   // JobRole is the role or title of the job.
	Salary    int    `json:"salary"`    // Salary is the salary associated with the job.
	CompanyId int    `json:"companyId"` // CompanyId is the identifier of the company associated with the job.
}
