package main

type LoginResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User
type CreateUserRequest struct {
	FullName string `json:"fullname"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
	GoldCard string `json:"goldcard"`
}

type User struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
	GoldCard string `json:"goldcard"`
}

// Employee
type Employee struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
	GoldCard string `json:"goldcard"`
	Service  string `json:"service"`
}

type CreateEmployeeRequest struct {
	FullName string `json:"fullname"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
	GoldCard string `json:"goldcard"`
	Service  string `json:"service"`
}

type jwtInupt struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type BookService struct {
	Id           string `json:"id"`
	UserId       string `json:"userid"`
	EmployeeId   string `json:"employeeId"`
	Service      string `json:"service"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Location     string `json:"location"`
	IsAuthorized string `json:"isaothorized"`
	Price        string `json:"price"`
}

type BookServiceRequest struct {
	UserId       string `json:"userid"`
	EmployeeId   string `json:"employeeId"`
	Service      string `json:"service"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Location     string `json:"location"`
	IsAuthorized string `json:"isaothorized"`
	Price        string `json:"price"`
}

type ID struct {
	Id string `json:"id"`
}
