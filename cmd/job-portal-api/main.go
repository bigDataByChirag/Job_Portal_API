package main

import (
	"fmt"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/database"
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func main() {
	// Loading the environment variables file
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("Error loading .env file")
	}

	// Create a new Chi router
	r := chi.NewRouter()

	// Use custom middleware for HTTP request logging
	r.Use(middleware.HttpLogger)

	cfg := database.DefaultPostgresConfig()
	db, err := database.Open(cfg)
	if err != nil {
		fmt.Println("not connected")
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("database Connected")
	defer db.Close()

	// Set up user service
	us, err := services.NewUserService(db)
	if err != nil {
		log.Panic(err)
	}

	// Set up company service
	cs, err := services.NewCompanyService(db)
	if err != nil {
		log.Panic(err)
	}

	// Set up job service
	js, err := services.NewJobService(db)
	if err != nil {
		log.Panic(err)
	}

	// Setup authentication using RSA keys
	privatePem, err := os.ReadFile("private.pem")
	if err != nil {
		log.Panic(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		log.Panic(err)
	}

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		log.Panic("not able to read pem file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		log.Panic(err)
	}

	a, err := auth.NewAuth(publicKey, privateKey)
	if err != nil {
		log.Panic(err)
	}

	// Setup middleware using the authentication service
	m, err := middleware.NewMid(a)
	if err != nil {
		log.Panic(err)
	}

	// Create handlers for user, company, and job operations
	usersC, err := handlers.NewUsers(us, a)
	if err != nil {
		log.Panic(err)
	}
	companyC, err := handlers.NewCompany(cs, a)
	if err != nil {
		log.Panic(err)
	}
	jobC, err := handlers.NewJob(js, a)
	if err != nil {
		log.Panic(err)
	}
	
	r.Post("/api/register", usersC.CreateUser)

	r.Post("/api/login", usersC.ProcessLoginIn)

	r.Post("/api/companies", m.JWTMiddlewareCookie(companyC.CreateCompany, auth.Admin))

	r.Get("/api/companies/user", m.JWTMiddlewareCookie(companyC.GetCompanyByUserID, auth.Admin))

	r.Get("/api/companies", m.JWTMiddlewareCookie(companyC.GetAllCompanies, auth.User))

	r.Get("/api/companies/{id}", m.JWTMiddlewareCookie(companyC.GetCompanyByID, auth.User))

	r.Post("/api/companies/{id}/jobs", m.JWTMiddlewareCookie(jobC.CreateJob, auth.Admin))

	r.Get("/api/companies/{id}/jobs", m.JWTMiddlewareCookie(jobC.GetJobByCompanyID, auth.User))

	r.Delete("/api/companies/user/{id}", m.JWTMiddlewareCookie(companyC.DeleteCompanyByUserID, auth.Admin))

	r.Patch("/api/companies/user/{id}", m.JWTMiddlewareCookie(companyC.UpdateCompanyByUserID, auth.Admin))

	r.Get("/api/jobs", m.JWTMiddlewareCookie(jobC.GetAllJob, auth.User))

	r.Get("/api/jobs/{id}", m.JWTMiddlewareCookie(jobC.GetJobByID, auth.User))

	r.Delete("/api/jobs/user/{id}", m.JWTMiddlewareCookie(jobC.DeleteJobByUserID, auth.Admin))

	r.Patch("/api/jobs/user/{id}", m.JWTMiddlewareCookie(jobC.UpdateJobByUserID, auth.Admin))

	http.ListenAndServe(":3030", r)
}
