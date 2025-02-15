package repository

import (
	"context"
	"errors"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
)

var ErrMerchItemNotFound = errors.New("merch item not found")

type MerchRepository struct {
	merchItems map[string]model.Merch
}

func NewMerchRepository() *MerchRepository {
	// Initialize with predefined merch items
	items := map[string]model.Merch{
		"t-shirt":    {Name: "t-shirt", Price: 80},
		"cup":        {Name: "cup", Price: 20},
		"book":       {Name: "book", Price: 50},
		"pen":        {Name: "pen", Price: 10},
		"powerbank":  {Name: "powerbank", Price: 200},
		"hoody":      {Name: "hoody", Price: 300},
		"umbrella":   {Name: "umbrella", Price: 200},
		"socks":      {Name: "socks", Price: 10},
		"wallet":     {Name: "wallet", Price: 50},
		"pink-hoody": {Name: "pink-hoody", Price: 500},
	}
	return &MerchRepository{merchItems: items}
}

func (r *MerchRepository) GetMerchItemByName(ctx context.Context, itemName string) (model.Merch, error) {
	item, exists := r.merchItems[itemName]
	if !exists {
		return model.Merch{}, ErrMerchItemNotFound
	}
	return item, nil
}

func (r *MerchRepository) ListMerchItems(ctx context.Context) ([]model.Merch, error) {
	items := make([]model.Merch, 0, len(r.merchItems))
	for _, item := range r.merchItems {
		items = append(items, item)
	}
	return items, nil
}
