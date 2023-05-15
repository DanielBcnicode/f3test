package f3account

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrorMarshalingAccount    = errors.New("error marshalling account data")
	ErrorCreatingUrl          = errors.New("url malformed")
	ErrorCreatingPostRequest  = errors.New("post request malformed")
	ErrorUnMarshalingResponse = errors.New("error unmarshalling")
	ErrorInAPICall            = errors.New("error in ret2 fake API call")
	ErrorCreatingGetRequest   = errors.New("get request malformed")
	ErrorEmptyBaseUrl         = errors.New("BaseUrl can not be empty")
)

// F3FAccountLib main library interface
type F3FAccountLib interface {
	Create(account *Account) (FakeLibResponse, error)
	Fetch(id string) (FakeLibResponse, error)
	Delete(id string, version int) (FakeLibResponse, error)
}

// FakeLib main library struct
type FakeLib struct {
	baseUrl string
}

// NewFakeLib is the FakeLib constructor
func NewFakeLib(baseUrl string) (*FakeLib, error) {
	if strings.TrimSpace(baseUrl) == "" {
		return nil, ErrorEmptyBaseUrl
	}
	lib := FakeLib{baseUrl: baseUrl}

	return &lib, nil
}

// Create is the bridge function to the F3 create account
// Returns the FakeLibResponse and error.
func (f *FakeLib) Create(a *Account) (FakeLibResponse, error) {
	fr := FakeLibRequest{Data: *a}
	body, err := json.Marshal(fr)
	if err != nil {
		return FakeLibResponse{}, ErrorMarshalingAccount
	}
	postUrl, err := url.JoinPath(f.baseUrl, "/v1/organisation/accounts")
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingUrl
	}

	// I need to use context and timeout
	r, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewBuffer(body))
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingPostRequest
	}
	r.Header.Add("Content-Type", "application/json")
	c := http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return FakeLibResponse{}, err
	}
	defer res.Body.Close()

	accountResponse := FakeLibResponse{}
	err = json.NewDecoder(res.Body).Decode(&accountResponse)
	if err != nil {
		return FakeLibResponse{}, ErrorUnMarshalingResponse
	}
	accountResponse.StatusCode = res.StatusCode
	if res.StatusCode >= 400 {
		return accountResponse, ErrorInAPICall
	}

	return accountResponse, nil
}

// Fetch is the bridge to the fetch F3 account
func (f *FakeLib) Fetch(Id string) (FakeLibResponse, error) {

	getUrl, err := url.JoinPath(f.baseUrl, "/v1/organisation/accounts", Id)
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingUrl
	}

	// I need to use context and timeout
	r, err := http.NewRequest(http.MethodGet, getUrl, nil)
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingGetRequest
	}
	r.Header.Add("Content-Type", "application/json")
	c := http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return FakeLibResponse{}, err
	}
	defer res.Body.Close()

	accountResponse := FakeLibResponse{}
	err = json.NewDecoder(res.Body).Decode(&accountResponse)
	if err != nil {
		return FakeLibResponse{}, ErrorUnMarshalingResponse
	}
	accountResponse.StatusCode = res.StatusCode
	if res.StatusCode >= 500 {
		return accountResponse, ErrorInAPICall
	}

	return accountResponse, nil
}

// Delete is the bridge function to the delete original F3 delete account
func (f *FakeLib) Delete(id string, version int) (FakeLibResponse, error) {

	deleteUrl, err := url.JoinPath(f.baseUrl, "/v1/organisation/accounts", id)
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingUrl
	}
	baseUrl, _ := url.Parse(deleteUrl)
	params := url.Values{}
	params.Add("version", strconv.Itoa(version))
	baseUrl.RawQuery = params.Encode()
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingUrl
	}

	// I need to use context and timeout
	r, err := http.NewRequest(http.MethodDelete, baseUrl.String(), nil)
	if err != nil {
		return FakeLibResponse{}, ErrorCreatingGetRequest
	}
	r.Header.Add("Content-Type", "application/json")
	c := http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return FakeLibResponse{}, err
	}
	defer res.Body.Close()

	accountResponse := FakeLibResponse{}
	if res.ContentLength > 0 {
		err = json.NewDecoder(res.Body).Decode(&accountResponse)
		if err != nil {
			return FakeLibResponse{}, ErrorUnMarshalingResponse
		}
	}
	accountResponse.StatusCode = res.StatusCode
	if res.StatusCode >= 400 {
		return accountResponse, ErrorInAPICall
	}

	return accountResponse, nil
}
