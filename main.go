package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
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

	// r.Use(authMiddleware)

	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomersByIdHandler)
	r.POST("/customers", addCustomerHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerByIDHandler)

	return r
}
