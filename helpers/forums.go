package helpers

import (
	"fmt"
	"reflect"

	"db_forum/database"

	"db_forum/models"

	"github.com/jackc/pgx"
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

const usersBySlug = `
SELECT DISTINCT 
			u.nickname, 
			u.email, 
			u.fullname,
			u.about
FROM 		users u
LEFT JOIN 	threads t on u.nickname = t.author
LEFT JOIN 	posts p on u.nickname = p.author
WHERE (p.forum = $1 OR t.forum = $1)
`
const descByUser = `
ORDER BY u.nickname DESC`

const ascByUser = `
ORDER BY u.nickname`

const sinceUserTrue = `
AND u.nickname < $3`

const sinceUserFalse = `
AND u.nickname > $3`

func GetUsersBySlug(slug string, limit []byte, desc []byte, since []byte) []*models.User {
	var users []*models.User
	transaction := database.StartTransaction()
	defer transaction.Commit()
	QueryString := usersBySlug
	if string(desc) == "true" {
		if len(since) > 0 {
			QueryString += sinceUserTrue
		}
		QueryString += descByUser
	} else {
		if len(since) > 0 {
			QueryString += sinceUserFalse
		}
		QueryString += ascByUser
	}
	if len(limit) > 0 {
		QueryString += limitQuery
	}
	var elements *pgx.Rows
	var err error
	if len(since) > 0 {

		if len(limit) > 0 {
			elements, err = transaction.Query(QueryString, slug, string(limit), string(since))
		} else {
			elements, err = transaction.Query(QueryString, slug, string(limit), string(since))
		}
	} else {
		if len(limit) > 0 {
			elements, err = transaction.Query(QueryString, slug, string(limit))
		} else {
			elements, err = transaction.Query(QueryString, slug)
		}

	}
	if err != nil {

		// fmt.Println(slug)
		// fmt.Println(string(limit))
		// fmt.Println("Я в ошибке")
		fmt.Println(reflect.TypeOf(elements))
		fmt.Println(err)
		// log.Fatal(err)
		return users
	} else {

		for elements.Next() {

			var user models.User
			if err := elements.Scan(
				&user.Nickname,
				&user.Email,
				&user.Fullname,
				&user.About); err != nil {
				fmt.Println(err)

				return users
			}

			users = append(users, &user)

		}

	}
	return users
}
