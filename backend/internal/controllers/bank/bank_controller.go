package bank

import (
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
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

	f, err := file.Open()
	if err != nil {
		bc.Logger.Error("Failed to open file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}

	defer f.Close() // Guarantee file closure

	content, err := readPdf(f, file.Size)
	if err != nil {
		bc.Logger.Error("Failed to read PDF", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}

func readPdf(f multipart.File, fileSize int64) (string, error) {
	r, err := pdf.NewReader(f, fileSize)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	var content string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

        rows, _ := p.GetTextByRow()
        for _, row := range rows {
            for _, word := range row.Content {
                content += word.S + " "
            }
            content += "\n"
        }
    }

    return content, nil
}