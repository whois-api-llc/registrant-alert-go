package registrantalert

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathRegistrantAlertResponseOK         = "/RegistrantAlert/ok"
	pathRegistrantAlertResponseError      = "/RegistrantAlert/error"
	pathRegistrantAlertResponse500        = "/RegistrantAlert/500"
	pathRegistrantAlertResponsePartial1   = "/RegistrantAlert/partial"
	pathRegistrantAlertResponsePartial2   = "/RegistrantAlert/partial2"
	pathRegistrantAlertResponseUnparsable = "/RegistrantAlert/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Registrant Alert API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathRegistrantAlertResponseOK:
		case pathRegistrantAlertResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathRegistrantAlertResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathRegistrantAlertResponsePartial1:
			response = response[:len(response)-10]
		case pathRegistrantAlertResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathRegistrantAlertResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Registrant Alert API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:             apiServer.Client(),
		RegistrantAlertBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestRegistrantAlertBasicPreview tests the BasicPreview function.
func TestRegistrantAlertBasicPreview(t *testing.T) {
	checkResultRec := func(res int) bool {
		return res != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":["Test error message."]}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory *BasicSearchTerms
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.include" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.include" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, []string{"1", "2", "3", "4", "5"}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.exclude" must have between 0 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.BasicPreview(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.BasicPreview() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("RegistrantAlert.BasicPreview() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != 0 {
					t.Errorf("RegistrantAlert.BasicPreview() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestRegistrantAlertBasicPurchase tests the BasicPurchase function.
func TestRegistrantAlertBasicPurchase(t *testing.T) {
	checkResultRec := func(res *RegistrantAlertResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":["Test error message."]}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory *BasicSearchTerms
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{}, nil},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.include" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.include" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, []string{"1", "2", "3", "4", "5"}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "basicSearchTerms.exclude" must have between 0 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.BasicPurchase(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.BasicPurchase() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("RegistrantAlert.BasicPurchase() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("RegistrantAlert.BasicPurchase() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestRegistrantAlertBasicRawData tests the BasicRawData function.
func TestRegistrantAlertBasicRawData(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory *BasicSearchTerms
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("xml"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{}, nil},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "basicSearchTerms.include" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "basicSearchTerms.include" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&BasicSearchTerms{[]string{"whois"}, []string{"1", "2", "3", "4", "5"}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "basicSearchTerms.exclude" must have between 0 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.BasicRawData(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.BasicRawData() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("RegistrantAlert.BasicRawData() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}

// TestRegistrantAlertAdvancedPreview tests the AdvancedPreview function.
func TestRegistrantAlertAdvancedPreview(t *testing.T) {
	checkResultRec := func(res int) bool {
		return res != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":["Test error message."]}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory []AdvancedSearchTerm
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"", "", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms.0.Field" is required.`,
		},
		{
			name: "invalid argument4",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{
						{"RegistrantContact.Organization", "whois1", false},
						{"RegistrantContact.Organization", "whois2", false},
						{"RegistrantContact.Organization", "whois3", false},
						{"RegistrantContact.Organization", "whois4", false},
						{"RegistrantContact.Organization", "whois5", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.AdvancedPreview(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.AdvancedPreview() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("RegistrantAlert.AdvancedPreview() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != 0 {
					t.Errorf("RegistrantAlert.AdvancedPreview() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestRegistrantAlertAdvancedPurchase tests the AdvancedPurchase function.
func TestRegistrantAlertAdvancedPurchase(t *testing.T) {
	checkResultRec := func(res *RegistrantAlertResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":["Test error message."]}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory []AdvancedSearchTerm
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"", "", false}},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "advancedSearchTerms.0.Field" is required.`,
		},
		{
			name: "invalid argument4",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{
						{"RegistrantContact.Organization", "whois1", false},
						{"RegistrantContact.Organization", "whois2", false},
						{"RegistrantContact.Organization", "whois3", false},
						{"RegistrantContact.Organization", "whois4", false},
						{"RegistrantContact.Organization", "whois5", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.AdvancedPurchase(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.AdvancedPurchase() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("RegistrantAlert.AdvancedPurchase() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("RegistrantAlert.AdvancedPurchase() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestRegistrantAlertAdvancedRawData tests the AdvancedRawData function.
func TestRegistrantAlertAdvancedRawData(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory []AdvancedSearchTerm
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathRegistrantAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathRegistrantAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathRegistrantAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathRegistrantAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathRegistrantAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("xml"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"RegistrantContact.Organization", "whois", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms" is required.`,
		},
		{
			name: "invalid argument3",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{{"", "", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms.0.Field" is required.`,
		},
		{
			name: "invalid argument4",
			path: pathRegistrantAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					[]AdvancedSearchTerm{
						{"RegistrantContact.Organization", "whois1", false},
						{"RegistrantContact.Organization", "whois2", false},
						{"RegistrantContact.Organization", "whois3", false},
						{"RegistrantContact.Organization", "whois4", false},
						{"RegistrantContact.Organization", "whois5", false}},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "advancedSearchTerms" must have between 1 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.AdvancedRawData(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("RegistrantAlert.AdvancedRawData() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("RegistrantAlert.AdvancedRawData() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
