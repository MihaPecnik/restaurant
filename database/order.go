package database

import (
	"errors"
	"github.com/MihaPecnik/restaurant/models"
	"gorm.io/gorm"
)

func (d *Database) CreateOrder(order models.Order) error {
	order.Paid = false
	return d.db.Table("orders").Create(&order).Error
}

func (d *Database) AddItemToOrder(orderId, itemId int64) error {
	var cost models.Cost
	err := d.db.Table("costs").
		Where("product_id = ? AND active = ?", itemId, true).Scan(&cost).Error
	if err != nil {
		return err
	}
	return d.db.Table("order_mappings").
		Create(&models.OrderMapping{
			OrderId:   orderId,
			ProductId: itemId,
			PriceId:   cost.ID,
		}).Error
}

func (d *Database) PayTheOrder(orderId int64, payment float64) (models.Receipt, error) {
	// Total sum of the order
	querySum := `
select sum(cost) from order_mappings
inner join costs c on order_mappings.price_id = c.id and c.active=true
where order_mappings.order_id = ?
`
	// All items in order grouped by cost
	queryOrder := `
select p.name, count(p.id) as quantity, sum(c.cost) as sum_cost from order_mappings
inner join products p on p.id = order_mappings.product_id
inner join costs c on order_mappings.price_id = c.id
where order_mappings.order_id = ?
group by p.id, c.cost
`
	stmt := d.db.Session(&gorm.Session{
		PrepareStmt: true,
	})
	reciept := models.Receipt{}

	err := stmt.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		// Check if order was already paid for
		err := tx.Table("orders").Where("id = ?",orderId).Scan(&order).Error
		if err != nil {
			return err
		}
		if order.Paid {
			return errors.New("this order was already paid for")
		}
		// Check if the funds are sufficient
		var sumOrder float64
		err = tx.
			Raw(querySum, orderId).
			Scan(&sumOrder).Error
		if err != nil {
			return err
		}
		if payment < sumOrder {
			return errors.New("insufficient funds.")
		}

		// Get all items in order
		var items []models.Items
		err = tx.
			Raw(queryOrder, orderId).
			Scan(&items).Error
		if err != nil {
			return err
		}
		reciept = models.Receipt{
			Items:          items,
			TotalCost:      sumOrder,
			ChangeReturned: payment - sumOrder,
		}
		// Set order status to paid
		err = tx.Table("orders").Where("id = ?", orderId).Update("paid", true).Error
		if err != nil {
			return err
		}
		return nil
	})
	return reciept, err
}
