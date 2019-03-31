package helpers

import (
	"github.com/iamfrommoscow/db_forum/database"
	"github.com/iamfrommoscow/db_forum/models"
)

const insertUser = `
INSERT INTO users (nickname, email, fullname, about) 
VALUES ($1, $2, $3, $4)`

func CreateUser(newUser *models.User) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()

	if _, err := transaction.Exec(insertUser, newUser.Nickname, newUser.Email, newUser.Fullname, newUser.About); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

const selectByNicknameOrEmail = `
SELECT nickname, email, fullname, about 
FROM users 
WHERE nickname = $1 or email = $2
`

func FindByNicknameOrEmail(nickname string, email string) []*models.User {
	var users []*models.User
	transaction := database.StartTransaction()
	defer transaction.Commit()

	if elements, err := transaction.Query(selectByNicknameOrEmail, nickname, email); err != nil {
		return users
	} else {
		for elements.Next() {
			var user models.User
			if err := elements.Scan(
				&user.Nickname,
				&user.Email,
				&user.Fullname,
				&user.About); err != nil {
				return users
			}
			users = append(users, &user)

		}
	}

	return users
}

const selectByNickname = `
SELECT nickname, email, fullname, about 
FROM users 
WHERE nickname = $1`

func FindByNickname(nickname string) *models.User {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var user models.User
	if err := transaction.QueryRow(selectByNickname, nickname).Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About); err != nil {

		return nil
	} else {
		return &user
	}
}

const selectByEmail = `
SELECT nickname, email, fullname, about 
FROM users 
WHERE email = $1`

func FindByEmail(email string) *models.User {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	var user models.User
	if err := transaction.QueryRow(selectByEmail, email).Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About); err != nil {

		return nil
	} else {
		return &user
	}
}

//  ИСПРАВИТЬ
const updateUsers = `
UPDATE users
SET
	email = coalesce(coalesce(nullif($2, ''), email)), 
	fullname = coalesce(coalesce(nullif($3, ''), fullname)), 
	about = coalesce(coalesce(nullif($4, ''), about))
WHERE
	nickname = $1
RETURNING
	email,
	fullname,
	about
`

func UpdateProfile(user *models.User) error {
	transaction := database.StartTransaction()
	defer transaction.Commit()
	if err := transaction.QueryRow(updateUsers, user.Nickname, user.Email, user.Fullname, user.About).Scan(&user.Email, &user.Fullname, &user.About); err != nil {
		// fmt.Println(err)

		return err
	} else {

		return nil
	}
}
