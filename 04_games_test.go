package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-2 Search Games
func TestSearchGames(t *testing.T) {
	t.Run("Test Search Games", func(t *testing.T) {
		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		gameQuery := url.QueryEscape(game.Title)

		// release req
		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/games/search?query=%s&offset=0&limit=10", ReleaseURL, gameQuery), "api-2")
		assert.NoError(t, err)

		var releaseRespBody map[string]interface{}
		err = json.NewDecoder(releaseResp.Body).Decode(&releaseRespBody)
		assert.NoError(t, err)

		releaseTotal := releaseRespBody["meta"].(map[string]interface{})["total"].(float64)

		// dev req
		devResp, err := SendGetRequest(fmt.Sprintf("%s/games/search?query=%s&offset=0&limit=10", DevURL, gameQuery), "api-2")
		assert.NoError(t, err)

		var devRespBody map[string]interface{}
		err = json.NewDecoder(devResp.Body).Decode(&devRespBody)
		assert.NoError(t, err)

		devTotal := devRespBody["meta"].(map[string]interface{})["total"].(float64)

		// Compare the total values for mismatch
		assert.NotEqual(t, releaseTotal, devTotal, "The 'total' values from Release and Dev should not match due to the search issue in Dev")
	})
}

// api-9 Fet a Game
func TestGetGame(t *testing.T) {
	t.Run("Test Get Game by UUID", func(t *testing.T) {
		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		// release req
		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/games/%s", ReleaseURL, game.UUID), "api-9")
		assert.NoError(t, err)
		assert.Equal(t, 200, releaseResp.StatusCode, "Release server should return 200")

		// dev req
		devResp, err := SendGetRequest(fmt.Sprintf("%s/games/%s", DevURL, game.UUID), "api-9")
		assert.NoError(t, err)
		assert.Equal(t, 404, devResp.StatusCode, "Dev server should return 404")

		// Compare the total values for mismatch
		assert.NotEqual(t, releaseResp.StatusCode, devResp.StatusCode, "The response codes should not match as Dev server isn't fetching game")
	})
}
