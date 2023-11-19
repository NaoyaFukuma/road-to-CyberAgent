package repositories

import (
	"database/sql"
	"log"
	"time"

	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

type UserScoresRepository interface {
	AddUserScore(userID entities.UserID, score entities.Score) error
	AddUserScoreTransaction(tx *sql.Tx, userID entities.UserID, score entities.Score) error
	GetUserScoreWithUserName(offset int64, limit entities.RankingListLimit) (*entities.UserScoresJoinedUserName, error)
}

func NewUserScoresRepository(db *sql.DB) UserScoresRepository {
	return &userScoresRepository{db}
}

type userScoresRepository struct {
	db *sql.DB
}

func (r *userScoresRepository) AddUserScore(userID entities.UserID, score entities.Score) error {
	query := "INSERT INTO user_scores (user_id, score) VALUES (?, ?)"
	_, err := r.db.Exec(query, userID, score)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *userScoresRepository) AddUserScoreTransaction(tx *sql.Tx, userID entities.UserID, score entities.Score) error {
	query := "INSERT INTO user_scores (user_id, score) VALUES (?, ?)"
	_, err := tx.Exec(query, userID, score)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *userScoresRepository) GetUserScoreWithUserName(offset int64, limit entities.RankingListLimit) (*entities.UserScoresJoinedUserName, error) {
	query := `
		SELECT user_scores.user_id, user.name, user_scores.score, user_scores.created_at
		FROM user_scores
		JOIN user ON user_scores.user_id = user.id
		ORDER BY user_scores.score DESC, user_scores.user_id ASC, user_scores.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var userScoresJoinedUserName entities.UserScoresJoinedUserName
	for rows.Next() {
		var data entities.UserScoreJoinedUserName
		var createdAt []byte
		if err := rows.Scan(&data.UserID, &data.UserName, &data.Score, &createdAt); err != nil {
			log.Println(err)
			return nil, err
		}
		if data.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt)); err != nil {
			log.Println(err)
			return nil, err
		}
		userScoresJoinedUserName = append(userScoresJoinedUserName, data)
	}
	return &userScoresJoinedUserName, nil
}
