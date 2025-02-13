package bank

import (
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BankController struct {
	Logger *slog.Logger
	DB     *sqlx.DB
}

func NewBankController(logger *slog.Logger, db *sqlx.DB) *BankController {
	return &BankController{
		Logger: logger,
		DB:     db,
	}
}

func (bc *BankController) UploadBankStatement(c *gin.Context) {
	var records [][]interface{} // two-dimensional slice to represent transaction history table

	// Bind the JSON data from the request body to the records interface
	if err := c.ShouldBindJSON(&records); err != nil {
		bc.Logger.Error("Failed to parse JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Regex pattern to validate date and amount values
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{4}$`) // e.g. "11/11/2011"
	amountPattern := regexp.MustCompile(`^-?\d+(\.\d{1,2})?$`)   // e.g. "-23.50", "11.40", or "100"

	// Create a slice to store the transactions
	var transactions []models.Transaction

	// Grab the index of each header to identify which column is date, amount and description
	header := records[0]
	dateIndex, amountIndex, descriptionIndex := -1, -1, -1
	for i, h := range header {
		switch h {
		case "Date":
			dateIndex = i
		case "Amount":
			amountIndex = i
		case "Description":
			descriptionIndex = i
		}
	}

	if dateIndex == -1 || amountIndex == -1 || descriptionIndex == -1 {
		bc.Logger.Error("Missing required headers")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required headers"})
		return
	}

	// Iterate over the records and store them in the slice
	for _, record := range records[1:] { // Skip the header row

		// Date data
		dateStr := record[dateIndex].(string)
		if !datePattern.MatchString(dateStr) { // User regex to validate date format
			bc.Logger.Error("Invalid date format", "date", dateStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		date, err := time.Parse("2/1/2006", dateStr) // DD/MM.YYYY
		if err != nil {
			bc.Logger.Error("Failed to parse date", "error", err) // User regex to validate amount format
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse date"})
			return
		}

		// Amount data
		amountStr := record[amountIndex].(string)
		if !amountPattern.MatchString(amountStr) {
			bc.Logger.Error("Invalid amount format", "amount", amountStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
			return
		}
		amountDollar, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			bc.Logger.Error("Failed to parse amount", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse amount"})
			return
		}
		amount := int(amountDollar * 100) // Convert dollars to cents

		// Description data
		description, success := record[descriptionIndex].(string)
		if !success {
			// If type assertion fails = record[descriptionIndex] is not a valid string
			description = "No description" // Default description
		}

		// TODO: [WEALT-15] Once Jalil completes authentication, replace the hardcoded User ID with the authenticated User ID
		transaction := models.Transaction{
			UserId:      1, // Hardcoded User ID for now
			Date:        date,
			AmountCents: amount,
			Description: description, // Default description is "No description"
		}

		// Insert the transaction into the database
		transactionId, err := bc.insertTransaction(transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}
		transaction.TransactionId = transactionId
		transactions = append(transactions, transaction)

		// [Optional] Log the transaction
		bc.Logger.Info("Transaction inserted successfully", "transaction", transaction)

	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "transactions": transactions})
}

func (bc *BankController) insertTransaction(transaction models.Transaction) (int, error) {
	// Insert the transaction into the database
	query := `INSERT INTO bank_transactions (user_id, date, amount_cents, description)
			  VALUES ($1, $2, $3, $4) RETURNING transaction_id`
	var transactionId int
	err := bc.DB.QueryRowx(query, transaction.UserId, transaction.Date, transaction.AmountCents, transaction.Description).Scan(&transactionId)
	if err != nil {
		bc.Logger.Error("Failed to insert transaction", "details", transaction, "error", err)
		return -1, err
	}

	// Return the transaction ID of the newly inserted transaction
	return transactionId, nil
}
