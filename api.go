package main

import (
	"encoding/json"
	"fmt"
	"time"

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
	router.HandleFunc("/UserLogin", corsMiddleware(makeHTTPHandleFunc(s.handleUserRegestration)))
	router.HandleFunc("/Employee", corsMiddleware(makeHTTPHandleFunc(s.handleEmployeeRegestration)))
	router.HandleFunc("/CreateBookSevice", corsMiddleware(makeHTTPHandleFunc(s.handleCreateBookService)))
	router.HandleFunc("/GetBookService", corsMiddleware(makeHTTPHandleFunc(s.handleGetBookService)))

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
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	enableCors(&w)

	return WriteJSON(w, http.StatusOK, users)
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
			Phone:    req.Phone,
			Email:    req.Email,
			Password: req.Password,
			GoldCard: req.GoldCard,
			Service:  req.Service,
		})
	return err
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, r *http.Request) error {
	employees, err := s.store.GetEmployee()
	if err != nil {
		return err
	}

	enableCors(&w)

	return WriteJSON(w, http.StatusOK, employees)
}

//-------------------------------

func (s *APIServer) handleUserRegestration(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Probleme in decoding")
		return err
	}

	// worker, err := s.store.UserRegister(req.Password, req.Email)
	user, err := s.store.UserRegister(req.Password, req.Email)
	if err != nil {
		WriteJSON(w, http.StatusNotAcceptable, err)
		fmt.Println("PREBLEME FROM REGISTRATION ///////// WORKER :", user)

	}

	userjwt := jwtInupt{user.ID, user.Email}

	token, err := createJWT(&userjwt)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "x-jwt-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Only send over HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleEmployeeRegestration(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// worker, err := s.store.UserRegister(req.Password, req.Email)
	user, err := s.store.EmployeeRegister(req.Password, req.Email)
	if err != nil {
		WriteJSON(w, http.StatusNotAcceptable, err)
		fmt.Println("PREBLEME FROM REGISTRATION ///////// WORKER :", user)

	}

	userjwt := jwtInupt{user.ID, user.Email}

	token, err := createJWT(&userjwt)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "x-jwt-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Only send over HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateBookService(w http.ResponseWriter, r *http.Request) error {
	//decode json
	req := new(BookServiceRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	//get account
	err := s.store.CreateBookSevice(
		&BookService{
			UserId:       req.UserId,
			EmployeeId:   req.EmployeeId,
			Service:      req.Service,
			Date:         req.Date,
			Time:         req.Time,
			Location:     req.Location,
			IsAuthorized: req.IsAuthorized,
			Price:        req.Price,
		})
	return err
}

func (s *APIServer) handleGetBookService(w http.ResponseWriter, r *http.Request) error {

	//calldb funtion
	bookservices, err := s.store.GetBookServices()
	if err != nil {
		return err
	}
	//return
	return WriteJSON(w, http.StatusAccepted, bookservices)
}

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
