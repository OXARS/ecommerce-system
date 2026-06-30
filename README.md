# Sistema de Gestión de E-commerce en Go

## Descripción

El presente proyecto consiste en el desarrollo de una aplicación web para la gestión de un e-commerce. El sistema fue desarrollado en el lenguaje de programación Go y utiliza el framework Gin para la creación de servicios web REST.

La aplicación permite administrar usuarios, productos, inventario y pedidos mediante solicitudes HTTP y respuestas en formato JSON.

## Justificación

Se seleccionó un sistema de gestión de e-commerce porque el comercio electrónico permite automatizar procesos importantes como el registro de productos, control de inventario, administración de usuarios y generación de pedidos.

Este proyecto permite aplicar conocimientos de programación orientada a objetos, estructuras de datos, interfaces, encapsulación, manejo de errores, concurrencia y pruebas de software.

## Objetivo general

Desarrollar una aplicación web de gestión de e-commerce utilizando Go, servicios web REST y programación orientada a objetos, permitiendo administrar usuarios, productos, inventario y pedidos de manera organizada y segura.

## Objetivos específicos

- Diseñar una arquitectura modular para la aplicación.
- Implementar servicios web utilizando Gin.
- Utilizar JSON para recibir y devolver información.
- Aplicar interfaces y encapsulación.
- Implementar manejo de errores.
- Proteger la información compartida mediante concurrencia segura.
- Crear pruebas unitarias, de integración y aceptación.
- Mantener el código fuente en un repositorio de GitHub.

## Tecnologías utilizadas

- Go
- Gin Web Framework
- JSON
- Git
- GitHub
- Visual Studio Code
- PowerShell
- bcrypt
- Paquete `testing` de Go
- `net/http/httptest`
- `sync.RWMutex`
- Goroutines
- Channels
- WaitGroup

## Arquitectura del proyecto

```text
ecommerce-system/
│
├── controllers/
│   ├── health_controller.go
│   ├── product_controller.go
│   ├── user_controller.go
│   └── order_controller.go
│
├── interfaces/
│   ├── product_interface.go
│   ├── user_interface.go
│   └── order_interface.go
│
├── middleware/
│   └── validation.go
│
├── models/
│   ├── product.go
│   ├── user.go
│   ├── order.go
│   └── cart.go
│
├── routes/
│   ├── routes.go
│   └── routes_test.go
│
├── services/
│   ├── product_service.go
│   ├── user_service.go
│   ├── order_service.go
│   ├── product_service_test.go
│   ├── user_service_test.go
│   └── order_service_test.go
│
├── main.go
├── go.mod
├── go.sum
└── README.md