package entities

import "time"

type (
	GameSettingID        int64
	GachaCoinConsumption int64
	RankingListLimit     int64
	Weight               int64
	MaxGachaTimes        int64

	GameSettings struct {
		ID                   GameSettingID        `json:"id"`
		GachaCoinConsumption GachaCoinConsumption `json:"gachaCoinConsumption"`
		RankingListLimit     RankingListLimit     `json:"rankingListLimit"`
		NWeight              Weight               `json:"nWeight"`
		RWeight              Weight               `json:"rWeight"`
		SrWeight             Weight               `json:"srWeight"`
		MaxGachaTimes        MaxGachaTimes        `json:"maxGachaTimes"`
		CreatedAt            time.Time            `json:"createdAt"`
		IsActive             bool                 `json:"isActive"`
		// and more...
	}
)
