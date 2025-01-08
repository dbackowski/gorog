package main

import (
	"crypto/rand"
	"math/big"
)

func GetRandomBetween(low int, high int) int {
	return GetDiceRoll(high-low) + high
}

func GetRandomInt(num int) int {
	x, _ := rand.Int(rand.Reader, big.NewInt(int64(num)))
	return int(x.Int64())

}

func GetDiceRoll(num int) int {
	x, _ := rand.Int(rand.Reader, big.NewInt(int64(num)))
	return int(x.Int64()) + 1
}
