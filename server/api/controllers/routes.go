package controllers

import "net/http"

// initializeRoutes sets up the routes for all the HTTP API endpoints.
func (s *Server) initializeRoutes(ApiPath string) {
	// Sets the home route.
	s.Router.HandleFunc(ApiPath, SetMiddlewareJSON(s.Home)).Methods("GET")

	// Sets the login route.
	s.Router.HandleFunc(ApiPath+"/login", SetMiddlewareJSON(s.Login)).Methods("POST")

	// Sets the routes for user endpoints.
	s.Router.HandleFunc(ApiPath+"/users", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.CreateUser, s, true,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/users/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetUserByID, s, false,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/users/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.UpdateUser, s, false,
	))).Methods("PUT")

	// Sets the routes for admin-only user endpoints.
	s.Router.HandleFunc(ApiPath+"/users", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetAllUsers,
		s, true,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/users/{id}", SetMiddlewareAuthentication(
		s.DeleteUser,
		s, true,
	)).Methods("DELETE")
	s.Router.HandleFunc(ApiPath+"/users", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.ReactivateUser, s, true,
	))).Methods("PUT")

	// Sets the routes for folder endpoints.
	s.Router.HandleFunc(ApiPath+"/folders", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.CreateFolder, s, false,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/folders/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetFolderByID, s, false,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/folders/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.UpdateFolder, s, false,
	))).Methods("PUT")
	s.Router.HandleFunc(ApiPath+"/folders/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.DeleteFolder, s, false,
	))).Methods("DELETE")

	// Sets the routes for file endpoints.
	s.Router.HandleFunc(ApiPath+"/files", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.CreateFile, s, false,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/files/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetFileByID, s, false,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/files/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.UpdateFile, s, false,
	))).Methods("PUT")
	s.Router.HandleFunc(ApiPath+"/files/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.DeleteFile, s, false,
	))).Methods("DELETE")

	// Sets the routes for file data endpoints.
	// These don't set the output header as JSON as they return raw file data.
	s.Router.HandleFunc(ApiPath+"/file-data/{id}", SetMiddlewareAuthentication(
		s.CreateFileData, s, false,
	)).Methods("PUT")
	s.Router.HandleFunc(ApiPath+"/file-data/{id}", SetMiddlewareAuthentication(
		s.GetFileData, s, false,
	)).Methods("GET")

	// Sets the routes for access role endpoints.
	s.Router.HandleFunc(ApiPath+"/access-roles", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.CreateAccessRole, s, true,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/access-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetAccessRoleByID, s, true,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/access-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.UpdateAccessRole, s, true,
	))).Methods("PUT")
	s.Router.HandleFunc(ApiPath+"/access-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.DeleteAccessRole, s, true,
	))).Methods("DELETE")

	// Sets the routes for user role endpoints.
	s.Router.HandleFunc(ApiPath+"/user-roles", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.CreateUserRole, s, true,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/user-roles", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetAllUserRoles, s, true,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/user-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.GetUserRoleByID, s, true,
	))).Methods("GET")
	s.Router.HandleFunc(ApiPath+"/user-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.UpdateUserRole, s, true,
	))).Methods("PUT")
	s.Router.HandleFunc(ApiPath+"/user-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.DeleteUserRole, s, true,
	))).Methods("DELETE")

	// Sets the routes for editing the user roles of a user.
	s.Router.HandleFunc(ApiPath+"/users/user-roles/{id}", SetMiddlewareJSON(SetMiddlewareAuthentication(
		s.AddUserRole, s, true,
	))).Methods("POST")
	s.Router.HandleFunc(ApiPath+"/users/user-roles/{id}", SetMiddlewareAuthentication(
		s.RemoveUserRole,
		s, true,
	)).Methods("DELETE")

	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/build/")))
}
