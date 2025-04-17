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

func setupClient(t *testing.T) *rd_station.Client {
	token := os.Getenv("RD_STATION_TOKEN")
	if token == "" {
		t.Skip("RD_STATION_TOKEN environment variable not set, skipping test")
	}
	return rd_station.NewClient(rd_station.WithToken(token))
}

func TestListContactsFilterBasic(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := rd_station.ListContactsFilterRequest{
		Limit:     "5",
		Page:      "1",
		Order:     "name",
		Direction: "asc",
	}

	response, err := client.ListContactsFilter(ctx, filter)
	require.NoError(t, err, "Should not return an error")
	require.NotNil(t, response, "Response should not be nil")

	assert.LessOrEqual(t, len(response.Contacts), 5, "Should return at most 5 contacts")
	assert.GreaterOrEqual(t, response.Total, 0.0, "Total should be non-negative")

	t.Logf("Successfully retrieved %d contacts out of %f total", len(response.Contacts), response.Total)
}

func TestListContactsFilterWithSearch(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := rd_station.ListContactsFilterRequest{
		Q: "Artur",
	}

	response, err := client.ListContactsFilter(ctx, filter)
	require.NoError(t, err, "Should not return an error")
	require.NotNil(t, response, "Response should not be nil")
	t.Logf("Search found %d contacts with name containing 'Artur'", len(response.Contacts))
}

func TestListContactsFilterByEmail(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := rd_station.ListContactsFilterRequest{
		Email: "test@example.com",
		Limit: "5",
	}

	response, err := client.ListContactsFilter(ctx, filter)
	require.NoError(t, err, "Should not return an error")
	require.NotNil(t, response, "Response should not be nil")

	t.Logf("Email search found %d contacts", len(response.Contacts))
}

func TestListContactsFilterByPhone(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := rd_station.ListContactsFilterRequest{
		Phone: "123456789",
		Limit: "5",
	}

	response, err := client.ListContactsFilter(ctx, filter)
	require.NoError(t, err, "Should not return an error")
	require.NotNil(t, response, "Response should not be nil")

	t.Logf("Phone search found %d contacts", len(response.Contacts))
}

func TestCreateContact(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uniqueEmail := "test_" + time.Now().Format("20060102150405") + "@example.com"

	birthday := rd_station.BirthdayData{
		Day:   1,
		Month: 1,
		Year:  2000,
	}

	contact := rd_station.CreateContactRequest{
		Contact: rd_station.CreateContactData{
			Birthday: &birthday,
			Emails: &[]rd_station.EmailData{
				{
					Email: uniqueEmail,
				},
			},
			Name: "Teste Automatizado",
			LegalBases: &[]rd_station.LegalBasis{
				{
					Category: "communications",
					Type:     "consent",
					Status:   "granted",
				},
			},
		},
	}

	response, err := client.CreateContact(ctx, contact)
	require.NoError(t, err, "Should not return an error")
	require.NotNil(t, response, "Response should not be nil")

	assert.NotEmpty(t, response.ID, "ID should not be empty")
	assert.Equal(t, contact.Contact.Name, response.Name, "Name should match")

	require.NotEmpty(t, response.Emails, "Emails should not be empty")
	assert.Equal(t, uniqueEmail, response.Emails[0].Email, "Email should match")

	t.Logf("Successfully created contact with ID: %s and email: %s", response.ID, uniqueEmail)
}

func TestUpdateContact(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingContactID := "67fffaa734a1ef0027cec987"

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newName := "Contato Atualizado em " + timestamp
	newTitle := "Título Atualizado " + timestamp

	updateContact := rd_station.UpdateContactRequest{
		Contact: rd_station.UpdateContactData{
			Name:  newName,
			Title: &newTitle,
		},
	}

	updatedContact, err := client.UpdateContact(ctx, existingContactID, updateContact)
	require.NoError(t, err, "Should update the contact without error")
	require.NotNil(t, updatedContact, "Updated contact should not be nil")

	assert.Equal(t, newName, updatedContact.Name, "Name should be updated")
	assert.Equal(t, newTitle, updatedContact.Title, "Title should be updated")
	assert.Equal(t, existingContactID, updatedContact.ID, "ID should match the existing contact ID")

	require.NotEmpty(t, updatedContact.Emails, "Contact should still have emails")

	t.Logf("Successfully updated contact with ID: %s", updatedContact.ID)
	t.Logf("Updated name: %s", updatedContact.Name)
	t.Logf("Updated title: %s", updatedContact.Title)
}
