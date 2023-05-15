package f3account

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ValidUUID = "eb0bd6f5-c3f5-44b2-b677-bcd23cdde73c"

func TestFlow(t *testing.T) {
	// Create an account
	t.Log("Test the flow ... Create, Fetch and Delete...")
	a := NewBasicAccount(
		ValidUUID,
		ValidUUID,
		"GB",
		"123456",
		"GBDSC",
		"EXMPLGB2XXX",
		"TestName",
	)
	flib, err := instanceFakeLib()
	assert.Nil(t, err)

	res1, err := flib.Create(&a)
	assert.Nil(t, err)
	assert.Equal(t, res1.StatusCode, http.StatusCreated)

	// find the account created
	aFetched, err := flib.Fetch(a.ID)
	assert.Nil(t, err)
	assert.Equal(t, a.ID, aFetched.Data.ID)

	// delete de account
	aDeleted, err := flib.Delete(a.ID, int(*aFetched.Data.Version))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, aDeleted.StatusCode)
	// can not find the account
	aReFetched, err := flib.Fetch(a.ID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, aReFetched.StatusCode)
}

func TestServerNotReachable(t *testing.T) {
	t.Log("Test when the server is not reachable")
	flib, err := NewFakeLib("http://www.nonexist3.not")
	assert.Nil(t, err)
	_, err = flib.Fetch("someID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Get \"http://www.nonexist3.not/v1/organisation/accounts/someID\": dial tcp: lookup www.nonexist3.not")
}

func TestBadParameters(t *testing.T) {
	flib, err := instanceFakeLib()
	assert.Nil(t, err)

	tests := []struct {
		name           string
		id             string
		organisationID string
		statusCode     int
		country        string
		bankID         string
		bankIDCode     string
		bic            string
		accountName    string
	}{
		{
			name:           "ID is not an UUID",
			id:             "xxxx",
			organisationID: ValidUUID,
			statusCode:     http.StatusBadRequest,
			country:        "GB",
			bankID:         "123456",
			bankIDCode:     "GBDSC",
			bic:            "EXMPLGB2XXX",
			accountName:    "testName",
		},
		{
			name:           "OrganisationID is not an UUID",
			id:             ValidUUID,
			organisationID: "xxxx",
			statusCode:     http.StatusBadRequest,
			country:        "GB",
			bankID:         "123456",
			bankIDCode:     "GBDSC",
			bic:            "EXMPLGB2XXX",
			accountName:    "testName",
		},
		{
			name:           "Country empty",
			id:             ValidUUID,
			organisationID: ValidUUID,
			statusCode:     http.StatusBadRequest,
			country:        "",
			bankID:         "123456",
			bankIDCode:     "GBDSC",
			bic:            "EXMPLGB2XXX",
			accountName:    "testName",
		},
		{
			name:           "Name is empty",
			id:             ValidUUID,
			organisationID: ValidUUID,
			statusCode:     http.StatusBadRequest,
			country:        "GB",
			bankID:         "123456",
			bankIDCode:     "GBDSC",
			bic:            "EXMPLGB2XXX",
			accountName:    "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := NewBasicAccount(
				test.id,
				test.organisationID,
				test.country,
				test.bankID,
				test.bankIDCode,
				test.bic,
				test.accountName,
			)
			res, err := flib.Create(&a)
			assert.Error(t, err)
			assert.Equal(t, test.statusCode, res.StatusCode)
			assert.Empty(t, res.Data)
		})
	}
}

func instanceFakeLib() (*FakeLib, error) {
	if baseURL := os.Getenv("F3_BASE_URL"); baseURL == "" {
		_ = os.Setenv("F3_BASE_URL", "http://localhost:8080")
	}
	flib, err := NewFakeLib(os.Getenv("F3_BASE_URL"))
	return flib, err
}
