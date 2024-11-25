package durazzo_test

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

func TestDurazzo_Select_All(t *testing.T) {
	newDurazzo := setupDatabase(t)
	defer tearDownDatabase(t, newDurazzo)

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
}

func TestDurazzo_Select_limit(t *testing.T) {
	newDurazzo := setupDatabase(t)
	defer tearDownDatabase(t, newDurazzo)

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
}

func Test_Durazzo_Select_Async(t *testing.T) {
	newDurazzo := setupDatabase(t)

	defer func() {
		err := newDurazzo.Close()
		assert.Nil(t, err)
	}()

	insertTestData(t, newDurazzo)
	defer tearDownDatabase(t, newDurazzo)
	queries := []struct {
		Name  string
		Model *User
		Field string
		Value string
	}{
		{Name: "Query1", Model: &User{}, Field: "name", Value: "kris"},
		{Name: "Query2", Model: &User{}, Field: "email", Value: "erald@yahoo.com"},
		{Name: "Query3", Model: &User{}, Field: "id", Value: "3"},
		{Name: "Query4", Model: &User{}, Field: "name", Value: "sara"},
		{Name: "Query5", Model: &User{}, Field: "email", Value: "kris@yahoo.com"},
	}

	var wg sync.WaitGroup
	wg.Add(len(queries))
	results := make([]error, len(queries))

	for i, q := range queries {
		go func(index int, query struct {
			Name  string
			Model *User
			Field string
			Value string
		}) {
			defer wg.Done()

			err := newDurazzo.Select(query.Model).
				Where(query.Field, query.Value).
				Run()

			results[index] = err

			log.Printf("%s completed. Result: %+v, Error: %v\n", query.Name, query.Model, err)
		}(i, q)
	}

	wg.Wait()

	for i, err := range results {
		assert.Nil(t, err, "Error in query %d: %v", i+1, err)
		log.Printf("Query %d passed successfully\n", i+1)
	}

}

func insertTestData(t *testing.T, db *durazzo.Durazzo) {
	users := []struct {
		Name  string
		Email string
		ID    int
	}{
		{Name: "kris", Email: "kris@yahoo.com", ID: 1},
		{Name: "erald", Email: "erald@yahoo.com", ID: 2},
		{Name: "jessie", Email: "jessie@gmail.com", ID: 3},
		{Name: "sara", Email: "sara@hotmail.com", ID: 4},
	}

	for _, user := range users {
		_, err := db.Db.Exec(`INSERT INTO "user" (id, name, email) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			user.ID, user.Name, user.Email)
		assert.Nil(t, err, "Error inserting user data: %v", err)
	}
}
