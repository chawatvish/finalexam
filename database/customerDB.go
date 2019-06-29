package database

import (
	"database/sql"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func CreateTable(db *sql.DB) error {

	row, err := db.Exec("SELECT 1 FROM customers LIMIT 1;")
	if row != nil {
		return nil
	}

	stmt := `
	CREATE TABLE customers(
		 id SERIAL PRIMARY KEY, 
		 name TEXT,
		 email TEXT,
		 status TEXT
		)
	`
	_, err = db.Exec(stmt)
	return err
}

func GetCustomers(db *sql.DB) ([]Customer, error) {
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Status)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func GetCustomerByID(db *sql.DB, id int) (Customer, error) {
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1;")
	if err != nil {
		return Customer{}, err
	}

	var customer Customer
	rows := stmt.QueryRow(id)
	err = rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Status)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func AddNewCustomer(db *sql.DB, customer Customer) (Customer, error) {
	query := `
	INSERT INTO customers (name, email, status) 
	VALUES ($1, $2, $3) 
	RETURNING id, name, email, status
	`

	var nCustomer Customer
	row := db.QueryRow(query, customer.Name, customer.Email, customer.Status)
	err := row.Scan(&nCustomer.ID, &nCustomer.Name, &nCustomer.Email, &nCustomer.Status)
	if err != nil {
		return Customer{}, err
	}

	return nCustomer, nil
}
