package helpers

import (
	"fmt"
	"time"

	"github.com/iamfrommoscow/db_forum/database"
	"github.com/iamfrommoscow/db_forum/models"
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
		fmt.Println(err)
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
WHERE forum = $1
`

const descByTime = `
ORDER BY created DESC`

const limitQuery = `
LIMIT $2`

func GetThreadsByForum(slug string, limit []byte, desc []byte) []*models.Thread {
	var threads []*models.Thread
	transaction := database.StartTransaction()
	defer transaction.Commit()
	if elements, err := transaction.Query(selectByLimit+descByTime+limitQuery, slug, limit); err != nil {
		fmt.Println(err)
		return threads
	} else {
		for elements.Next() {
			var thread models.Thread
			if err := elements.Scan(
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Message,
				&thread.Title,
				&thread.Slug,
				&thread.ID); err != nil {
				return threads
			}
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
