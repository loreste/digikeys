package mrz

import (
	"testing"
	"time"

	"github.com/digikeys/backend/internal/domain"
)

func TestTransliterateBasicASCII(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"JOHN", "JOHN"},
		{"john", "JOHN"},
		{"Jean Pierre", "JEAN<PIERRE"},
		{"O'Brien", "O<BRIEN"},
		{"Smith-Jones", "SMITH<JONES"},
	}
	for _, tt := range tests {
		got := Transliterate(tt.input)
		if got != tt.expected {
			t.Errorf("Transliterate(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestTransliterateAccentedCharacters(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"Ouédraogo", "OUEDRAOGO"},           // é -> E
		{"François", "FRANCOIS"},             // ç -> C
		{"Müller", "MUELLER"},                 // ü -> UE
		{"Ångström", "AANGSTROEM"},             // Å -> AA, ö -> OE
		{"Ñoño", "NONO"},                      // ñ -> N
		{"Héloïse", "HELOISE"},                // é -> E, ï -> I
		{"Bäcker", "BAECKER"},                 // ä -> AE
		{"Ólafsson", "OLAFSSON"},              // Ó -> O
		{"André", "ANDRE"},                    // é -> E
		{"Dağıstan", "DAISTAN"},               // ğ and ı are non-ASCII non-mapped, dropped
	}
	for _, tt := range tests {
		got := Transliterate(tt.input)
		if got != tt.expected {
			t.Errorf("Transliterate(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestTransliterateDropsNonMRZChars(t *testing.T) {
	// Numbers, special chars (except space, hyphen, apostrophe) should be dropped
	got := Transliterate("John123!@#")
	if got != "JOHN" {
		t.Errorf("Transliterate with non-MRZ chars: got %q, want %q", got, "JOHN")
	}
}

func TestPadRightPadding(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"ABC", 5, "ABC<<"},
		{"AB", 2, "AB"},
		{"ABCDE", 3, "ABC"},
		{"", 4, "<<<<"},
		{"X", 1, "X"},
	}
	for _, tt := range tests {
		got := PadRight(tt.input, tt.length)
		if got != tt.expected {
			t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.length, got, tt.expected)
		}
		if len(got) != tt.length {
			t.Errorf("PadRight(%q, %d) length = %d, want %d", tt.input, tt.length, len(got), tt.length)
		}
	}
}

func TestCheckDigitICAO(t *testing.T) {
	// Known ICAO 9303 test values
	tests := []struct {
		input    string
		expected int
	}{
		// Weights cycle: 7, 3, 1, 7, 3, 1, ...
		// '0' = 0, '1' = 1, ..., '9' = 9, 'A' = 10, 'B' = 11, ..., 'Z' = 35, '<' = 0
		{"520727", 3},     // 5*7 + 2*3 + 0*1 + 7*7 + 2*3 + 7*1 = 35+6+0+49+6+7 = 103. 103%10 = 3
		{"AB", 3},         // A(10)*7 + B(11)*3 = 70+33 = 103. 103%10 = 3
		{"<<<<<<", 0},     // all zeros
		{"7", 9},          // 7*7 = 49. 49%10 = 9
		{"L898902C3", 6},  // ICAO example: L(21)*7+8*3+9*1+8*7+9*3+0*1+2*7+C(12)*3+3*1 = 147+24+9+56+27+0+14+36+3 = 316. 316%10 = 6
	}
	for _, tt := range tests {
		got := CheckDigit(tt.input)
		if got != tt.expected {
			t.Errorf("CheckDigit(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestCheckDigitWeightsPattern(t *testing.T) {
	// Verify the 7,3,1 cycling pattern explicitly
	// "A" at position 0: 10 * 7 = 70
	if CheckDigit("A") != 0 { // 70 % 10 = 0
		t.Errorf("CheckDigit(A) = %d, want 0", CheckDigit("A"))
	}
	// "B" at position 0: 11 * 7 = 77
	if CheckDigit("B") != 7 { // 77 % 10 = 7
		t.Errorf("CheckDigit(B) = %d, want 7", CheckDigit("B"))
	}
}

func TestFormatDateYYMMDD(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC), "900515"},
		{time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "000101"},
		{time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), "251231"},
		{time.Date(1985, 7, 27, 0, 0, 0, 0, time.UTC), "850727"},
	}
	for _, tt := range tests {
		got := FormatDate(tt.input)
		if got != tt.expected {
			t.Errorf("FormatDate(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGenerateTD1NilInputs(t *testing.T) {
	g := NewGenerator()

	_, _, _, err := g.GenerateTD1(nil, &domain.Card{}, &domain.Embassy{})
	if err == nil {
		t.Error("expected error for nil citizen")
	}
	_, _, _, err = g.GenerateTD1(&domain.Citizen{}, nil, &domain.Embassy{})
	if err == nil {
		t.Error("expected error for nil card")
	}
	_, _, _, err = g.GenerateTD1(&domain.Citizen{}, &domain.Card{}, nil)
	if err == nil {
		t.Error("expected error for nil embassy")
	}
}

func TestGenerateTD1LineLengths(t *testing.T) {
	g := NewGenerator()

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
	embassy := &domain.Embassy{CardPrefix: "BF"}

	line1, line2, line3, err := g.GenerateTD1(citizen, card, embassy)
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

func TestGenerateTD1Structure(t *testing.T) {
	g := NewGenerator()

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
	embassy := &domain.Embassy{CardPrefix: "CD"}

	line1, line2, line3, err := g.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}

	// Line 1: positions 0-1 = "ID", 2-4 = "BFA"
	if line1[0:2] != "ID" {
		t.Errorf("line1[0:2] = %q, want ID", line1[0:2])
	}
	if line1[2:5] != "BFA" {
		t.Errorf("line1[2:5] = %q, want BFA", line1[2:5])
	}

	// Line 2: DOB "850727" at positions 0-5
	if line2[0:6] != "850727" {
		t.Errorf("line2 DOB = %q, want 850727", line2[0:6])
	}
	// DOB check digit at position 6
	expectedDOBCheck := CheckDigit("850727")
	gotDOBCheck := int(line2[6] - '0')
	if gotDOBCheck != expectedDOBCheck {
		t.Errorf("line2 DOB check digit = %d, want %d", gotDOBCheck, expectedDOBCheck)
	}

	// Gender at position 7
	if line2[7:8] != "F" {
		t.Errorf("line2 gender = %q, want F", line2[7:8])
	}

	// Expiry "291231" at positions 8-13
	if line2[8:14] != "291231" {
		t.Errorf("line2 expiry = %q, want 291231", line2[8:14])
	}

	// Nationality "BFA" at positions 15-17
	if line2[15:18] != "BFA" {
		t.Errorf("line2 nationality = %q, want BFA", line2[15:18])
	}

	// Line 3: name starts with TRAORE (transliterated from Traoré)
	if line3[0:6] != "TRAORE" {
		t.Errorf("line3 starts with %q, want TRAORE", line3[0:6])
	}
	// Should contain << separator between surname and given name
	if line3[6:8] != "<<" {
		t.Errorf("line3[6:8] = %q, want <<", line3[6:8])
	}
}

func TestGenerateTD1NoExpiry(t *testing.T) {
	g := NewGenerator()

	citizen := &domain.Citizen{
		FirstName:   "Ali",
		LastName:    "Diallo",
		DateOfBirth: time.Date(1995, 1, 10, 0, 0, 0, 0, time.UTC),
		Gender:      "M",
	}
	card := &domain.Card{
		CardNumber: "BF2024003",
		ExpiresAt:  nil,
	}
	embassy := &domain.Embassy{CardPrefix: "BF"}

	_, line2, _, err := g.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}

	// No expiry -> "<<<<<<" at positions 8-13
	if line2[8:14] != "<<<<<<" {
		t.Errorf("expected no-expiry '<<<<<<', got %q", line2[8:14])
	}
}

func TestGenerateTD1GenderMapping(t *testing.T) {
	g := NewGenerator()

	expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	card := &domain.Card{CardNumber: "BF2024099", ExpiresAt: &expiry}
	embassy := &domain.Embassy{CardPrefix: "BF"}

	cases := []struct {
		gender   string
		expected string
	}{
		{"M", "M"},
		{"MALE", "M"},
		{"F", "F"},
		{"FEMALE", "F"},
		{"X", "<"},
		{"other", "<"},
	}

	for _, tc := range cases {
		citizen := &domain.Citizen{
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:      tc.gender,
		}
		_, line2, _, err := g.GenerateTD1(citizen, card, embassy)
		if err != nil {
			t.Fatalf("GenerateTD1 gender=%q: %v", tc.gender, err)
		}
		if line2[7:8] != tc.expected {
			t.Errorf("gender=%q: MRZ gender = %q, want %q", tc.gender, line2[7:8], tc.expected)
		}
	}
}

func TestGenerateTD1CompositeCheckDigit(t *testing.T) {
	g := NewGenerator()

	expiry := time.Date(2030, 6, 15, 0, 0, 0, 0, time.UTC)
	citizen := &domain.Citizen{
		FirstName:   "Amadou",
		LastName:    "Ouedraogo",
		DateOfBirth: time.Date(1990, 3, 20, 0, 0, 0, 0, time.UTC),
		Gender:      "M",
	}
	card := &domain.Card{CardNumber: "BF2024001", ExpiresAt: &expiry}
	embassy := &domain.Embassy{CardPrefix: "BF"}

	line1, line2, _, err := g.GenerateTD1(citizen, card, embassy)
	if err != nil {
		t.Fatalf("GenerateTD1 failed: %v", err)
	}

	// Verify composite check digit manually
	docNum := line1[5:14]     // 9 chars
	docNumCheckStr := line1[14:15]
	optional1 := line1[15:30] // 15 chars
	dob := line2[0:6]
	dobCheckStr := line2[6:7]
	expiryStr := line2[8:14]
	expiryCheckStr := line2[14:15]
	optional2 := line2[18:29] // 11 chars

	compositeInput := docNum + docNumCheckStr + optional1 + dob + dobCheckStr + expiryStr + expiryCheckStr + optional2
	expectedComposite := CheckDigit(compositeInput)
	gotComposite := int(line2[29] - '0')

	if gotComposite != expectedComposite {
		t.Errorf("composite check digit = %d, want %d (input: %q)", gotComposite, expectedComposite, compositeInput)
	}
}
