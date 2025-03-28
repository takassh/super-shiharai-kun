package repository

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	e "github.com/takassh/super-shiharai-kun/entity"
	util "github.com/takassh/super-shiharai-kun/util"
)

func PopulateCompany(repo CompanyRepository) {
	csvPath := "data/company.csv"
	rows := returnCSV(csvPath)
	// id,name,ceo,phone,postal_code,address
	for i, row := range rows {
		if i == 0 {
			continue
		}

		company := e.Company{
			Name:       row[1],
			CEO:        row[2],
			Phone:      row[3],
			PostalCode: row[4],
			Address:    row[5],
		}

		if err := repo.Save(company); err != nil {
			message := fmt.Sprintf("Failed to save company: %s", company.Name)
			log.Fatal(message)
		}
	}
}

func PopulateUser(repo UserRepository) {
	csvPath := "data/user.csv"
	rows := returnCSV(csvPath)
	// id,company_id,name,email,password
	for i, row := range rows {
		if i == 0 {
			continue
		}

		user := e.User{
			Name:     row[0],
			Email:    row[1],
			Password: row[2],
		}

		companyID, err := strconv.Atoi(row[1])
		if err != nil {
			message := fmt.Sprintf("Failed to convert company ID: %s", row[1])
			log.Fatal(message)
		}
		user.CompanyID = uint(companyID)

		if err := repo.Save(user); err != nil {
			message := fmt.Sprintf("Failed to save user: %s", user.Name)
			log.Fatal(message)
		}
	}
}

func PopulateClient(repo CompanyRepository) {
	csvPath := "data/client.csv"
	rows := returnCSV(csvPath)
	// id,company_id,name,ceo,phone,postal_code,address
	for i, row := range rows {
		if i == 0 {
			continue
		}

		client := e.Company{
			Name:       row[2],
			CEO:        row[3],
			Phone:      row[4],
			PostalCode: row[5],
			Address:    row[6],
		}

		companyID, err := strconv.Atoi(row[1])
		if err != nil {
			message := fmt.Sprintf("Failed to convert company ID: %s", row[1])
			log.Fatal(message)
		}

		company, err := repo.FindByID(uint(companyID))
		if err != nil {
			message := fmt.Sprintf("Failed to find company: %s", row[1])
			log.Fatal(message)
		}

		if err := repo.AppendClient(company, client); err != nil {
			message := fmt.Sprintf("Failed to append client: %s", client.Name)
			log.Fatal(message)
		}
	}
}

func PopulateBankAccount(repo BankAccountRepository) {
	csvPath := "data/bank_account.csv"
	rows := returnCSV(csvPath)
	// id,client_id,bank_name,branch_name,account_number
	for i, row := range rows {
		if i == 0 {
			continue
		}

		bankAccount := e.BankAccount{
			BankName:    row[2],
			Branch:      row[3],
			AccountNo:   row[4],
			AccountName: row[5],
		}

		companyID, err := strconv.Atoi(row[1])
		if err != nil {
			message := fmt.Sprintf("Failed to convert company ID: %s", row[1])
			log.Fatal(message)
		}

		bankAccount.CompanyID = uint(companyID)

		if err := repo.Save(bankAccount); err != nil {
			message := fmt.Sprintf("Failed to append bank account: %s", bankAccount.BankName)
			log.Fatal(message)
		}
	}
}

func PopulateInvoice(repo InvoiceRepository) {
	csvPath := "data/invoice.csv"
	rows := returnCSV(csvPath)
	// id,company_id,client_id,issued_at,amount,fee,fee_rate,tax,tax_rate,total_amount,due_date,status
	for i, row := range rows {
		if i == 0 {
			continue
		}
		issueAt, err := util.Atot(row[3])
		if err != nil {
			message := fmt.Sprintf("Failed to convert issued at: %s", row[3])
			log.Fatal(message)
		}

		invoice := e.Invoice{
			IssuedAt: issueAt,
			Status:   e.Status(row[11]),
		}

		amount, err := util.Atof[float64](row[4])
		if err != nil {
			message := fmt.Sprintf("Failed to convert amount: %s", row[4])
			log.Fatal(message)
		}
		invoice.Amount = amount

		fee, err := util.Atof[float64](row[5])
		if err != nil {
			message := fmt.Sprintf("Failed to convert fee: %s", row[5])
			log.Fatal(message)
		}
		invoice.Fee = fee

		feeRate, err := util.Atof[float64](row[6])
		if err != nil {
			message := fmt.Sprintf("Failed to convert fee rate: %s", row[6])
			log.Fatal(message)
		}
		invoice.FeeRate = feeRate

		tax, err := util.Atof[float64](row[7])
		if err != nil {
			message := fmt.Sprintf("Failed to convert tax: %s", row[7])
			log.Fatal(message)
		}
		invoice.Tax = tax

		taxRate, err := util.Atof[float64](row[8])
		if err != nil {
			message := fmt.Sprintf("Failed to convert tax rate: %s", row[8])
			log.Fatal(message)
		}
		invoice.TaxRate = taxRate

		totalAmount, err := util.Atof[float64](row[9])
		if err != nil {
			message := fmt.Sprintf("Failed to convert total amount: %s", row[9])
			log.Fatal(message)
		}
		invoice.TotalAmount = totalAmount

		dueDate, err := util.Atot(row[10])
		if err != nil {
			message := fmt.Sprintf("Failed to convert due date: %s", row[10])
			log.Fatal(message)
		}
		invoice.DueDate = dueDate

		companyID, err := strconv.Atoi(row[1])
		if err != nil {
			message := fmt.Sprintf("Failed to convert company ID: %s", row[1])
			log.Fatal(message)
		}

		clientID, err := strconv.Atoi(row[2])
		if err != nil {
			message := fmt.Sprintf("Failed to convert client ID: %s", row[2])
			log.Fatal(message)
		}

		invoice.CompanyID = uint(companyID)
		invoice.ClientID = uint(clientID)

		if err := repo.Save(&invoice); err != nil {
			message := fmt.Sprintf("Failed to save invoice: client %v company %v", invoice.ClientID, invoice.CompanyID)
			log.Fatal(message)
		}
	}
}

func returnCSV(csvPath string) [][]string {
	file, err := os.Open(csvPath)
	if err != nil {
		message := fmt.Sprintf("Failed to open CSV file: %s", csvPath)
		log.Fatal(message)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		message := fmt.Sprintf("Failed to read CSV file: %s", csvPath)
		log.Fatal(message)
	}

	return rows
}
