package bank

import (
	"encoding/csv"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BankController struct {
	Logger *slog.Logger
}

func NewBankController(logger *slog.Logger) *BankController {
	return &BankController{
		Logger: logger,
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

	// Create a map to store transactions
	transactions := make(map[string]models.Transaction)

	// Iterate over the records and store them in the map
	for _, record := range records {
		date, err := time.Parse("02/01/2006", record[0]) // DD/MM/YYYY
		if err != nil {
			bc.Logger.Error("Failed to parse date", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse date"})
			return
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			bc.Logger.Error("Failed to parse amount", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse amount"})
			return
		}

		transaction := models.Transaction{
			ID:          uuid.New().String(),
			Date:        date,
			Amount:      amount,
			Description: record[2],
		}

		transactions[transaction.ID] = transaction
	}

	// Log the transactions for debugging [Optional]
	for id, transaction := range transactions {
		bc.Logger.Info("Transaction stored", "id", id, "transaction", transaction)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "transactions": transactions})

	// TODO: Save the transactions to the database with User ID
}
