package models

// Tables
type Product struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

// History of all prices
type Cost struct {
	ID        int64   `json:"id" gorm:"primary_key"`
	Cost      float64 `json:"cost,string" sql:"type:decimal(10,2);"`
	Active    bool    `json:"active"`
	ProductId int64
	Product   Product `gorm:"ForeignKey:ProductId;References:ID"`
}

type Order struct {
	ID         int64  `json:"id" gorm:"primary_key"`
	ClientName string `json:"client_name"`
	Paid       bool   `json:"paid"`
}

type OrderMapping struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	OrderId   int64
	ProductId int64
	PriceId   int64
	Order     Order   `gorm:"ForeignKey:OrderId;References:ID"`
	Product   Product `gorm:"ForeignKey:ProductId;References:ID"`
	Cost      Cost    `gorm:"ForeignKey:PriceId;References:ID"`
}

// Responses, Requests
type ProductPrice struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost,string" sql:"type:decimal(10,2);"`
}

type Items struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	SumCost  float64 `json:"sum_cost"`
}

type Receipt struct {
	Items          []Items `json:"items"`
	TotalCost      float64 `json:"total_cost"`
	ChangeReturned float64 `json:"change_returned"`
}
