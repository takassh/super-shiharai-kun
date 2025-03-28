package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	e "github.com/takassh/super-shiharai-kun/entity"
	repo "github.com/takassh/super-shiharai-kun/repository"
	req "github.com/takassh/super-shiharai-kun/request"
	"github.com/takassh/super-shiharai-kun/util"
)

var (
	factory     = repo.NewRepositoryFactory(repo.DialectSQLite, true)
	companyRepo = factory.NewCompanyRepository()
	userRepo    = factory.NewUserRepository()
	invoiceRepo = factory.NewInvoiceRepository()
)

func LoginCompany1(c echo.Context) error {
	// always login as admin and user 1
	claims := &util.Claims{
		"1",
		"1",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv(util.SIGNNING_KEY)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorCodes.GetError("500-000", err))
	}

	return c.String(http.StatusOK, t)
}

func LoginCompany2(c echo.Context) error {
	// always login as user 2
	claims := &util.Claims{
		"2",
		"2",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv(util.SIGNNING_KEY)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorCodes.GetError("500-000", err))
	}

	return c.String(http.StatusOK, t)
}

func CreateInvoice(c echo.Context) error {
	_user := c.Get("user").(*jwt.Token)
	claims := _user.Claims.(*util.Claims)
	userID, _ := strconv.Atoi(claims.UserID)
	user, err := userRepo.FindByID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorCodes.GetError("400-001", err))
	}

	invoiceReq := new(req.CreateInvoice)
	if err := req.BindAndValidate(c, invoiceReq); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorCodes.GetError("400-002", err))
	}

	if err := invoiceReq.ValidateValue(); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorCodes.GetError("400-006", err))
	}

	if user.CompanyID == invoiceReq.ClientID {
		return c.JSON(http.StatusBadRequest, util.ErrorCodes.GetError("400-003", nil))
	}

	invoice, err := invoiceReq.ToEntity(user.CompanyID, e.StatusUnprocessed, util.FEE_RATE, util.TAX_RATE)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorCodes.GetError("400-004", err))
	}

	if err := invoiceRepo.Save(invoice); err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorCodes.GetError("500-000", err))
	}

	return c.JSON(http.StatusOK, invoice)
}

func GetInvoices(c echo.Context) error {
	_user := c.Get("user").(*jwt.Token)
	claims := _user.Claims.(*util.Claims)
	companyID, _ := strconv.Atoi(claims.CompanyID)
	if _, err := companyRepo.FindByID(uint(companyID)); err != nil {
		return c.JSON(http.StatusNotFound, util.ErrorCodes.GetError("400-005", err))
	}

	startDate := c.QueryParam("start_date")
	dueDate := c.QueryParam("due_date")

	invoices, err := invoiceRepo.FindByDate(uint(companyID), startDate, dueDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorCodes.GetError("500-001", err))
	}

	return c.JSON(http.StatusOK, invoices)
}
