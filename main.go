package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal terhubung dengan DB:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	categoriesRepo := repositories.NewCategoryRepository(db)
	categoriesService := services.NewCategoryService(categoriesRepo)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)

	// Setup routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	http.HandleFunc("/categories", categoriesHandler.HandleCategories)
	http.HandleFunc("/categories/", categoriesHandler.HandleCategoriesByID)

	//Root BaseURL Handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Selamat Datang!")
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Berjalan dengan sempurna",
		})
	})

	// http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		getKategoriByID(w, r)
	// 	case http.MethodPut:
	// 		updateKategori(w, r)
	// 	case http.MethodDelete:
	// 		deleteKategori(w, r)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	// http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(kategori)
	// 	case http.MethodPost:
	// 		var kategoriBaru Kategori
	// 		err := json.NewDecoder(r.Body).Decode(&kategoriBaru)
	// 		if err != nil {
	// 			http.Error(w, "Invalid request", http.StatusBadRequest)
	// 			return
	// 		}
	// 		kategoriBaru.ID = len(kategori) + 1
	// 		kategori = append(kategori, kategoriBaru)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusCreated) // 201
	// 		json.NewEncoder(w).Encode(kategoriBaru)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	fmt.Println("Server Bisa Diakses Pada localhost:", +config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal Menjalankan API")
	}
}

// Bagian Kategori
type Kategori struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}

var kategori = []Kategori{
	{
		ID:        1,
		Nama:      "Makanan",
		Deskripsi: "Makanan Halal",
	},
	{
		ID:        2,
		Nama:      "Minuman",
		Deskripsi: "Minuman Halal",
	},
}

func getKategoriName(kategoriID int) string {
	for _, k := range kategori {
		if k.ID == kategoriID {
			return k.Nama
		}
	}
	return "Unknown"
}

func getKategoriByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	for _, k := range kategori {
		if k.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func updateKategori(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	var updateKategori Kategori
	err = json.NewDecoder(r.Body).Decode(&updateKategori)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for i := range kategori {
		if kategori[i].ID == id {
			updateKategori.ID = id
			kategori[i] = updateKategori
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateKategori)
			return
		}
	}
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func deleteKategori(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	for i, k := range kategori {
		if k.ID == id {
			kategori = append(kategori[:i], kategori[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}
