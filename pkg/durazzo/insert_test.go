package durazzo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_Insert(t *testing.T) {
	newDurazzo := setupDatabase(t)
	defer tearDownDatabase(t, newDurazzo)

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)

	userBody2 := User{ID: 3, Name: "edgar", Email: "edgar@gmail.com"}
	err = newDurazzo.Insert(&userBody2).Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "edgar", users[0].Name)
	assert.Equal(t, "edgar@gmail.com", users[0].Email)
}
