package log

import (
	"context"
	"math/rand"
	"time"
)

type id string

const keyId = id("request_ID")

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GetRequestId(ctx context.Context) string {
	s, _ := ctx.Value(keyId).(string)
	return s
}

func generateRequestID() string {
	requestID := make([]rune, 32)
	for i := range requestID {
		requestID[i] = letters[random.Intn(len(letters))]
	}
	return string(requestID)
}
