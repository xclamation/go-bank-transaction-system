package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/xclamation/go-bank-transaction-system/internal/database"
	"github.com/xclamation/go-bank-transaction-system/internal/server"
	"github.com/xclamation/go-bank-transaction-system/internal/worker"
)

// type apiConfig struct {
// 	DB *database.Queries
// }

func main() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// Задержка для ожидания запуска базы данных
	time.Sleep(20 * time.Second)

	// Подключение к базе данных
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	// Проверка подключения к базе данных
	if err = conn.Ping(); err != nil {
		log.Fatal("Can't ping database: ", err)
	}

	// Инициализация Queries

	db := database.New(conn)

	go worker.StartWorker(db)

	// apiCfg := apiConfig{
	// 	DB: db,
	// }

	// Настройка роутера
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/send-transaction", server.HandlerTransaction)

	log.Printf("Server starting on port %v", portString)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
