package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chawatvish/finalexam/database"

	"github.com/gin-gonic/gin"
)

func main() {
	p := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r := setupRouter()
	setupDatabase()
	r.Run(p)
}

/*
POST /customers
GET /customers/:id
GET /customers
PUT /customers/:id
DELETE /customers/:id

*/

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(authenMiddleware)

	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomersByIDHandler)
	r.POST("/customers", addCustomerHandler)
	// r.PUT("/customers/:id", updateCustomerHandler)
	// r.DELETE("/customers/:id", deleteCustomerByIDHandler)

	return r
}

func setupDatabase() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Can't connect DB", err.Error())
	}

	err = database.CreateTable(db)
	if err != nil {
		log.Fatal("Can't create table", err.Error())
	}
}

func authenMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}
}

// func responseError(c *gin.Context, errNumber int, err error) {
// 	//Internal error print
// 	fmt.Println("Error number : %d | info : %s", errNumber, err.Error())

// 	//Response to external
// 	//TODO
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// }

func getCustomersHandler(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	customers, err := database.GetCustomers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, customers)
}

func getCustomersByIDHandler(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo, err := database.GetCustomerByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todo)
}

func addCustomerHandler(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var customer database.Customer
	c.BindJSON(&customer)
	customer, err = database.AddNewCustomer(db, customer)
	fmt.Println(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, customer)
}
