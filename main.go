package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chawatvish/finalexam/database"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	p := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r := setupRouter()
	db = setupDatabase()
	r.Run(p)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(authenMiddleware)

	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomersByIDHandler)
	r.POST("/customers", addCustomerHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerByIDHandler)

	return r
}

func setupDatabase() *sql.DB {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Can't connect DB", err.Error())
		defer db.Close()
	}

	err = database.CreateTable(db)
	if err != nil {
		log.Fatal("Can't create table", err.Error())
		defer db.Close()
	}

	return db
}

func authenMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}
}

func responseError(c *gin.Context, errNumber int, err error, text string) {
	//Internal error print
	fmt.Printf("Error number : %d | info : %s", errNumber, err.Error())

	//Response to external
	//TODO
	c.JSON(http.StatusInternalServerError, gin.H{"error": text})
}

func getCustomersHandler(c *gin.Context) {
	customers, err := database.GetCustomers(db)
	if err != nil {
		responseError(c,
			http.StatusInternalServerError,
			err,
			"Can't find customer")
		return
	}

	c.JSON(200, customers)
}

func getCustomersByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responseError(c,
			http.StatusBadRequest,
			err,
			"Wrong parameter")
		return
	}

	todo, err := database.GetCustomerByID(db, id)
	if err != nil {
		responseError(c,
			http.StatusInternalServerError,
			err,
			fmt.Sprintf("Can't find customer of id : %d", id))
		return
	}

	c.JSON(200, todo)
}

func addCustomerHandler(c *gin.Context) {
	var customer database.Customer
	c.BindJSON(&customer)
	customer, err := database.AddNewCustomer(db, customer)
	if err != nil {
		responseError(c,
			http.StatusInternalServerError,
			err,
			"Can't create new customer")
		return
	}

	c.JSON(201, customer)
}

func updateCustomerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responseError(c,
			http.StatusBadRequest,
			err,
			"Wrong parameter")
		return
	}

	var customer database.Customer
	c.BindJSON(&customer)
	customer.ID = id
	err = database.UpdateCustomerInfo(db, customer)
	if err != nil {
		responseError(c,
			http.StatusInternalServerError,
			err,
			fmt.Sprintf("Can't update customer of id : %d", id))
		return
	}

	c.JSON(200, customer)
}

func deleteCustomerByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responseError(c,
			http.StatusBadRequest,
			err,
			"Wrong parameter")
		return
	}

	_, err = database.DeleteTodoByID(db, id)
	if err != nil {
		responseError(c,
			http.StatusInternalServerError,
			err,
			fmt.Sprintf("Can't delete customer of id : %d", id))
		return
	}

	c.JSON(200, gin.H{"message": "customer deleted"})
}
