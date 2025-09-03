package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/rauan06/realtime-map/seeder/domain"
)

func genCord() float64 {
	return rand.Float64() * 100
}

func genOBUData() domain.OBUData {
	return domain.OBUData{
		ID:   uuid.New(),
		Long: genCord(),
		Lat:  genCord(),
	}
}

func main() {
	for {
		time.Sleep(time.Duration(genCord()))

		obuData := genOBUData()
		fmt.Printf("%+v\n", obuData)
	}
}
