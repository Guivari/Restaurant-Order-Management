# management system
A "quick" implementation of a system that I was supposed to implement in 90 minutes but took 3 hours instead.


## Data Structures
### Shop
```
type Shop struct {
	Menu         map[string]MenuItem
	DiscountCode map[string]func(float64) float64
	Orders       map[string]*Order
	Tax          float64
	OrderCounter int64
}
```
>Our main struct. This is responsible for creating the menu (and their IDs), holding relevant constants, and order creation + order management. 
>
>DiscountCode maps a discount code to a function that returns a raw discount amount. It supports taking a float64, meant to be used as an order's total, and used to calculate absolute discount. You are in charge of writing and providing respective discount functions. Can return a negative discount if you want :)
>
>Orders mapping maps to a pointer to simplify modification in Shop's functions.
>
>Uses ordercounter to ensure that all orderIDs are unique (when considering scaling, number of order exceeds max int64. Alternative is to use UUID)

### MenuItem
```
type MenuItem struct {
	ID       string
	Name     string
	Category category.Category
	Price    float64
}
```
>Represents a menu item. Holds the items name, category and price.
>
>Category: "Appetizer", "Main", "Dessert", "Beverage"

### Order
```
type Order struct {
	ID           string 
	Customer     string
	Status       status.Status
	Items        map[string]int
	DiscountCode string
}
```
>Represents customer orders. Holds customer metadata, the items and their quantity in their order, the order status, and the discount code used.
>
>Items only holds the itemID and quantity, not an actual MenuItem object.
>
>Statuses: PENDING, CONFIRMED, COMPLETED, CANCELLED
>
>Discount code is stored even if user uses a wrong discount code.
>
>Currently only customer name is held as metadata.



## Shop Functions

### CreateShop
```
CreateShop(
    tax int64,
    discountCode map[string]func(float64) float64,
    menu []model.UseItem,
) *Shop
```
> **tax**: Will be converted to percentage when printing.
>
> **discountCode**: Have to write func(float64)float64 yourself.
>
> **menu**: User only provides name, category, and price. 

### ViewMenu
```
Shop.ViewMenu()
```
> Nothing much to explain. Prints whole menu.


### CreateOrder
```
Shop.CreateOrder(
    name string, 
    orders []model.OrderItem,
) (string, error)
```
> **name**: The one in your passport.
>
> **orders**: User has to provide itemID and quantity.
>
> **return**: OrderID string.
>
> **errors**:
> -  Item in orders not in menu.
> -  Order item quantity less than 0.

### ModifyOrder
```
Shop.ModifyOrder(
    orderID string, 
    orders []model.OrderItem,
) (error)
```
> **orderID**: Using the format ORD{id} (eg. ORD123, ORD1, ORD999)
>
> **orders**: User has to provide itemID and quantity. Incremental to old order.
>
> **errors**:
> -  Order not found in the system.
> -  Order has been completed/cancelled.
> -  Item in orders not in menu.
> -  Old item quantity + new order quantity = less than 0.


### ViewOrder
```
Shop.ViewOrder(
    orderID string, 
) (error)
```
> **orderID**: Using the format ORD{id} (eg. ORD123, ORD1, ORD999)
>
> **errors**:
> -  Order not found in the system.

### CompleteOrder
```
Shop.ViewOrder(
    orderID string, 
    discCode string,
) (error)
```
> **orderID**: Using the format ORD{id} (eg. ORD123, ORD1, ORD999)
>
> **discCode**: if no match, assumes no discount. DOES NOT CANCEL COMPLETION!
>
> **errors**:
> -  Order not found in the system.
> -  Order has been completed/cancelled.
> -  Order item has qtt < 0. If implementation is correct, this should never hit...

### CancelOrder
```
Shop.ViewOrder(
    orderID string, 
) (error)
```
> **orderID**: Using the format ORD{id} (eg. ORD123, ORD1, ORD999)
>
> **errors**:
> -  Order not found in the system.

## Enums
I made enums instead of making "magic values" for code resillience. Prevents one misspelling from causing half an hour of debugging.
### type Category string

```
const (
    Appetizer Category = "Appetizer"
    Main      Category = "Main"
    Dessert   Category = "Dessert"
    Beverage  Category = "Beverage"
)
```

### type Status string
```
const (
    Pending   Status = "PENDING"
    Confirm   Status = "CONFIRMED"
    Complete  Status = "COMPLETED"
    Cancelled Status = "CANCELLED"
)
```
