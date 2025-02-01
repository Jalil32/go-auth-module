package bank

import (
	"encoding/csv"
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
	file, err := c.FormFile("file") // Single File Upload
	if err != nil {
		bc.Logger.Error("Failed to retrieve file from request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		bc.Logger.Error("Failed to open file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close() // Guarantee the file is closed

	// Parse the CSV file
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		bc.Logger.Error("Failed to parse CSV file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	// Slice to store the transactions
	var transactions []models.Transaction

	// Regex pattern to match date, amount and description values
	datePattern := regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`)
	amountPattern := regexp.MustCompile(`^-?\d+\.\d{2}$`)

	// Iterate over the records and store them in the map
	for _, record := range records {
		var date time.Time
		var amount int                            // Amount in cents
		var description string = "No description" // Default description

		dateFound := false
		amountFound := false
		descriptionFound := false

		for _, field := range record {
			if datePattern.MatchString(field) && !dateFound {
				date, err = time.Parse("02/01/2006", field) // DD/MM/YYYY
				if err != nil {
					bc.Logger.Error("Failed to parse date", "error", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse date"})
					return
				}
				dateFound = true
			} else if amountPattern.MatchString(field) && !amountFound {
				parsedAmount, err := strconv.ParseFloat(field, 64)
				if err != nil {
					bc.Logger.Error("Failed to parse amount", "error", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse amount"})
					return
				}
				amount = int(parsedAmount * 100) // Convert dollars to cents
				amountFound = true
			} else if !descriptionFound {
				description = field
				descriptionFound = true
			}
		}

		// Create a transaction object if at least the date and amount are found
		if dateFound && amountFound {
			transaction := models.Transaction{
				UserID:      1, // Hardcoded User ID for now
				Date:        date,
				AmountCents: amount,
				Description: description, // Default description is "No description"
			}

			transactions = append(transactions, transaction)
		}
	}

	// Insert the transactions into the database
	if err := bc.insertTransactions(transactions); err != nil {
		bc.Logger.Error("Failed to insert transactions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "transactions": transactions})
}

func (bc *BankController) insertTransactions(transactions []models.Transaction) error {
	tx, err := bc.DB.Beginx()
	if err != nil {
		bc.Logger.Error("Failed to start transaction", "error", err)
		return err
	}

	for _, transaction := range transactions {
		query := `INSERT INTO bank_transactions (user_id, date, amount_cents, description)
                  VALUES ($1, $2, $3, $4) RETURNING transaction_id`
		err := tx.QueryRowx(query, transaction.UserID, transaction.Date, transaction.AmountCents, transaction.Description).Scan(&transaction.TransactionID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
