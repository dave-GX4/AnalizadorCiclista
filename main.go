package main

import (
	"compilerciclista/src/database"
	"compilerciclista/src/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// 1. Cargar las variables de entorno desde el archivo .env
	// Esto debe ser lo primero que se ejecute.
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo encontrar el archivo .env, se usarán las variables de entorno del sistema.")
	}

	// 2. Construir la cadena de conexión (DSN) para MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 3. Inicializar la conexión a la base de datos
	if err := database.InitDB(dsn); err != nil {
		log.Fatalf("Error fatal: No se pudo conectar a la base de datos: %v", err)
	}
	// Nos aseguramos de cerrar la conexión cuando la aplicación termine
	defer database.DB.Close()
	log.Println("Conexión a la base de datos establecida exitosamente.")

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterParticipantHandler)

	// 5. Configurar el middleware de CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // <-- La URL de tu frontend de Vite
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	})
	handler := c.Handler(mux) // Envuelve tu mux con el manejador de CORS

	// 6. Iniciar el servidor HTTP usando el manejador con CORS
	port := ":8080"
	log.Printf("Servidor API escuchando en http://localhost%s", port)
	if err := http.ListenAndServe(port, handler); err != nil { // <-- USA 'handler' en lugar de 'mux' o 'nil'
		log.Fatalf("Error fatal: No se pudo iniciar el servidor: %v", err)
	}
}