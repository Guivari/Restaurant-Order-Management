package main

import (
	"log"
	"restaurant/model"
	"restaurant/model/category"
	"restaurant/resto"
)

var (
	TAX       float64 = 0.1 // INFO: in percentage
	DISCOUNTS       = map[string]func(float64) float64{
		"WELCOME10": func(cost float64) float64 {
			return cost * 0.1
		},
		"FEAST20": func(cost float64) float64 {
			if cost > 50 {
				return cost * 0.2
			}
			return 0
		},
		"SAVE5": func(cost float64) float64 {
			return 5
		},
	}
	MENU = []model.UserItem{
		{Name: "Spring Rolls", Category: category.Appetizer, Cost: 6},
		{Name: "Soup", Category: category.Appetizer, Cost: 5},

		{Name: "Pasta", Category: category.Main, Cost: 12},
		{Name: "Burger", Category: category.Main, Cost: 10},
		{Name: "Pizza", Category: category.Main, Cost: 15},

		{Name: "Ice Cream", Category: category.Dessert, Cost: 4},
		{Name: "Cake", Category: category.Dessert, Cost: 6},

		{Name: "Soda", Category: category.Beverage, Cost: 2},
		{Name: "Juice", Category: category.Beverage, Cost: 3},
	}
)

func main() {
  var err error
	
  shop := resto.CreateShop(
		TAX,
		DISCOUNTS,
		MENU,
	)

	shop.ViewMenu()

  order1 := shop.CreateOrder(
    "James Brown",
    []model.OrderItem{
      {Item: "M8",Quantity: 4,},
      {Item: "M1",Quantity: 2,},
      {Item: "M4",Quantity: 3,},
      {Item: "M3",Quantity: 1,},
      {Item: "M7",Quantity: 2,},
    },
  )

  
  order2 := shop.CreateOrder(
    "John Smith",
    []model.OrderItem{
      {Item: "M3",Quantity: 2,},
      {Item: "M8",Quantity: 2,},
      {Item: "M6",Quantity: 1,},
    },
  )

  err = shop.ViewOrder(order1)
  if err != nil {
    log.Fatal(err)
  }

  err = shop.ViewOrder(order2)
  if err != nil {
    log.Fatal(err)
  }
  
  err = shop.CompleteOrder(order2,"WELCOME10")
  if err != nil {
    log.Fatal(err)
  }
}
