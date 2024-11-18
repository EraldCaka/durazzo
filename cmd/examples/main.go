package main

import (
	"fmt"
	"github.com/EraldCaka/durazzo/pkg/orm"

	"log"
)

func main() {

	dbType := orm.Postgres
	config := orm.Config{
		Driver: dbType,
		DSN:    "postgresql://postgres:1234@localhost:5432/keeper?sslmode=disable",
	}

	type User struct {
		ID    int    `orm:"primary_key"`
		Name  string `orm:"size:100"`
		Email string `orm:"unique"`
	}

	durazzo := orm.NewDurazzo(config)

	var users []*User
	err := durazzo.Select(users).Run()
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	fmt.Println(users) // TODO : FIX THIS ISSUE

}
