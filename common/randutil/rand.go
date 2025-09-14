package randutil

import (
	"math/rand"
	"time"
)

var randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
