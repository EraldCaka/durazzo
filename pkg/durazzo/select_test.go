package durazzo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_Select_All(t *testing.T) {

	newDurazzo := setupDatabase(t)

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)

	insertQuery := `INSERT INTO "user" (name, email) VALUES ($1, $2)`
	_, err = newDurazzo.Db.Exec(insertQuery, "edgar", "edgar@gmail.com")
	assert.Nil(t, err)

	_, err = newDurazzo.Db.Exec(insertQuery, "ermelinda", "ermelinda@gmail.com")
	assert.Nil(t, err)
	_, err = newDurazzo.Db.Exec(insertQuery, "kris", "kris@yahoo.com")
	assert.Nil(t, err)
	var users []User
	err = newDurazzo.Select(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 3, len(users))
	assert.Equal(t, "edgar", users[0].Name)
	assert.Equal(t, "edgar@gmail.com", users[0].Email)
	assert.Equal(t, "ermelinda", users[1].Name)
	assert.Equal(t, "ermelinda@gmail.com", users[1].Email)
	assert.Equal(t, "kris", users[2].Name)
	assert.Equal(t, "kris@yahoo.com", users[2].Email)

	tearDownDatabase(t, newDurazzo)
	assert.Nil(t, err)
}

func TestDurazzo_Select_limit(t *testing.T) {

	newDurazzo := setupDatabase(t)

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)

	insertQuery := `INSERT INTO "user" (name, email) VALUES ($1, $2)`
	_, err = newDurazzo.Db.Exec(insertQuery, "edgar", "edgar@gmail.com")
	assert.Nil(t, err)

	_, err = newDurazzo.Db.Exec(insertQuery, "ermelinda", "ermelinda@gmail.com")
	assert.Nil(t, err)
	_, err = newDurazzo.Db.Exec(insertQuery, "kris", "kris@yahoo.com")
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Limit(2).Run()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users))
	assert.Equal(t, "edgar", users[0].Name)
	assert.Equal(t, "edgar@gmail.com", users[0].Email)
	assert.Equal(t, "ermelinda", users[1].Name)
	assert.Equal(t, "ermelinda@gmail.com", users[1].Email)

	dropTableQuery := `DROP TABLE IF EXISTS "user"`
	_, err = newDurazzo.Db.Exec(dropTableQuery)
	assert.Nil(t, err)
}

func TestDurazzo_Select_Where(t *testing.T) {
	newDurazzo := setupDatabase(t)

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

	userBody1 := User{ID: 2, Name: "ermelinda", Email: "ermelinda@gmail.com"}
	err = newDurazzo.Insert(&userBody1).Run()
	assert.Nil(t, err)

	userBody2 := User{ID: 3, Name: "kris", Email: "kris@yahoo.com"}
	err = newDurazzo.Insert(&userBody2).Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Where("name", "kris").Run()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "kris", users[0].Name)
	assert.Equal(t, "kris@yahoo.com", users[0].Email)

	tearDownDatabase(t, newDurazzo)
}
