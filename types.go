package main

type LoginResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateCommandRequest struct {
	FullName    string `json:"fullname"`
	Number      string `json:"number"`
	Flor        string `json:"flor"`
	Itemtype    string `json:"itemtype"`
	Service     string `json:"service"`
	Workers     string `json:"workers"`
	Start       string `json:"start"`
	Distination string `json:"distination"`
	Prix        string `json:"prise"`
	IsAccepted  string `json:"isaccepted"`
}

type Command struct {
	ID          string `json:"id"`
	FullName    string `json:"fullname"`
	Number      string `json:"number"`
	Flor        string `json:"flor"`
	Itemtype    string `json:"itemtype"`
	Service     string `json:"service"`
	Workers     string `json:"workers"`
	Start       string `json:"start"`
	Distination string `json:"distination"`
	Prix        string `json:"prise"`
	IsAccepted  string `json:"isaccepted"`
}

func NewCommand(fullname, number, flor, itemtype, service, workers, start, distination, prix, isaccepted string) (*Command, error) {
	return &Command{
		FullName:    fullname,
		Number:      number,
		Flor:        flor,
		Itemtype:    itemtype,
		Service:     service,
		Workers:     workers,
		Start:       start,
		Distination: distination,
		Prix:        prix,
		IsAccepted:  isaccepted,
	}, nil
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
