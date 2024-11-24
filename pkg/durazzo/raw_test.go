package durazzo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	ID    int    `durazzo:"primary_key"`
	Name  string `durazzo:"size:100"`
	Email string `durazzo:"unique"`
}

type Post struct {
	ID     int    `durazzo:"primary_key"`
	Title  string `durazzo:"size:255"`
	Body   string `durazzo:"type:text"`
	UserID int    `durazzo:"foreign_key"`
}

func TestRaw_Insert(t *testing.T) {
	newDurazzo := setupDatabase(t)

	insertQuery := `INSERT INTO user (name, email) VALUES ($1, $2)`
	err := newDurazzo.Raw(insertQuery, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)

	err = newDurazzo.Raw(insertQuery, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users))
	assert.Equal(t, "edgar", users[0].Name)
	assert.Equal(t, "ermelinda", users[1].Name)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_Select(t *testing.T) {
	newDurazzo := setupDatabase(t)

	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Raw("SELECT * FROM user").Model(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users))
	assert.Equal(t, "edgar", users[0].Name)
	assert.Equal(t, "ermelinda", users[1].Name)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_Update(t *testing.T) {
	newDurazzo := setupDatabase(t)

	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)

	updateQuery := `UPDATE user SET email = $1 WHERE name = $2`
	err = newDurazzo.Raw(updateQuery, "edgar.updated@gmail.com", "edgar").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Raw("SELECT * FROM user WHERE name = $1", "edgar").Model(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "edgar.updated@gmail.com", users[0].Email)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_Delete(t *testing.T) {
	newDurazzo := setupDatabase(t)

	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)

	deleteQuery := `DELETE FROM user WHERE name = $1`
	err = newDurazzo.Raw(deleteQuery, "edgar").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Raw("SELECT * FROM user").Model(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "ermelinda", users[0].Name)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_ComplexQueries(t *testing.T) {
	newDurazzo := setupDatabase(t)

	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "kris", "kris@yahoo.com").Run()
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Raw("SELECT * FROM user WHERE email LIKE $1 ORDER BY id DESC", "%@gmail.com").Model(&users).Run()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users))
	assert.Equal(t, "ermelinda", users[0].Name)
	assert.Equal(t, "edgar", users[1].Name)
	var users1 []User
	err = newDurazzo.Raw("SELECT * FROM user ORDER BY id ASC LIMIT $1", 2).Model(&users1).Run()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users1))
	assert.Equal(t, "edgar", users1[0].Name)
	assert.Equal(t, "ermelinda", users1[1].Name)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_AggregateFunctions(t *testing.T) {
	newDurazzo := setupDatabase(t)

	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "kris", "kris@yahoo.com").Run()
	assert.Nil(t, err)

	tearDownDatabase(t, newDurazzo)
}
