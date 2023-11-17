package entities

import "time"

type (
	Score int64

	UserScore struct {
		UserID    UserID
		Score     Score
		CreatedAt time.Time
	}

	UserScoreJoinedUserName struct {
		UserID    UserID
		UserName  UserName
		Score     Score
		CreatedAt time.Time
	}

	UserScores               []UserScore
	UserScoresJoinedUserName []UserScoreJoinedUserName
)
