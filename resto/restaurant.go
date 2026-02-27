package resto

import (
	"fmt"
	"restaurant/model"
	"restaurant/model/category"
	"restaurant/model/status"
	"strconv"
)

type MenuItem struct {
	ID       string
	Name     string
	Category category.Category // Appetizer/Main/Dessert/Beverage
	Price    float64
}

type Order struct {
	ID           string // ORD-001
	Customer     string
	Status       status.Status // (PENDING, CONFIRMED, COMPLETED, CANCELLED)
	Items        map[string]int
	DiscountCode string
}

type Shop struct {
	Tax          float64
	DiscountCode map[string]func(float64) float64 //enum or map string -> integer
	Menu         map[string]MenuItem              //use order ID for map's key
	Orders       map[string]*Order
	OrderCounter int64
}

func CreateShop(
	tax float64,
	discountCode map[string]func(float64) float64,
	menu []model.UserItem,
) *Shop {
	newMenu := make(map[string]MenuItem)
	counter := 1
	for _, item := range menu {
		menuId := "M" + strconv.Itoa(counter)
		newMenu[menuId] = MenuItem{
			ID:       menuId,
			Name:     item.Name,
			Category: item.Category,
			Price:    item.Cost,
		}
		counter++
	}

	return &Shop{
		Tax:          tax,
		DiscountCode: discountCode,
		Menu:         newMenu,
		Orders:       make(map[string]*Order),
		OrderCounter: 0,
	}
}

func (this *Shop) ViewMenu() {
	fmtCategories := make(map[category.Category][]MenuItem, len(category.AllCat))
	for _, item := range this.Menu {
		fmtCategories[item.Category] = append(fmtCategories[item.Category], item)
	}
	fmt.Printf("MENU:\n")

	for _, cat := range category.AllCat {
		items := fmtCategories[cat]
		fmt.Printf("%ss:\n", cat)
		for _, item := range items {
			fmt.Printf("  [%s] %s - $%.2f\n", item.ID, item.Name, item.Price)
		}
	}
}

func (this *Shop) CreateOrder(name string, orders []model.OrderItem) (string, error) {
	newOrder := Order{
		ID:       "ORD" + strconv.Itoa(int(this.OrderCounter)),
		Customer: name,
		Status:   status.Pending,
		Items:    make(map[string]int),
	}
	for _, order := range orders {
		_, ok := this.Menu[order.ItemID]
		if !ok {
			return "", fmt.Errorf("[createOrder] %s not found in menu", order.ItemID)
		}
		if order.Quantity < 0 {
			return "", fmt.Errorf("[createOrder] invalid quantity %d for %s, cannot be less than 0", order.Quantity, order.ItemID)
		}
		newOrder.Items[order.ItemID] = order.Quantity
	}

	this.OrderCounter++
	this.Orders[newOrder.ID] = &newOrder
	return newOrder.ID, nil
}

func (this *Shop) ModifyOrder(orderID string, orders []model.OrderItem) error {
	orderOrig, ok := this.Orders[orderID]
	if !ok {
		return fmt.Errorf("[modifyOrder] order %s not found", orderID)
	}
	if orderOrig.Status == status.Cancelled {
		return fmt.Errorf("[modifyOrder] invalid orderID %s, order previously cancelled", orderID)
	}
	if orderOrig.Status == status.Complete {
		return fmt.Errorf("[modifyOrder] invalid orderID %s, order previously completed", orderID)
	}

	orderCopy := *orderOrig

	for _, change := range orders {
		_, ok := this.Menu[change.ItemID]
		if !ok {
			return fmt.Errorf("[modifyOrder] %s not found in menu", change.ItemID)
		}
		oldQtt, ok := orderCopy.Items[change.ItemID]
		if !ok {
			oldQtt = 0
			orderCopy.Items[change.ItemID] = oldQtt
		}
		newQtt := oldQtt + change.Quantity
		if newQtt < 0 {
			return fmt.Errorf("[modifyOrder] invalid change in %s, %s cannot have (%d)+(%d)=%d", orderID, change.ItemID,
				orderCopy.Items[change.ItemID], change.Quantity, newQtt)
		}
		orderCopy.Items[change.ItemID] = newQtt
	}
	this.Orders[orderID] = &orderCopy
	return nil
}

func (this *Shop) ViewOrder(orderID string) error {
	order, ok := this.Orders[orderID]
	if !ok {
		return fmt.Errorf("[viewOrder] order %s not found", orderID)
	}
	var totalCost float64 = 0
	fmt.Printf("Order Details:\n")
	fmt.Printf("Customer: %s\n", order.Customer)
	fmt.Printf("Status: %s\n", order.Status)
	fmt.Printf("Items:\n")
	for itemID, qtt := range order.Items {
		item := this.Menu[itemID]
		itemCost := (item.Price * float64(qtt))
		totalCost += itemCost
		fmt.Printf(" - %s x%d: $%.2f\n", item.Name, qtt, itemCost)
	}
	var orderTax = totalCost * this.Tax
	fmt.Printf("Subtotal: $%.2f\n", totalCost)
	fmt.Printf("Tax (%.0f%%): $%.2f\n", this.Tax*100, orderTax)
	fmt.Printf("Total: $%.2f\n", totalCost+orderTax)
	return nil
}

func (this *Shop) CompleteOrder(orderID string, discCode string) error {
	order, ok := this.Orders[orderID]
	if !ok {
		return fmt.Errorf("[completeOrder] invalid order %s not found", orderID)
	}
	if order.Status == status.Cancelled {
		return fmt.Errorf("[completeOrder] invalid orderID %s, order previously cancelled", orderID)
	}
	if order.Status == status.Complete {
		return fmt.Errorf("[completeOrder] invalid orderID %s, order previously completed", orderID)
	}

	order.Status = status.Confirm
	this.Orders[orderID] = order

	var totalCost float64 = 0
	for itemID, qtt := range order.Items {
		if qtt < 0 {
			return fmt.Errorf("[completeOrder] invalid order, %s has quantity of %d how did you get here??\n", itemID, qtt)
			// INFO: if all goes well this should never be triggered
		}
		item := this.Menu[itemID]
		itemCost := (item.Price * float64(qtt))
		totalCost += itemCost
	}

	var orderDiscount float64 = 0
	discountFunc, ok := this.DiscountCode[discCode]
	if ok {
		orderDiscount = discountFunc(totalCost)
		fmt.Printf("Discount applied: -$%.2f\n", orderDiscount)
	} else {
		fmt.Printf("Discount code not valid\n")
	}

	fmt.Printf("Final bill:\n")
	fmt.Printf("Subtotal: $%.2f\n", totalCost)
	fmt.Printf("Discount: -$%.2f\n", orderDiscount)

	totalCost -= orderDiscount
	var orderTax = totalCost * float64(this.Tax)
	fmt.Printf("Tax: $%.2f\n", orderTax)

	totalCost += orderTax
	fmt.Printf("Total: $%.2f\n", totalCost)

	fmt.Printf("Order completed. Thank you!\n")

	order.Status = status.Complete
	order.DiscountCode = discCode
	this.Orders[orderID] = order

	return nil
}

func (this *Shop) CancelOrder(orderID string) error {
	order, ok := this.Orders[orderID]
	if !ok {
		return fmt.Errorf("[cancelOrder] order %s not found", orderID)
	}
	order.Status = status.Cancelled
	this.Orders[orderID] = order
	return nil
}
