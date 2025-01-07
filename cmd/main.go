package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"auther/configs"
	database "auther/internal/db"
	"auther/internal/server/handlers"
	"auther/internal/server/middlewares"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to config file")
	flag.Parse()

	if err := run(*configPath); err != nil {
		log.Fatalf("Failed to run the application: %v", err)
	}
}

func run(configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	db, err := initDatabase(cfg.DB)
	if err != nil {
		return err
	}

	mux := initRouter(cfg, db)

	log.Printf("Server is running on port %d", cfg.Server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mux)
}

func loadConfig(configPath string) (*configs.Config, error) {
	configsManager := configs.ConfigManager{}
	cfg, err := configsManager.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	log.Println("Config loaded successfully")
	return cfg, nil
}

func initDatabase(DBConfig *configs.DBConfig) (*gorm.DB, error) {
	db, err := database.ConnectDB(DBConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	log.Println("Database connected successfully")
	return db, nil
}

func initRouter(cfg *configs.Config, db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginHandler(cfg.JWT, db)).Methods("POST")
	router.HandleFunc("/refresh", handlers.RefreshTokenHandler(cfg.JWT, db)).Methods("POST")
	router.Handle("/users", middlewares.AdminTokenMiddleware(cfg.Admin)(http.HandlerFunc(handlers.CreateUserHandler(db)))).Methods("POST")

	router.Handle("/users", middlewares.AdminTokenMiddleware(cfg.Admin)(http.HandlerFunc(handlers.DeleteUserHandler(db)))).Methods("DELETE")
	router.Handle("/users/id/{id}", middlewares.AdminTokenMiddleware(cfg.Admin)(http.HandlerFunc(handlers.DeleteUserByIDHandler(db)))).Methods("DELETE")
	router.Handle("/users/login", middlewares.AdminTokenMiddleware(cfg.Admin)(http.HandlerFunc(handlers.DeleteUserByLoginHandler(db)))).Methods("DELETE")

	log.Println("Router initialized successfully")

	return router
}
