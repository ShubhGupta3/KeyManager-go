package helpers

import (
	"crypto/rand"
	"math/big"
)

const KEY_LENGTH = 16

// func GenerateRandomKey() (models.Key, error) {
// 	tempBytes := make([]byte, 32)
// 	_, err := rand.Read(tempBytes)

// 	if err != nil {
// 		return models.Key{}, err
// 	}

// 	var finalKey models.Key
// 	keyStr := base32.StdEncoding.EncodeToString(tempBytes)[:KEY_LENGTH]
// 	finalKey.Key = keyStr
// 	return finalKey, nil
// }

func GenerateRandNo(rangeNo int) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	return n
}
