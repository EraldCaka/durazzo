package durazzo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_Update(t *testing.T) {
	newDurazzo := setupDatabase(t)
	defer tearDownDatabase(t, newDurazzo)

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)

	userBody := User{ID: 1, Name: "edgar", Email: "edgar@gmail.com"}
	err = newDurazzo.Insert(&userBody).Run()
	assert.Nil(t, err)

	err = newDurazzo.Update("user").
		Set("name", "kris").
		Set("email", "kris@yahoo.com").
		Where("id", "1").
		Run()

	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Where("name", "kris").Run()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "kris", users[0].Name)
	assert.Equal(t, "kris@yahoo.com", users[0].Email)
}
