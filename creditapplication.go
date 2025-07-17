package main

import "time"

type Applicant struct {
	Type string `json:"type"`

	UserId string `json:"user_id" aml:"user_id"`
	IINBIN string `json:"iinbin" aml:"iinbin"`

	IDCardNumber         string    `json:"id_card_number"`
	IDCardIssueDate      time.Time `json:"id_card_issue_date"`
	IDCardExpirationDate time.Time `json:"id_card_expiration_date"`

	Fullname   string `json:"fullname"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Patronymic string `json:"patronymic"`

	BirthDate   string `json:"birth_date"`
	Nationality string `json:"nationality"`
	Citizenship string `json:"citizenship"`
	Gender      string `json:"gender"`

	PhoneNumber string `json:"phone_number" aml:"phone_number"`
	Email       string `json:"email"`

	RegisteredAddress  string `json:"registered_address"`
	ResidentialAddress string `json:"residential_address"`
}

type CreditApplication struct {
	Id       string    `json:"application_id"`
	Date     time.Time `json:"application_date"`
	Time     string    `json:"application_time"`
	Channel  string    `json:"application_channel"`
	Duration int       `json:"application_duration"`
	Region   string    `json:"application_region"`

	Applicant Applicant `json:"applicant"`

	CreditType     string `json:"credit_type"`
	CreditAmount   string `json:"credit_amount"`
	CreditCurrency string `json:"credit_currency"`
	CreditTerm     int    `json:"credit_term"`

	IncomeAmount string `json:"income_amount"`
	JobType      string `json:"job_type"`
	JobDuration  int    `json:"job_duration"`

	SocialStatus           string `json:"social_status,optional"`
	MilitaryService        string `json:"military_service"`
	DriversLicenseCategory string `json:"drivers_license_category"`
	DriversLicenseNumber   string `json:"drivers_license_number"`

	SpouseIIN     string `json:"spouce_iin"`
	MaritalStatus string `json:"marital_status"`
	Children      int    `json:"children"`

	DebtBurdenRatio         float64 `json:"debt_burden_ratio"`
	CreditScoring           float64 `json:"credit_scoring"`
	ActiveObligationsCount  int     `json:"active_obligations_count"`
	OverduePaymentsCount    int     `json:"overdue_payments_count"`
	OverduePaymentsCount90  int     `json:"overdue_payments_count_90"`
	OverdueAmount           string  `json:"overdue_amount"`
	OutstandingDebt         string  `json:"outstanding_debt"`
	CreditApplicationsCount int     `json:"credit_applications_count"`

	GamblingTotalNumberOfPayments int    `json:"gambling_total_number_of_payments"`
	GamblingPaymentAmount         string `json:"gambling_total_payment_amount"`

	IpAddress string `json:"ip_address"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`

	DeviceId       string `json:"device_id,optional"`
	DeviceModel    string `json:"device_model"`
	DevicePlatform string `json:"device_platform"`
	OsId           string `json:"os_id"`
	DeviceLanguage string `json:"device_language"`

	BankruptcyStatus          bool `json:"bankruptcy_status"`
	CreditIssuanceRestriction bool `json:"credit_issuance_restriction"`
	Biometrics                bool `json:"biometrics"`
	DriversLicense            bool `json:"drivers_license"`
	SpousesConsent            bool `json:"spouses_consent"`
}
