package main

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"

	"log"
)

func main() {

	config := durazzo.Config{
		Driver: durazzo.Postgres,
		DSN:    "postgresql://postgres:1234@localhost:5432/keeper?sslmode=disable",
	}

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}
	type Product struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"unique size:100"`
		Price int
	}

	newDurazzo := durazzo.NewDurazzo(config)
	defer func(newDurazzo *durazzo.Durazzo) {
		err := newDurazzo.Close()
		if err != nil {
			log.Fatalf("Error closing Durazzo:%v", err)
		}
	}(newDurazzo)
	err := newDurazzo.AutoMigrate(&User{}, &Product{})
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	var userPtr User

	if err := newDurazzo.Select(&userPtr).
		Where("name", "erald").
		Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(userPtr)

	var users []User
	if err := newDurazzo.Select(&users).Where("name", "kris").Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	for _, user := range users {
		log.Println(user)
	}

	var users1 []*User
	if err := newDurazzo.Select(&users1).Limit(1).Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	for _, user := range users1 {
		log.Println(user)
	}

	var user *User
	if err := newDurazzo.Select(&user).Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(user)

}
