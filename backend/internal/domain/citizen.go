package domain

import "time"

type Citizen struct {
	ID                 string    `json:"id" db:"id"`
	FirstName          string    `json:"firstName" db:"first_name"`
	LastName           string    `json:"lastName" db:"last_name"`
	MaidenName         string    `json:"maidenName,omitempty" db:"maiden_name"`
	DateOfBirth        time.Time `json:"dateOfBirth" db:"date_of_birth"`
	PlaceOfBirth       string    `json:"placeOfBirth" db:"place_of_birth"`
	Gender             string    `json:"gender" db:"gender"`
	Nationality        string    `json:"nationality" db:"nationality"`
	NationalID         string    `json:"nationalId,omitempty" db:"national_id"`
	UniqueIdentifier   string    `json:"uniqueIdentifier,omitempty" db:"unique_identifier"`
	PassportNumber     string    `json:"passportNumber,omitempty" db:"passport_number"`
	Phone              string    `json:"phone,omitempty" db:"phone"`
	Email              string    `json:"email,omitempty" db:"email"`
	CountryOfResidence string    `json:"countryOfResidence" db:"country_of_residence"`
	CityOfResidence    string    `json:"cityOfResidence" db:"city_of_residence"`
	AddressAbroad      string    `json:"addressAbroad,omitempty" db:"address_abroad"`
	ProvinceOfOrigin   string    `json:"provinceOfOrigin,omitempty" db:"province_of_origin"`
	CommuneOfOrigin    string    `json:"communeOfOrigin,omitempty" db:"commune_of_origin"`
	EmbassyID          string    `json:"embassyId" db:"embassy_id"`
	RegisteredBy       string    `json:"registeredBy" db:"registered_by"`
	PhotoKey           string    `json:"photoKey,omitempty" db:"photo_key"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`
}

type CitizenFilter struct {
	EmbassyID          string
	CountryOfResidence string
	Query              string
	Status             string
	Page               int
	PageSize           int
}
