package util

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/spf13/afero"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/invincibot/penn-spark-server/api/controllers"
	"github.com/invincibot/penn-spark-server/api/filesystem"
	"github.com/invincibot/penn-spark-server/api/models"
)

var Users = []models.User{
	{
		Username:  "vincetiu8",
		FirstName: "Vince",
		LastName:  "Tiu",
		Password:  "password",
		IsAdmin:   true,
	},
	{
		Username:  "qtaro",
		FirstName: "jotaro",
		LastName:  "kujo",
		Password:  "star platinum",
		IsAdmin:   false,
	},
	{
		Username:  "Pet",
		FirstName: "Nugget",
		LastName:  "Tiu",
		Password:  "catto",
		IsAdmin:   false,
	},
	{
		Username:  "aleckgs",
		FirstName: "Alec",
		LastName:  "See",
		Password:  "12345",
		IsAdmin:   false,
	},
}

var Folders = []models.Folder{
	{
		Name: "root",
	},
	{
		Name: "folder1",
	},
	{
		Name: "folder2",
	},
}

var Files = []models.File{
	{
		Name:        "file1",
		IsPublished: true,
	},
	{
		Name:        "file2",
		IsPublished: false,
	},
	{
		Name:        "file3",
		IsPublished: true,
	},
}

var AccessRoles = []models.AccessRole{
	{
		AccessLevel: models.Publisher,
	},
	{
		AccessLevel: models.Uploader,
	},
	{
		AccessLevel: models.Viewer,
	},
	{
		AccessLevel: models.None,
	},
}

var UserRoles = []models.UserRole{
	{
		Name: "role1",
	},
	{
		Name: "role2",
	},
	{
		Name: "role3",
	},
	{
		Name: "role4",
	},
}

type TestServer struct {
	Server controllers.Server
	Data   SeedData
}

type SeedData struct {
	Users       []models.User
	Folders     []models.Folder
	Files       []models.File
	AccessRoles []models.AccessRole
	UserRoles   []models.UserRole
}

func NewTestServer() *TestServer {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	s := &TestServer{}

	s.RefreshFileSystem()

	s.Server.DB, err = gorm.Open(sqlite.Open(os.Getenv("TEST_DB_PATH")), &gorm.Config{})
	if err != nil {
		fmt.Printf("cannot connect to the database\n")
		log.Fatal("this is the error:", err)
	} else {
		fmt.Printf("connected to the database\n")
	}

	return s
}

func (s *TestServer) RefreshTable(i interface{}) error {
	err := s.Server.DB.Migrator().DropTable(i)
	if err != nil {
		return err
	}
	err = s.Server.DB.AutoMigrate(i)
	return err
}

func (s *TestServer) RefreshTables() error {
	v := reflect.ValueOf(s.Data)

	if s.Server.DB.Migrator().HasTable("assigned_access_roles") {
		err := s.Server.DB.Migrator().DropTable("assigned_access_roles")
		if err != nil {
			return err
		}
	}

	if s.Server.DB.Migrator().HasTable("assigned_user_roles") {
		err := s.Server.DB.Migrator().DropTable("assigned_user_roles")
		if err != nil {
			return err
		}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()
		err := s.RefreshTable(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TestServer) RefreshFileSystem() {
	s.Server.FileSystem = filesystem.FileSystem{
		Fs:       afero.NewMemMapFs(),
		FilePath: os.Getenv("TEST_FS_PATH"),
	}
}

func (s *TestServer) SeedData() error {
	err := s.RefreshTables()
	if err != nil {
		return err
	}

	users := make([]models.User, len(Users))
	for i, user := range Users {
		user, err = models.CreateUser(s.Server.DB, user)
		if err != nil {
			return err
		}
		users[i] = user
	}
	s.Data.Users = users

	folders := make([]models.Folder, len(Folders))
	for i, folder := range Folders {
		folder.LastEditorID = users[i].ID
		parentFolderID := uint(i)
		folder.ParentFolderID = &parentFolderID
		if i == 0 {
			err = s.Server.DB.Create(&folder).Take(&folder).Error
		} else {
			folder, err = models.CreateFolder(s.Server.DB, folder)
		}
		if err != nil {
			return err
		}
		folders[i] = folder
		if i > 0 {
			folders[i-1].ChildFolders = append(folders[i-1].ChildFolders, folder)
		}
	}
	s.Data.Folders = folders

	files := make([]models.File, len(Files))
	for i, file := range Files {
		file.LastEditorID = users[i].ID
		file.FolderID = folders[i].ID
		file, err = models.CreateFile(s.Server.DB, file)
		if err != nil {
			return err
		}
		files[i] = file
		folders[i].Files = append(folders[i].Files, file)
	}
	s.Data.Files = files

	userRoles := make([]models.UserRole, len(UserRoles))
	var accessRoles []models.AccessRole
	for i, userRole := range UserRoles {
		userRole, err = models.CreateUserRole(s.Server.DB, userRole)
		if err != nil {
			return err
		}
		userRoles[i] = userRole
		if i < len(UserRoles)-1 {
			for j := i; j < len(Folders); j++ {
				accessRole := AccessRoles[j]
				accessRole.FolderID = folders[j].ID
				accessRole.UserRoleID = userRole.ID
				accessRole, err = models.CreateAccessRole(s.Server.DB, accessRole)
				if err != nil {
					return err
				}
				accessRoles = append(accessRoles, accessRole)
				userRoles[i].AccessRoles = append(userRoles[i].AccessRoles, accessRole)
				folders[j].AccessRoles = append(folders[j].AccessRoles, accessRole)
			}
		} else {
			accessRole := AccessRoles[3]
			accessRole.FolderID = folders[2].ID
			accessRole.UserRoleID = userRole.ID
			accessRole, err = models.CreateAccessRole(s.Server.DB, accessRole)
			if err != nil {
				return err
			}
			accessRoles = append(accessRoles, accessRole)
			userRoles[i].AccessRoles = append(userRoles[i].AccessRoles, accessRole)
			folders[2].AccessRoles = append(folders[2].AccessRoles, accessRole)
		}
	}
	s.Data.UserRoles = userRoles
	s.Data.AccessRoles = accessRoles

	for i, user := range users {
		user.UserRoles = append(user.UserRoles, userRoles[i])
		user, err = models.UpdateUser(s.Server.DB, user, true)
		if err != nil {
			return err
		}
		users[i] = user
	}

	return nil
}
