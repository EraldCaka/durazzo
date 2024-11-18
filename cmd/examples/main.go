package main

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"

	"log"
)

func main() {

	dbType := durazzo.Postgres
	config := durazzo.Config{
		Driver: dbType,
		DSN:    "postgresql://postgres:1234@localhost:5432/keeper?sslmode=disable",
	}

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	newDurazzo := durazzo.NewDurazzo(config)

	err := newDurazzo.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("ptr user")
	var userPtr User
	if err := newDurazzo.Select(&userPtr).Where("name", "erald2").Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(userPtr)

	var users []User
	if err := newDurazzo.Select(&users).Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	for _, user := range users {
		log.Println(user)
	}

	log.Println("users ptr")
	var users1 []*User
	if err := newDurazzo.Select(&users1).Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	for _, user := range users1 {
		log.Println(user)
	}

	log.Println("user")
	var user User
	if err := newDurazzo.Select(&user).Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(user)

}
