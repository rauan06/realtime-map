package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	client "github.com/rauan06/realtime-map/seeder/controller"
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
	client := client.NewHTTPClient(nil, "")

	for {
		time.Sleep(time.Duration(genCord()))

		obuData := genOBUData()
		encodedData, err  := json.Marshal(obuData)
		if err != nil {
			log.Fatal(err)
		}

		client.NewRequest("POST", "/v1/locate", encodedData)
		fmt.Printf("%+v\n", obuData)
	}
}
