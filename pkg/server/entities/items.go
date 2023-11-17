package entities

const (
	N  Rarity = 1
	R  Rarity = 2
	SR Rarity = 3
)

type (
	ItemID   int64
	ItemName string
	Rarity   int64 // 1: N, 2: R, 3: SR
	HasItem  bool

	Item struct {
		ID     ItemID   `json:"collectionID"`
		Name   ItemName `json:"name"`
		Rarity Rarity   `json:"rarity"` // 1: N, 2: R, 3: SR
	}

	ItemWithWeight struct {
		Item
		Weight Weight
	}

	Items []Item
)
