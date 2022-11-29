package registrantalert

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values *registrantAlertRequest
		option Option
		want   string
	}{
		{
			name:   "responseFormat",
			values: &registrantAlertRequest{},
			option: OptionResponseFormat("json"),
			want:   "json",
		},
		{
			name:   "sinceDate",
			values: &registrantAlertRequest{},
			option: OptionSinceDate(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
			want:   "2021-01-01",
		},
		{
			name:   "punycode",
			values: &registrantAlertRequest{},
			option: OptionPunycode(false),
			want:   "false",
		},
		{
			name:   "createdDateFrom",
			values: &registrantAlertRequest{},
			option: OptionCreatedDateFrom(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)),
			want:   "2022-01-01",
		},
		{
			name:   "createdDateTo",
			values: &registrantAlertRequest{},
			option: OptionCreatedDateTo(time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)),
			want:   "2021-02-01",
		},
		{
			name:   "updatedDateFrom",
			values: &registrantAlertRequest{},
			option: OptionUpdatedDateFrom(time.Date(2021, 01, 03, 0, 0, 0, 0, time.UTC)),
			want:   "2021-01-03",
		},
		{
			name:   "updatedDateTo",
			values: &registrantAlertRequest{},
			option: OptionUpdatedDateTo(time.Date(2021, 01, 01, 1, 0, 0, 0, time.UTC)),
			want:   "2021-01-01",
		},
		{
			name:   "expiredDateFrom",
			values: &registrantAlertRequest{},
			option: OptionExpiredDateFrom(time.Date(2021, 01, 01, 0, 2, 0, 0, time.UTC)),
			want:   "2021-01-01",
		},
		{
			name:   "expiredDateTo",
			values: &registrantAlertRequest{},
			option: OptionExpiredDateTo(time.Date(2021, 01, 01, 0, 0, 3, 0, time.UTC)),
			want:   "2021-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			tt.option(tt.values)

			switch tt.name {
			case "responseFormat":
				got = tt.values.ResponseFormat
			case "sinceDate":
				got = tt.values.SinceDate
			case "createdDateFrom":
				got = tt.values.CreatedDateFrom
			case "createdDateTo":
				got = tt.values.CreatedDateTo
			case "updatedDateFrom":
				got = tt.values.UpdatedDateFrom
			case "updatedDateTo":
				got = tt.values.UpdatedDateTo
			case "expiredDateFrom":
				got = tt.values.ExpiredDateFrom
			case "expiredDateTo":
				got = tt.values.ExpiredDateTo
			case "punycode":
				got = strconv.FormatBool(tt.values.Punycode)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
