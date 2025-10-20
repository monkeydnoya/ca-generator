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
	Timeout: 30 * time.Second,
}

// var bufPool = sync.Pool{
// 	New: func() any {
// 		return new(bytes.Buffer)
// 	},
// }

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
				IINBIN:               "830622350419",
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

func doPostCARequest(requestUrl string, data []byte) error {
	// Slice of bytes -> buffer -> http.Response
	_, err := client.Post(requestUrl, "application/json", bytes.NewReader(data))
	if err != nil {
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

func GenerateLoadCAs(ctx context.Context, count int) error {
	url := "https://caf.baraiq.io/api/gtwsvc/async/credit-application"

	reqChan := make(chan []byte, count)
	applications, err := generateRandomCAs(count)
	if err != nil {
		return err
	}

	for _, application := range applications {
		reqChan <- application
	}

	startTime := time.Now().UTC()
	slog.Info("load ca-test:", "start at", startTime)

	workers := 1000
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for application := range reqChan {
				_ = doPostCARequest(url, application)
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
