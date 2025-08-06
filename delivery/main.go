package main

import (
	"a2sv-project/Delivery/routers"
	"a2sv-project/Infrastructure/database"
)

func main() {
	// Connect to the database
	database.ConnectDB()
	
	r := routers.SetupRoutes()
	r.Run(":8080")
}
