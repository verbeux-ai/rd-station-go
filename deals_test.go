package rd_station_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rd_station "github.com/verbeux-ai/rd-station-go"
)

func TestListDealsFilterBasic(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := rd_station.ListDealsFilterRequest{
		Limit:     "5",
		Page:      "1",
		Order:     "name",
		Direction: "asc",
	}

	response, err := client.ListDealsFilter(ctx, filter)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.LessOrEqual(t, len(response.Deals), 5)
	assert.GreaterOrEqual(t, response.Total, 0)

	t.Logf("Successfully retrieved %d deals out of %d total", len(response.Deals), response.Total)
}

func TestCreateDeal(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deal := rd_station.CreateDealRequest{
		Deal: rd_station.CreateDealData{
			Name: "Automated Deal Test " + time.Now().Format("20060102150405"),
		},
	}

	response, err := client.CreateDeal(ctx, deal)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, deal.Deal.Name, response.Name)

	t.Logf("Successfully created deal with ID: %s", response.ID)
}

func TestUpdateDeal(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingDealID := os.Getenv("RD_TEST_DEAL_ID")
	if existingDealID == "" {
		t.Skip("RD_TEST_DEAL_ID environment variable not set, skipping test")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	updatedName := "Updated Deal" + timestamp
	updatedNote := "Automated Deal note" + timestamp

	updateReq := rd_station.UpdateDealRequest{
		Deal: rd_station.UpdateDealRequestData{
			Name:         &updatedName,
			DealLostNote: &updatedNote,
		},
	}

	updatedDeal, err := client.UpdateDeal(ctx, existingDealID, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updatedDeal)

	assert.Equal(t, updatedName, updatedDeal.Name)
	assert.Equal(t, existingDealID, updatedDeal.ID)

	t.Logf("Successfully updated deal with ID: %s", updatedDeal.ID)
	t.Logf("Updated name: %s", updatedDeal.Name)
	if updatedDeal.DealLostNote != nil {
		t.Logf("Updated deal lost note: %s", *updatedDeal.DealLostNote)
	}
}
