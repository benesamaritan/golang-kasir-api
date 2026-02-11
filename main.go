package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
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
	APIKey string `mapstructure:"API_KEY"`
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
		APIKey: viper.GetString("API_KEY"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	apiKeyMiddleware := middlewares.APIKey(config.APIKey)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Dynamic Route Tracking
	var registeredRoutes []handlers.RouteInfo
	register := func(path string, handler http.HandlerFunc) {
		registeredRoutes = append(registeredRoutes, handlers.RouteInfo{Path: path})
		http.HandleFunc(path, handler)
	}

	// API Discovery
	apiHandler := handlers.NewAPIHandler(&registeredRoutes)
	register("/", apiHandler.Welcome)
	register("/api", middlewares.CORS(middlewares.Logger(apiHandler.ListRoutes)))
	http.Handle("/api/", http.RedirectHandler("/api", http.StatusMovedPermanently))

	// Setup routes
	register("/api/product", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	register("/api/product/{id}", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))
	http.Handle("/api/produk", http.RedirectHandler("/api/product", http.StatusMovedPermanently))
	http.Handle("/api/produk/", http.RedirectHandler("/api/product/", http.StatusMovedPermanently))

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	register("/api/checkout", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.Checkout))))

	// Categories Route
	categoriesRepo := repositories.NewCategoryRepository(db)
	categoriesService := services.NewCategoryService(categoriesRepo)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)
	register("/categories", categoriesHandler.HandleCategories)
	register("/categories/{id}", categoriesHandler.HandleCategoryByID)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)
	register("/api/report", reportHandler.GetReport)
	register("/api/report/today", reportHandler.GetReportToday)
	http.Handle("/api/report/hari-ini", http.RedirectHandler("/api/report/today", http.StatusMovedPermanently))
	http.Handle("/api/report/range", http.RedirectHandler("/api/report", http.StatusMovedPermanently))

	// localhost:8080/health
	register("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("Server berjalan di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal menjalankan server:", err)
	}
}
