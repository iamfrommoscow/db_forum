package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/iamfrommoscow/db_forum/database"
	"github.com/iamfrommoscow/db_forum/models"
	"github.com/jackc/pgx"
)

const insertPost = `
INSERT into posts(
	author, 
	created, 
	forum, 
	message, 
	parent, 
	thread,
	id
	)
VALUES (
	$1, 
	$2, 
	(SELECT slug FROM forums WHERE slug = $3), 
	$4, 
	$5,
	(SELECT id FROM threads WHERE id = $6),
	$7)
`
const selectPostsCount = `SELECT COUNT(*) FROM posts`

func InsertPosts(posts []*models.Post) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var postID int
	if err := transaction.QueryRow(selectPostsCount).Scan(&postID); err != nil {
		// fmt.Println(err)
		return err
	}
	for _, post := range posts {
		postID++
		post.ID = postID
		if _, err := transaction.Exec(insertPost, post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread, postID); err != nil {
			fmt.Println("Post:")
			fmt.Println(post.Author)
			fmt.Println(post.Created)
			fmt.Println(post.Forum)
			fmt.Println(post.Message)
			fmt.Println(post.Parent)
			fmt.Println(post.Thread)
			fmt.Println(postID)
			fmt.Println("")
			fmt.Println(err)
			// log.Fatal(err)

			return err
		}
	}

	return nil
}

const selectPostsByThread = `
SELECT 	author, 
		created, 
		forum, 
		message, 
		parent, 
		thread,
		id
FROM posts 
WHERE thread = $1`

const withParent = `
AND parent = $3`

const descByTimePost = `
ORDER BY created, id DESC`

const descByIdPost = `
ORDER BY id DESC`

const sinceQueryPost = `
AND id > `

func GetPostsTree(slug int, limit []byte, sort []byte, desc []byte, since []byte, parent int) []*models.Post {
	var posts []*models.Post
	transaction := database.StartTransaction()
	defer transaction.Commit()
	QueryString := selectPostsByThread + withParent
	if len(since) > 0 {
		QueryString += sinceQueryPost + ` $4`
	}

	if string(desc) == "true" && (parent == 0 || string(sort) == "tree") {
		QueryString += descByIdPost
	} else {
		QueryString += ascByTime
	}
	if len(limit) > 0 {
		QueryString += limitQuery
	}
	var elements *pgx.Rows
	var err error
	if len(since) > 0 {
		fmt.Println(QueryString)
		fmt.Println("$1", slug)
		fmt.Println("$2", string(limit))
		fmt.Println("$3", parent)
		fmt.Println("$4", string(since))

		elements, err = transaction.Query(QueryString, slug, string(limit), parent, string(since))
		if elements.Next() == false {

		}
	} else {
		elements, err = transaction.Query(QueryString, slug, string(limit), parent)
	}

	if err != nil {

		// fmt.Println(slug)
		// fmt.Println(string(limit))
		// fmt.Println("Я в ошибке")
		fmt.Println(err)
		// log.Fatal(err)
		return posts
	} else {

		for elements.Next() {

			var post models.Post
			var created time.Time
			if err := elements.Scan(
				&post.Author,
				&created,
				&post.Forum,
				&post.Message,
				&post.Parent,
				&post.Thread,
				&post.ID); err != nil {
				fmt.Println(err)

				return posts
			}
			post.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
			if string(desc) == "true" && string(sort) == "tree" {
				posts = append(posts, GetPostsTree(slug, limit, sort, desc, since, post.ID)...)
				posts = append(posts, &post)
			} else {
				posts = append(posts, &post)
				posts = append(posts, GetPostsTree(slug, limit, sort, desc, since, post.ID)...)
			}
			lim, _ := strconv.Atoi(string(limit))
			if len(limit) > 0 {
				if len(posts) >= lim {
					posts = posts[:lim]
					return posts
				}
			}

		}

	}
	return posts
}

func GetPostsByThread(slug int, limit []byte, sort []byte, since []byte, desc []byte) []*models.Post {
	var posts []*models.Post
	if string(sort) == "tree" {
		posts = GetPostsTree(slug, limit, sort, desc, since, 0)
		return posts
	}
	if string(sort) == "parent_tree" {
		posts = GetPostsTree(slug, limit, sort, desc, since, 0)
		return posts
	}
	transaction := database.StartTransaction()
	defer transaction.Commit()

	QueryString := selectPostsByThread
	if len(since) > 0 {
		QueryString += sinceQueryPost + ` $3`

	}
	if string(desc) == "true" {
		// if len(since) > 0 {
		// 	QueryString += sinceQueryTrue
		// }
		QueryString += descByIdPost
	} else {
		QueryString += ascByTime
	}
	if len(limit) > 0 {
		QueryString += limitQuery
	}
	var elements *pgx.Rows
	var err error
	if len(since) > 0 {

		elements, err = transaction.Query(QueryString, slug, string(limit), string(since))
	} else {
		elements, err = transaction.Query(QueryString, slug, string(limit))
	}

	if err != nil {

		// fmt.Println(slug)
		// fmt.Println(string(limit))
		// fmt.Println("Я в ошибке")
		fmt.Println("")
		fmt.Println(QueryString)
		fmt.Println(err)
		// log.Fatal(err)
		return posts
	} else {

		for elements.Next() {

			var post models.Post
			var created time.Time
			if err := elements.Scan(
				&post.Author,
				&created,
				&post.Forum,
				&post.Message,
				&post.Parent,
				&post.Thread,
				&post.ID); err != nil {
				fmt.Println(err)

				return posts
			}
			post.Created = created.Format("2006-01-02T15:04:05.000Z07:00")

			posts = append(posts, &post)

		}

	}
	return posts
}
