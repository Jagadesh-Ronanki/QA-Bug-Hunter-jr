package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-19 get a payment
func TestGetPayment(t *testing.T) {
	t.Run("Get payment and compare Release and Dev", func(t *testing.T) {
		// create an order in Release
		user, err := FetchExistingUser(9)
		assert.NoError(t, err)

		// Fetch an existing game
		game, err := FetchExistingGame(2)
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

		// creating an order
		releaseOrderResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/orders", ReleaseURL, user.UUID), orderData, "api-16")
		assert.NoError(t, err)

		var releaseOrderRespBody map[string]interface{}
		err = json.NewDecoder(releaseOrderResp.Body).Decode(&releaseOrderRespBody)
		assert.NoError(t, err)

		orderUUIDRelease, _ := releaseOrderRespBody["uuid"].(string)

		// create payment data for Release
		paymentDataRelease := struct {
			OrderUUID     string `json:"order_uuid"`
			PaymentMethod string `json:"payment_method"`
		}{
			OrderUUID:     orderUUIDRelease,
			PaymentMethod: "mir_pay",
		}

		// creating a payment in release
		releasePaymentResp, err := SendPostRequest(fmt.Sprintf("%s/users/%s/payments", ReleaseURL, user.UUID), paymentDataRelease, "api-20")
		assert.NoError(t, err)

		var releasePaymentRespBody map[string]interface{}
		err = json.NewDecoder(releasePaymentResp.Body).Decode(&releasePaymentRespBody)
		assert.NoError(t, err)

		paymentUUIDRelease, _ := releasePaymentRespBody["uuid"].(string)

		// get payment details from release
		releaseGetResp, err := SendGetRequest(fmt.Sprintf("%s/payments/%s", ReleaseURL, paymentUUIDRelease), "api-19")
		assert.NoError(t, err)

		var releaseGetRespBody map[string]interface{}
		err = json.NewDecoder(releaseGetResp.Body).Decode(&releaseGetRespBody)
		assert.NoError(t, err)

		// get payment details from dev
		devGetResp, err := SendGetRequest(fmt.Sprintf("%s/payments/%s", DevURL, paymentUUIDRelease), "api-19")
		assert.NoError(t, err)

		var devGetRespBody map[string]interface{}
		err = json.NewDecoder(devGetResp.Body).Decode(&devGetRespBody)
		assert.NoError(t, err)

		// Release response contains both created_at and updated_at
		assert.Contains(t, releaseGetRespBody, "created_at", "Release response should contain created_at")
		assert.Contains(t, releaseGetRespBody, "updated_at", "Release response should contain updated_at")

		// dev response doesn't contains both created_at and updated_at
		assert.NotContains(t, devGetRespBody, "created_at", "Dev response should not contain created_at")
		assert.NotContains(t, devGetRespBody, "updated_at", "Dev response should not contain updated_at")
	})
}
