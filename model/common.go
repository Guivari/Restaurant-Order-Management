package model

import "restaurant/model/category"

type UserItem struct {
	Name     string
	Category category.Category
	Cost     float64
}

type OrderItem struct {
	ItemID   string
	Quantity int
}
