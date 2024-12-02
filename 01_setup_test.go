package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sendSetupRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/setup", url), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// /setup endpoint
func TestSetupRelease(t *testing.T) {
	t.Run("Test Setup on Release Environment", func(t *testing.T) {
		// Send the setup request
		resp, err := sendSetupRequest(ReleaseURL)
		assert.NoError(t, err)

		// Check that the response status code is 205 (ResetContent)
		assert.Equal(t, http.StatusResetContent, resp.StatusCode)
	})
}

// /users endpoint
func TestGetAllUsers(t *testing.T) {
	t.Run("Test Get All Users", func(t *testing.T) {
		var taskID = "api-6"
		resp, err := SendGetRequest(fmt.Sprintf("%s/users", ReleaseURL), taskID)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]interface{}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		meta, ok := responseBody["meta"].(map[string]interface{})
		assert.True(t, ok, "meta field is missing or incorrect format")

		total, ok := meta["total"].(float64)
		assert.True(t, ok, "total field is missing or incorrect format")

		// If fetched successfully initial users returned are 11
		assert.Equal(t, float64(11), total)
	})
}
