package userManagement

import (
	"net/http"
	"github.com/gorilla/mux"
)

func UserManagementRoutes(r *mux.Router) {
  r.HandleFunc("/signup", handleUserSignup).Methods(http.MethodPost)
  r.HandleFunc("/get-user", handleUserSignup).Methods(http.MethodGet)
  r.HandleFunc("/remove-user", handleUserDeletion).Methods(http.MethodDelete)
}
