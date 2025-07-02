package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/CreateUser", corsMiddleware(makeHTTPHandleFunc(s.handleCreateUser)))
	router.HandleFunc("/CreateEmployee", corsMiddleware(makeHTTPHandleFunc(s.handleCreateEmployee)))
	router.HandleFunc("/GetUser", corsMiddleware(makeHTTPHandleFunc(s.handleGetUsers)))
	router.HandleFunc("/GetEmployee", corsMiddleware(makeHTTPHandleFunc(s.handleGetEmployee)))
	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Or your frontend origin
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allowed Origin
		w.Header().Set("Access-Control-Allow-Origin", "*") // Change this to your frontend origin
		// Allowed Methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		// Allowed Headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle Preflight Request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed to actual request
		next.ServeHTTP(w, r)
	}
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {

	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	//get account
	err := s.store.CreateUsers(
		&User{
			FullName: req.FullName,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: req.Password,
			GoldCard: req.GoldCard,
		})
	return err
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	workers, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	enableCors(&w)

	return WriteJSON(w, http.StatusOK, workers)
}
func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {

	fmt.Println("get to handleGetWorkerByID ")
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)

}

// ------------ EMPLOYEE
func (s *APIServer) handleCreateEmployee(w http.ResponseWriter, r *http.Request) error {

	req := new(CreateEmployeeRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	//get account
	err := s.store.CreateEmployee(
		&Employee{
			FullName: req.FullName,
			Email:    req.Email,
			Password: req.Password,
			Phone:    req.Phone,
			GoldCard: req.GoldCard,
		})
	return err
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, r *http.Request) error {
	workers, err := s.store.GetEmployee()
	if err != nil {
		return err
	}

	enableCors(&w)

	return WriteJSON(w, http.StatusOK, workers)
}
func (s *APIServer) handleGetEmployeeByID(w http.ResponseWriter, r *http.Request) error {

	fmt.Println("get to handleGetWorkerByID ")
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetEmployeeByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)

}

//-------------------------------

// func (s *APIServer) handleRegestration(w http.ResponseWriter, r *http.Request) error {
// 	enableCors(&w)
// 	req := new(LoginRequest)
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return err
// 	}

// 	if req.Email == "Krixo" || req.Password == "Nasro1234" {
// 		return WriteJSON(w, http.StatusOK, "Welcome Admin")
// 	}

// 	worker, err := s.store.Register(req.Password, req.Email)
// 	if err != nil {
// 		WriteJSON(w, http.StatusNotAcceptable, err)
// 		fmt.Println("PREBLEME FROM REGISTRATION ///////// WORKER :", worker)

// 	}

// 	token, err := createJWT(worker)
// 	if err != nil {
// 		return err
// 	}

// 	resp := LoginResponse{
// 		Token: token,
// 		ID:    worker.ID,
// 	}
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "x-jwt-token",
// 		Value:    token,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true, // Prevent JavaScript access
// 		Secure:   true, // Only send over HTTPS
// 		SameSite: http.SameSiteStrictMode,
// 		Path:     "/",
// 	})

// 	return WriteJSON(w, http.StatusOK, resp)
// }

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
