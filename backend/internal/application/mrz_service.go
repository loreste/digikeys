package application

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/digikeys/backend/internal/domain"
)

// MRZService generates ICAO 9303 TD1 format MRZ (3 lines x 30 characters).
type MRZService struct{}

func NewMRZService() *MRZService {
	return &MRZService{}
}

// transliteration maps accented characters to their MRZ equivalents.
var transliteration = map[rune]string{
	'\u00C0': "A", '\u00C1': "A", '\u00C2': "A", '\u00C3': "A", '\u00C4': "AE", '\u00C5': "AA",
	'\u00C6': "AE", '\u00C7': "C", '\u00C8': "E", '\u00C9': "E", '\u00CA': "E", '\u00CB': "E",
	'\u00CC': "I", '\u00CD': "I", '\u00CE': "I", '\u00CF': "I",
	'\u00D0': "D", '\u00D1': "N", '\u00D2': "O", '\u00D3': "O", '\u00D4': "O", '\u00D5': "O",
	'\u00D6': "OE", '\u00D8': "OE", '\u00D9': "U", '\u00DA': "U", '\u00DB': "U", '\u00DC': "UE",
	'\u00DD': "Y", '\u00DE': "TH", '\u00DF': "SS",
	'\u00E0': "A", '\u00E1': "A", '\u00E2': "A", '\u00E3': "A", '\u00E4': "AE", '\u00E5': "AA",
	'\u00E6': "AE", '\u00E7': "C", '\u00E8': "E", '\u00E9': "E", '\u00EA': "E", '\u00EB': "E",
	'\u00EC': "I", '\u00ED': "I", '\u00EE': "I", '\u00EF': "I",
	'\u00F0': "D", '\u00F1': "N", '\u00F2': "O", '\u00F3': "O", '\u00F4': "O", '\u00F5': "O",
	'\u00F6': "OE", '\u00F8': "OE", '\u00F9': "U", '\u00FA': "U", '\u00FB': "U", '\u00FC': "UE",
	'\u00FD': "Y", '\u00FE': "TH", '\u00FF': "Y",
}

// transliterate converts a name string to MRZ-compatible uppercase ASCII.
func transliterate(name string) string {
	name = strings.ToUpper(name)
	var result strings.Builder
	for _, r := range name {
		if repl, ok := transliteration[r]; ok {
			result.WriteString(repl)
		} else if unicode.IsLetter(r) && r < 128 {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '\'' {
			result.WriteRune('<')
		}
		// Drop any other character
	}
	return result.String()
}

// padRight pads a string with '<' to the given length, or truncates it.
func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat("<", length-len(s))
}

// checkDigit computes ICAO 9303 check digit using mod 10, weights 7,3,1.
func checkDigit(s string) int {
	weights := []int{7, 3, 1}
	total := 0
	for i, r := range s {
		var val int
		switch {
		case r >= '0' && r <= '9':
			val = int(r - '0')
		case r >= 'A' && r <= 'Z':
			val = int(r-'A') + 10
		default: // '<' or filler
			val = 0
		}
		total += val * weights[i%3]
	}
	return total % 10
}

// formatDate formats a time.Time as YYMMDD for MRZ.
func formatDate(t time.Time) string {
	return t.Format("060102")
}

// GenerateTD1 generates a TD1 MRZ (ID card, 3 lines x 30 chars) per ICAO 9303.
//
// Line 1: document type (2) + issuing state (3) + document number (9) + check (1) + optional (15)
// Line 2: DOB (6) + check (1) + sex (1) + expiry (6) + check (1) + nationality (3) + optional (11) + composite check (1)
// Line 3: surname<<given<names (30)
func (s *MRZService) GenerateTD1(citizen *domain.Citizen, card *domain.Card, embassy *domain.Embassy) (string, string, string, error) {
	if citizen == nil || card == nil || embassy == nil {
		return "", "", "", fmt.Errorf("citizen, card, and embassy are required")
	}

	docType := "ID"
	country := padRight(transliterate("BFA"), 3) // Burkina Faso ISO 3166-1 alpha-3

	// Document number: use first 9 chars of card number
	docNum := padRight(transliterate(card.CardNumber), 9)
	docNumCheck := checkDigit(docNum)

	// Optional data for line 1 (15 chars)
	optional1 := padRight("", 15)

	line1 := fmt.Sprintf("%s%s%s%d%s",
		padRight(docType, 2),
		country,
		docNum,
		docNumCheck,
		optional1,
	)

	// Line 2
	dob := formatDate(citizen.DateOfBirth)
	dobCheck := checkDigit(dob)

	gender := "M"
	switch strings.ToUpper(citizen.Gender) {
	case "F", "FEMALE":
		gender = "F"
	case "M", "MALE":
		gender = "M"
	default:
		gender = "<"
	}

	expiry := "<<<<<<" // If no expiry date set
	if card.ExpiresAt != nil {
		expiry = formatDate(*card.ExpiresAt)
	}
	expiryCheck := checkDigit(expiry)

	nationality := padRight(transliterate("BFA"), 3)

	optional2 := padRight("", 11)

	// Composite check digit: over doc number + check + optional1 + dob + check + expiry + check + optional2
	compositeData := fmt.Sprintf("%s%d%s%s%d%s%d%s",
		docNum, docNumCheck, optional1,
		dob, dobCheck,
		expiry, expiryCheck,
		optional2,
	)
	compositeCheck := checkDigit(compositeData)

	line2 := fmt.Sprintf("%s%d%s%s%d%s%s%d",
		dob, dobCheck,
		gender,
		expiry, expiryCheck,
		nationality,
		optional2,
		compositeCheck,
	)

	// Line 3: name
	surname := transliterate(citizen.LastName)
	givenNames := transliterate(citizen.FirstName)
	nameLine := surname + "<<" + givenNames
	line3 := padRight(nameLine, 30)

	// Validate lengths
	if len(line1) != 30 || len(line2) != 30 || len(line3) != 30 {
		return "", "", "", fmt.Errorf("MRZ line length mismatch: L1=%d L2=%d L3=%d", len(line1), len(line2), len(line3))
	}

	return line1, line2, line3, nil
}
