package service

import (
	"errors"
	"inventory-service/internal/model"
	"sync"
	"time"
)

var (
	ErrSKUNotFound    = errors.New("sku not found")
	ErrSetNotFound    = errors.New("set product not found")
	ErrStockNotEnough = errors.New("insufficient stock")
)

type InventoryService struct {
	mu        sync.RWMutex
	skus      map[string]*model.SKU
	sets      map[string]*model.SetProduct
	skuToSets map[string][]string
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		skus:      make(map[string]*model.SKU),
		sets:      make(map[string]*model.SetProduct),
		skuToSets: make(map[string][]string),
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
	for _, item := range sp.Items {
		s.skuToSets[item.SKUID] = append(s.skuToSets[item.SKUID], sp.ID)
	}
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

	return s.calcSetStockLocked(set)
}

func (s *InventoryService) DeductSKUStock(skuID string, quantity int) (*model.StockInfo, []model.StockInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sku, ok := s.skus[skuID]
	if !ok {
		return nil, nil, ErrSKUNotFound
	}

	if sku.Stock < quantity {
		return nil, nil, ErrStockNotEnough
	}

	sku.Stock -= quantity
	sku.UpdatedAt = time.Now()

	result := &model.StockInfo{
		ID:    sku.ID,
		Name:  sku.Name,
		Type:  model.TypeSKU,
		Stock: sku.Stock,
	}

	var affectedSets []model.StockInfo
	setIDs := s.skuToSets[skuID]
	for _, setID := range setIDs {
		set, ok := s.sets[setID]
		if !ok {
			continue
		}
		setStock, err := s.calcSetStockLocked(set)
		if err != nil {
			continue
		}
		affectedSets = append(affectedSets, *setStock)
	}

	return result, affectedSets, nil
}

func (s *InventoryService) DeductSetStock(setID string, quantity int) (*model.StockInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	set, ok := s.sets[setID]
	if !ok {
		return nil, ErrSetNotFound
	}

	for _, item := range set.Items {
		sku, ok := s.skus[item.SKUID]
		if !ok {
			return nil, ErrSKUNotFound
		}
		need := item.Quantity * quantity
		if sku.Stock < need {
			return nil, ErrStockNotEnough
		}
	}

	for _, item := range set.Items {
		sku := s.skus[item.SKUID]
		sku.Stock -= item.Quantity * quantity
		sku.UpdatedAt = time.Now()
	}

	set.UpdatedAt = time.Now()

	return s.calcSetStockLocked(set)
}

func (s *InventoryService) UpdateSKUStock(skuID string, stock int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sku, ok := s.skus[skuID]
	if !ok {
		return ErrSKUNotFound
	}

	sku.Stock = stock
	sku.UpdatedAt = time.Now()
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
