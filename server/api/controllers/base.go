// Package controllers handles the API server and processes incoming requests.
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/vincetiu8/penn-spark-server/api/filesystem"
	"github.com/vincetiu8/penn-spark-server/api/models"
)

// Server provides a struct that houses all aspects of the backend server.
// The Server Mutex protects against concurrency issues.
type Server struct {
	Mutex      sync.RWMutex
	FileSystem filesystem.FileSystem
	DB         *gorm.DB
	Router     *mux.Router
}

// Initialize sets up the server.
func (s *Server) Initialize(DbPath, ApiPath, FsPath string) {
	var err error

	// Setting up the filesystem.
	s.FileSystem = filesystem.FileSystem{
		Fs:       afero.NewOsFs(),
		FilePath: FsPath,
	}

	// Make a folder to host files
	err = s.FileSystem.MkdirAll(FsPath, 0666)
	if err != nil {
		log.Fatalln("can't initialize filesystem:", err)
	} else {
		fmt.Println("linked to filesystem")
	}

	// Setting up the database
	s.DB, err = gorm.Open(sqlite.Open(DbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln("can't connect to the database:", err)
	} else {
		fmt.Println("connected to the database")
	}

	// Creating tables for all structs in the database
	err = s.DB.AutoMigrate(&models.User{}, &models.Folder{}, &models.File{}, &models.UserRole{}, &models.AccessRole{})
	if err != nil {
		log.Fatalln("can't migrate tables", err)
	}

	// Seed the database with the minimum amount of information to be usable
	s.SeedDatabase()

	// Setting up the router
	s.Router = mux.NewRouter()

	// Initialize the api routes
	s.initializeRoutes(ApiPath)
}

// SeedDatabase seeds the database with default information, if not present.
func (s *Server) SeedDatabase() {

	// Creates a user if none are present
	user, err := models.GetUserByID(s.DB, 1)
	if err == models.ErrUserNotFound {
		user, err = models.CreateUser(s.DB, models.User{
			Username:  "admin",
			FirstName: "admin",
			LastName:  "admin",
			Password:  "password",
			IsAdmin:   true,
		})
		if err != nil {
			log.Fatalf("cannot seed user: %v", err)
		}
	} else if err != nil {
		log.Fatalf("cannot get user: %v", err)
	}

	// Creates the root folder if it doesn't exist
	folder, err := models.GetFolderByID(s.DB, 1)
	if err == models.ErrFolderNotFound {
		parentFolderID := uint(0)

		//
		err = s.DB.Create(&models.Folder{
			Name:           "root",
			LastEditorID:   user.ID,
			ParentFolderID: &parentFolderID,
		}).Take(&folder).Error
		if err != nil {
			log.Fatalf("cannot seed folder: %v", err)
		}
	} else if err != nil {
		log.Fatalf("cannot get folder: %v", err)
	}
}

// Run runs the server on the specified address.
func (s *Server) Run(addr string) {
	headersOk := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	credentialsOk := handlers.AllowCredentials()
	methodsOk := handlers.AllowedMethods([]string{
		"CREATE",
		"GET",
		"POST",
		"PUT",
		"DELETE",
	})

	fmt.Printf("Listening to port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(s.Router)))
}
