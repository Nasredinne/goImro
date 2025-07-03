package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	CreateUsers(*User) error
	GetUsers() ([]*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UserRegister(password string, email string) (*User, error)
	CreateEmployee(*Employee) error
	GetEmployee() ([]*Employee, error)
	GetEmployeeByID(id string) (*Employee, error)
	EmployeeRegister(password string, email string) (*Employee, error)
	GetEmployeeByEmail(email string) (*Employee, error)
	CreateBookSevice(*BookService) error
	GetBookServices() ([]*BookService, error)
	GetBookServiceByEmployee(*ID) ([]*BookService, error)
	AutoriseBookService(bookservice *BookService) error
	UpdatePrice(bookservice *BookService) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	// connStr := "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=goImro sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return nil, err
	// }

	// DO THIS BEFORE PUSH
	dsn := os.Getenv("DB_HOST")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database!")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {

	if err := s.createUserTable(); err != nil {
		fmt.Println("Proble in Creation user table")
		return err
	}
	if err := s.createEmployeeTable(); err != nil {
		return err
	}
	if err := s.createEmployeeTable(); err != nil {
		return err
	}
	if err := s.createBookServiceTable(); err != nil {
		return err
	}

	return nil
}

// USER
func (s *PostgresStore) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(100) NOT NULL ,
    phone VARCHAR(20) NOT NULL ,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL ,
	goldcard VARCHAR(100) NOT NULL 
);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUsers(user *User) error {
	hashedpassword, err := s.CreateUser(user.Email, user.Password)
	if err != nil {
		return err
	} else {
		query := `INSERT INTO users (fullname, phone, email, password, goldcard) 
								VALUES ($1, $2, $3, $4, $5); `

		_, err := s.db.Query(query, user.FullName, user.Phone, user.Email, hashedpassword, user.GoldCard)

		return err
	}
}

func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("select * from users ")
	if err != nil {
		return nil, err
	}

	Users := []*User{}
	for rows.Next() {
		worker, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		Users = append(Users, worker)
	}

	return Users, nil
}

func (s *PostgresStore) UserRegister(password string, email string) (*User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password mismatch
		return nil, nil
	}
	return user, nil
}

// EMPLOYEE
func (s *PostgresStore) createEmployeeTable() error {
	query := `CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(100) NOT NULL ,
    phone VARCHAR(20) NOT NULL ,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL ,
	services VARCHAR(100) NOT NULL ,
	goldcard VARCHAR(100) NOT NULL 
);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateEmployee(employee *Employee) error {
	hashedpassword, err := s.CreateEmp(employee.Email, employee.Password)
	if err != nil {
		return err
	} else {
		query := `INSERT INTO employee (fullname, phone,  email, password, services, goldcard) 
								VALUES ($1, $2, $3, $4, $5, $6); `

		_, err := s.db.Query(query, employee.FullName, employee.Phone, employee.Email, hashedpassword, employee.Service, employee.GoldCard)

		return err
	}
}

func (s *PostgresStore) GetEmployee() ([]*Employee, error) {
	rows, err := s.db.Query("select * from employee")
	if err != nil {
		return nil, err
	}

	Employee := []*Employee{}
	for rows.Next() {
		employee, err := scanIntoEmployee(rows)
		if err != nil {
			return nil, err
		}
		Employee = append(Employee, employee)
	}

	return Employee, nil
}

func (s *PostgresStore) EmployeeRegister(password string, email string) (*Employee, error) {
	employee, err := s.GetEmployeeByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password))
	if err != nil {
		// Password mismatch
		return nil, nil
	}
	return employee, nil
}

// func (s *PostgresStore) Register(password string, email string) (*Worker, error) {
// 	worker, err := s.GetWorkerByEmail(email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = bcrypt.CompareHashAndPassword([]byte(worker.Password), []byte(password))
// 	if err != nil {
// 		// Password mismatch
// 		return nil, nil
// 	}

// 	return worker, nil
// }

func (s *PostgresStore) GetUserByID(id string) (*User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("account %s not found", id)
}

func (s *PostgresStore) GetEmployeeByID(id string) (*Employee, error) {
	rows, err := s.db.Query("select * from employee where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoEmployee(rows)
	}

	return nil, fmt.Errorf("account %s not found", id)
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	rows, err := s.db.Query("select * from users where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("Worker %s not found", email)
}

func (s *PostgresStore) GetEmployeeByEmail(email string) (*Employee, error) {
	rows, err := s.db.Query("select * from employee where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoEmployee(rows)
	}

	return nil, fmt.Errorf("Worker %s not found", email)
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.FullName,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.GoldCard,
	)

	return user, err
}

func scanIntoEmployee(rows *sql.Rows) (*Employee, error) {
	employee := new(Employee)
	err := rows.Scan(
		&employee.ID,
		&employee.FullName,
		&employee.Phone,
		&employee.Email,
		&employee.Password,
		&employee.Service,
		&employee.GoldCard,
	)

	return employee, err
}

func (s *PostgresStore) createBookServiceTable() error {
	query := `CREATE TABLE IF NOT EXISTS book_service (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    service VARCHAR(255) NOT NULL,
    date VARCHAR(50) NOT NULL,
    time VARCHAR(50) NOT NULL,
    location VARCHAR(255) NOT NULL,
    is_authorized BOOLEAN NOT NULL,
    price VARCHAR(50) NOT NULL,

    -- Foreign key constraints
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_employee FOREIGN KEY (employee_id) REFERENCES employee (id) ON DELETE CASCADE
);
`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateBookSevice(bookservice *BookService) error {
	query := `INSERT INTO book_service (
    user_id, employee_id, service, date, time, location, is_authorized, price
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ); `

	_, err := s.db.Query(query, bookservice.UserId, bookservice.EmployeeId, bookservice.Service, bookservice.Date, bookservice.Time, bookservice.Location, bookservice.IsAuthorized, bookservice.Price)

	return err
}

func (s *PostgresStore) GetBookServices() ([]*BookService, error) {
	rows, err := s.db.Query("select * from book_service")
	if err != nil {
		return nil, err
	}

	BookServices := []*BookService{}
	for rows.Next() {
		bookservice, err := scanIntoBookServices(rows)
		if err != nil {
			return nil, err
		}
		BookServices = append(BookServices, bookservice)
	}

	return BookServices, nil
}

func (s *PostgresStore) GetBookServiceByEmployee(req *ID) ([]*BookService, error) {
	rows, err := s.db.Query("SELECT * FROM book_service WHERE employee_id = $1;", req.Id)
	if err != nil {
		return nil, err
	}

	BookServices := []*BookService{}
	for rows.Next() {
		bookservice, err := scanIntoBookServices(rows)
		if err != nil {
			return nil, err
		}
		BookServices = append(BookServices, bookservice)
	}

	return BookServices, nil
}
func (s *PostgresStore) AutoriseBookService(bookservice *BookService) error {

	query := `
		UPDATE book_service
		SET is_authorized = $1
		WHERE id = $2
	`
	result, err := s.db.Exec(query, bookservice.IsAuthorized, bookservice.Id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no command found with ID %s", bookservice.Id)
	}

	return nil
}

func (s *PostgresStore) UpdatePrice(bookservice *BookService) error {

	query := `
		UPDATE book_service
		SET price = $1
		WHERE id = $2
	`
	result, err := s.db.Exec(query, bookservice.Price, bookservice.Id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no command found with ID %s", bookservice.Id)
	}

	return nil
}

func scanIntoBookServices(rows *sql.Rows) (*BookService, error) {
	services := new(BookService)
	err := rows.Scan(
		&services.Id,
		&services.UserId,
		&services.EmployeeId,
		&services.Service,
		&services.Date,
		&services.Time,
		&services.Location,
		&services.IsAuthorized,
		&services.Price,
	)

	return services, err
}
