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

	// Raw COUNT query

	var count int
	err = newDurazzo.Raw(`SELECT COUNT(*) FROM user WHERE email LIKE $1`, "%@gmail.com").Model(&count).Run()
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	// Raw AVG query (example)
	//var avgAge float64
	//err = newDurazzo.Raw(`SELECT AVG(age) FROM user`).Model(&avgAge).Run()
	//assert.Nil(t, err)
	//
	//assert.Equal(t, 0.0, avgAge)

	tearDownDatabase(t, newDurazzo)
}

func TestRaw_ComplexJoinQuery(t *testing.T) {
	newDurazzo := setupDatabase(t)

	// Insert user data
	err := newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "edgar", "edgar@gmail.com").Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO user (name, email) VALUES ($1, $2)`, "ermelinda", "ermelinda@gmail.com").Run()
	assert.Nil(t, err)

	// Insert post data
	err = newDurazzo.Raw(`INSERT INTO post (title, body, userid) VALUES ($1, $2, $3)`, "Post 1", "Body of post 1", 1).Run()
	assert.Nil(t, err)
	err = newDurazzo.Raw(`INSERT INTO post (title, body, userid) VALUES ($1, $2, $3)`, "Post 2", "Body of post 2", 2).Run()
	assert.Nil(t, err)

	type result struct {
		UserName  string
		PostTitle string
	}

	var results []result
	err = newDurazzo.Raw(`
		SELECT u.name AS UserName, p.title AS PostTitle
		FROM user u
		LEFT JOIN post p ON u.id = p.userid
		WHERE u.email LIKE $1
		ORDER BY u.id DESC
	`, "%@gmail.com").Model(&results).Run()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "ermelinda", results[0].UserName)
	assert.Equal(t, "Post 2", results[0].PostTitle)
	assert.Equal(t, "edgar", results[1].UserName)
	assert.Equal(t, "Post 1", results[1].PostTitle)

	tearDownDatabase(t, newDurazzo)
}
