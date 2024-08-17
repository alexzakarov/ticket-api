# Ticket-App

- ## General Information
  - ### The application is written on hexagonal architecture
  - ### Fiber is used for http server. 
  - ### Embedded-postgres in-memory db library is used to handle db queries for test cases.
  - ### I18n implementation is working with "x-lang-code" header. 

- ## Service File System


```
    ├── cmd                                     # Root of main
    │   ├── grpc                                # ticket_service main
    ├── config                                  # Configuration parser, structure and .yaml file
    │   ├── config.go                           # Config file reader and parser 
    │   ├── config.yaml                         # Config .yaml file
    │   └── models.go                           # Config model structure
    ├── docs                                    # Swagger documentation
    │   ├── docs.go                             # Document parser
    │   ├── swagger.json                        # Auto generated swagger JSON file 
    │   └── swagger.yaml                        # Auto generated swagger YAML file 
    ├── internal                                # Hexagonal - DDD Layer
    │   └── auth                                # Auth Domain
    │       ├── application                     # Application Layer
    │       │   ├──service                      # Usecase Layer
    │       │   │   └── service.go              # Usecase file of current domain
    │       │   └── tests                       # Service Test Layer
    │       │       └── service_test.go         # Service Test file of current domain
    │       ├── domain                          # Domain Layer
    │       │   └── entities                    # Models/Structures
    │       │       └── entities.go             # Models/Structures file of current domain
    │       │   └── errors                      # Custom Errors
    │       │       └── errors.go               # Custom Errors file of current domain
    │       │   └── ports                       # Ports / Interfaces
    │       │       └── http_handler.go         # HTTP Handlers Interfaces file of current domain
    │       │       └── postgresql.go           # Postgresql Interfaces file of current domain
    │       │       └── service.go              # Usecase Interface file of current domain
    │       │   └── utils                       # Specific Utilities
    │       │       └── utils.go                # Specific Utilities file of current domain
    │       └── handler                         # Handler Layer
    │       │   └── http                        # HTTP Handlers
    │       │       └── handlers.go             # HTTP Handler file of current domain  
    │       │       └── routes.go               # HTTP routes file of current domain  
    │       ├── infrastructure                  # Infrastructure Layer
    │       │   └── repository                  # Repositories 
    │       │       └── postgresql.go           # PostgresqlRepository file of current domain
    │       │   └── tests                       # Postgresql Test Layer
    │       │       └── postgresql_test.go      # Postgresql Test file of current domain
    ├── pkg                                     # Custom packages for general common usage
    └── ...
```
- ## Tech Stack
  - ### Golang (v1.23.0)
  - ### Fiber
  - ### Postgresql
  - ### Docker

- ## Makefile commands
```shell

make swag # for generate swagger documentations

```