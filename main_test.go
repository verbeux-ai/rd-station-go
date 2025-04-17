package rd_station_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	rd_station "github.com/verbeux-ai/rd-station-go"
)

var client *rd_station.Client

func TestMain(m *testing.M) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Warning: could not load .env file. Ensure RD_STATION_TOKEN is set via environment.")
	}

	apiToken := os.Getenv("RD_STATION_TOKEN")
	if apiToken == "" {
		log.Fatal("Error: RD_STATION_TOKEN environment variable not set.")
	}

	client = rd_station.NewClient(
		rd_station.WithToken(apiToken),
	)

	exitCode := m.Run()
	os.Exit(exitCode)
}
