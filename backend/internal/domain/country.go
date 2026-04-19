package domain

const (
	CountryBF = "BF" // Burkina Faso
	CountryCD = "CD" // Democratic Republic of Congo
)

type Country struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	NameEN    string `json:"nameEn"`
	Currency  string `json:"currency"`
	Motto     string `json:"motto"`
	FlagEmoji string `json:"flagEmoji"`
}

var Countries = map[string]Country{
	CountryBF: {
		Code: "BF", Name: "Burkina Faso", NameEN: "Burkina Faso",
		Currency: "XOF", Motto: "Unité - Progrès - Justice", FlagEmoji: "🇧🇫",
	},
	CountryCD: {
		Code: "CD", Name: "République Démocratique du Congo", NameEN: "Democratic Republic of the Congo",
		Currency: "CDF", Motto: "Justice - Paix - Travail", FlagEmoji: "🇨🇩",
	},
}

// NationalityCode returns the 3-letter ICAO nationality code for MRZ generation.
func NationalityCode(country string) string {
	switch country {
	case CountryBF:
		return "BFA"
	case CountryCD:
		return "COD"
	default:
		return "XXX"
	}
}

// CardPrefix returns the consular card document prefix.
func CardPrefix(country string) string {
	switch country {
	case CountryBF:
		return "CC" // Carte Consulaire
	case CountryCD:
		return "CD" // Carte Diaspora
	default:
		return "CC"
	}
}

// FundName returns the solidarity fund name for the country.
func FundName(country string) string {
	switch country {
	case CountryBF:
		return "Fonds de Solidarité Burkinabè (FSB)"
	case CountryCD:
		return "Fonds de Solidarité Congolais (FSC)"
	default:
		return "Fonds de Solidarité"
	}
}

// ContributionAmount returns the solidarity fund contribution in smallest currency unit.
func ContributionAmount(country string) int64 {
	switch country {
	case CountryBF:
		return 150000 // 1,500 FCFA (XOF centimes)
	case CountryCD:
		return 500000 // 5,000 CDF (Congolese Franc centimes)
	default:
		return 150000
	}
}

// MobileMoneyProviders returns available payment providers for a country.
func MobileMoneyProviders(country string) []string {
	switch country {
	case CountryBF:
		return []string{"ORANGE_MONEY", "MOOV_MONEY"}
	case CountryCD:
		return []string{"ORANGE_MONEY_RDC", "AIRTEL_MONEY", "AFRICELL_MONEY", "VODACOM"}
	default:
		return []string{"ORANGE_MONEY"}
	}
}

// PartnerBanks returns partner banks for account opening per country.
func PartnerBanks(country string) []string {
	switch country {
	case CountryBF:
		return []string{"Coris Bank International", "Bank of Africa Burkina", "Ecobank Burkina"}
	case CountryCD:
		return []string{"Rawbank", "Equity BCDC", "TMB (Trust Merchant Bank)", "FBN Bank RDC"}
	default:
		return []string{}
	}
}
