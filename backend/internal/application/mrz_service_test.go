package application

import (
	"testing"
	"time"

	"github.com/digikeys/backend/internal/domain"
)

func TestCheckDigitKnownValues(t *testing.T) {
	// ICAO 9303 check digit algorithm: weights 7, 3, 1, mod 10
	tests := []struct {
		input    string
		expected int
	}{
		{"520727", 3},   // date check digit example
		{"D231458907", 6}, // D=13*7+2*3+3*1+1*7+4*3+5*1+8*7+9*3+0*1+7*7=256%10=6
		{"<<<<<<", 0},   // all fillers
		{"0", 0},        // single zero
		{"AB", 3},       // A(10)*7 + B(11)*3 = 70+33 = 103 mod 10 = 3
	}
	for _, tt := range tests {
		got := checkDigit(tt.input)
		if got != tt.expected {
			t.Errorf("checkDigit(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestTransliterateAccented(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"Ouédraogo", "OUEDRAOGO"},
		{"François", "FRANCOIS"},
		{"André", "ANDRE"},
		{"Hélène", "HELENE"},
		{"Jean-Pierre", "JEAN<PIERRE"},
		{"N'Diaye", "N<DIAYE"},
		{"Müller", "MUELLER"},     // ü -> UE
		{"Bäcker", "BAECKER"},     // ä -> AE
		{"Ångström", "AANGSTROEM"}, // Å -> AA, ö -> OE
	}
	for _, tt := range tests {
		got := transliterate(tt.input)
		if got != tt.expected {
			t.Errorf("transliterate(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"ABC", 5, "ABC<<"},
		{"ABCDE", 5, "ABCDE"},
		{"ABCDEF", 5, "ABCDE"}, // truncate
		{"", 3, "<<<"},
	}
	for _, tt := range tests {
		got := padRight(tt.input, tt.length)
		if got != tt.expected {
			t.Errorf("padRight(%q, %d) = %q, want %q", tt.input, tt.length, got, tt.expected)
		}
	}
}

func TestFormatDate(t *testing.T) {
	dt := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	got := formatDate(dt)
	if got != "900515" {
		t.Errorf("formatDate = %q, want %q", got, "900515")
	}

	dt2 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	got2 := formatDate(dt2)
	if got2 != "000101" {
		t.Errorf("formatDate = %q, want %q", got2, "000101")
	}
}

func TestGenerateTD1NilInputs(t *testing.T) {
	svc := NewMRZService()

	_, _, _, err := svc.GenerateTD1(nil, &domain.Card{}, &domain.Embassy{})
	if err == nil {
		t.Error("expected error for nil citizen")
	}
	_, _, _, err = svc.GenerateTD1(&domain.Citizen{}, nil, &domain.Embassy{})
	if err == nil {
		t.Error("expected error for nil card")
	}
	_, _, _, err = svc.GenerateTD1(&domain.Citizen{}, &domain.Card{}, nil)
	if err == nil {
		t.Error("expected error for nil embassy")
	}
}

func TestGenerateTD1LineLength(t *testing.T) {
	svc := NewMRZService()

	expiry := time.Date(2030, 6, 15, 0, 0, 0, 0, time.UTC)
	citizen := &domain.Citizen{
		FirstName:   "Amadou",
		LastName:    "Ouedraogo",
		DateOfBirth: time.Date(1990, 3, 20, 0, 0, 0, 0, time.UTC),
		Gender:      "M",
	}
	card := &domain.Card{
		CardNumber: "BF2024001",
		ExpiresAt:  &expiry,
	}
	embassy := &domain.Embassy{
		CardPrefix: "BF",
	}

	line1, line2, line3, err := svc.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}

	if len(line1) != 30 {
		t.Errorf("line1 length = %d, want 30: %q", len(line1), line1)
	}
	if len(line2) != 30 {
		t.Errorf("line2 length = %d, want 30: %q", len(line2), line2)
	}
	if len(line3) != 30 {
		t.Errorf("line3 length = %d, want 30: %q", len(line3), line3)
	}
}

func TestGenerateTD1Content(t *testing.T) {
	svc := NewMRZService()

	expiry := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)
	citizen := &domain.Citizen{
		FirstName:   "Fatimata",
		LastName:    "Traoré",
		DateOfBirth: time.Date(1985, 7, 27, 0, 0, 0, 0, time.UTC),
		Gender:      "F",
	}
	card := &domain.Card{
		CardNumber: "CD2024002",
		ExpiresAt:  &expiry,
	}
	embassy := &domain.Embassy{
		CardPrefix: "CD",
	}

	line1, line2, line3, err := svc.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}

	// Line 1 starts with "ID" + "BFA"
	if line1[:2] != "ID" {
		t.Errorf("line1 should start with ID, got %q", line1[:2])
	}
	if line1[2:5] != "BFA" {
		t.Errorf("line1[2:5] should be BFA, got %q", line1[2:5])
	}

	// Line 2 starts with DOB "850727"
	if line2[:6] != "850727" {
		t.Errorf("line2 DOB should be 850727, got %q", line2[:6])
	}
	// Gender at position 7
	if line2[7:8] != "F" {
		t.Errorf("line2 gender should be F, got %q", line2[7:8])
	}

	// Line 3: transliterated name
	if line3[:6] != "TRAORE" {
		t.Errorf("line3 should start with TRAORE, got %q", line3[:6])
	}
}

func TestGenerateTD1NoExpiry(t *testing.T) {
	svc := NewMRZService()

	citizen := &domain.Citizen{
		FirstName:   "Ali",
		LastName:    "Diallo",
		DateOfBirth: time.Date(1995, 1, 10, 0, 0, 0, 0, time.UTC),
		Gender:      "M",
	}
	card := &domain.Card{
		CardNumber: "BF2024003",
		ExpiresAt:  nil, // no expiry
	}
	embassy := &domain.Embassy{
		CardPrefix: "BF",
	}

	line1, line2, line3, err := svc.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}
	if len(line1) != 30 || len(line2) != 30 || len(line3) != 30 {
		t.Errorf("MRZ line lengths incorrect: %d, %d, %d", len(line1), len(line2), len(line3))
	}

	// Expiry should be "<<<<<<" at positions 8-13 in line2
	if line2[8:14] != "<<<<<<" {
		t.Errorf("expected no-expiry marker '<<<<<<' at line2[8:14], got %q", line2[8:14])
	}
}

func TestGenerateTD1GenderVariants(t *testing.T) {
	svc := NewMRZService()

	expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	base := &domain.Citizen{
		FirstName:   "Test",
		LastName:    "User",
		DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	card := &domain.Card{CardNumber: "BF2024099", ExpiresAt: &expiry}
	embassy := &domain.Embassy{CardPrefix: "BF"}

	genderTests := []struct {
		input    string
		expected string
	}{
		{"M", "M"},
		{"MALE", "M"},
		{"F", "F"},
		{"FEMALE", "F"},
		{"X", "<"},
		{"", "<"},
	}

	for _, tt := range genderTests {
		base.Gender = tt.input
		_, line2, _, err := svc.GenerateTD1(base, card, embassy)
		if err != nil {
			t.Fatalf("GenerateTD1 gender=%q: %v", tt.input, err)
		}
		if line2[7:8] != tt.expected {
			t.Errorf("gender=%q: expected MRZ gender %q, got %q", tt.input, tt.expected, line2[7:8])
		}
	}
}
