package main

import (
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	repo "github.com/takassh/super-shiharai-kun/repository"
)

// hard coded values
const (
	FEE_RATE     = 0.04
	TAX_RATE     = 0.10
	SIGNNING_KEY = "SIGNING_KEY"
)

var (
	factory     = repo.NewRepositoryFactory(repo.DialectSQLite, true)
	companyRepo = factory.NewCompanyRepository()
	userRepo    = factory.NewUserRepository()
	bankRepo    = factory.NewBankAccountRepository()
	invoiceRepo = factory.NewInvoiceRepository()
)

func loadEnv() {
	godotenv.Load(".env")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	loadErrors()
	loadEnv()

	// initialize data (this is just for testing)
	repo.PopulateCompany(companyRepo)
	repo.PopulateUser(userRepo)
	repo.PopulateClient(companyRepo)
	repo.PopulateBankAccount(bankRepo)
	repo.PopulateInvoice(invoiceRepo)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// login
	e.GET("/api/login/company/1", LoginCompany1)
	e.GET("/api/login/company/2", LoginCompany2)

	// JWT middleware setup
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(Claims)
		},
		SigningKey: []byte(os.Getenv("SIGNING_KEY")),
	}

	// invoices
	i := e.Group("/api/invoices")
	i.Use(echojwt.WithConfig(config))
	i.POST("", CreateInvoice)
	i.GET("", GetInvoices)

	e.Logger.Fatal(e.Start(":8080"))

}
