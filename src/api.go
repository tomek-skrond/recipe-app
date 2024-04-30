package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type APIServer struct {
	DB         *Storage
	listenPort string
}

func NewAPIServer(db *Storage, listenPort string) (*APIServer, error) {
	return &APIServer{DB: db, listenPort: listenPort}, nil
}

func (s *APIServer) Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", makeFrontendHandler(s.handleHome))
	router.HandleFunc("GET /recipe/{id}", makeFrontendHandler(s.handleViewRecipeByID))

	router.HandleFunc("GET /v1/recipes", makeHandler(s.handleGetRecipe))
	router.HandleFunc("GET /v1/recipes/{id}", makeHandler(s.handleGetRecipeByID))
	router.HandleFunc("POST /v1/recipes", makeHandler(s.handleCreateRecipe))
	router.HandleFunc("POST /v1/recipes/bulk", makeHandler(s.handleCreateRecipeBulk))
	router.HandleFunc("PUT /v1/recipes/{id}", makeHandler(s.handleUpdateRecipeByID))
	router.HandleFunc("DELETE /v1/recipes/{id}", makeHandler(s.handleDeleteRecipeByID))

	log.Println("Server listening on address", s.listenPort)
	log.Println(http.ListenAndServe(s.listenPort, router))
}

// Frontend handlers
func (s *APIServer) handleHome(w http.ResponseWriter, r *http.Request) error {
	log.Printf("%s %s", r.RemoteAddr, r.URL.Path)
	workdir, err := os.Getwd()
	if err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}
	tmpl, err := template.New("index.html").ParseFiles(workdir + "/template/index.html")
	if err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}

	if err := tmpl.Execute(w, nil); err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}
	return nil
}

func (s *APIServer) handleViewRecipeByID(w http.ResponseWriter, r *http.Request) error {
	log.Printf("%s %s", r.RemoteAddr, r.URL.Path)
	workdir, err := os.Getwd()
	if err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}
	tmpl, err := template.New("recipe.html").ParseFiles(workdir + "/template/recipe.html")
	if err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}

	id := getId(w, r)
	idInt, _ := strconv.Atoi(id)

	recipeData, err := s.DB.GetByID(idInt)
	if err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}

	if err := tmpl.Execute(w, recipeData); err != nil {
		return apiError{message: err.Error(), statuscode: http.StatusInternalServerError}
	}
	return nil
}

// Backend handlers
func (s *APIServer) handleGetRecipe(w http.ResponseWriter, r *http.Request) error {
	var recipes []*Recipe
	recipes, err := s.DB.Get()
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}

	return writeJSON(w, http.StatusOK, recipes)
}

func getId(w http.ResponseWriter, r *http.Request) string {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	return id
}

func (s *APIServer) handleGetRecipeByID(w http.ResponseWriter, r *http.Request) error {
	id := getId(w, r)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}

	recipe, err := s.DB.GetByID(idInt)
	return writeJSON(w, http.StatusOK, recipe)
}
func (s *APIServer) handleCreateRecipe(w http.ResponseWriter, r *http.Request) error {
	var newRecipe *Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipe); err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}
	if err := s.DB.Create(newRecipe); err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}
	return writeJSON(w, http.StatusOK, "successful")
}

func (s *APIServer) handleCreateRecipeBulk(w http.ResponseWriter, r *http.Request) error {
	var newRecipes []*Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipes); err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}
	if err := s.DB.CreateBulk(newRecipes); err != nil {
		return writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return writeJSON(w, http.StatusOK, "successful")
}

func (s *APIServer) handleUpdateRecipeByID(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (s *APIServer) handleDeleteRecipeByID(w http.ResponseWriter, r *http.Request) error {

	return nil
}
