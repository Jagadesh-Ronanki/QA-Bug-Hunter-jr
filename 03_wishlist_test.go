package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-5 Add an item to user's wishlist
func TestAddItemToWishlist(t *testing.T) {
	t.Run("Test Add Item to Wishlist", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(1)
		assert.NoError(t, err)

		userUUID := user.UUID
		gameUUID := game.UUID
		game2UUID := game2.UUID

		releaseBody := WishlistBody{
			ItemUUID: gameUUID,
		}

		devBody := WishlistBody{
			ItemUUID: game2UUID,
		}

		// release req
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", ReleaseURL, userUUID), releaseBody, "api-5")
		assert.NoError(t, err)
		assert.Equal(t, 200, releaseResp.StatusCode, "Release server should return 200")

		// dev req
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", DevURL, userUUID), devBody, "api-5")
		assert.NoError(t, err)
		assert.Equal(t, 422, devResp.StatusCode, "Dev server should return 422 when wishlist limit is reached")

		// release req whishlist items
		releaseWishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", ReleaseURL, userUUID), "api-5")
		assert.NoError(t, err)

		var releaseWishlist map[string]interface{}
		err = json.NewDecoder(releaseWishlistResp.Body).Decode(&releaseWishlist)
		assert.NoError(t, err)

		// get items from wishlist
		items := releaseWishlist["items"].([]interface{})
		// Check if the gameUUID is in the wishlist items
		var itemUUIDs []string
		for _, item := range items {
			itemUUIDs = append(itemUUIDs, item.(map[string]interface{})["uuid"].(string))
		}
		assert.Contains(t, itemUUIDs, gameUUID, "Wishlist should contain the added item on release")

		// dev req wishlist items
		devWishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", DevURL, userUUID), "api-5")
		assert.NoError(t, err)

		var devWishlist map[string]interface{}
		err = json.NewDecoder(devWishlistResp.Body).Decode(&devWishlist)
		assert.NoError(t, err)

		// get items
		devItems := devWishlist["items"].([]interface{})
		var devItemUUIDs []string
		for _, item := range devItems {
			devItemUUIDs = append(devItemUUIDs, item.(map[string]interface{})["uuid"].(string))
		}
		assert.NotContains(t, devItemUUIDs, game2UUID, "Wishlist should not contain the item on dev due to error")

		// the number of items in Release server's wishlist doesn't exceed 10 still dev returns 422 response code
		assert.LessOrEqual(t, len(items), 10, "Wishlist should not contain more than 10 items")
	})
}

// api-25 add an item to user's wishlist
func TestAddItemToWishlistAPI25(t *testing.T) {
	t.Run("Test Add Item to Wishlist (API-25)", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(1)
		assert.NoError(t, err)

		userUUID := user.UUID
		gameUUID := game.UUID
		game2UUID := game2.UUID

		releaseBody := WishlistBody{
			ItemUUID: gameUUID,
		}

		devBody := WishlistBody{
			ItemUUID: game2UUID,
		}

		wishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", ReleaseURL, userUUID), "api-25")
		assert.NoError(t, err)

		var wishlist map[string]interface{}
		err = json.NewDecoder(wishlistResp.Body).Decode(&wishlist)
		assert.NoError(t, err)

		items := wishlist["items"].([]interface{})
		for _, item := range items {
			itemUUID := item.(map[string]interface{})["uuid"].(string)
			removeBody := WishlistBody{
				ItemUUID: itemUUID,
			}
			_, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/remove", ReleaseURL, userUUID), removeBody, "api-25")
			assert.NoError(t, err)
		}

		// release req - Add game to wishlist
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", ReleaseURL, userUUID), releaseBody, "api-25")
		assert.NoError(t, err)
		assert.Equal(t, 200, releaseResp.StatusCode, "Release server should return 200")

		// dev req - Add game to wishlist
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", DevURL, userUUID), devBody, "api-25")
		assert.NoError(t, err)
		assert.Equal(t, 200, devResp.StatusCode, "Dev server should return 200")

		// Verify wishlist items in Release server
		releaseWishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", ReleaseURL, userUUID), "api-25")
		assert.NoError(t, err)

		var releaseWishlist map[string]interface{}
		err = json.NewDecoder(releaseWishlistResp.Body).Decode(&releaseWishlist)
		assert.NoError(t, err)

		items = releaseWishlist["items"].([]interface{})
		var releaseItemUUIDs []string
		for _, item := range items {
			releaseItemUUIDs = append(releaseItemUUIDs, item.(map[string]interface{})["uuid"].(string))
		}

		// Verify the item was added to the wishlist
		assert.Contains(t, releaseItemUUIDs, gameUUID, "Wishlist should contain the added item on release")

		// Verify the item on dev request was not added to the wishlist
		assert.NotContains(t, releaseItemUUIDs, game2UUID, "Wishlist should not contain the item on dev due to temporary addition")

		// Verify wishlist items in Dev server
		devWishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", DevURL, userUUID), "api-25")
		assert.NoError(t, err)

		var devWishlist map[string]interface{}
		err = json.NewDecoder(devWishlistResp.Body).Decode(&devWishlist)
		assert.NoError(t, err)

		devItems := devWishlist["items"].([]interface{})
		var devItemUUIDs []string
		for _, item := range devItems {
			devItemUUIDs = append(devItemUUIDs, item.(map[string]interface{})["uuid"].(string))
		}

		// Verify that the item is NOT in the dev wishlist (since it is not actually saved)
		assert.NotContains(t, devItemUUIDs, game2UUID, "Wishlist should not contain the item on dev due to temporary addition")
	})
}

// api-8 Add an item to user's wishlist
func TestRemoveItemFromWishlistAPI8(t *testing.T) {
	t.Run("Test Remove Item from Wishlist (API-8)", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game1, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(1)
		assert.NoError(t, err)

		userUUID := user.UUID
		game1UUID := game1.UUID
		game2UUID := game2.UUID

		releaseBody := WishlistBody{
			ItemUUID: game1UUID,
		}

		devBody := WishlistBody{
			ItemUUID: game2UUID,
		}

		// Add two games to the wishlist on release and Dev servers
		_, err = SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", ReleaseURL, userUUID), releaseBody, "api-8")
		assert.NoError(t, err)
		_, err = SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/add", ReleaseURL, userUUID), devBody, "api-8")
		assert.NoError(t, err)

		// Try removing the first game on Release server
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/remove", ReleaseURL, userUUID), releaseBody, "api-8")
		assert.NoError(t, err)
		assert.Equal(t, 200, releaseResp.StatusCode, "Release server should return 200 when removing an existing item")

		// Try removing the first game again on Release server -  should return 404
		releaseResp2, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/remove", ReleaseURL, userUUID), releaseBody, "api-8")
		assert.NoError(t, err)
		assert.Equal(t, 404, releaseResp2.StatusCode, "Release server should return 404 after item has been removed")

		// Try removing the second game on dev server
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/remove", DevURL, userUUID), devBody, "api-8")
		assert.NoError(t, err)
		assert.Equal(t, 200, devResp.StatusCode, "Dev server should return 200 when attempting to remove an item")

		// Try removing the second game again on Dev server - returning 200 - indicates not removed in previous call
		devResp2, err := SendPostRequest(fmt.Sprintf("%s/users/%s/wishlist/remove", DevURL, userUUID), devBody, "api-8")
		assert.NoError(t, err)
		assert.Equal(t, 200, devResp2.StatusCode, "Dev server should still return 200 after an item is not actually removed")

		// Verify wishlist items on server after removal
		releaseWishlistResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/wishlist", ReleaseURL, userUUID), "api-8")
		assert.NoError(t, err)

		var releaseWishlist map[string]interface{}
		err = json.NewDecoder(releaseWishlistResp.Body).Decode(&releaseWishlist)
		assert.NoError(t, err)

		items := releaseWishlist["items"].([]interface{})
		var releaseItemUUIDs []string
		for _, item := range items {
			releaseItemUUIDs = append(releaseItemUUIDs, item.(map[string]interface{})["uuid"].(string))
		}

		// Assert that the first game was removed (not in wishlist)
		assert.NotContains(t, releaseItemUUIDs, game1UUID, "Wishlist should not contain the first game after removal on release")

		// Assert that the second game was not removed
		assert.Contains(t, releaseItemUUIDs, game2UUID, "Wishlist should still contain the second game")
	})
}
