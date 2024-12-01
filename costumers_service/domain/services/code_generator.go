package services

import (
	"fmt"
	"math/rand"
	"time"
)

type CodeGenerator struct{}

func (gen CodeGenerator) GeneratePhoneCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (gen CodeGenerator) GenerateEmailCode() string {
	length := 16
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomString)
}
