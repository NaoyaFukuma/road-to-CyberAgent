package entities

type (
	Rank int64

	RankInfo struct {
		UserID   UserID   `json:"userId"`
		UserName UserName `json:"userName"`
		Rank     Rank     `json:"rank"`
		Score    Score    `json:"score"`
	}

	RankingListResponse struct {
		RankInfoList []RankInfo `json:"ranks"`
	}
)
