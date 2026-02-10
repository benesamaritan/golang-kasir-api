package main

import (
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
	"github.com/gorilla/mux"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization, X-Requested-With, Accept",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
		log.Fatal("Gagal Terhubung Dengan DB:", err)
	}
	defer db.Close()

	r := mux.newRouter()

	// BaseURI dan Health
	http.HandleFunc("/", handlers.WelcomeHandler())
	http.HandleFunc("/health", handlers.HealthCheckHandler())

	// Products Route
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Categories Route
	categoriesRepo := repositories.NewCategoryRepository(db)
	categoriesService := services.NewCategoryService(categoriesRepo)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)
	http.HandleFunc("/categories", categoriesHandler.HandleCategories)
	http.HandleFunc("/categories/", categoriesHandler.HandleCategoryByID)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	r.HandleFunc("/api/checkout", transactionHandler.HandleCheckout).Methods("POST")

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)
	r.HandleFunc("/api/report", reportHandler.GetReport).Methods("GET")
	r.HandleFunc("/api/report/today", reportHandler.GetReportToday).Methods("GET")

	// Pesan Status Aplikasi
	fmt.Println("Aplikasi Kasir Berhasil Jalan Pada Port:", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal Menjalankan API")
	}

	handler := enableCORS(r)
	// fmt.Println("Server running di localhost:" + port)
	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
