package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-12 Get a Cart
func TestCartTotalPriceDifference(t *testing.T) {
	t.Run("Compare Cart Total Price Between Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game1, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(1)
		assert.NoError(t, err)

		item1UUID := game1.UUID
		item2UUID := game2.UUID

		_, err = AddItemToCart(user.UUID, item1UUID, 2, ReleaseURL, "api-12")
		assert.NoError(t, err)

		_, err = AddItemToCart(user.UUID, item2UUID, 1, ReleaseURL, "api-12")
		assert.NoError(t, err)

		// release req
		releaseCartResp, err := GetUserCart(user.UUID, ReleaseURL, "api-12")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseCartResp.StatusCode)

		var releaseCart CartResponse
		err = json.NewDecoder(releaseCartResp.Body).Decode(&releaseCart)
		assert.NoError(t, err)

		releaseTotalPrice := releaseCart.TotalPrice

		// dev req
		devCartResp, err := GetUserCart(user.UUID, DevURL, "api-12")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devCartResp.StatusCode)

		var devCart CartResponse
		err = json.NewDecoder(devCartResp.Body).Decode(&devCart)
		assert.NoError(t, err)

		devTotalPrice := devCart.TotalPrice

		// compare that the total prices are not equal between Release and dev
		assert.NotEqual(t, releaseTotalPrice, devTotalPrice, "Total price between Release and Dev should not be equal")
	})
}

// api-13 Change an item in user's cart
func TestChangeItemQuantity(t *testing.T) {
	t.Run("Change Item Quantity and Compare Release vs Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		itemUUID := game.UUID
		_, err = AddItemToCart(user.UUID, itemUUID, 1, ReleaseURL, "api-13")
		assert.NoError(t, err)

		// release req
		releaseCartResp, err := GetUserCart(user.UUID, ReleaseURL, "api-13")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseCartResp.StatusCode)

		var releaseCart CartResponse
		err = json.NewDecoder(releaseCartResp.Body).Decode(&releaseCart)
		assert.NoError(t, err)

		initialReleaseTotalPrice := releaseCart.TotalPrice

		response, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/change", ReleaseURL, user.UUID), ChangeItemQuantityRequest{ItemUUID: itemUUID, Quantity: 2}, "api-13")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode, "Expected a 200 OK response from Release server")

		// again release req
		releaseUpdatedCartResp, err := GetUserCart(user.UUID, ReleaseURL, "api-13")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseUpdatedCartResp.StatusCode)

		var releaseUpdatedCart CartResponse
		err = json.NewDecoder(releaseUpdatedCartResp.Body).Decode(&releaseUpdatedCart)
		assert.NoError(t, err)

		// The total price should increase after updating quantity to 2
		assert.Greater(t, releaseUpdatedCart.TotalPrice, initialReleaseTotalPrice, "Total price should increase after changing quantity")

		var itemFound bool
		for _, item := range releaseUpdatedCart.Items {
			if item.ItemUUID == itemUUID {
				itemFound = true
				break
			}
		}
		assert.True(t, itemFound, "Item with the specified UUID should be present in the updated cart")

		// Change the item quantity to 1 on Dev
		devResponse, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/change", DevURL, user.UUID), ChangeItemQuantityRequest{ItemUUID: itemUUID, Quantity: 1}, "api-13")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResponse.StatusCode, "Expected a 200 OK response from Release server")

		var devCart CartResponse
		err = json.NewDecoder(devResponse.Body).Decode(&devCart)
		assert.NoError(t, err)

		fmt.Print(devCart)

		// The total price should decrease after changing quantity to 1
		assert.Less(t, devCart.TotalPrice, releaseUpdatedCart.TotalPrice, "Total price should decrease after changing quantity to 1")

		// Assert if the items list is empty in Dev due to the bug
		assert.Empty(t, devCart.Items, "The items list should be empty in Dev due to the bug")

		// Fetch the cart again from Release to confirm the issue persists
		cartAgainResp, err := GetUserCart(user.UUID, ReleaseURL, "api-13")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, cartAgainResp.StatusCode)

		var devCartAgain CartResponse
		err = json.NewDecoder(cartAgainResp.Body).Decode(&devCartAgain)
		assert.NoError(t, err)

		// Confirm the total price matches the previous Dev response
		assert.Equal(t, devCart.TotalPrice, devCartAgain.TotalPrice, "The total price should remain the same when fetching cart again in Dev")
	})
}

// api-14 Remove an item from user's cart
func TestRemoveItemFromCart(t *testing.T) {
	t.Run("Remove item from cart and compare Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game1, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game3, err := FetchExistingGame(0)
		assert.NoError(t, err)

		// Add 3 items to the cart for the user on Release
		_, err = AddItemToCart(user.UUID, game1.UUID, 1, ReleaseURL, "api-14") // Add item 1
		assert.NoError(t, err)

		_, err = AddItemToCart(user.UUID, game2.UUID, 1, ReleaseURL, "api-14") // Add item 2
		assert.NoError(t, err)

		_, err = AddItemToCart(user.UUID, game3.UUID, 1, ReleaseURL, "api-14") // Add item 3
		assert.NoError(t, err)

		// Fetch the cart from Release to get the initial total price and items
		releaseCartResp, err := GetUserCart(user.UUID, ReleaseURL, "api-14")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseCartResp.StatusCode)

		var releaseCart CartResponse
		err = json.NewDecoder(releaseCartResp.Body).Decode(&releaseCart)
		assert.NoError(t, err)

		initialReleaseTotalPrice := releaseCart.TotalPrice
		initialReleaseItemsCount := len(releaseCart.Items)

		// remove an item from the Release cart
		removeItemData := RemoveItemRequest{ItemUUID: game2.UUID}
		releaseRemoveResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/remove", ReleaseURL, user.UUID), removeItemData, "api-14")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseRemoveResp.StatusCode)

		// Check release cart after removal, assert total price decreases, and item count changes
		releaseCartRespAfterRemove, err := GetUserCart(user.UUID, ReleaseURL, "api-14")
		assert.NoError(t, err)

		err = json.NewDecoder(releaseCartRespAfterRemove.Body).Decode(&releaseCart)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseCartRespAfterRemove.StatusCode)

		assert.Less(t, releaseCart.TotalPrice, initialReleaseTotalPrice, "Total price should decrease after removal")
		assert.Equal(t, initialReleaseItemsCount-1, len(releaseCart.Items), "Item list should decrease by 1")

		// remove the same item from the Dev cart
		devRemoveResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/remove", DevURL, user.UUID), removeItemData, "api-14")
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusNotFound, devRemoveResp.StatusCode)

		// Decode the response from Dev into CartResponse struct
		var devCart CartResponse
		err = json.NewDecoder(devRemoveResp.Body).Decode(&devCart)
		assert.NoError(t, err)

		assert.Equal(t, 0, devCart.TotalPrice, "Total price on Dev should be 0")
		assert.Equal(t, 0, len(devCart.Items), "Item list on Dev should be empty indicates every item removed from cart")

		// Check Release cart after removal, assert total price is reset to 0, and item list is empty
		devCartRespAfterRemove, err := GetUserCart(user.UUID, ReleaseURL, "api-14")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devCartRespAfterRemove.StatusCode)

		err = json.NewDecoder(devCartRespAfterRemove.Body).Decode(&devCart)
		assert.NoError(t, err)
		assert.Equal(t, 0, devCart.TotalPrice, "Total price on Release should be 0")
		assert.Equal(t, 0, len(devCart.Items), "Item list on Release should be empty due to the bug")
	})
}

// api-15 Clear user's cart
func TestClearCart(t *testing.T) {
	t.Run("Clear cart and compare Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		// add a single item to the cart for the user on Release
		_, err = AddItemToCart(user.UUID, game.UUID, 1, ReleaseURL, "api-15") // Add item
		assert.NoError(t, err)

		var releaseCart CartResponse

		// clear the cart on Release
		clearCartData := struct{}{} // Empty body for the clear request
		releaseClearResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/clear", ReleaseURL, user.UUID), clearCartData, "api-15")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseClearResp.StatusCode)

		// Check Release cart after clear, assert items are empty, and total price is 0
		releaseCartRespAfterClear, err := GetUserCart(user.UUID, ReleaseURL, "api-15")
		assert.NoError(t, err)

		err = json.NewDecoder(releaseCartRespAfterClear.Body).Decode(&releaseCart)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseCartRespAfterClear.StatusCode)

		assert.Equal(t, 0, len(releaseCart.Items), "Item list should be empty after clear on Release")
		assert.Equal(t, 0, releaseCart.TotalPrice, "Total price should be 0 after clear on Release")

		// add a single item to the cart for the user on Release
		_, err = AddItemToCart(user.UUID, game.UUID, 1, ReleaseURL, "api-15") // Add item
		assert.NoError(t, err)

		// clear the cart on Dev
		devClearResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/cart/clear", DevURL, user.UUID), clearCartData, "api-15")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devClearResp.StatusCode)

		var devCart CartResponse
		err = json.NewDecoder(devClearResp.Body).Decode(&devCart)
		assert.NoError(t, err)

		assert.NotZero(t, devCart.TotalPrice, "Total price on Dev should remain the same")
		assert.NotEqual(t, 0, len(devCart.Items), "Item list on Dev should not be empty, bug in Dev cart clear")

		// get the cart from Dev after clear request to verify if the bug persists
		cartRespAfterClear, err := GetUserCart(user.UUID, ReleaseURL, "api-15")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, cartRespAfterClear.StatusCode)

		err = json.NewDecoder(cartRespAfterClear.Body).Decode(&devCart)
		assert.NoError(t, err)
		assert.NotZero(t, devCart.TotalPrice, "Total price on Dev should still be the same after clear")
		assert.NotEqual(t, 0, len(devCart.Items), "Item list on Dev should still not be empty after clear due to the bug")
	})
}
