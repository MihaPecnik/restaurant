package database

import (
	"github.com/MihaPecnik/restaurant/models"
	"gorm.io/gorm"
)

func (d *Database) CreateProduct(product models.ProductPrice) error {
	p := models.Product{
		Name: product.Name,
	}
	stmt := d.db.Session(&gorm.Session{
		PrepareStmt: true,
	})
	err := stmt.Transaction(func(tx *gorm.DB) error {
		// Creates a new item
		err := d.db.Table("products").Create(&p).Error
		if err != nil {
			return err
		}
		cost := models.Cost{
			Cost:      product.Cost,
			Active:    true,
			ProductId: p.ID,
		}
		// Creates a price for newly created item
		err = d.db.Table("costs").Create(&cost).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) ListProducts() ([]models.ProductPrice, error) {
	var response []models.ProductPrice
	query := `
	select products.*, c.cost from products
		inner join costs c on products.id = c.product_id and c.active=true
	`
	err := d.db.
		Raw(query).
		Scan(&response).Error
	if err != nil {
		return []models.ProductPrice{}, nil
	}
	return response, nil
}

func (d *Database) UpdateThePrice(id int64, cost float64) error {
	stmt := d.db.Session(&gorm.Session{
		PrepareStmt: true,
	})
	err := stmt.Transaction(func(tx *gorm.DB) error {
		// Set the current price to inactive
		err := d.db.Table("costs").
			Where("product_id = ? AND active = ?", id, true).
			Update("active", false).Error
		if err != nil {
			return err
		}

		// Create a new price for the updated item
		err =d.db.Table("costs").Create(&models.Cost{
			Cost:      cost,
			Active:    true,
			ProductId: id,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
