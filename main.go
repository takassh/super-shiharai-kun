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
	"github.com/takassh/super-shiharai-kun/handler"
	repo "github.com/takassh/super-shiharai-kun/repository"
	"github.com/takassh/super-shiharai-kun/util"
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
	util.LoadErrors()
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
	e.GET("/api/login/company/1", handler.LoginCompany1)
	e.GET("/api/login/company/2", handler.LoginCompany2)

	// JWT middleware setup
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.Claims)
		},
		SigningKey: []byte(os.Getenv(util.SIGNNING_KEY)),
	}

	// invoices
	i := e.Group("/api/invoices")
	i.Use(echojwt.WithConfig(config))
	i.POST("", handler.CreateInvoice)
	i.GET("", handler.GetInvoices)

	e.Logger.Fatal(e.Start(":8080"))

}
