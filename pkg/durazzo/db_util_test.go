package durazzo_test

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Setup function to initialize the database connection for testing purposes
func setupDatabase(t *testing.T) *durazzo.Durazzo {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatalf("DATABASE_URL is not set in environment")
	}
	newDurazzo := durazzo.NewDurazzo(durazzo.Config{
		Driver: durazzo.Postgres,
		DSN:    databaseURL,
	})
	err := newDurazzo.AutoMigrate(&User{}, &Post{})
	assert.Nil(t, err)
	return newDurazzo
}

// Helper function to clean up database
func tearDownDatabase(t *testing.T, d *durazzo.Durazzo) {
	dropTableQuery := `DROP TABLE IF EXISTS "user"`
	_, err := d.Db.Exec(dropTableQuery)
	assert.Nil(t, err)
	dropTableQueryPost := `DROP TABLE IF EXISTS "post"`
	_, err = d.Db.Exec(dropTableQueryPost)
	assert.Nil(t, err)
}
