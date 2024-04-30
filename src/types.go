package main

type Recipe struct {
	ID          int           `gorm:"primaryKey" json:"id"`
	Name        string        `json:"recipe_name"`
	Ingredients []*Ingredient `gorm:"foreignKey:RecipeID" json:"ingredients"`
}

type Ingredient struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	RecipeID int    `json:"recipe_id"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
}
