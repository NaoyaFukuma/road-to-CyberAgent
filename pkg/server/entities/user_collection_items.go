package entities

type (
	CollectionItem struct {
		ID      ItemID   `json:"collectionID"`
		Name    ItemName `json:"name"`
		Rarity  Rarity   `json:"rarity"` // 1: N, 2: R, 3: SR
		HasItem HasItem  `json:"hasItem"`
	}

	CollectionItemList struct {
		Items []CollectionItem `json:"collections"`
	}

	GachaResult struct {
		ID     ItemID   `json:"collectionID"`
		Name   ItemName `json:"name"`
		Rarity Rarity   `json:"rarity"` // 1: N, 2: R, 3: SR
		IsNew  bool     `json:"isNew"`
	}

	GachaResultList struct {
		Items []GachaResult `json:"results"`
	}
)
