package helpers

import (
	"fmt"
	"strconv"
	"time"

	"db_forum/database"

	"db_forum/models"
)

const plus2 = `
UPDATE threads
SET
	votes = votes + 2
WHERE
	id = $1
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id,
	votes
`

const minus2 = `
UPDATE threads
SET
	votes = votes - 2
WHERE
	id = $1
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id,
	votes
`
const like = `
UPDATE threads
SET
	votes = votes + 1 
WHERE
	id = $1
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id,
	votes
`

const dislike = `
UPDATE threads
SET
	votes = votes - 1
WHERE
	id = $1
RETURNING
	author,
	created,
	forum,
	message,
	title,
	slug,
	id,
	votes
`

const addVote = `
INSERT into votes(
	thread, 
	nickname, 
	voice
	)
VALUES (
	$1,
	$2, 
	$3
)
`
const findVoteByUserAndId = `
SELECT 	thread, 
		nickname, 
		voice
FROM 	votes
WHERE 	thread = $1 AND
		nickname = $2`

const dislikeVote = `
UPDATE votes
SET
	voice = -1
	WHERE 	thread = $1 AND
	nickname = $2
`

const likeVote = `
UPDATE votes
SET
	voice = 1
	WHERE 	thread = $1 AND
	nickname = $2
`

func VoteFound(id int, vote *models.Vote) *models.Vote {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var prevVote models.Vote
	err := transaction.QueryRow(findVoteByUserAndId, id, vote.Nickname).Scan(&prevVote.Thread, &prevVote.Nickname, &prevVote.Voice)
	fmt.Println("vf ", err)
	return &prevVote
}

func likeVoteF(id int, vote *models.Vote) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	_, err := transaction.Exec(likeVote, id, vote.Nickname)
	fmt.Println("lvf ", err)
	return err
}

func dislikeVoteF(id int, vote *models.Vote) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	_, err := transaction.Exec(dislikeVote, id, vote.Nickname)
	fmt.Println("dvf ", err)
	return err
}

func CreateVote(id int, vote *models.Vote) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	_, err := transaction.Exec(addVote, id, vote.Nickname, vote.Voice)
	return err
}

func VoteForThread(id int, vote *models.Vote) *models.Thread {
	transaction := database.StartTransaction()
	defer transaction.Commit()

	fmt.Println("0")
	prevVote := VoteFound(id, vote)
	var thread models.Thread
	var created time.Time
	fmt.Println(prevVote)
	if prevVote.Voice == 0 {
		if err := CreateVote(id, vote); err != nil {

			fmt.Println(err)

			return nil
		}
		fmt.Println("1")
		if vote.Voice > 0 {
			if err := transaction.QueryRow(like, id).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID, &thread.Votes); err != nil {
				fmt.Println("11", err)
				return nil
			} else {
				thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
				return &thread
			}
		} else {
			if err := transaction.QueryRow(dislike, id).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID, &thread.Votes); err != nil {
				fmt.Println("12", err)
				return nil
			} else {
				thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
				return &thread
			}
		}
	}
	fmt.Println("prev:", prevVote.Voice, "new:", vote.Voice)
	fmt.Println(prevVote.Thread)

	if prevVote.Voice > 0 && vote.Voice < 0 {
		err := dislikeVoteF(id, vote)
		if err != nil {
			return nil
		}
		if err := transaction.QueryRow(minus2, id).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID, &thread.Votes); err != nil {
			fmt.Println(err)
			return nil
		} else {
			thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
			return &thread
		}
	} else if prevVote.Voice < 0 && vote.Voice > 0 {
		err := likeVoteF(id, vote)
		if err != nil {
			return nil
		}
		if err := transaction.QueryRow(plus2, id).Scan(&thread.Author, &created, &thread.Forum, &thread.Message, &thread.Title, &thread.Slug, &thread.ID, &thread.Votes); err != nil {
			fmt.Println(err)
			return nil
		} else {
			thread.Created = created.Format("2006-01-02T15:04:05.000Z07:00")
			return &thread
		}
	} else {
		return GetThreadByID(strconv.Itoa(id))
	}

}
