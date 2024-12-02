package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-16 Create a new Order
func TestCreateOrderWithDuplicateItems(t *testing.T) {
	t.Run("Create order with duplicate items and compare Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game, err := FetchExistingGame(0)
		assert.NoError(t, err)

		// invalid body to orders
		orderData := struct {
			Items []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		}{
			Items: []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			}{
				{ItemUUID: game.UUID, Quantity: 2}, // Add item first time
				{ItemUUID: game.UUID, Quantity: 1}, // Duplicate item
			},
		}

		// creating the order on release
		releaseOrderResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/orders", ReleaseURL, user.UUID), orderData, "api-16")
		assert.NoError(t, err)

		// expecting 400 on release when duplicate items are added
		assert.Equal(t, http.StatusBadRequest, releaseOrderResp.StatusCode)

		// creating the same order on dev
		devOrderResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/orders", DevURL, user.UUID), orderData, "api-16")
		assert.NoError(t, err)

		// Expect a 400 response but returns 200 for duplicate items
		assert.NotEqual(t, http.StatusBadRequest, devOrderResp.StatusCode)
	})
}

// api-17 List all orders for a user
func TestListOrdersWithLimitAndOffset(t *testing.T) {
	t.Run("List orders with limit and offset and compare Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		for i := 0; i < 5; i++ {
			game, err := FetchExistingGame(int32(i))
			assert.NoError(t, err)

			orderData := struct {
				Items []struct {
					ItemUUID string `json:"item_uuid"`
					Quantity int    `json:"quantity"`
				} `json:"items"`
			}{
				Items: []struct {
					ItemUUID string `json:"item_uuid"`
					Quantity int    `json:"quantity"`
				}{
					{ItemUUID: game.UUID, Quantity: 2},
				},
			}

			_, err = SendPostRequest(fmt.Sprintf("%s/users/%s/orders", ReleaseURL, user.UUID), orderData, "api-16")
			assert.NoError(t, err)
		}

		offset := 1
		limit := 1

		// list orders in Release with offset and limit
		releaseOrdersResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/orders?offset=%d&limit=%d", ReleaseURL, user.UUID, offset, limit), "api-17")
		assert.NoError(t, err)

		var releaseOrdersRespBody map[string]interface{}
		err = json.NewDecoder(releaseOrdersResp.Body).Decode(&releaseOrdersRespBody)
		assert.NoError(t, err)

		// Expect only 1 order due to the limit value being considered
		assert.Equal(t, 1, len(releaseOrdersRespBody["orders"].([]interface{})), "Release should return only 1 order based on the limit")

		// list orders in Dev with offset and limit
		devOrdersResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s/orders?offset=%d&limit=%d", DevURL, user.UUID, offset, limit), "api-17")
		assert.NoError(t, err)

		var devOrdersRespBody map[string]interface{}
		err = json.NewDecoder(devOrdersResp.Body).Decode(&devOrdersRespBody)
		assert.NoError(t, err)

		// in Dev - the limit is ignored, and more orders than expected are returned so it returns more than one item
		devOrders := devOrdersRespBody["orders"].([]interface{})
		assert.Greater(t, len(devOrders), 1, "Dev should return more than 1 order as the limit is not considered due to the bug")

		// compare the number of orders returned by Release and Dev
		assert.NotEqual(t, len(releaseOrdersRespBody["orders"].([]interface{})), len(devOrders), "The number of orders should differ between Release and Dev due to the bug in Dev not considering limit")
	})
}

// api-18 Update an order status
func TestUpdateOrderStatus(t *testing.T) {
	t.Run("Update order status and compare Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		game1, err := FetchExistingGame(0)
		assert.NoError(t, err)

		game2, err := FetchExistingGame(0)
		assert.NoError(t, err)

		orderData1 := struct {
			Items []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		}{
			Items: []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			}{
				{ItemUUID: game1.UUID, Quantity: 2},
			},
		}

		orderData2 := struct {
			Items []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		}{
			Items: []struct {
				ItemUUID string `json:"item_uuid"`
				Quantity int    `json:"quantity"`
			}{
				{ItemUUID: game2.UUID, Quantity: 2},
			},
		}

		// Creating an order in release
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/orders", ReleaseURL, user.UUID), orderData1, "api-18")
		assert.NoError(t, err)

		var releaseRespBody map[string]interface{}
		err = json.NewDecoder(releaseResp.Body).Decode(&releaseRespBody)
		assert.NoError(t, err)

		orderUUIDRelease, _ := releaseRespBody["uuid"].(string)

		// another order 2 in Release
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/orders", ReleaseURL, user.UUID), orderData2, "api-18")
		assert.NoError(t, err)

		var devRespBody map[string]interface{}
		err = json.NewDecoder(devResp.Body).Decode(&devRespBody)
		assert.NoError(t, err)

		orderUUIDDev, _ := devRespBody["uuid"].(string)

		statusUpdate := OrderStatusUpdateRequest{
			Status: "canceled",
		}

		// update order status to "canceled" in release
		releaseUpdateResp, err := SendPatchRequest(fmt.Sprintf("%s/orders/%s/status", ReleaseURL, orderUUIDRelease), statusUpdate, "api-18")
		assert.NoError(t, err)

		// update order status to "canceled" in dev
		devUpdateResp, err := SendPatchRequest(fmt.Sprintf("%s/orders/%s/status", DevURL, orderUUIDDev), statusUpdate, "api-18")
		assert.NoError(t, err)

		// compare the response from Release and Dev for mismatch as dev restricting to update orders status even with "open" status and returns 422
		assert.NotEqual(t, releaseUpdateResp.StatusCode, devUpdateResp.StatusCode, "The status of the order should differ between Release and Dev")
	})
}
