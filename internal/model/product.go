package model

import "time"

type ProductType string

const (
	TypeSKU  ProductType = "sku"
	TypeSet  ProductType = "set"
)

type SKU struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SetItem struct {
	SKUID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
}

type SetProduct struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Items     []SetItem `json:"items"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockInfo struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     ProductType `json:"type"`
	Stock    int         `json:"stock"`
	Detail   interface{} `json:"detail,omitempty"`
}

type SetStockDetail struct {
	SetID string           `json:"set_id"`
	Items []SetItemStockDetail `json:"items"`
}

type SetItemStockDetail struct {
	SKUID      string `json:"sku_id"`
	SKUName    string `json:"sku_name"`
	Required   int    `json:"required"`
	Available  int    `json:"available"`
	MaxSets    int    `json:"max_sets"`
}
