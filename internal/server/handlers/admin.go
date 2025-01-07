package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"auther/internal/auth"
	"auther/internal/db/models"

	"github.com/gorilla/mux"
)

func CreateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := auth.CreateUser(db, &user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := auth.DeleteUser(db, &user); err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteUserByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		if err := auth.DeleteUserByID(db, userID); err != nil {
			http.Error(w, "Failed to delete user by ID", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
func DeleteUserByLoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var deleteByLogin struct {
			Login string `json:"login"`
		}

		if err := json.NewDecoder(r.Body).Decode(&deleteByLogin); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := auth.DeleteUserByLogin(db, deleteByLogin.Login); err != nil {
			http.Error(w, "Failed to delete user by login", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
