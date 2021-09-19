package common

import (
	"math/rand"
	"strconv"
	"time"
)



func GenEmailCode() (s string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(90000000) + 10000000)

}
