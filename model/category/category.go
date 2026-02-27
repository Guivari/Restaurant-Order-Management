package category

type Category string

const (
	Appetizer Category = "Appetizer"
	Main      Category = "Main"
	Dessert   Category = "Dessert"
	Beverage  Category = "Beverage"
)

var AllCat = []Category{
	Appetizer,
	Main,
	Dessert,
	Beverage,
}
