package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-10 Get games by category
func TestGetGamesByCategory(t *testing.T) {
	t.Run("Test Get Games by Category", func(t *testing.T) {
		category, err := FetchCategory()
		assert.NoError(t, err)

		// release req
		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/categories/%s/games", ReleaseURL, category.UUID), "api-10")
		assert.NoError(t, err)
		assert.Equal(t, 200, releaseResp.StatusCode, "Release server should return 200")

		// dev req
		devResp, err := SendGetRequest(fmt.Sprintf("%s/categories/%s/games", DevURL, category.UUID), "api-10")
		assert.NoError(t, err)
		assert.Equal(t, 200, devResp.StatusCode, "Dev server should return 200")

		// Compare the release and dev server response bodies
		var releaseBody map[string]interface{}
		err = json.NewDecoder(releaseResp.Body).Decode(&releaseBody)
		assert.NoError(t, err)

		var devBody map[string]interface{}
		err = json.NewDecoder(devResp.Body).Decode(&devBody)
		assert.NoError(t, err)

		// Check if release server's game category matches the passed category UUID
		releaseGames, ok := releaseBody["games"].([]interface{})
		assert.True(t, ok, "Release response should contain 'games' field")

		for _, game := range releaseGames {
			gameObj := game.(map[string]interface{})
			categoryUUIDs := gameObj["category_uuids"].([]interface{})
			assert.Contains(t, categoryUUIDs, category.UUID, "Release server returned incorrect category UUID")
		}

		// checking if dev server's game category does NOT match the passed category UUID
		devGames, ok := devBody["games"].([]interface{})
		assert.True(t, ok, "Dev response should contain 'games' field")

		for _, game := range devGames {
			gameObj := game.(map[string]interface{})
			categoryUUIDs := gameObj["category_uuids"].([]interface{})
			assert.NotContains(t, categoryUUIDs, category.UUID, "Dev server returned incorrect category UUID")
		}
	})
}
