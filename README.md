[![GoDoc](https://pkg.go.dev/badge/github.com/EraldCaka/durazzo)](https://pkg.go.dev/github.com/EraldCaka/durazzo)
[![GitHub Stars](https://img.shields.io/github/stars/EraldCaka/durazzo)](https://github.com/EraldCaka/durazzo/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/EraldCaka/durazzo)](https://github.com/EraldCaka/durazzo/network/members)
[![GitHub License](https://img.shields.io/github/license/EraldCaka/durazzo)](https://opensource.org/licenses/MIT)
[![GitHub Issues](https://img.shields.io/github/issues/EraldCaka/durazzo)](https://github.com/EraldCaka/durazzo/issues)

# Durazzo ORM

Durazzo is a simple and efficient ORM library for Go that provides a powerful, flexible, and easy-to-use interface for interacting with relational databases using raw SQL queries and model-based migrations. Still under development new features will come up soon.

---



## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [CRUD Operations](#crud-operations)
4. [Testing](#testing)

---

## Installation

To install Durazzo, run the following Go command:

```bash
  $ go get github.com/EraldCaka/durazzo
```

Once installed, you can begin using Durazzo in your Go projects to interact with databases easily.

---

## Quick Start

Durazzo works by first setting up a connection to your database, and then interacting with the database using models. Here’s how to get started with Durazzo:

1. **Create a model** (e.g., `User`):

```go
    type User struct {
        ID    int    `durazzo:"primary_key"`
        Name  string `durazzo:"size:100"`
        Email string `durazzo:"unique"`
    }
```

2. **Initialize Durazzo**:

```go
    package main
    
    import (
        "github.com/EraldCaka/durazzo/pkg/durazzo"
        "log"
    )
    
    func main() {
        db := durazzo.NewDurazzo(durazzo.Config{
            Driver: durazzo.Postgres,
            DSN:    "postgresql://username:password@localhost:5432/yourdbname?sslmode=disable",
        })
    
        err := db.AutoMigrate(&User{})
        if err != nil {
            log.Fatal("Error during migration:", err)
        }
		
        userBody := User{ID: 1, Name: "erald", Email: "erald@gmail.com"}
        if err := newDurazzo.Insert(&userBody).Run(); err != nil {
            log.Fatalf("Error inserting user: %v", err)
        }   
    }
```

---

## CRUD Operations

Durazzo makes it easy to work with your database through interfaces or raw SQL queries. Here are the common operations:

---
### Create


To insert a new record:

```go
    user := User{Name: "erald", Email: "erald@yahoo.com"}
    err := db.Insert(&user).Run()
```
---
### Select


Fetching records from the database:

```go
    var users []User
    err := db.Select(&users).Where("email", "erald@yahoo.com").Run()
```
---
### Update

To update a record:
```go
    err := db.Update(&User{}).Set("email", "erald@gmail.com").Where("name", "erald").Run()
```
---
### Delete

To delete a record:

```go
    err := db.Delete(&User{}).Where("name", "erald").Run()
```
---
### Raw SQL Queries

Durazzo allows you to execute raw SQL queries directly, for operations such as joins, complex selects, and more. 

```go
    var users1 []User
    err = newDurazzo.Raw("SELECT * FROM user ORDER BY id ASC LIMIT $1", 2).Model(&users1).Run()
```

---

## Testing

Durazzo comes with built-in support for testing your database operations.


### Running Tests

Once your database is set up for testing, you can execute tests using the following command:

```bash
  $ make test
```


---

## License

© EraldCaka, 2024

Released under the [MIT License](https://github.com/EraldCaka/durazzo/blob/main/license)