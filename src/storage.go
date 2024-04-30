package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage() (*Storage, error) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("SSLMODE")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s TimeZone=Europe/Warsaw",
		host,
		user,
		pass,
		db_name,
		sslmode)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	if db.AutoMigrate(&Recipe{}, &Ingredient{}) != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) Configure(db *gorm.DB) {
	s.db = db
}

func (s *Storage) Create(r *Recipe) error {
	if err := s.db.Create(r).Error; err != nil {
		return err
	}
	return nil
}
func (s *Storage) CreateBulk(recipes []*Recipe) error {

	for _, r := range recipes {
		if err := s.db.Create(r).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Get() ([]*Recipe, error) {
	recipes := []*Recipe{}

	res := s.db.Find(&recipes)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, recipe := range recipes {
		if err := s.db.Model(recipe).Association("Ingredients").Find(&recipe.Ingredients); err != nil {
			return recipes, err
		}
	}
	fmt.Printf("%+v\n", recipes)

	return recipes, nil
}

func (s *Storage) GetByID(id int) (*Recipe, error) {

	recipe := &Recipe{}
	err := s.db.First(recipe, id).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(recipe).Association("Ingredients").Find(&recipe.Ingredients); err != nil {
		return recipe, err
	}

	return recipe, nil
}
func (s *Storage) Update(r *Recipe) error {

	if err := s.db.Model(r).Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(id int) error {
	var recipeToDelete Recipe
	if err := s.db.Where("id = ?", id).First(&recipeToDelete).Error; err != nil {
		return err
	}

	s.db.Delete(&recipeToDelete)
	return nil
}

func (s *Storage) SeedRecipes() error {
	recipes := []Recipe{
		{
			Name: "Chocolate Cake",
			Ingredients: []*Ingredient{
				{Name: "Flour", Amount: "2 cups"},
				{Name: "Sugar", Amount: "1.5 cups"},
				{Name: "Cocoa powder", Amount: "0.75 cups"},
				{Name: "Baking powder", Amount: "1.5 tsp"},
				{Name: "Baking soda", Amount: "1.5 tsp"},
				{Name: "Salt", Amount: "1 tsp"},
				{Name: "Eggs", Amount: "2"},
				{Name: "Milk", Amount: "1 cup"},
				{Name: "Vegetable oil", Amount: "0.5 cup"},
				{Name: "Vanilla extract", Amount: "2 tsp"},
				{Name: "Boiling water", Amount: "1 cup"},
			},
		},
		{
			Name: "Chicken Alfredo Pasta",
			Ingredients: []*Ingredient{
				{Name: "Fettuccine pasta", Amount: "8 oz"},
				{Name: "Chicken breast", Amount: "2"},
				{Name: "Butter", Amount: "0.5 cup"},
				{Name: "Heavy cream", Amount: "1 cup"},
				{Name: "Parmesan cheese", Amount: "1 cup"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Salt", Amount: "1 tsp"},
				{Name: "Black pepper", Amount: "0.5 tsp"},
				{Name: "Parsley", Amount: "2 tbsp"},
			},
		},
		{
			Name: "Spaghetti Bolognese",
			Ingredients: []*Ingredient{
				{Name: "Spaghetti", Amount: "8 oz"},
				{Name: "Ground beef", Amount: "1 lb"},
				{Name: "Onion", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Tomato sauce", Amount: "2 cups"},
				{Name: "Dried oregano", Amount: "1 tsp"},
				{Name: "Dried basil", Amount: "1 tsp"},
				{Name: "Salt", Amount: "1 tsp"},
				{Name: "Black pepper", Amount: "0.5 tsp"},
				{Name: "Parmesan cheese", Amount: "0.5 cup"},
			},
		},
		{
			Name: "Classic Margherita Pizza",
			Ingredients: []*Ingredient{
				{Name: "Pizza dough", Amount: "1 lb"},
				{Name: "Tomato sauce", Amount: "1 cup"},
				{Name: "Fresh mozzarella cheese", Amount: "8 oz"},
				{Name: "Fresh basil leaves", Amount: "1/2 cup"},
				{Name: "Extra-virgin olive oil", Amount: "2 tbsp"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
			},
		},
		{
			Name: "Caesar Salad",
			Ingredients: []*Ingredient{
				{Name: "Romaine lettuce", Amount: "1 head"},
				{Name: "Croutons", Amount: "1 cup"},
				{Name: "Parmesan cheese", Amount: "1/2 cup"},
				{Name: "Chicken breast", Amount: "2"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Lemon juice", Amount: "2 tbsp"},
				{Name: "Garlic clove", Amount: "1"},
				{Name: "Worcestershire sauce", Amount: "1 tsp"},
				{Name: "Dijon mustard", Amount: "1 tsp"},
				{Name: "Anchovy fillets", Amount: "4"},
				{Name: "Salt", Amount: "1/4 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
			},
		},
		{
			Name: "Classic Beef Burger",
			Ingredients: []*Ingredient{
				{Name: "Ground beef", Amount: "1 lb"},
				{Name: "Hamburger buns", Amount: "4"},
				{Name: "Cheddar cheese", Amount: "4 slices"},
				{Name: "Lettuce leaves", Amount: "4"},
				{Name: "Tomato slices", Amount: "4"},
				{Name: "Red onion slices", Amount: "4"},
				{Name: "Ketchup", Amount: "1/4 cup"},
				{Name: "Mustard", Amount: "1/4 cup"},
				{Name: "Salt", Amount: "1 tsp"},
				{Name: "Black pepper", Amount: "1/2 tsp"},
			},
		},
		{
			Name: "Chicken Alfredo Pasta",
			Ingredients: []*Ingredient{
				{Name: "Fettuccine pasta", Amount: "1 lb"},
				{Name: "Chicken breast", Amount: "2"},
				{Name: "Heavy cream", Amount: "1 cup"},
				{Name: "Parmesan cheese", Amount: "1 cup"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Butter", Amount: "2 tbsp"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
				{Name: "Fresh parsley", Amount: "2 tbsp"},
			},
		},
		{
			Name: "Vegetable Stir Fry",
			Ingredients: []*Ingredient{
				{Name: "Broccoli florets", Amount: "2 cups"},
				{Name: "Bell peppers", Amount: "1 cup"},
				{Name: "Carrots", Amount: "1 cup"},
				{Name: "Snap peas", Amount: "1 cup"},
				{Name: "Onion", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Ginger", Amount: "1 tbsp"},
				{Name: "Soy sauce", Amount: "1/4 cup"},
				{Name: "Sesame oil", Amount: "1 tbsp"},
				{Name: "Honey", Amount: "1 tbsp"},
				{Name: "Cornstarch", Amount: "1 tsp"},
				{Name: "Vegetable oil", Amount: "2 tbsp"},
				{Name: "Salt", Amount: "1/4 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
			},
		},
		{
			Name: "Spaghetti Bolognese",
			Ingredients: []*Ingredient{
				{Name: "Spaghetti pasta", Amount: "1 lb"},
				{Name: "Ground beef", Amount: "1 lb"},
				{Name: "Tomato sauce", Amount: "2 cups"},
				{Name: "Onion", Amount: "1"},
				{Name: "Carrot", Amount: "1"},
				{Name: "Celery stalk", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Red wine", Amount: "1/2 cup"},
				{Name: "Beef broth", Amount: "1/2 cup"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
				{Name: "Fresh parsley", Amount: "2 tbsp"},
				{Name: "Grated Parmesan cheese", Amount: "1/4 cup"},
			},
		},
		{
			Name: "Vegetarian Chili",
			Ingredients: []*Ingredient{
				{Name: "Kidney beans", Amount: "2 cans"},
				{Name: "Black beans", Amount: "1 can"},
				{Name: "Diced tomatoes", Amount: "1 can"},
				{Name: "Bell peppers", Amount: "2"},
				{Name: "Onion", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Vegetable broth", Amount: "2 cups"},
				{Name: "Chili powder", Amount: "2 tbsp"},
				{Name: "Cumin", Amount: "1 tbsp"},
				{Name: "Paprika", Amount: "1 tbsp"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
				{Name: "Fresh cilantro", Amount: "2 tbsp"},
				{Name: "Sour cream", Amount: "1/4 cup"},
			},
		},
		{
			Name: "Shrimp Scampi",
			Ingredients: []*Ingredient{
				{Name: "Shrimp", Amount: "1 lb"},
				{Name: "Linguine pasta", Amount: "1/2 lb"},
				{Name: "Butter", Amount: "4 tbsp"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Garlic cloves", Amount: "4"},
				{Name: "White wine", Amount: "1/2 cup"},
				{Name: "Lemon juice", Amount: "2 tbsp"},
				{Name: "Red pepper flakes", Amount: "1/4 tsp"},
				{Name: "Salt", Amount: "1/4 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
				{Name: "Fresh parsley", Amount: "2 tbsp"},
				{Name: "Grated Parmesan cheese", Amount: "1/4 cup"},
			},
		},
		{
			Name: "Beef Tacos",
			Ingredients: []*Ingredient{
				{Name: "Ground beef", Amount: "1 lb"},
				{Name: "Taco shells", Amount: "12"},
				{Name: "Lettuce", Amount: "1 cup"},
				{Name: "Tomatoes", Amount: "2"},
				{Name: "Cheddar cheese", Amount: "1 cup"},
				{Name: "Onion", Amount: "1"},
				{Name: "Taco seasoning", Amount: "1 packet"},
				{Name: "Sour cream", Amount: "1/2 cup"},
				{Name: "Salsa", Amount: "1/2 cup"},
				{Name: "Guacamole", Amount: "1/2 cup"},
			},
		},
		{
			Name: "Chicken Tikka Masala",
			Ingredients: []*Ingredient{
				{Name: "Chicken thighs", Amount: "1 lb"},
				{Name: "Tomato sauce", Amount: "1 cup"},
				{Name: "Greek yogurt", Amount: "1/2 cup"},
				{Name: "Onion", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Ginger", Amount: "1 tbsp"},
				{Name: "Garam masala", Amount: "1 tbsp"},
				{Name: "Paprika", Amount: "1 tsp"},
				{Name: "Cumin", Amount: "1 tsp"},
				{Name: "Coriander", Amount: "1 tsp"},
				{Name: "Cayenne pepper", Amount: "1/4 tsp"},
				{Name: "Butter", Amount: "2 tbsp"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
				{Name: "Fresh cilantro", Amount: "2 tbsp"},
			},
		},
		{
			Name: "Vegetable Lasagna",
			Ingredients: []*Ingredient{
				{Name: "Lasagna noodles", Amount: "12"},
				{Name: "Zucchini", Amount: "2"},
				{Name: "Yellow squash", Amount: "2"},
				{Name: "Carrots", Amount: "2"},
				{Name: "Onion", Amount: "1"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Spinach", Amount: "2 cups"},
				{Name: "Ricotta cheese", Amount: "1 1/2 cups"},
				{Name: "Mozzarella cheese", Amount: "2 cups"},
				{Name: "Parmesan cheese", Amount: "1/2 cup"},
				{Name: "Tomato sauce", Amount: "3 cups"},
				{Name: "Olive oil", Amount: "2 tbsp"},
				{Name: "Italian seasoning", Amount: "1 tbsp"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
			},
		},
		{
			Name: "Chicken Caesar Wrap",
			Ingredients: []*Ingredient{
				{Name: "Grilled chicken breast", Amount: "2"},
				{Name: "Flour tortillas", Amount: "2"},
				{Name: "Romaine lettuce", Amount: "1 cup"},
				{Name: "Caesar dressing", Amount: "1/4 cup"},
				{Name: "Parmesan cheese", Amount: "1/4 cup"},
				{Name: "Tomato", Amount: "1"},
			},
		},
		{
			Name: "Vegetarian Pizza",
			Ingredients: []*Ingredient{
				{Name: "Pizza dough", Amount: "1 lb"},
				{Name: "Tomato sauce", Amount: "1 cup"},
				{Name: "Mozzarella cheese", Amount: "8 oz"},
				{Name: "Bell peppers", Amount: "1"},
				{Name: "Red onion", Amount: "1/2"},
				{Name: "Black olives", Amount: "1/4 cup"},
				{Name: "Fresh basil leaves", Amount: "1/4 cup"},
				{Name: "Extra-virgin olive oil", Amount: "2 tbsp"},
				{Name: "Garlic cloves", Amount: "2"},
				{Name: "Salt", Amount: "1/2 tsp"},
				{Name: "Black pepper", Amount: "1/4 tsp"},
			},
		},
	}

	for i := range recipes {
		// Check if the recipe already exists
		var existingRecipe Recipe
		if err := s.db.Where("name = ?", recipes[i].Name).First(&existingRecipe).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Recipe doesn't exist, so create it
				if err := s.db.Create(&recipes[i]).Error; err != nil {
					return fmt.Errorf("failed to create recipe: %w", err)
				}
				log.Printf("Created recipe: %s", recipes[i].Name)
			} else {
				return fmt.Errorf("failed to check if recipe exists: %w", err)
			}
		} else {
			log.Printf("Recipe already exists: %s", recipes[i].Name)
		}
	}

	log.Println("Recipes seeded successfully")
	return nil
}
