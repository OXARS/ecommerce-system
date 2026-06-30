package main

import (
	"fmt"
	"log"

	"ecommerce-system/routes"
)

func main() {
	router := routes.SetupRouter()

	fmt.Println("Servidor e-commerce iniciado")
	fmt.Println("Dirección: http://localhost:8080")
	fmt.Println("Servicio de prueba: http://localhost:8080/api/health")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("No fue posible iniciar el servidor: ", err)
	}
}
