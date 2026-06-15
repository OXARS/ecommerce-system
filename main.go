package main

import (
	"ecommerce-system/controllers"
	"ecommerce-system/middleware"
	"ecommerce-system/models"
	"ecommerce-system/routes"
	"ecommerce-system/services"
)

func main() {

	routes.StartRoutes()
	middleware.Validate()

	productService := services.ProductService{}
	userService := services.UserService{}

	productController := controllers.ProductController{
		Service: productService,
	}

	userController := controllers.UserController{
		Service: userService,
	}

	product := models.Product{
		ID:          1,
		Name:        "Laptop",
		Description: "Laptop Gamer",
		Price:       1200.50,
		Stock:       5,
	}

	user := models.User{
		ID:       1,
		Name:     "Oscar",
		Email:    "oscar@gmail.com",
		Password: "123456",
		Role:     "Admin",
	}

	productController.AddProduct(product)
	userController.AddUser(user)
}
