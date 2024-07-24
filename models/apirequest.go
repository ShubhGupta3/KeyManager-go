package models

type Key struct {
	ID         int    `json:"_id"`
	Key        string `json:"key"`
	CreationTS int64  `json:"cts"`
	DeathTS    int64  `json:"dts"`
	IsBlocked  bool   `json:"blkd"`
	BlockTs    int64  `json:"bts"`
	IsRemoved  bool   `json:"rmvd"`
}

type GenerateKeyReq struct {
	NumberOfKeys int `json:"count"`
}
