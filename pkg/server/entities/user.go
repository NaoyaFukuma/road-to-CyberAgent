package entities

type (
	UserID    int64
	UserName  string
	Coin      int64
	AuthToken string

	User struct {
		ID        UserID    `json:"id"`
		Name      UserName  `json:"name"`
		HighScore Score     `json:"highScore"`
		Coin      Coin      `json:"coin"`
		AuthToken AuthToken `json:"-"`
		// and more...
	}

	Users []User
)
