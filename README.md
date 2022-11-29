[![registrant-alert-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![registrant-alert-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/registrant-alert-go)
[![registrant-alert-go test](https://github.com/whois-api-llc/registrant-alert-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/registrant-alert-go/actions/)

# Overview

The client library for
[Registrant Alert API](https://registrant-alert.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/registrant-alert-go
```

# Examples

Full API documentation available [here](https://registrant-alert.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := registrantalert.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := registrantalert.NewClient(apiKey, registrantalert.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Registrant Alert API lets you monitor specific domain registrants to be alerted whenever their information is linked to a newly-registered or just-expired domain name.

Basic search requires less configuration and produces broader results.

```go

// Make request to get a list of all domains matching the criteria.
registrantAlertResp, _, err := client.BasicPurchase(ctx,
    &registrantalert.BasicSearchTerms{[]string{"Airbnb", "US"}, []string{"Europe", "EU"}})

for _, obj := range registrantAlertResp.DomainsList {
    log.Println(obj.DomainName)
}


// Make request to get only domains count.
domainsCount, _, err := client.BasicPreview(ctx,
    &registrantalert.BasicSearchTerms{[]string{"Google"}})

log.Println(domainsCount)

// Make request to get raw data in XML.
resp, err := client.BasicRawData(ctx,
    &registrantalert.SearchTerms{"google", "blog"},
    &registrantalert.SearchTerms{"analytics"},
    registrantalert.OptionResponseFormat("XML"))

log.Println(string(resp.Body))

```

## Advanced usage
Advanced search allows searching through specific WHOIS fields.

```go
registrantAlertResp, resp, err := client.AdvancedPurchase(ctx,
    []registrantalert.AdvancedSearchTerm{{"RegistrantContact.Organization", "Airbnb, Inc.", true}},
    registrantalert.OptionSinceDate(time.Date(2022, 11, 01, 0, 0, 0, 0, time.UTC)))

```