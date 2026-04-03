package helpers

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jasonsavage/financeflow/internal/models"
)

func ParseTransactionsCSV(r io.Reader, accountID string, bankAccountName string) ([]models.Transaction, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	// Accept variable row lengths if banks are messy
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("empty csv file")
		}
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	dateIdx := -1
	descIdx := -1
	creditIdx := -1
	debitIdx := -1
	typeIdx := -1
	amountIdx := -1

	// Normalize and find indices
	for i, header := range headers {
		lower := strings.ToLower(strings.TrimSpace(header))
		if strings.Contains(lower, "date") && dateIdx == -1 {
			dateIdx = i
		} else if strings.Contains(lower, "description") && descIdx == -1 {
			descIdx = i
		} else if lower == "credit" {
			creditIdx = i
		} else if lower == "debit" {
			debitIdx = i
		} else if strings.Contains(lower, "type") && typeIdx == -1 {
			typeIdx = i
		} else if strings.Contains(lower, "amount") && amountIdx == -1 {
			amountIdx = i
		}
	}

	if dateIdx == -1 {
		return nil, errors.New("could not find 'date' column")
	}
	if descIdx == -1 {
		return nil, errors.New("could not find 'description' column")
	}

	var format string
	if creditIdx != -1 && debitIdx != -1 {
		format = "credit_debit"
	} else if typeIdx != -1 && amountIdx != -1 {
		format = "type_amount"
	} else if amountIdx != -1 {
		format = "pos_neg"
	} else {
		return nil, errors.New("could not detect amount format")
	}

	var transactions []models.Transaction
	now := time.Now()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed reading row: %w", err)
		}

		// Ensure we don't query out of bounds
		if len(record) <= dateIdx || len(record) <= descIdx {
			continue // skip malformed row
		}

		// 1. Transaction ID - hash the raw row data + accountID for deduping
		rowStr := strings.Join(record, ",")
		hash := sha256.Sum256([]byte(rowStr + accountID))
		txID := hex.EncodeToString(hash[:])

		dateStr := strings.TrimSpace(record[dateIdx])
		desc := strings.TrimSpace(record[descIdx])

		// Parse date
		parsedDate, err := parseDate(dateStr)
		if err != nil {
			continue // skip rows with bad dates
		}

		var deposit, withdrawal float64
		switch format {
		case "credit_debit":
			if len(record) > creditIdx && record[creditIdx] != "" {
				deposit, _ = parseAmount(record[creditIdx])
			}
			if len(record) > debitIdx && record[debitIdx] != "" {
				withdrawal, _ = parseAmount(record[debitIdx])
				withdrawal = abs(withdrawal)
			}
		case "type_amount":
			if len(record) > amountIdx && len(record) > typeIdx {
				val, _ := parseAmount(record[amountIdx])
				tType := strings.ToLower(strings.TrimSpace(record[typeIdx]))
				if tType == "credit" {
					deposit = abs(val)
				} else if tType == "debit" {
					withdrawal = abs(val)
				}
			}
		case "pos_neg":
			if len(record) > amountIdx {
				val, _ := parseAmount(record[amountIdx])
				if val >= 0 {
					deposit = val
				} else {
					withdrawal = abs(val)
				}
			}
		}

		// Only save valid rows where deposit or withdrawal exists
		if deposit == 0 && withdrawal == 0 {
			continue
		}

		transactions = append(transactions, models.Transaction{
			TransactionID:   txID,
			AccountID:       accountID,
			Date:            parsedDate,
			Description:     desc,
			Category:        nil,
			Deposit:         deposit,
			Withdrawal:      withdrawal,
			BankAccountName: bankAccountName,
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	}

	return transactions, nil
}

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
		"2006/01/02",
		"01/02/06",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unknown date format")
}

func parseAmount(amtStr string) (float64, error) {
	amtStr = strings.ReplaceAll(amtStr, "$", "")
	amtStr = strings.ReplaceAll(amtStr, ",", "")
	amtStr = strings.ReplaceAll(amtStr, "\"", "")
	amtStr = strings.TrimSpace(amtStr)
	if amtStr == "" {
		return 0, nil
	}
	return strconv.ParseFloat(amtStr, 64)
}

func abs(val float64) float64 {
	if val < 0 {
		return -val
	}
	return val
}
