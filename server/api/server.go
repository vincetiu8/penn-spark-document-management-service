package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/vincetiu8/penn-spark-server/api/controllers"
)

var server = controllers.Server{}

// Run sets up and runs the server using the specified environment variables.
// Database path set with DB_PATH environment variable.
// URL extension path set with API_PATH environment variable.
// Path to stored documents relative to server folder set with FS_PATH environment variable.
func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("error getting env, not comming through %v", err)
	} else {
		fmt.Println("loaded env values")
	}

	server.Initialize(
		os.Getenv("DB_PATH"),
		os.Getenv("API_PATH"),
		os.Getenv("FS_PATH"),
	)

	server.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
}
