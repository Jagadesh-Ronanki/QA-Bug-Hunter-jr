package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-11 Update user's avatar
func TestUpdateUserAvatar(t *testing.T) {
	t.Run("Update User Avatar and Test Login", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		avatarFilePath := "/path/to/valid/avatar.jpg"

		// release req
		releaseAvatarResp, err := SendPutRequestWithFile(fmt.Sprintf("%s/users/%s/avatar", ReleaseURL, user.UUID), avatarFilePath, "api-11")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseAvatarResp.StatusCode)

		// storing the avatar_url
		var releaseAvatar map[string]interface{}
		err = json.NewDecoder(releaseAvatarResp.Body).Decode(&releaseAvatar)
		assert.NoError(t, err)
		releaseAvatarURL := releaseAvatar["avatar_url"].(string)

		loginData := LoginRequest{
			Email:    user.Email,
			Password: "password",
		}

		// Test on Release server - login to fetch the avatar
		releaseLoginResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", ReleaseURL), loginData, "api-11")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseLoginResp.StatusCode)

		var releaseLoginData map[string]interface{}
		err = json.NewDecoder(releaseLoginResp.Body).Decode(&releaseLoginData)
		assert.NoError(t, err)
		releaseLoginAvatarURL := releaseLoginData["avatar_url"].(string)

		// Assert that the avatar URL matches the one returned by the avatar update
		assert.Equal(t, releaseAvatarURL, releaseLoginAvatarURL, "Avatar URL on release should match after update")

		// dev req
		devAvatarResp, err := SendPutRequestWithFile(fmt.Sprintf("%s/users/%s/avatar", DevURL, user.UUID), avatarFilePath, "api-11")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devAvatarResp.StatusCode)

		// Store the avatar_url
		var devAvatar map[string]interface{}
		err = json.NewDecoder(devAvatarResp.Body).Decode(&devAvatar)
		assert.NoError(t, err)
		devAvatarURL := devAvatar["avatar_url"].(string)

		// Login to fetch the user details after avatar update on dev
		devLoginResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", DevURL), loginData, "api-11")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devLoginResp.StatusCode)

		var devLoginData map[string]interface{}
		err = json.NewDecoder(devLoginResp.Body).Decode(&devLoginData)
		assert.NoError(t, err)
		devLoginAvatarURL := devLoginData["avatar_url"].(string)

		// Assert that the avatar URL on dev is NOT updated correctly (since the bug exists)
		assert.NotEqual(t, devAvatarURL, devLoginAvatarURL, "Avatar URL on dev should not match after update, indicating the bug")

		// Compare avatar URL between release and dev
		assert.Equal(t, releaseLoginAvatarURL, devLoginAvatarURL, "Avatar URLs between release and dev should be equal as dev fails to update in db")
	})
}
