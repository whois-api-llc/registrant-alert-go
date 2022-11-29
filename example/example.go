package example

import (
	"context"
	"errors"
	registrantalert "github.com/whois-api-llc/registrant-alert-go"
	"log"
	"net/http"
	"time"
)

func BasicPreview(apikey string) {
	client := registrantalert.NewBasicClient(apikey)

	// Get a number of domains matching the criteria.
	domainsCount, _, err := client.BasicPreview(context.Background(),
		// specify the including and excluding search terms
		&registrantalert.BasicSearchTerms{[]string{"Airbnb", "US"}, []string{"Europe", "EU"}},
	)

	if err != nil {
		// Handle error message returned by server
		var apiErr *registrantalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	log.Println(domainsCount)
}

func BasicPurchase(apikey string) {
	client := registrantalert.NewClient(apikey, registrantalert.ClientParams{
		HTTPClient: &http.Client{
			Transport: nil,
			Timeout:   40 * time.Second,
		},
	})

	// Get parsed Registrant Alert API response as a model instance.
	registrantAlertResp, resp, err := client.BasicPurchase(context.Background(),
		// specify the including search terms, excluding search terms can be unspecified
		&registrantalert.BasicSearchTerms{[]string{"Airbnb", "US"}, nil},
		// this option is ignored, as the inner parser works with JSON only
		registrantalert.OptionResponseFormat("XML"),
		// this option results in domain names in the response will be encoded to Punycode
		registrantalert.OptionPunycode(true),
		// this option results in search through activities discovered since the given date
		registrantalert.OptionSinceDate(time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC)))

	if err != nil {
		// Handle error message returned by server
		var apiErr *registrantalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	// Then print all "added" domains.
	for _, obj := range registrantAlertResp.DomainsList {
		if obj.Action == registrantalert.Added {
			log.Println(obj.DomainName, obj.Action, time.Time(obj.Date).Format("2006-01-02"))
		}
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func BasicRawData(apikey string) {
	client := registrantalert.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.BasicRawData(context.Background(),
		// specify the including search terms
		&registrantalert.BasicSearchTerms{[]string{"Google"}, nil},
		// specify the domain-related dates to search through
		registrantalert.OptionCreatedDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionCreatedDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionUpdatedDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionUpdatedDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionExpiredDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionExpiredDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
	)

	if err != nil {
		// Handle error message returned by server
		log.Println(err)
	}

	if resp != nil {
		log.Println(string(resp.Body))
	}
}

func AdvancedPreview(apikey string) {
	client := registrantalert.NewBasicClient(apikey)

	// Get a number of domains matching the criteria.
	domainsCount, _, err := client.AdvancedPreview(context.Background(),
		// specify the advanced search terms
		[]registrantalert.AdvancedSearchTerm{
			{"RegistrantContact.Organization", "Airbnb, Inc.", true},
			{"RegistrantContact.Country", "UNITED STATES", false}})

	if err != nil {
		// Handle error message returned by server
		var apiErr *registrantalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	log.Println(domainsCount)
}

func AdvancedPurchase(apikey string) {
	client := registrantalert.NewClient(apikey, registrantalert.ClientParams{
		HTTPClient: &http.Client{
			Transport: nil,
			Timeout:   40 * time.Second,
		},
	})

	// Get parsed Registrant Alert API response as a model instance.
	registrantAlertResp, resp, err := client.AdvancedPurchase(context.Background(),
		// specify the advanced search terms
		[]registrantalert.AdvancedSearchTerm{{"RegistrantContact.Organization", "Airbnb, Inc.", true}},
		// this option is ignored, as the inner parser works with JSON only
		registrantalert.OptionResponseFormat("XML"),
		// this option results in search through activities discovered since the given date
		registrantalert.OptionSinceDate(time.Date(2022, 10, 31, 0, 0, 0, 0, time.UTC)))

	if err != nil {
		// Handle error message returned by server
		var apiErr *registrantalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	// Then print all "dropped" domains.
	for _, obj := range registrantAlertResp.DomainsList {
		if obj.Action == registrantalert.Dropped {
			log.Println(obj.DomainName, obj.Action, time.Time(obj.Date).Format("2006-01-02"))
		}
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func AdvancedRawData(apikey string) {
	client := registrantalert.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.AdvancedRawData(context.Background(),
		// specify the including search terms
		[]registrantalert.AdvancedSearchTerm{{"RegistrantContact.Organization", "Airbnb", false}},
		// this option results in search through activities discovered since the given date
		registrantalert.OptionSinceDate(time.Date(2022, 10, 01, 0, 0, 0, 0, time.UTC)),
		// specify the domain-related dates to search through
		registrantalert.OptionCreatedDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionCreatedDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionUpdatedDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionUpdatedDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionExpiredDateFrom(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		registrantalert.OptionExpiredDateTo(time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)),
	)

	if err != nil {
		// Handle error message returned by server
		log.Println(err)
	}

	if resp != nil {
		log.Println(string(resp.Body))
	}
}
