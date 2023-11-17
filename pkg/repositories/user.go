package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

type UserRepository interface {
	GetUsers() ([]*entities.User, error)
	GetUserByID(ID entities.UserID) (*entities.User, error)
	GetUserIDAuthToken(token entities.AuthToken) (entities.UserID, error)
	CreateUser(user *entities.User) error
	UpdateUserNameByID(ID entities.UserID, name entities.UserName) error
	UpdateUserCoinsByID(ID entities.UserID, coin entities.Coin) error
	UpdateUserCoinsByIDTransaction(tx *sql.Tx, ID entities.UserID, coin entities.Coin) error
	UpdateUserHighScoreByID(ID entities.UserID, score entities.Score) error
	DeleteUserByID(ID entities.UserID) error
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) GetUsers() ([]*entities.User, error) {
	query := "SELECT * FROM user"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.AuthToken); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ID entities.UserID) (*entities.User, error) {
	query := "SELECT * FROM user WHERE id = ? LIMIT 1"
	row := r.db.QueryRow(query, ID)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var user entities.User
	if err := row.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.AuthToken); err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserIDAuthToken(token entities.AuthToken) (entities.UserID, error) {
	query := "SELECT id FROM user WHERE auth_token = ? LIMIT 1"
	row := r.db.QueryRow(query, token)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var ID entities.UserID
	if err := row.Scan(&ID); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found with provided auth token")
			log.Println(err)
			return 0, err
		}
		log.Println(err)
		return 0, err
	}

	return ID, nil
}

func (r *userRepository) CreateUser(user *entities.User) error {
	query := "INSERT INTO user (name, auth_token) VALUES (?, ?)"
	_, err := r.db.Exec(query, user.Name, user.AuthToken)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserNameByID(ID entities.UserID, name entities.UserName) error {
	query := "UPDATE user SET name = ? WHERE id = ?"
	_, err := execQueryAndReturnAffectedRows(r.db, query, name, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserCoinsByID(ID entities.UserID, coin entities.Coin) error {
	query := "UPDATE user SET coin = ? WHERE id = ?"
	_, err := execQueryAndReturnAffectedRows(r.db, query, coin, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserCoinsByIDTransaction(tx *sql.Tx, ID entities.UserID, coin entities.Coin) error {
	query := "UPDATE user SET coin = ? WHERE id = ?"
	_, err := execQueryAndReturnAffectedRows(tx, query, coin, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserHighScoreByID(ID entities.UserID, highScore entities.Score) error {
	query := "UPDATE user SET high_score = ? WHERE id = ?"
	_, err := execQueryAndReturnAffectedRows(r.db, query, highScore, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *userRepository) DeleteUserByID(ID entities.UserID) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := execQueryAndReturnAffectedRows(r.db, query, ID)
	if err != nil {
		log.Println(err)
		return err
	}
	if result == 0 {
		return fmt.Errorf("user id %d not found", ID)
	}

	return nil
}
