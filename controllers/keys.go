package controllers

import (
	"crypto/rand"
	"edra/helpers"
	"edra/models"
	"encoding/base32"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const KEY_LENGTH = 16

var ActiveKeys int64

var InactiveKeys int64

var KeyStore map[int64]models.Key

// Function to generate keys
func GenerateKeys(c *gin.Context) {
	// Read request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, "Error in loading request")
		return
	}

	// Unmarshal the request
	var req models.GenerateKeyReq

	err = json.Unmarshal(body, &req)

	if err != nil {
		c.JSON(400, "Error in unmarshalling request")
		return
	}

	var keyArr []models.Key
	count := req.NumberOfKeys

	for i := 0; i < count; i++ {
		key, err := GenerateRandomKey()
		if err != nil {
			c.JSON(400, "Error in generating key")
		}
		keyArr = append(keyArr, key)
	}
	c.JSON(200, keyArr)
}

// Function to retrieve a key
func RetrieveKey(c *gin.Context) {
	l := len(KeyStore)
	if l == 0 {
		c.JSON(404, "No keys available")
		return
	}
	var key models.Key
	for {
		randInt := helpers.GenerateRandNo(l)
		key = KeyStore[randInt]
		if !key.IsBlocked && !key.IsRemoved {
			key.IsBlocked = true
			key.BlockTs = time.Now().Unix()
			KeyStore[randInt] = key
			break
		}
	}
	c.JSON(200, key)
}

// Function to retrieve a key by ID
func RetrieveKeyByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, "Invalid ID")
	}
	finalKey := KeyStore[int64(idInt)]
	if finalKey.ID == 0 {
		c.JSON(400, "Key not present")
		return
	}

	c.JSON(200, finalKey)
}

// Function to delete a key from system
func DeleteKeyByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, "Invalid ID")
	}
	finalKey := KeyStore[int64(idInt)]
	if finalKey.ID == 0 {
		c.JSON(400, "Key not present")
		return
	}
	delete(KeyStore, int64(idInt))
	c.JSON(200, "Succesfully deleted!")
}

// Function to unblock a key for further use
func UnblockKeyByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, "Invalid ID")
	}
	finalKey := KeyStore[int64(idInt)]
	if !finalKey.IsBlocked {
		c.JSON(404, "Key is not blocked")
	}
	finalKey.IsBlocked = false
	finalKey.BlockTs = 0
	KeyStore[int64(idInt)] = finalKey
	c.JSON(200, "Succesfully unblocked!")
}

// Function to keep a key alive for T+5
func KeepKeyAliveByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, "Invalid ID")
	}
	finalKey := KeyStore[int64(idInt)]
	if finalKey.IsBlocked {
		c.JSON(400, "Key is blocked, cannot keep alive")
		return
	}
	finalKey.DeathTS = time.Now().Unix() + 300
	KeyStore[int64(idInt)] = finalKey
	c.JSON(200, finalKey)
}

// func DeleteKeysCRON(c *gin.Context) {
// 	go BlacklistKeys()
// }

func BlacklistKeys() {
	interval := 2 * time.Second

	ts := time.Now().Unix() - 300
	for {
		log.Println("Running delete CRON")
		for i, val := range KeyStore {
			if val.DeathTS <= ts {
				val.IsRemoved = true
				KeyStore[i] = val
			}
		}
		time.Sleep(interval)
	}
}

// func UnblockKeysCRON(c *gin.Context) {
// 	go unblockKeys()
// }

func UnblockKeys() {
	interval := 2 * time.Second

	ts := time.Now().Unix() - 60
	for {
		log.Println("Running unblock CRON")
		for i, val := range KeyStore {
			if val.IsBlocked {
				if val.BlockTs <= ts {
					val.BlockTs = 0
					val.IsBlocked = false
					KeyStore[i] = val
				}
			}
		}
		time.Sleep(interval)
	}

}

func GenerateRandomKey() (models.Key, error) {
	// if KeyStore == nil {
	// 	KeyStore = make(map[int64]models.Key, 1)
	// }
	ActiveKeys++
	tempBytes := make([]byte, 32)
	_, err := rand.Read(tempBytes)

	if err != nil {
		return models.Key{}, err
	}

	var finalKey models.Key
	keyStr := base32.StdEncoding.EncodeToString(tempBytes)[:KEY_LENGTH]
	tsINT := time.Now().Unix()

	finalKey.ID = len(KeyStore) + 1
	finalKey.Key = keyStr
	finalKey.CreationTS = tsINT
	finalKey.DeathTS = tsINT + 300
	finalKey.IsBlocked = false
	finalKey.BlockTs = 0
	finalKey.IsRemoved = false
	KeyStore[int64(finalKey.ID)] = finalKey
	return finalKey, nil
}
