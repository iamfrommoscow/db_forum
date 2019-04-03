package helpers

import (
	"fmt"

	"github.com/iamfrommoscow/db_forum/database"
	"github.com/iamfrommoscow/db_forum/models"
)

const insertForum = `
INSERT INTO forums (slug, title, "user") 
VALUES ($1, $2, (SELECT nickname
	FROM users 
	WHERE nickname = $3
	)) `

func CreateForum(newForum *models.Forum) (error, string) {
	transaction := database.StartTransaction()
	defer transaction.Commit()

	if _, err := transaction.Exec(insertForum, newForum.Slug, newForum.Title, newForum.User); err != nil {
		transaction.Rollback()
		return err, newForum.Title
	}
	nickname := FindByNickname(newForum.User)

	return nil, nickname.Nickname
}

const selectBySlug = `
SELECT slug, title, "user"
FROM forums 
WHERE slug = $1`

func FindBySlug(slug string) *models.Forum {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var forum models.Forum
	if err := transaction.QueryRow(selectBySlug, slug).Scan(&forum.Slug, &forum.Title, &forum.User); err != nil {

		return nil
	} else {
		return &forum
	}
}

const postsCount = `
SELECT COUNT(*) FROM posts 
WHERE forum = $1`

func GetPostsCountByForum(slug string) int {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var count int
	if err := transaction.QueryRow(postsCount, slug).Scan(&count); err != nil {
		fmt.Println(err)
		return 0
	} else {
		return count
	}
}

const threadsCount = `
SELECT COUNT(*) FROM threads 
WHERE forum = $1`

func GetThreadsCountByForum(slug string) int {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var count int
	if err := transaction.QueryRow(threadsCount, slug).Scan(&count); err != nil {
		fmt.Println(err)
		return 0
	} else {
		return count
	}
}
