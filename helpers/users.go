package helpers

import (
	"db_forum/database"
	"fmt"

	"db_forum/models"
)

const insertUser = `
INSERT INTO users (nickname, email, fullname, about) 
VALUES ($1, $2, $3, $4)`

func CreateUser(newUser *models.User) error {
	transaction := database.StartTransaction()
	// defer transaction.Commit()

	if _, err := transaction.Exec(insertUser, newUser.Nickname, newUser.Email, newUser.Fullname, newUser.About); err != nil {
		transaction.Rollback()
		fmt.Println(newUser.Nickname, "<-user")
		fmt.Println("newUser", err)
		return err
	}
	transaction.Commit()
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
	// defer transaction.Rollback()
	fmt.Println(nickname)
	fmt.Println(email)

	if elements, err := database.Connection.Query(selectByNicknameOrEmail, nickname, email); err != nil {
		fmt.Println("selectByNicknameOrEmail", err)
		fmt.Println("length", elements.Next())
		transaction.Rollback()
		elements.Close()

		return users
	} else {
		defer elements.Close()

		for elements.Next() {
			var user models.User
			if err := elements.Scan(
				&user.Nickname,
				&user.Email,
				&user.Fullname,
				&user.About); err != nil {
				fmt.Println("selectByNicknameOrEmail2", err)
				transaction.Rollback()
				return users
			}
			users = append(users, &user)

		}
	}

	transaction.Commit()
	return users
}

const selectByNickname = `
SELECT nickname, email, fullname, about 
FROM users 
WHERE nickname = $1`

func FindByNickname(nickname string) *models.User {

	transaction := database.StartTransaction()
	var user models.User

	if elements, err := database.Connection.Query(selectByNickname, nickname); err != nil {
		fmt.Println("selectByNickname", err)
		fmt.Println("length", elements.Next())
		transaction.Rollback()
		elements.Close()
		transaction := database.StartTransaction()
		var user models.User
		if err := database.Connection.QueryRow(selectByNickname, nickname).Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About); err != nil {
			transaction.Rollback()
			fmt.Println("FUCK IT", err)
			return nil
		} else {
			transaction.Commit()

			return &user
		}
		return &user
	} else {
		defer elements.Close()

		for elements.Next() {
			if err := elements.Scan(
				&user.Nickname,
				&user.Email,
				&user.Fullname,
				&user.About); err != nil {
				fmt.Println("selectByNicknameOrEmail2", err)

				return &user
			}
			transaction.Commit()

			return &user

		}
	}
	return nil
}

const selectByEmail = `
SELECT nickname, email, fullname, about 
FROM users 
WHERE email = $1`

func FindByEmail(email string) *models.User {
	fmt.Println("FINDBYEMAIL")

	transaction := database.StartTransaction()
	defer transaction.Commit()
	var user models.User
	if err := database.Connection.QueryRow(selectByEmail, email).Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About); err != nil {
		fmt.Println("fbye", err)
		return nil
	} else {
		return &user
	}
}

//  ИСПРАВИТЬ
const updateUsers = `
UPDATE users
SET
	email = $2, 
	fullname = $3, 
	about = $4
WHERE
	nickname = $1
`

func UpdateProfile(user *models.User) error {
	transaction := database.StartTransaction()
	if _, err := transaction.Exec(updateUsers, user.Nickname, user.Email, user.Fullname, user.About); err != nil {
		fmt.Println("update:", err)
		transaction.Rollback()
		return err
	} else {
		transaction.Commit()
		return nil
	}
}
