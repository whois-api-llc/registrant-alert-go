package registrantalert

import (
	"time"
)

// Option adds parameters to the query.
type Option func(v *registrantAlertRequest)

var _ = []Option{
	OptionResponseFormat("JSON"),
	OptionSinceDate(time.Time{}),
	OptionPunycode(true),
	OptionCreatedDateFrom(time.Time{}),
	OptionCreatedDateTo(time.Time{}),
	OptionUpdatedDateFrom(time.Time{}),
	OptionUpdatedDateTo(time.Time{}),
	OptionExpiredDateFrom(time.Time{}),
	OptionExpiredDateTo(time.Time{}),
}

// OptionResponseFormat sets Response output format json | xml. Default: json.
func OptionResponseFormat(outputFormat string) Option {
	return func(v *registrantAlertRequest) {
		v.ResponseFormat = outputFormat
	}
}

// OptionSinceDate results in search through activities discovered since the given date.
func OptionSinceDate(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.SinceDate = date.Format(dateFormat)
	}
}

// OptionPunycode sets the punycode option.
// If true, domain names in the response will be encoded to Punycode. Default: true.
func OptionPunycode(punycode bool) Option {
	return func(v *registrantAlertRequest) {
		v.Punycode = punycode
	}
}

// OptionCreatedDateFrom searches through domains created after the given date.
func OptionCreatedDateFrom(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.CreatedDateFrom = date.Format(dateFormat)
	}
}

// OptionCreatedDateTo searches through domains created before the given date.
func OptionCreatedDateTo(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.CreatedDateTo = date.Format(dateFormat)
	}
}

// OptionUpdatedDateFrom searches through domains updated after the given date.
func OptionUpdatedDateFrom(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.UpdatedDateFrom = date.Format(dateFormat)
	}
}

// OptionUpdatedDateTo searches through domains updated before the given date.
func OptionUpdatedDateTo(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.UpdatedDateTo = date.Format(dateFormat)
	}
}

// OptionExpiredDateFrom searches through domains expired after the given date.
func OptionExpiredDateFrom(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.ExpiredDateFrom = date.Format(dateFormat)
	}
}

// OptionExpiredDateTo searches through domains expired before the given date.
func OptionExpiredDateTo(date time.Time) Option {
	return func(v *registrantAlertRequest) {
		v.ExpiredDateTo = date.Format(dateFormat)
	}
}
