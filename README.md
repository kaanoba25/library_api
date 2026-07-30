# Library Management API

A RESTful API built with Go, PostgreSQL, JWT Authentication, and Gorilla Mux using layered architecture.

## Features

- **Layered Architecture:** Handler, Service, Repository, Model pattern.
- **PostgreSQL & DB Transactions:** Safe stock tracking and atomic loan processes.
- **Authentication:** JWT tokens and password hashing via Bcrypt.
- **Role-Based Authorization:** Separate access for `admin` and `member` roles.
- **API Documentation:** Interactive Swagger UI documentation.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL (Docker container)
- **Router:** Gorilla Mux
- **Docs:** Swaggo / Swagger UI
