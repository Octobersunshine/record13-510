package main

import (
	"fmt"
	"inventory-service/internal/handler"
	"inventory-service/internal/model"
	"inventory-service/internal/service"
	"log"
	"net/http"
	"time"
)

func main() {
	svc := service.NewInventoryService()

	initTestData(svc)

	h := handler.NewInventoryHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /api/stock?id={id}          - query stock by id")
	fmt.Println("  GET /api/stock/sku/{sku_id}     - query sku stock")
	fmt.Println("  GET /api/stock/set/{set_id}     - query set stock")
	fmt.Println("  GET /api/stock/all              - list all stock")

	log.Fatal(http.ListenAndServe(addr, mux))
}

func initTestData(svc *service.InventoryService) {
	now := time.Now()

	skus := []*model.SKU{
		{ID: "sku_001", Name: "白色T恤", Stock: 100, CreatedAt: now, UpdatedAt: now},
		{ID: "sku_002", Name: "蓝色牛仔裤", Stock: 50, CreatedAt: now, UpdatedAt: now},
		{ID: "sku_003", Name: "运动袜子", Stock: 200, CreatedAt: now, UpdatedAt: now},
		{ID: "sku_004", Name: "黑色运动鞋", Stock: 30, CreatedAt: now, UpdatedAt: now},
		{ID: "sku_005", Name: "红色卫衣", Stock: 80, CreatedAt: now, UpdatedAt: now},
	}

	for _, sku := range skus {
		svc.AddSKU(sku)
	}

	sets := []*model.SetProduct{
		{
			ID:   "set_001",
			Name: "休闲套装（T恤+牛仔裤）",
			Items: []model.SetItem{
				{SKUID: "sku_001", Quantity: 1},
				{SKUID: "sku_002", Quantity: 1},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:   "set_002",
			Name: "运动三件套（T恤+袜子+运动鞋）",
			Items: []model.SetItem{
				{SKUID: "sku_001", Quantity: 1},
				{SKUID: "sku_003", Quantity: 2},
				{SKUID: "sku_004", Quantity: 1},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:   "set_003",
			Name: "卫衣套装（卫衣+牛仔裤）",
			Items: []model.SetItem{
				{SKUID: "sku_005", Quantity: 1},
				{SKUID: "sku_002", Quantity: 1},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, set := range sets {
		svc.AddSetProduct(set)
	}

	fmt.Println("Test data initialized:")
	fmt.Println("  SKUs: 5 items")
	fmt.Println("  Set products: 3 sets")
	fmt.Println()
}
