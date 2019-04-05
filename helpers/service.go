package helpers

import (
	"fmt"

	"db_forum/database"

	"db_forum/models"
)

const countUsers = `
SELECT COUNT(*) FROM users`
const countThreads = `
SELECT COUNT(*) FROM threads`
const countForums = `
SELECT COUNT(*) FROM forums`
const countPosts = `
SELECT COUNT(*) FROM posts`

func GetCount() *models.Service {
	transaction := database.StartTransaction()
	defer transaction.Rollback()
	var service models.Service
	if err := transaction.QueryRow(countUsers).Scan(&service.Users); err != nil {
		fmt.Println(err)
	}
	if err := transaction.QueryRow(countThreads).Scan(&service.Threads); err != nil {
		fmt.Println(err)
	}
	if err := transaction.QueryRow(countForums).Scan(&service.Forums); err != nil {
		fmt.Println(err)
	}
	if err := transaction.QueryRow(countPosts).Scan(&service.Posts); err != nil {
		fmt.Println(err)
	}
	return &service
}

const deleteAll = `
TRUNCATE users CASCADE;
`

func DropDatabase() error {

	if err := database.Exec(deleteAll); err != nil {
		fmt.Println("drop database: ", err)
	}
	return nil
}
