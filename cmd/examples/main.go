package main

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"

	"log"
)

func main() {

	config := durazzo.Config{
		Driver: durazzo.Postgres,
		DSN:    "postgresql://postgres:postgres@localhost:5432/testdb?sslmode=disable",
	}

	type User struct {
		ID    int    `durazzo:"primary_key""`
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

	userBody := User{ID: 1, Name: "emir", Email: "emir@gmail.com"}
	if err := newDurazzo.Insert(&userBody).Run(); err != nil {
		log.Fatalf("Error inserting user: %v", err)
	}
	var user User

	if err := newDurazzo.Select(&user).
		Where("name", "emir").
		Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(user)

	if err := newDurazzo.Update("user").
		Set("name", "kris").
		Set("email", "kris@gmail.com").
		Where("id", "1").
		Run(); err != nil {
		log.Fatalf("Error updating user: %v", err)
	}

	var userPtr User

	if err := newDurazzo.Select(&userPtr).
		Where("name", "kris").
		Run(); err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Println(userPtr)

	if err := newDurazzo.Delete("user").Where("id", "1").Run(); err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
	//var users []User
	//err = newDurazzo.Raw("SELECT * FROM user").
	//	Model(&users).
	//	Run()
	//
	//if err != nil {
	//	log.Fatalf("Error executing raw query: %v", err)
	//}
	//
	//for _, user := range users {
	//	log.Printf("User: %+v", user)
	//}

}
