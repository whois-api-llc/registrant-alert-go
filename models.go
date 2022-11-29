package registrantalert

import (
	"encoding/json"
	"fmt"
	"time"
)

func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time
type Time time.Time

var emptyTime Time

const dateFormat = "2006-01-02"

// UnmarshalJSON decodes time as Registrant Alert API does.
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as Registrant Alert API does.
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format(dateFormat) + `"`), nil
}

// BasicSearchTerms is a part of the Registrant Alert API request.
type BasicSearchTerms struct {
	// Include is an array of search strings.
	// All of them should be present in the domain's registrant details.
	Include []string `json:"include,omitempty"`

	// Exclude is an array of search strings.
	// All of them should NOT be present in the domain's registrant details.
	Exclude []string `json:"exclude,omitempty"`
}

// AdvancedSearchTerm is a part of the Registrant Alert API request.
type AdvancedSearchTerm struct {
	// Field is the WHOIS field to search in.
	Field string `json:"field"`

	// Term is the search string. Case insensitive.
	Term string `json:"term"`

	// ExactMatch defines whether the field should exactly match the search term.
	// If false, the field is allowed to contain a search term as a substring.
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// registrantAlertRequest is the request struct for Registrant Alert API.
type registrantAlertRequest struct {
	// APIKey is the user's API key.
	APIKey string `json:"apiKey"`

	// BasicSearchTerms is the set of search terms for the Basic search.
	BasicSearchTerms *BasicSearchTerms `json:"basicSearchTerms,omitempty"`

	// BasicSearchTerms is the set of search terms for the Advanced search.
	AdvancedSearchTerms []AdvancedSearchTerm `json:"advancedSearchTerms,omitempty"`

	// SinceDate If present, search through activities discovered since the given date.
	SinceDate string `json:"sinceDate,omitempty"`

	// Mode is the mode of the API call. Acceptable values: preview | purchase.
	Mode string `json:"mode,omitempty"`

	// Punycode If true, domain names in the response will be encoded to punycode.
	Punycode bool `json:"punycode,omitempty"`

	// ResponseFormat is the response output format JSON | XML.
	ResponseFormat string `json:"responseFormat,omitempty"`

	// CreatedDateFrom If present, search through domains created after the given date.
	CreatedDateFrom string `json:"createdDateFrom,omitempty"`

	// CreatedDateTo If present, search through domains created before the given date.
	CreatedDateTo string `json:"createdDateTo,omitempty"`

	// UpdatedDateFrom If present, search through domains updated after the given date.
	UpdatedDateFrom string `json:"updatedDateFrom,omitempty"`

	// UpdatedDateTo If present, search through domains updated before the given date.
	UpdatedDateTo string `json:"updatedDateTo,omitempty"`

	// ExpiredDateFrom If present, search through domains expired after the given date.
	ExpiredDateFrom string `json:"expiredDateFrom,omitempty"`

	// ExpiredDateFrom If present, search through domains expired before the given date.
	ExpiredDateTo string `json:"expiredDateTo,omitempty"`
}

// Action is a wrapper on string.
type Action string

// List of possible actions.
const (
	Added      Action = "added"
	Updated    Action = "updated"
	Dropped    Action = "dropped"
	Discovered Action = "discovered"
)

var _ = []Action{
	Added,
	Updated,
	Dropped,
	Discovered,
}

// DomainItem is a part of the Registrant Alert API response.
type DomainItem struct {
	// DomainName is the full domain name.
	DomainName string `json:"domainName"`

	// Action is the related action. Possible actions: added | updated | dropped | discovered.
	Action Action `json:"action"`

	// Date is the event date.
	Date Time `json:"date"`
}

// RegistrantAlertResponse is a response of Registrant Alert API.
type RegistrantAlertResponse struct {
	// DomainsList is the list of domains matching the criteria.
	DomainsList []DomainItem `json:"domainsList"`

	// DomainsCount is the number of domains matching the criteria.
	DomainsCount int `json:"domainsCount"`
}

// Messages is a wrapper on []string.
type Messages []string

// UnmarshalJSON decodes the error messages returned by Registrant Alert API.
func (m *Messages) UnmarshalJSON(b []byte) error {
	var msgs []string

	if err := json.Unmarshal(b, &msgs); err == nil {
		*m = msgs
		return nil
	}

	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	*m = append(*m, fmt.Sprintf("%s", x))
	return nil
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    int      `json:"code"`
	Message Messages `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
