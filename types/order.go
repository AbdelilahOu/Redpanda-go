package types

type Order struct {
    OrderID     int         `json:"OrderID"`
    CustomerID  int         `json:"CustomerID"`
    OrderDate   interface{}   `json:"OrderDate"`
    TotalAmount float64     `json:"TotalAmount"`
    Status      string      `json:"Status"`
}
