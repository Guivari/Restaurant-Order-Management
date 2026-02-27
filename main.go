package main

import (
	"fmt"
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

  order1, err := shop.CreateOrder(
    "James Brown",
    []model.OrderItem{
      {ItemID: "M8",Quantity: 4,},
      {ItemID: "M1",Quantity: 2,},
      {ItemID: "M4",Quantity: 3,},
      {ItemID: "M3",Quantity: 1,},
      {ItemID: "M7",Quantity: 2,},
    },
  )  
  order2, err := shop.CreateOrder(
    "Jeremiah Singh",
    []model.OrderItem{
      {ItemID: "M8",Quantity: 4,},
      {ItemID: "M1",Quantity: 2,},
      {ItemID: "M4",Quantity: 3,},
      {ItemID: "M3",Quantity: 1,},
      {ItemID: "M7",Quantity: 2,},
    },
  )
  order3, err := shop.CreateOrder(
    "Jane Doe",
    []model.OrderItem{
      {ItemID: "M1",Quantity: 4,},
      {ItemID: "M2",Quantity: 5,},
      {ItemID: "M3",Quantity: 6,},
    },
  )
  orderZ, err := shop.CreateOrder("Lorem Ipsum",[]model.OrderItem{
    {ItemID: "M1", Quantity: 1},
  })
  if err != nil {log.Println(err)}

  //No modification, cancelled
  err = shop.ViewOrder(order1)
  if err != nil {log.Println(err)}
  err = shop.ModifyOrder(order1, []model.OrderItem{
    {ItemID: "M4",Quantity: -1,},
    {ItemID: "M3",Quantity: 2,},
  })
  if err != nil {log.Println(err)}
  err = shop.CancelOrder(order1)
  if err != nil {log.Println(err)}

  //Good modification, good discount
  err = shop.ViewOrder(order2)
  if err != nil {log.Println(err)}
  err = shop.ModifyOrder(order2, []model.OrderItem{
    {ItemID: "M4",Quantity: -1,},
    {ItemID: "M3",Quantity: 2,},
    {ItemID: "M6",Quantity: 5,},
  })
  if err != nil {log.Println(err)}
  err = shop.ViewOrder(order2)
  if err != nil {log.Println(err)}
  err = shop.CompleteOrder(order2,"WELCOME10")
  if err != nil {log.Println(err)}
  
  fmt.Println()

  //bad orderID
  err = shop.ModifyOrder("badID",[]model.OrderItem{})
  if err != nil {log.Println(err)}
  err = shop.ViewOrder("badID")
  if err != nil {log.Println(err)}
  err = shop.CompleteOrder("badID","")
  if err != nil {log.Println(err)}
  err = shop.CancelOrder("badID")
  if err != nil {log.Println(err)}

  fmt.Println()

  //bad create, wrong itemid
  _, err = shop.CreateOrder("Lorem Ipsum",[]model.OrderItem{
    {ItemID: "M6969", Quantity: 0},
  })
  if err != nil {log.Println(err)}
  //bad modification, wrong itemid
  err = shop.ModifyOrder(orderZ, []model.OrderItem{
    {ItemID: "M420", Quantity: 1},
  })
  if err != nil {log.Println(err)}

  
  fmt.Println()

  //bad modification, invalid quantity
  err = shop.ModifyOrder(order3, []model.OrderItem{
    {ItemID: "M4",Quantity: -10,},
  })
  if err != nil {log.Println(err)}
  
  fmt.Println()
  
  //order expired (either completed or cancelled)
  err = shop.ModifyOrder(order1, []model.OrderItem{})
  if err != nil {log.Println(err)}
  err = shop.ModifyOrder(order2, []model.OrderItem{})
  if err != nil {log.Println(err)}
  err = shop.CompleteOrder(order1,"BADCODE")
  if err != nil {log.Println(err)}
  err = shop.CompleteOrder(order2,"BADCODE")
  if err != nil {log.Println(err)}

  fmt.Println()
  
  //bad discount
  err = shop.CompleteOrder(order3,"BADCODE")
  if err != nil {log.Println(err)}




}
