package main

import (
	"edra/controllers"
	"edra/models"
	"log"

	"github.com/gin-gonic/gin"
)

// func Init() {
// 	controllers.Initialize()
// }

func main() {
	// Create a server with a port
	r := gin.Default()

	log.Println("Starting goroutines")
	controllers.KeyStore = make(map[int64]models.Key, 1)

	go controllers.BlacklistKeys()
	go controllers.UnblockKeys()
	// fmt.Println("Key generated: ", key)

	public := r.Group("/api")
	// Log initial KeyStore state
	log.Println("Initial KeyStore state:")
	for i, val := range controllers.KeyStore {
		log.Printf("Key %d: %+v\n", i, val)
	}

	// POST api /keys to retrieve an available key for client use
	public.POST("/keys", controllers.GenerateKeys)

	// GET /keys: Retrieve an available key for client use.
	public.GET("/keys", controllers.RetrieveKey)

	// GET /keys/:id: Provide information (e.g., assignment timestamps) about a specific key.
	public.GET("/keys/:id", controllers.RetrieveKeyByID)

	// DELETE /keys/:id: Remove a specific key, identified by :id, from the system.
	public.DELETE("keys/:id", controllers.DeleteKeyByID)

	// PUT /keys/:id: Unblock a key for further use.
	public.PUT("keys/:id", controllers.UnblockKeyByID)

	// PUT /keepalive/:id: Signal the server to keep the specified key, identified by :id, from being deleted.
	public.PUT("/keepalive/:id", controllers.KeepKeyAliveByID)

	// public.POST("/unblock", controllers.UnblockKeysCRON)

	// public.POST("/blacklist", controllers.DeleteKeysCRON)

	r.Run("127.0.0.1:1918")

}
