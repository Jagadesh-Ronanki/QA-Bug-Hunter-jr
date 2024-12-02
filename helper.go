// helpers.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ReleaseURL = "https://release-gs.qa-playground.com/api/v1"
	DevURL     = "https://dev-gs.qa-playground.com/api/v1"
	AuthHeader = "Bearer qahack2024:jagadeshc0891@gmail.com"
)

// login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// creating a user
type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

// update a user
type UserUpdateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

// user
type User struct {
	UUID      string `json:"uuid"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// game
type Game struct {
	CategoryUUIDs []string `json:"category_uuids"`
	Price         int      `json:"price"`
	Title         string   `json:"title"`
	UUID          string   `json:"uuid"`
}

// category partial
type Category struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

// add to wishlist
type WishlistBody struct {
	ItemUUID string `json:"item_uuid"`
}

// add to cart
type AddItemRequest struct {
	ItemUUID string `json:"item_uuid"`
	Quantity int    `json:"quantity"`
}

// change quantity in cart
type ChangeItemQuantityRequest struct {
	ItemUUID string `json:"item_uuid"`
	Quantity int    `json:"quantity"`
}

// delete item in cart
type RemoveItemRequest struct {
	ItemUUID string `json:"item_uuid"`
	Quantity int    `json:"quantity"`
}

// item in the user's cart
type CartItem struct {
	ItemUUID   string `json:"item_uuid"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
}

// structure of the cart response
type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice int        `json:"total_price"`
	UserUUID   string     `json:"user_uuid"`
}

// update order status
type OrderStatusUpdateRequest struct {
	Status string `json:"status"`
}

// Helper - POST requests
func SendPostRequest(url string, body interface{}, taskID string) (*http.Response, error) {
	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper - GET requests
func SendGetRequest(url string, taskID string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)

	// No Cache headers
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper - DELETE requests
func SendDeleteRequest(url string, taskID string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper - PATCH requests
func SendPatchRequest(url string, body interface{}, taskID string) (*http.Response, error) {
	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper - PUT requests
func ParseJSONResponse(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()
	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

// Helper - PUT requests with file
func SendPutRequestWithFile(url string, filePath string, taskID string) (*http.Response, error) {
	// Create a buffer to write the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// If file doesn't exist, create a temporary empty file
	var file *os.File
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		// Create an empty temporary file
		file, err = os.CreateTemp("", "empty-avatar-*.jpg")
		if err != nil {
			return nil, fmt.Errorf("failed to create empty file: %v", err)
		}
		defer file.Close()
	} else {
		// Open the existing file (if it exists)
		file, err = os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()
	}

	// Create a form file field named "avatar_file" and attach the file
	part, err := writer.CreateFormFile("avatar_file", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form file (empty or real file content)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create the PUT request with the appropriate headers
	req, err := http.NewRequest("PUT", url, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}

	// Add necessary headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer qahack2024:jagadeshc0891@gmail.com") // Replace with actual authorization token
	req.Header.Set("X-Task-Id", taskID)

	// Send the request using the http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %v", err)
	}

	return resp, nil
}

// ----------- Other Helpers Functions --------------

func FetchAllUsers(url, taskID string) ([]map[string]interface{}, error) {
	// Send the GET request to fetch all users
	resp, err := SendGetRequest(fmt.Sprintf("%s/users", url), taskID)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	// Ensure "users" field exists in the response
	users, ok := body["users"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'users' key in the response")
	}

	// Convert []interface{} to []map[string]interface{}
	var userList []map[string]interface{}
	for _, user := range users {
		userList = append(userList, user.(map[string]interface{}))
	}

	return userList, nil
}

func FetchExistingUser(index int32) (*User, error) {
	// Get the user list from the server
	resp, err := SendGetRequest(fmt.Sprintf("%s/users", ReleaseURL), "api-6")
	if err != nil {
		return nil, err
	}

	// Parse and extract the users
	body, err := ParseJSONResponse(resp)
	if err != nil {
		return nil, err
	}

	// Assuming there is at least one user in the response
	if users, ok := body["users"].([]interface{}); ok && len(users) > 0 {
		firstUser := users[index].(map[string]interface{})
		// Populate the User struct with the necessary fields
		user := &User{
			UUID:      firstUser["uuid"].(string),
			Email:     firstUser["email"].(string),
			Nickname:  firstUser["nickname"].(string),
			Name:      firstUser["name"].(string),
			AvatarURL: firstUser["avatar_url"].(string),
		}
		return user, nil
	}

	return nil, fmt.Errorf("no users found")
}

func FetchExistingGame(index int32) (*Game, error) {
	// Send GET request to fetch all games
	resp, err := SendGetRequest(fmt.Sprintf("%s/games", ReleaseURL), "api-9")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response body to extract the game data
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	// Check if "games" is available in the response and extract the first game
	if games, ok := body["games"].([]interface{}); ok && len(games) > 0 {
		// Extract the first game and unmarshal it into the Game struct
		firstGame := games[index].(map[string]interface{})
		gameJSON, err := json.Marshal(firstGame)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal game data: %w", err)
		}

		var game Game
		err = json.Unmarshal(gameJSON, &game)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal game data: %w", err)
		}

		// Return the game struct
		return &game, nil
	}

	return nil, fmt.Errorf("no games found")
}

func FetchCategory() (*Category, error) {
	// Send GET request to fetch all categories
	resp, err := SendGetRequest(fmt.Sprintf("%s/categories", ReleaseURL), "api-10")
	if err != nil {
		return nil, err
	}

	// Parse the response body to extract the category UUID
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	// Check if categories exist in the response body
	if categories, ok := body["categories"].([]interface{}); ok && len(categories) > 0 {
		// Extract UUID of the first category
		firstCategory := categories[0].(map[string]interface{})

		// Check if UUID field exists
		if uuid, exists := firstCategory["uuid"].(string); exists {
			return &Category{
				UUID: uuid,
			}, nil
		} else {
			return nil, fmt.Errorf("category UUID is missing or not a string")
		}
	}

	return nil, fmt.Errorf("no categories found in the response")
}

func AddItemToCart(userUUID string, itemUUID string, quantity int, environmentURL string, taskID string) (*http.Response, error) {
	// Prepare the request body using the AddItemRequest struct
	requestBody := AddItemRequest{
		ItemUUID: itemUUID,
		Quantity: quantity,
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	// Create a POST request to add the item to the cart
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/cart/add", environmentURL, userUUID), bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %v", err)
	}

	return resp, nil
}

func GetUserCart(userUUID string, environmentURL string, taskID string) (*http.Response, error) {
	// Send GET request to fetch the cart
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/cart", environmentURL, userUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", AuthHeader)
	req.Header.Set("X-Task-Id", taskID)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}

	return resp, nil
}
