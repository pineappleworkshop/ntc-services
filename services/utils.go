package services

import (
	"math/rand"
	"ntc-services/models"
	"time"
)

func scramble(arr []*models.BlockRaw) {
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
