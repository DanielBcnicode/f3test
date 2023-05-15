package f3account

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type Account struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func NewBasicAccount(
	ID string,
	organisationID string,
	country string,
	bankID string,
	bankIDCode string,
	bic string,
	name ...string,
) Account {
	countryLocal := country
	att := AccountAttributes{}
	acc := Account{}
	acc.Attributes = &att
	acc.ID = ID
	acc.Type = "accounts"
	acc.OrganisationID = organisationID
	acc.Attributes.Country = &countryLocal
	acc.Attributes.BankID = bankID
	acc.Attributes.BankIDCode = bankIDCode
	acc.Attributes.Bic = bic
	acc.Attributes.Name = name

	return acc
}

type ResponseLinks struct {
	Self string `json:"self,omitempty"`
}
type FakeLibResponse struct {
	Data         Account        `json:"data,omitempty"`
	Links        *ResponseLinks `json:"links,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	StatusCode   int            `json:"status_code,omitempty"`
}

type FakeLibRequest struct {
	Data Account `json:"data,omitempty"`
}
