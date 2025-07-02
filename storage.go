package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUsers(*User) error
	GetUsers() ([]*User, error)
	GetUserByID(id string) (*User, error)
	CreateEmployee(*Employee) error
	GetEmployee() ([]*Employee, error)
	GetEmployeeByID(id string) (*Employee, error)
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
	hashedpassword, err := s.CreateUser(employee.Email, employee.Password)
	if err != nil {
		return err
	} else {
		query := `INSERT INTO employee (fullname, phone,  email, password, services, goldcard) 
								VALUES ($1, $2, $3, $4, $5); `

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

// func (s *PostgresStore) Register(password string, email string) (bool, error) {
// 	var hashedPassword string

// 	// Fetch the hashed password from the database
// 	query := `SELECT password FROM worker WHERE email = $1`
// 	err := s.db.QueryRow(query, email).Scan(&hashedPassword)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			// User not found
// 			return false, nil
// 		}
// 		// Other DB error
// 		return false, err
// 	}

// 	// Compare the plaintext password with the stored hash
// 	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	if err != nil {
// 		// Password mismatch
// 		return false, nil
// 	}

// 	// Password matches
// 	return true, nil
// }

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

// func (s *PostgresStore) GetCommandByID(id string) (*Command, error) {
// 	rows, err := s.db.Query("select * from commandsss where id = $1", id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		return scanIntoAccount(rows)
// 	}

// 	return nil, fmt.Errorf("command %s not found", id)
// }

// func (s *PostgresStore) UpdateCommand(command *Command) error {

// 	query := `
// 		UPDATE commandsss
// 		SET isaccepted = $1
// 		WHERE id = $2
// 	`
// 	result, err := s.db.Exec(query, command.IsAccepted, command.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute update query: %w", err)
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve affected rows: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no command found with ID %s", command.ID)
// 	}

// 	return nil
// }

// func (s *PostgresStore) UpdateWorker(worker *Worker) error {

// 	query := `
// 		UPDATE worker
// 		SET isaccepted = $1
// 		WHERE id = $2
// 	`
// 	result, err := s.db.Exec(query, worker.IsAccepted, worker.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute update query: %w", err)
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve affected rows: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no worker found with ID %s", worker.ID)
// 	}

// 	return nil
// }

// func scanIntoAccount(rows *sql.Rows) (*Command, error) {
// 	command := new(Command)
// 	err := rows.Scan(
// 		&command.ID,
// 		&command.FullName,
// 		&command.Number,
// 		&command.Flor,
// 		&command.Itemtype,
// 		&command.Service,
// 		&command.Workers,
// 		&command.Start,
// 		&command.Distination,
// 		&command.IsAccepted,
// 		&command.Prix,
// 	)

// 	return command, err
// }

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
