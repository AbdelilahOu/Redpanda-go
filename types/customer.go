package types

type Customer struct {
	CustomerID int    `json:"CustomerID"`
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Email      string `json:"Email"`
	Phone      string `json:"Phone"`
	Address    string `json:"Address"`
	City       string `json:"City"`
	State      string `json:"State"`
	ZipCode    string `json:"ZipCode"`
}
