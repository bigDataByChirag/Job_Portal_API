# Job Portal API

## Overview

The Job Portal API is a backend system for managing user accounts, companies, and job postings. It provides a set of endpoints for user registration, login, company creation, job posting, and various other functionalities related to job management. The API is built in Go and utilizes the Chi router for routing HTTP requests.

## Features

### User Management
- **User Registration**: Endpoint to allow users to register for an account.
- **User Login**: Endpoint for user authentication and login.

### Company Management
- **Create Company**: Allows authorized users (admin) to create a new company.
- **Get Company by UserID**: Retrieves company details associated with a specific user ID.
- **Get All Companies**: Retrieves a list of all companies.
- **Get Company by ID**: Retrieves details of a specific company by its ID.
- **Update Company by UserID**: Allows authorized users (admin) to update company details associated with a specific user ID.
- **Delete Company by UserID**: Allows authorized users (admin) to delete a company associated with a specific user ID.

### Job Management
- **Create Job**: Allows authorized users (admin) to create a new job posting for a specific company.
- **Get Job by Company ID**: Retrieves a list of job postings associated with a specific company ID.
- **Get All Jobs**: Retrieves a list of all job postings.
- **Get Job by ID**: Retrieves details of a specific job posting by its ID.
- **Update Job by UserID**: Allows authorized users (admin) to update job details associated with a specific user ID.
- **Delete Job by UserID**: Allows authorized users (admin) to delete a job posting associated with a specific user ID.

## Authentication and Authorization

The API implements RSA key-based authentication using JWT (JSON Web Tokens). Public and private keys are used to sign and verify the tokens. Middleware is applied to specific routes to enforce role-based access control, allowing certain operations only for authorized users (admin).

## Database

The API connects to a PostgreSQL database to store and retrieve user, company, and job data. The connection is established using the database/sql package.

## Middleware

Custom middleware is implemented for HTTP request logging and JWT validation. JWTMiddlewareCookie ensures that certain routes are accessible only with a valid JWT cookie, enforcing authentication and authorization.

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/your-username/job-portal-api.git
cd job-portal-api
```

2. Set up environment variables:

   Create a `.env` file and specify the required environment variables, such as database connection details and RSA key file paths.

3. Build and run the application:

```bash
go build
./job-portal-api
```

The API will be accessible at `http://localhost:3030`.

## Dependencies

- [Chi Router](https://github.com/go-chi/chi): Lightweight and flexible HTTP router for Go.
- [Golang JWT](https://github.com/golang-jwt/jwt): JSON Web Token implementation for Go.
- [Joho Godotenv](https://github.com/joho/godotenv): GoDotEnv loads environment variables from a .env file.

## Contributors

- Chirag Singh

## License

This project is licensed under the [MIT License](LICENSE).
