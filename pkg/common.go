package pkg

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func TokenGenerator(num int) string {
	b := make([]byte, num)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func IsInt(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}

	return false
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ContainsInt(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ExtractUsernames(text string) []string {
	// Define the regular expression pattern
	pattern := `@([a-zA-Z0-9_]+)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find all matches in the text
	matches := re.FindAllString(text, -1)

	return matches
}

func ExtractSentencesAfterWord(text, keyword string) []string {
	// Split the text into sentences
	sentences := strings.Split(text, ".")

	var result []string

	// Iterate through each sentence
	for _, sentence := range sentences {
		// Find the position of the keyword in the sentence
		index := strings.Index(sentence, keyword)
		if index != -1 {
			// If the keyword is found, extract the part of the sentence after the keyword
			result = append(result, strings.TrimSpace(sentence[index+len(keyword):]))
		}
	}

	return result
}

func IsTimeInBetween(startEpoch, endEpoch int64) bool {
	// Convert epoch values to time.Time
	startTime := time.Unix(startEpoch, 0)
	endTime := time.Unix(endEpoch, 0)

	// Get the current time
	currentTime := time.Now()

	return currentTime.After(startTime) && currentTime.Before(endTime)
}

func IsLetter(s string) bool {
	var regex, _ = regexp.Compile(`^[a-zA-z\s]+$`)
	return regex.MatchString(s)
}

func ShortUUID(uuid string) string {
	// Hash the UUID using MD5
	hash := md5.Sum([]byte(uuid))

	// Convert the hash to base64 encoding
	base64Str := base64.StdEncoding.EncodeToString(hash[:])

	// Take the first 6 characters from the base64 string
	return base64Str[:6]
}

func ReplaceChars(str string, chars []string, replacement string) string {
	for _, char := range chars {
		str = strings.Replace(str, char, replacement, -1)
	}
	return str
}

func FormatDate(date time.Time) string {
	// Define the months in Indonesian
	months := []string{
		"Januari", "Februari", "Maret", "April",
		"Mei", "Juni", "Juli", "Agustus",
		"September", "Oktober", "November", "Desember",
	}

	// Extract the day, month, and year
	day := date.Day()
	month := months[date.Month()-1]
	year := date.Year()

	// Format the date
	formattedDate := fmt.Sprintf("%d %s %d", day, month, year)

	return formattedDate
}

func MaskSuffix(input string) string {
	if len(input) <= 3 {
		return "***"
	}
	return input[:len(input)-3] + "***"
}

func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))  // Generate MD5 hash (returns [16]byte)
	return hex.EncodeToString(hash[:]) // Convert to hexadecimal string
}

func PrizeData(prize string) (prizeType string) {
	// Map to define prize categories
	prizes := map[string]string{
		"belum beruntung":  "reguler",
		"pulsa 25000":      "reguler",
		"pulsa 50000":      "reguler",
		"pulsa 100000":     "reguler",
		"jaket":            "reguler",
		"helm":             "reguler",
		"handphone":        "grand",
		"smart tv 40 inch": "grand",
		"sepeda motor":     "grand",
	}

	// Normalize the input
	prize = strings.ToLower(prize)

	// Check the map for the prize type
	if val, ok := prizes[prize]; ok {
		return val
	}

	// Return a default value for unknown prizes
	return "unknown"
}

func WrapError(err error, msg string) error {
	if err != nil {
		return errors.Wrap(err, msg)
	}
	return nil
}

func WriteXLS(data []map[string]interface{}) (path string, rows int, err error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		err = errors.Wrap(err, "[WriteXLS] error create sheet")
		return
	}

	if len(data) > 0 {
		headers := []string{}
		for key := range data[0] {
			if key == "rNum" {
				continue
			}
			headers = append(headers, strings.ToUpper(key))
		}

		// Write headers to the first row
		for col, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			if err = f.SetCellValue(sheetName, cell, header); err != nil {
				err = errors.Wrap(err, "[WriteXLS] error writing headers")
				return
			}
		}

		// Write data rows
		for rowIndex, row := range data {
			for colIndex, header := range headers {
				cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
				if err = f.SetCellValue(sheetName, cell, row[header]); err != nil {
					err = errors.Wrap(err, "[WriteXLS] error writing data")
					return
				}
			}
		}

		rows = len(data) - 1
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	path = fmt.Sprintf("%s/%d.xlsx", "download", time.Now().Unix())
	if err = f.SaveAs(path); err != nil {
		err = errors.Wrap(err, "[WriteXLS]")
		return
	}

	return
}

func GenerateRandomCode(length int) (code string, err error) {
	bytes := make([]byte, length/2) // Length/2 because hex encoding doubles the size
	if _, err = rand.Read(bytes); err != nil {
		err = errors.Wrap(err, "[GenerateRandomCode]")
		return
	}
	code = hex.EncodeToString(bytes)
	return code[:length], nil
}

func ReverseStrings(slice []string) []string {
	length := len(slice)
	reversed := make([]string, length)
	for i, v := range slice {
		reversed[length-i-1] = v
	}
	return reversed
}
