package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
	Timeout: 10 * time.Second,
}

var txnTypes = [...]string{"C2C", "C2A", "P2P", "A2C", "A2A", "CORPC2A", "CORPC2C", "CORPA2A", "CORPA2C"}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var channels = []string{"mobile", "office", "web"}
var applicationTypes = [...]string{"Автомобильный", "Потребительский", "Ипотечный"}

func generateRandomCAs(count int) ([][]byte, error) {
	applications := make([][]byte, count)
	for i := 0; i < count; i++ {
		creditAmountFloat64 := gofakeit.Float64Range(float64(1000000), float64(10000000))
		creditAmountString := fmt.Sprintf("%f", creditAmountFloat64)
		incomeAmountFloat64 := gofakeit.Float64Range(float64(100000), float64(900000))
		applicationDate := time.Now().UTC()
		applicationTime := applicationDate.Format("15:04:05")

		application := CreditApplication{
			Id:       uuid.New().String(),
			Date:     applicationDate,
			Time:     applicationTime,
			Channel:  channels[r.Intn(len(channels))],
			Duration: gofakeit.IntRange(60, 120),
			Region:   "Алматы",
			Applicant: Applicant{
				Type:                 "person",
				Age:                  gofakeit.IntRange(18, 90),
				UserId:               uuid.New().String(),
				IINBIN:               "001000001000",
				IDCardNumber:         gofakeit.SSN(),
				IDCardIssueDate:      gofakeit.PastDate(),
				IDCardExpirationDate: gofakeit.FutureDate(),
				Firstname:            gofakeit.FirstName(),
				Lastname:             gofakeit.LastName(),
				Patronymic:           gofakeit.MiddleName(),
				BirthDate:            gofakeit.PastDate().String(),
				Nationality:          "казах",
				Citizenship:          "Казахстан",
				Gender:               gofakeit.Gender(),
				PhoneNumber:          gofakeit.PhoneFormatted(),
				Email:                gofakeit.Email(),
				RegisteredAddress:    gofakeit.Address().Address,
				ResidentialAddress:   gofakeit.Address().Address,
			},
			CreditType:                    applicationTypes[r.Intn(len(applicationTypes))],
			CreditAmount:                  creditAmountString,
			CreditCurrency:                "KZT",
			CreditTerm:                    gofakeit.IntRange(12, 72),
			IncomeAmount:                  fmt.Sprintf("%f", incomeAmountFloat64),
			JobType:                       gofakeit.JobTitle(),
			JobDuration:                   gofakeit.IntRange(1, 5),
			BankruptcyStatus:              gofakeit.Bool(),
			CreditIssuanceRestriction:     gofakeit.Bool(),
			MilitaryService:               "Военнобязанный",
			Biometrics:                    true,
			DriversLicense:                gofakeit.Bool(),
			DriversLicenseCategory:        "",
			DriversLicenseNumber:          "",
			SpousesConsent:                true,
			SpouseIIN:                     "",
			MaritalStatus:                 "",
			Children:                      0,
			DebtBurdenRatio:               gofakeit.Float64(),
			CreditScoring:                 gofakeit.Float64(),
			ActiveObligationsCount:        gofakeit.IntRange(0, 10),
			OverduePaymentsCount:          gofakeit.IntRange(0, 5),
			OverduePaymentsCount90:        gofakeit.IntRange(0, 3),
			OverdueAmount:                 "",
			OutstandingDebt:               "",
			CreditApplicationsCount:       gofakeit.IntRange(0, 5),
			GamblingTotalNumberOfPayments: 0,
			GamblingPaymentAmount:         "",
			IpAddress:                     gofakeit.IPv4Address(),
			Longitude:                     "",
			Latitude:                      "",
			DeviceId:                      gofakeit.UUID(),
			DeviceModel:                   "IPhone 15",
			DevicePlatform:                "IOS",
			OsId:                          "",
			DeviceLanguage:                "ru",
		}

		applicationJSON, err := json.Marshal(application)
		if err != nil {
			return nil, err
		}

		applications[i] = applicationJSON
	}

	return applications, nil
}

func generateRandomTxns(count int) ([][]byte, error) {
	txns := make([][]byte, count)
	for i := 0; i < count; i++ {
		date := time.Now().UTC()
		txnTime := date.Format("15:04:05")

		txn := Transaction{
			Id:                   uuid.New().String(),
			SourceUserId:         uuid.New().String(),
			SourceIdentifier:     strconv.Itoa(gofakeit.Number(100000000000, 900000000000)),
			SourceFullname:       gofakeit.LastName() + " " + gofakeit.FirstName() + " " + gofakeit.MiddleName(),
			SourceCardNumber:     gofakeit.CreditCardNumber(nil),
			SourceAccount:        "KZ" + gofakeit.SSN() + gofakeit.SSN(),
			TargetUserId:         uuid.New().String(),
			TargetIdentifier:     strconv.Itoa(gofakeit.Number(100000000000, 900000000000)),
			TargetFullname:       gofakeit.LastName() + " " + gofakeit.FirstName() + " " + gofakeit.MiddleName(),
			TargetCardNumber:     gofakeit.CreditCardNumber(nil),
			TargetAccount:        "KZ" + gofakeit.SSN() + gofakeit.SSN(),
			MerchantId:           uuid.New().String(),
			MerchantTerminalId:   uuid.New().String(),
			MerchantMCCCode:      gofakeit.CreditCardCvv(),
			Date:                 date.Format(time.RFC3339),
			Time:                 txnTime,
			Amount:               strconv.Itoa(gofakeit.Number(100000, 10000000)),
			Currency:             "KZT",
			PaymentMode:          gofakeit.CreditCardCvv(),
			TransactionType:      txnTypes[r.Intn(len(txnTypes))],
			TransactionCountry:   "Казахстан",
			TransactionCity:      "Алматы",
			TransactionChannel:   "Y",
			TransactionRRN:       uuid.New().String(),
			TransactionStatus:    "Non-3DS",
			RegistrationDate:     gofakeit.Date().Format(time.RFC3339),
			CardType:             "debit",
			NewRecipient:         "yes",
			NewTerminal:          "false",
			DeviceId:             uuid.New().String(),
			LastDeviceUpdateDate: gofakeit.Date().Format(time.RFC3339),
			RemoteAccess:         "false",
			ScreenSharing:        "false",
			HardwareId:           uuid.New().String(),
			OSID:                 uuid.New().String(),
			IsTokenized:          "false",
			CookieEnabled:        "false",
			LastLoginDate:        gofakeit.Date().Format(time.RFC3339Nano),
			LastRegistrationDate: gofakeit.Date().Format(time.RFC3339),
			LastDenyEvent:        gofakeit.Date().Format(time.RFC3339Nano),
			LastReviewEvent:      gofakeit.Date().Format(time.RFC3339),
			LastLimitsUpdateDate: gofakeit.Date().Format(time.RFC3339),
			PinUpdateDate:        gofakeit.Date().Format(time.RFC3339),
		}

		txnBytes, err := json.Marshal(txn)
		if err != nil {
			return nil, err
		}

		txns[i] = txnBytes
	}

	return txns, nil
}

func doPostCARequest(requestUrl string, data []byte) error {
	// TODO: Handle response and log fails
	reqId := uuid.New()

	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", reqId.String())

	_, err = client.Do(req)
	if err != nil {
		var object Transaction

		if err := json.Unmarshal(data, &object); err != nil {
			slog.Error("failed to unmarshal", "err", err.Error())
			return err
		}

		slog.Info("failed object", "id", object.Id, "requestId", reqId)

		slog.Error("request failed", "err", err)
		return err
	}

	return nil
}

func GenerateManualCAs(ctx context.Context, applications []CreditApplication) error {
	var url = os.Getenv("CAF_GATEWAY_URL")

	mu := sync.Mutex{}
	errApplications := make([]string, 0)
	var wg sync.WaitGroup

	for _, application := range applications {
		wg.Add(1)
		go func(application CreditApplication) {
			defer wg.Done()
			jsonData, err := json.Marshal(application)
			if err != nil {
				mu.Lock()
				errApplications = append(errApplications, application.Id)
				mu.Unlock()
				slog.Error("failed to marshal json", "id", application.Id)
				return
			}

			if err := doPostCARequest(url, jsonData); err != nil {
				mu.Lock()
				errApplications = append(errApplications, application.Id)
				mu.Unlock()
				slog.Error("request failed", "err", err)
				return
			}
		}(application)
	}

	wg.Wait()

	if len(errApplications) != 0 {
		return fmt.Errorf("given applications was not sent: %v", errApplications)
	}

	return nil
}

func GenerateLoadCAs(count int, objectType string) error {
	var url string

	reqChan := make(chan []byte, count)
	objects := make([][]byte, count)

	var err error

	switch objectType {
	case "credit-application":
		url = "https://caf.baraiq.io/api/gtwsvc/async/credit-application"
		objects, err = generateRandomCAs(count)
		if err != nil {
			return err
		}
	case "transaction":
		url = "https://taf.baraiq.io/api/gtwsvc/async/transaction"
		objects, err = generateRandomTxns(count)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported object type")
	}

	for _, object := range objects {
		reqChan <- object
	}

	startTime := time.Now().UTC()
	slog.Info("load ca-test:", "type", objectType, "start at", startTime)

	workers := 80
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for object := range reqChan {
				_ = doPostCARequest(url, object)
			}
		}()
	}

	close(reqChan)
	wg.Wait()

	duration := int(time.Since(startTime).Nanoseconds())

	var rps float64
	if duration != 0 {
		rps = float64(count) / (float64(duration) / float64(1000000000))
	}

	slog.Info("load ca-test:", "count", count, "duration(ns)", duration, "rps", rps)
	return nil
}
