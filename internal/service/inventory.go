package service

import (
	"errors"
	"inventory-service/internal/model"
	"sync"
)

var (
	ErrSKUNotFound  = errors.New("sku not found")
	ErrSetNotFound  = errors.New("set product not found")
)

type InventoryService struct {
	mu       sync.RWMutex
	skus     map[string]*model.SKU
	sets     map[string]*model.SetProduct
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		skus: make(map[string]*model.SKU),
		sets: make(map[string]*model.SetProduct),
	}
}

func (s *InventoryService) AddSKU(sku *model.SKU) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.skus[sku.ID] = sku
}

func (s *InventoryService) AddSetProduct(sp *model.SetProduct) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sets[sp.ID] = sp
}

func (s *InventoryService) GetSKUStock(skuID string) (*model.StockInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sku, ok := s.skus[skuID]
	if !ok {
		return nil, ErrSKUNotFound
	}

	return &model.StockInfo{
		ID:    sku.ID,
		Name:  sku.Name,
		Type:  model.TypeSKU,
		Stock: sku.Stock,
	}, nil
}

func (s *InventoryService) GetSetStock(setID string) (*model.StockInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set, ok := s.sets[setID]
	if !ok {
		return nil, ErrSetNotFound
	}

	minStock := -1
	var details []model.SetItemStockDetail

	for _, item := range set.Items {
		sku, ok := s.skus[item.SKUID]
		if !ok {
			return nil, ErrSKUNotFound
		}

		maxSets := sku.Stock / item.Quantity

		details = append(details, model.SetItemStockDetail{
			SKUID:     sku.ID,
			SKUName:   sku.Name,
			Required:  item.Quantity,
			Available: sku.Stock,
			MaxSets:   maxSets,
		})

		if minStock == -1 || maxSets < minStock {
			minStock = maxSets
		}
	}

	if minStock < 0 {
		minStock = 0
	}

	return &model.StockInfo{
		ID:    set.ID,
		Name:  set.Name,
		Type:  model.TypeSet,
		Stock: minStock,
		Detail: model.SetStockDetail{
			SetID: set.ID,
			Items: details,
		},
	}, nil
}

func (s *InventoryService) UpdateSKUStock(skuID string, stock int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sku, ok := s.skus[skuID]
	if !ok {
		return ErrSKUNotFound
	}

	sku.Stock = stock
	return nil
}

func (s *InventoryService) ListAllStock() ([]*model.StockInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.StockInfo

	for _, sku := range s.skus {
		result = append(result, &model.StockInfo{
			ID:    sku.ID,
			Name:  sku.Name,
			Type:  model.TypeSKU,
			Stock: sku.Stock,
		})
	}

	for _, set := range s.sets {
		stockInfo, err := s.calcSetStockLocked(set)
		if err != nil {
			return nil, err
		}
		result = append(result, stockInfo)
	}

	return result, nil
}

func (s *InventoryService) calcSetStockLocked(set *model.SetProduct) (*model.StockInfo, error) {
	minStock := -1
	var details []model.SetItemStockDetail

	for _, item := range set.Items {
		sku, ok := s.skus[item.SKUID]
		if !ok {
			return nil, ErrSKUNotFound
		}

		maxSets := sku.Stock / item.Quantity

		details = append(details, model.SetItemStockDetail{
			SKUID:     sku.ID,
			SKUName:   sku.Name,
			Required:  item.Quantity,
			Available: sku.Stock,
			MaxSets:   maxSets,
		})

		if minStock == -1 || maxSets < minStock {
			minStock = maxSets
		}
	}

	if minStock < 0 {
		minStock = 0
	}

	return &model.StockInfo{
		ID:    set.ID,
		Name:  set.Name,
		Type:  model.TypeSet,
		Stock: minStock,
		Detail: model.SetStockDetail{
			SetID: set.ID,
			Items: details,
		},
	}, nil
}

func (s *InventoryService) GetAllSKUs() map[string]*model.SKU {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*model.SKU)
	for k, v := range s.skus {
		result[k] = v
	}
	return result
}

func (s *InventoryService) GetAllSets() map[string]*model.SetProduct {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*model.SetProduct)
	for k, v := range s.sets {
		result[k] = v
	}
	return result
}
