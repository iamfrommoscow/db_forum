package helpers

import (
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
