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
		log.Fatal("Gagal menghubungi database:", err)
	}
	defer db.Close()

	//Products
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	//Categories
	categoriesRepo := repositories.NewCategoryRepository(db)
	categoriesService := services.NewCategoryService(categoriesRepo)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)

	//Transactions
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes

	// BaseURI dan Health
	http.HandleFunc("/", handlers.WelcomeHandler())
	http.HandleFunc("/health", handlers.HealthCheckHandler())

	// Products Route
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Transactions Route
	http.HandleFunc("/api/checkout", transactionHandler.Checkout)

	// Categories Route
	http.HandleFunc("/categories", categoriesHandler.HandleCategories)
	http.HandleFunc("/categories/", categoriesHandler.HandleCategoryByID)

	// Pesan Status Aplikasi
	fmt.Println("Aplikasi Kasir Berhasil Jalan Pada Port:", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal Menjalankan API")
	}
}
