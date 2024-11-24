package durazzo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_Delete(t *testing.T) {
	newDurazzo := setupDatabase(t)

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

	err = newDurazzo.Delete("user").Where("name", "edgar").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Where("name", "edgar").Run()
	assert.Nil(t, err)

	assert.Equal(t, 0, len(users))

	tearDownDatabase(t, newDurazzo)
}
