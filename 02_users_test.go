package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// api-21: List All Users
func TestListAllUsers(t *testing.T) {
	t.Run("List All Users on Release and Dev", func(t *testing.T) {

		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/users", ReleaseURL), "api-21")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseResp.StatusCode)

		releaseBody, err := ParseJSONResponse(releaseResp)
		assert.NoError(t, err)
		assert.Contains(t, releaseBody, "meta")
		assert.Contains(t, releaseBody["meta"], "total")
		releaseTotal := releaseBody["meta"].(map[string]interface{})["total"].(float64)
		assert.Equal(t, releaseTotal, 11.0)

		devResp, err := SendGetRequest(fmt.Sprintf("%s/users", DevURL), "api-21")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResp.StatusCode)

		devBody, err := ParseJSONResponse(devResp)
		assert.NoError(t, err)
		assert.Contains(t, devBody, "meta")
		assert.Contains(t, devBody["meta"], "total")
		devTotal := devBody["meta"].(map[string]interface{})["total"].(float64)

		// mismatch in the 'total' between Release and Dev
		assert.NotEqual(t, releaseTotal, devTotal, "Release and Dev 'total' values should not match")
	})
}

// api-7 Get a user by email and pass
func TestUserLogin(t *testing.T) {
	t.Run("User Login on Release and Dev", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		email := user.Email

		loginData := LoginRequest{
			Email:    email,
			Password: "password",
		}

		// release req
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", ReleaseURL), loginData, "api-7")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseResp.StatusCode)

		// dev req
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", DevURL), loginData, "api-7")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, devResp.StatusCode)

		// The test compares the response codes between Release and Dev
		assert.NotEqual(t, releaseResp.StatusCode, devResp.StatusCode, "The response codes should differ between Release and Dev")
	})
}

// api-3 create a new user
func TestCreateUser(t *testing.T) {
	t.Run("Create User on Release and Dev with Existing Nickname", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		createUserData := UserCreateRequest{
			Email:    "new.1@gmail.com",
			Password: "password",
			Name:     "new.1 user",
			Nickname: user.Nickname, // existing nickname
		}

		// release req
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users", ReleaseURL), createUserData, "api-3")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, releaseResp.StatusCode) // Expect 409 due to conflict (duplicate nickname)

		// dev req
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users", DevURL), createUserData, "api-3")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResp.StatusCode) // Expect 200 even though the nickname is duplicated

		// should Mismatch codes between Release and Dev
		assert.NotEqual(t, releaseResp.StatusCode, devResp.StatusCode, "The response codes should differ between Release and Dev")

		// reset server as it bricks next calls
		setupResp, err := SendPostRequest("https://release-gs.qa-playground.com/api/v1/setup", "", "api-6")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusResetContent, setupResp.StatusCode, "Setup should return a 205 Content Reset")
	})
}

// api-22 create a new user
func TestCreateUser2(t *testing.T) {
	t.Run("Create User with Duplicate Nickname, Existing Email, Short Password, and Valid Data", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		existingNicknameUserData := UserCreateRequest{
			Email:    "max@gmail.com",
			Password: "password",
			Name:     "max",
			Nickname: user.Nickname, // Reusing the existing nickname
		}

		// release req
		releaseResp, err := SendPostRequest(fmt.Sprintf("%s/users", ReleaseURL), existingNicknameUserData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, releaseResp.StatusCode) // Expect 409 due to conflict (duplicate nickname)

		// dev req
		devResp, err := SendPostRequest(fmt.Sprintf("%s/users", DevURL), existingNicknameUserData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, devResp.StatusCode) // Expect 500 (server issue)

		// release: Test with existing Email
		existingEmailData := UserCreateRequest{
			Email:    user.Email,
			Password: "password",
			Name:     "new user",
			Nickname: "newNickname",
		}

		releaseRespEmail, err := SendPostRequest(fmt.Sprintf("%s/users", ReleaseURL), existingEmailData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, releaseRespEmail.StatusCode) // Expect 409

		// dev: Test with Existing Email
		devRespEmail, err := SendPostRequest(fmt.Sprintf("%s/users", DevURL), existingEmailData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, devRespEmail.StatusCode) // Expect 500 (server issue)

		// release: Test with short Pass
		shortPasswordData := UserCreateRequest{
			Email:    "shortpassword@example.com",
			Password: "123", // Password less than 5 characters
			Name:     "Short Password User",
			Nickname: "shortpass",
		}

		releaseRespPassword, err := SendPostRequest(fmt.Sprintf("%s/users", ReleaseURL), shortPasswordData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, releaseRespPassword.StatusCode) // Expect 400

		// dev: Test with short Pass
		devRespPassword, err := SendPostRequest(fmt.Sprintf("%s/users", DevURL), shortPasswordData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, devRespPassword.StatusCode) // Expect 400

		// perfect data to create new user
		validData := UserCreateRequest{
			Email:    "valid.email@example.com", // Non-existing email
			Password: "password",                // Valid password
			Name:     "Valid User",
			Nickname: "validNickname", // Non-existing nickname
		}

		// release req
		releaseRespValid, err := SendPostRequest(fmt.Sprintf("%s/users", ReleaseURL), validData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseRespValid.StatusCode) // Expect 200 due to valid input

		// dev req
		devRespValid, err := SendPostRequest(fmt.Sprintf("%s/users", DevURL), validData, "api-22")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, devRespValid.StatusCode) // Expect 500 due to server issue

		// reset server as it bricks next calls
		setupResp, err := SendPostRequest("https://release-gs.qa-playground.com/api/v1/setup", "", "api-6")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusResetContent, setupResp.StatusCode, "Setup should return a 205 Content Reset")
	})
}

// api-4 update a user
func TestUpdateUserWithDuplicateEmail(t *testing.T) {
	t.Run("Update User with Duplicate Email", func(t *testing.T) {
		// Fetch an existing user to reuse their email for the test
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		userLast, err := FetchExistingUser(1)
		assert.NoError(t, err)

		updateUserData := UserUpdateRequest{
			Email:    user.Email, // Reusing the existing email
			Password: "password",
			Name:     "new name",
			Nickname: "newNickname",
		}

		// release req
		releaseResp, err := SendPatchRequest(fmt.Sprintf("%s/users/%s", ReleaseURL, userLast.UUID), updateUserData, "api-4")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, releaseResp.StatusCode)

		// dev req
		devResp, err := SendPatchRequest(fmt.Sprintf("%s/users/%s", DevURL, userLast.UUID), updateUserData, "api-4")
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusConflict, devResp.StatusCode)

		// Compare response codes between Release and Dev for mismatch
		assert.NotEqual(t, releaseResp.StatusCode, devResp.StatusCode, "The response codes should differ between Release and Dev")
	})
}

// api-24 update a user
func TestUpdateUserAndLogin(t *testing.T) {
	t.Run("Update User and Test Login", func(t *testing.T) {
		user, err := FetchExistingUser(0)
		assert.NoError(t, err)

		updateUserData := UserCreateRequest{
			Email:    user.Email,
			Password: "passwords", // changing password
			Name:     user.Name,
			Nickname: user.Nickname,
		}

		// release req
		releaseResp, err := SendPatchRequest(fmt.Sprintf("%s/users/%s", ReleaseURL, user.UUID), updateUserData, "api-24")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseResp.StatusCode)

		// login
		loginData := LoginRequest{
			Email:    user.Email,
			Password: "passwords",
		}

		// release req login
		releaseLoginResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", ReleaseURL), loginData, "api-24")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseLoginResp.StatusCode)

		// Update in dev -----------------------------------

		userLast, err := FetchExistingUser(0)
		assert.NoError(t, err)

		updateUserLastData := UserCreateRequest{
			Email:    userLast.Email,
			Password: "passwords", // changing password
			Name:     userLast.Name,
			Nickname: userLast.Nickname,
		}

		// dev req
		devResp, err := SendPatchRequest(fmt.Sprintf("%s/users/%s", DevURL, userLast.UUID), updateUserLastData, "api-24")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResp.StatusCode)

		// login
		loginLastData := LoginRequest{
			Email:    user.Email,
			Password: "passwords",
		}

		// release req login
		devLoginResp, err := SendPostRequest(fmt.Sprintf("%s/users/login", ReleaseURL), loginLastData, "api-24")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, devLoginResp.StatusCode)

		// Compare update response codes between Release and Dev for match
		assert.Equal(t, releaseResp.StatusCode, devResp.StatusCode, "The response codes are same between Release and Dev for Update User Details request")

		// Compare Login response codes between Release and Dev for mismatch
		assert.NotEqual(t, releaseLoginResp.StatusCode, devLoginResp.StatusCode, "The response codes should differ between Release and Dev for Login request")
	})
}

// api-6 List all Users
func TestOffsetHandlingMismatch(t *testing.T) {
	t.Run("Test Offset Handling Mismatch", func(t *testing.T) {
		offset := 15
		limit := 15

		// release req
		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/users?offset=%d&limit=%d", ReleaseURL, offset, limit), "api-6")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseResp.StatusCode)

		releaseBody, err := ParseJSONResponse(releaseResp)
		assert.NoError(t, err)

		releaseUsers := releaseBody["users"].([]interface{})

		// dev req
		devResp, err := SendGetRequest(fmt.Sprintf("%s/users?offset=%d&limit=%d", DevURL, offset, limit), "api-6")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResp.StatusCode)

		devBody, err := ParseJSONResponse(devResp)
		assert.NoError(t, err)

		devUsers := devBody["users"].([]interface{})

		// Compare the lengths: Release should return users based on offset, Dev should return all users
		assert.NotEqual(t, len(releaseUsers), len(devUsers), "The user list length mismatch: Release and Dev servers handle offset differently.")
	})
}

// api-23 Get a user
func TestFetchUserByUUID(t *testing.T) {
	t.Run("Test Fetch User by UUID Mismatch", func(t *testing.T) {
		user, err := FetchExistingUser(2)
		assert.NoError(t, err)

		// release req
		releaseResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s", ReleaseURL, user.UUID), "api-23")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, releaseResp.StatusCode)

		releaseBody, err := ParseJSONResponse(releaseResp)
		assert.NoError(t, err)

		releaseUUID := releaseBody["uuid"].(string)
		assert.Equal(t, user.UUID, releaseUUID, "The UUID from Release server does not match the requested UUID")

		// dev req
		devResp, err := SendGetRequest(fmt.Sprintf("%s/users/%s", DevURL, user.UUID), "api-23")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, devResp.StatusCode)

		devBody, err := ParseJSONResponse(devResp)
		assert.NoError(t, err)

		devUUID := devBody["uuid"].(string)

		// compare uuids returned in response for same user.UUID for mismatch
		assert.NotEqual(t, user.UUID, devUUID, "The UUID from Dev server matches the requested UUID")
	})
}

// api-1 Delete User
func TestDeleteUserByUUID(t *testing.T) {
	t.Run("Test DELETE User by UUID", func(t *testing.T) {
		existingUser, err := FetchExistingUser(0)
		assert.NoError(t, err)

		// release req existing uuid
		releaseResp, err := SendDeleteRequest(fmt.Sprintf("%s/users/%s", ReleaseURL, existingUser.UUID), "api-1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, releaseResp.StatusCode, "Release server should return 204 for existing UUID")

		nonExistingUUID := existingUser.UUID

		// release req non existing uuid
		releaseRespNonExist, err := SendDeleteRequest(fmt.Sprintf("%s/users/%s", ReleaseURL, nonExistingUUID), "api-1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, releaseRespNonExist.StatusCode, "Release server should return 404 for non-existing UUID")

		// dev req with non-existing UUID
		devRespNonExist, err := SendDeleteRequest(fmt.Sprintf("%s/users/%s", DevURL, nonExistingUUID), "api-1")
		assert.NoError(t, err)

		// dev server incapable of handling this request
		assert.Equal(t, http.StatusInternalServerError, devRespNonExist.StatusCode, "Dev server should return 500 for non-existing UUID")
	})
}
