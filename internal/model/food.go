package model

type FoodCategory int

// PERF: i guess the dish types could be mixed, but i would just
// skip this feature and provide only one type to each dish.
const (
	FoodCategoryStaple FoodCategory = iota
	FoodCategoryVegetable
	FoodCategoryMeat
	FoodCategorySeafood
	FoodCategorySoup
	FoodCategoryDessert
	FoodCategoryDrink
	FoodCategoryOther
)

type Food struct {
	ID           int
	FoodID       string
	Name         string
	Detail       string
	Prize        float32
	Rate         float32
	RequiredTime int64
	SoldCount    int
	ImageURIs    []string
	Ingredients  []string
	Category     FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
	CreatedAt    int64
	DeletedAt    int64
}

type FoodStep struct {
	ID     int
	FoodID string
	Sort   int
	Detail string
}

// DTO.
type FoodDetail struct {
	Food     Food
	Steps    []FoodStep
	Comments []FoodComment
}

// Comment.
type FoodComment struct {
	ID        int
	CommentID string
	FoodID    string
	Detail    string
	CreatedAt int64
	DeletedAt int64
}
