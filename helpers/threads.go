package helpers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/iamfrommoscow/db_forum/database"
	"github.com/iamfrommoscow/db_forum/models"
	"github.com/jackc/pgx"
)

const insertThread = `
INSERT INTO threads (
	author, 
	created, 
	forum, 
	message, 
	title, 
	slug,
	id
	) 
VALUES (
	$1, 
	$2, 
	(SELECT slug FROM forums WHERE slug = $3), 
	$4, 
	$5, 
	$6,
	$7)
`

const iterateThreads = `
UPDATE forums
SET threads = threads + 1
WHERE slug = $1`
const selectCount = `SELECT COUNT(*) FROM threads`

//int лишняя
func CreateThread(thread *models.Thread) (int, error) {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var threadID int
	if err := transaction.QueryRow(selectCount).Scan(&threadID); err != nil {
		fmt.Println(err)
		return -1, err
	}
	threadID++
	var created string
	if thread.Created == "" {
		created = time.Now().Format("2006-01-02 15:04:05")
	} else {
		created = thread.Created
	}

	if _, err := transaction.Exec(insertThread, thread.Author, created, thread.Forum, thread.Message, thread.Title, thread.Slug, threadID); err != nil {
		// fmt.Println("Thread:")
		// fmt.Println(thread.Author)
		// fmt.Println(created)
		// fmt.Println(thread.Forum)
		// fmt.Println(thread.Message)
		// fmt.Println(thread.Title)
		// fmt.Println(thread.Slug)
		// fmt.Println(threadID)
		// fmt.Println("")
		fmt.Println(err)
		// log.Fatal(err)

		return threadID, err
	}
	if _, err := transaction.Exec(iterateThreads, thread.Forum); err != nil {
		fmt.Println(err)
		return threadID, err
	}
	return threadID, nil

}

const selectByLimit = `
SELECT 	author, 
		created, 
		forum, 
		message, 
		title, 
		slug,
		id 
FROM threads 
WHERE forum = $1`

const sinceQueryTrue = `
AND created <= $3`

const sinceQueryFalse = `
AND created >= $3`

const descByTime = `
ORDER BY created DESC`

const ascByTime = `
ORDER BY created, id`

const limitQuery = `
LIMIT $2`

func GetThreadsByForum(slug string, limit []byte, desc []byte, since []byte) []*models.Thread {
	var threads []*models.Thread
	transaction := database.StartTransaction()
	defer transaction.Commit()
	QueryString := selectByLimit

	if string(desc) == "true" {
		if len(since) > 0 {
			QueryString += sinceQueryTrue
		}
		QueryString += descByTime
	} else {
		if len(since) > 0 {
			QueryString += sinceQueryFalse
		}
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
		fmt.Println(reflect.TypeOf(elements))
		fmt.Println(err)
		// log.Fatal(err)
		return threads
	} else {

		for elements.Next() {

			var thread models.Thread
			var created time.Time
			if err := elements.Scan(
				&thread.Author,
				&created,
				&thread.Forum,
				&thread.Message,
				&thread.Title,
				&thread.Slug,
				&thread.ID); err != nil {
				fmt.Println(err)

				return threads
			}
			thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")

			threads = append(threads, &thread)

		}

	}
	return threads
}

const selectThreadBySlug = `
SELECT author,
		created,
		forum,
		message,
		title,
		slug,
		id
FROM threads
WHERE lower(slug) = lower($1)`

func GetThreadBySlug(slug string) *models.Thread {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var thread models.Thread
	var created time.Time
	if err := transaction.QueryRow(selectThreadBySlug, slug).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID); err != nil {

		return nil
	} else {
		thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
		return &thread
	}
}

const selectThreadByID = `
SELECT author,
		created,
		forum,
		message,
		title,
		slug,
		id
FROM threads
WHERE id = $1`

func GetThreadByID(id string) *models.Thread {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var thread models.Thread
	var created time.Time
	if err := transaction.QueryRow(selectThreadByID, id).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID); err != nil {
		fmt.Println(err)
		return nil
	} else {
		thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
		return &thread
	}
}

const UpdateThreadQuery = `UPDATE threads
SET
	message = $2, 
	title = $3
WHERE
	lower(slug) = lower($1)
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id
`

func UpdateThreadBySlug(slug string, message string, title string) *models.Thread {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var thread models.Thread
	var created time.Time
	if err := transaction.QueryRow(UpdateThreadQuery, slug, message, title).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID); err != nil {
		fmt.Println(err)
		// fmt.Println(UpdateThreadQuery)
		// fmt.Println(slug)
		// fmt.Println(message)
		// fmt.Println(title)
		// fmt.Println("")
		return nil
	} else {

		// thread.Message = message
		// fmt.Println(thread.Message)
		// thread.Title = title
		thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
		return &thread
	}
}

const UpdateThreadQueryID = `UPDATE threads
SET
	message = $2, 
	title = $3
WHERE
	id = $1
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id
`

func UpdateThreadByID(slug string, message string, title string) *models.Thread {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var thread models.Thread
	var created time.Time
	if err := transaction.QueryRow(UpdateThreadQueryID, slug, message, title).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID); err != nil {
		fmt.Println(err)
		// fmt.Println(UpdateThreadQuery)
		// fmt.Println(slug)
		// fmt.Println(message)
		// fmt.Println(title)
		// fmt.Println("")
		return nil
	} else {

		// thread.Message = message
		// fmt.Println(thread.Message)
		// thread.Title = title
		thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
		return &thread
	}
}
