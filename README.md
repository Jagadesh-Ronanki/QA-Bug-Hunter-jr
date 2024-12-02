# QA Bug Hunter - Test Suite

This repository contains the test suite for the "QA Bug Hunter" project, which tests various endpoints of a web application, including user management, games, orders, payments, and more.

[Installation](#installation)  
[Run Tests](#running-tests)  
[Bugs Description](#bug-notes)

## File Structure

Here is the structure of the test files in the repository:

```bash
├── 01_setup_test.go
├── 02_users_test.go
├── 03_wishlist_test.go
├── 04_games_test.go
├── 05_categories_test.go
├── 06_avatar_test.go
├── 07_cart_test.go
├── 08_orders_test.go
├── 09_payment_test.go
├── go.mod
├── go.sum
├── helper.go
└── README.md --> You are Here
```

### Test Files

Each test file corresponds to a specific functionality of the application.


- `02_users_test.go`: Tests related to user creation and user-specific endpoints.

## Prerequisites

Before you begin, make sure you have the following installed on your system:

- **Go**: Go version 1.18 (1.23.3 prefered) or above is required.
- **Go Modules**: Ensure Go Modules are enabled (enabled by default in Go 1.11+).

## Installation

1. Clone this repository to your local machine.

    ```bash
    https://github.com/Jagadesh-Ronanki/QA-Bug-Hunter-jr/
    cd QA-Bug-Hunter-jr
    ```

2. Install the required dependencies using Go Modules.

    ```bash
    go mod tidy
    ```

## Running Tests

To run the entire test suite, you can use the following command:

```bash
go test -v
```

This will run all the tests in the `*_test.go` files in the repository.

### Running a Specific Test File

If you want to run tests from a specific file, use:

```bash
go test -v 02_users_test.go
```

This will run only the tests defined in the `02_users_test.go` file.<br/>
But make sure you run `01_setup_test.go` before running individual file


## Bug Notes

Note: Both `helpers.go` and `01_setup_test.go` doesn't contain any tests but essential to run all the test cases. If you are running individual testcases run `go test -v 01_setup_test.go` to complete the setup. 

### Users (9/9) | [Tests](./02_users_test.go)

1. **API-21: List All Users**  
   - **Bug:** Mismatch in the 'total' value between Release and Dev servers.
   - **Release vs. Dev:** Release returns 11 total users, but Dev returns a different count.
  
2. **API-7: User Login**  
   - **Bug:** Inconsistent response codes for login requests.  
   - **Release vs. Dev:** Release returns `200 OK` while Dev returns `404 Not Found` for the same login request.

3. **API-3: Create User with Existing Nickname**  
   - **Bug:** Conflict in handling user creation with an existing nickname.  
   - **Release vs. Dev:** Release returns `409 Conflict` when attempting to create a user with an existing nickname, while Dev allows it (`200 OK`)

4. **API-22: Create User with Duplicate Nickname, Existing Email, and Short Password**  
   - **Bug:** Inconsistent handling of duplicate nickname, email, and short passwords.  
   - **Release vs. Dev:** 
     - Duplicate nickname: Release returns `409 Conflict`, while Dev does not handle it properly.
     - Existing email: Release returns `409 Conflict`, while Dev returns `500 Internal Server Error`.
     - Short password: Both servers return `400 Bad Request` for invalid short passwords.
     - Valid user creation: Release processes it correctly (`200 OK`), but Dev fails (`500 Internal Server Error`).

5. **API-4: Update User with Duplicate Email**  
   - **Bug:** Mismatch in handling updates with duplicate email.  
   - **Release vs. Dev:** Release returns `409 Conflict` when trying to update a user with an existing email, but Dev server handles it incorrectly and updates with duplicate email.

6. **API-24: Update User and Login**  
   - **Bug:** Mismatch in login after updating user details.  
   - **Release vs. Dev:**  
     - Update request: Both servers return `200 OK` for updating user details.
     - Login: Release allows login with the updated password, while Dev server returns `404 Not Found` for the same login credentials.

7. **API-6: List Users with Offset Handling**  
   - **Bug:** Inconsistent handling of the offset and limit in pagination.
   - **Release vs. Dev:** Release returns the correct user list based on offset and limit, while Dev returns all users, leading to a mismatch in the number of users returned.

8. **API-23: Fetch User by UUID**  
   - **Bug:** Mismatch in UUIDs returned for the same user on Release and Dev servers.  
   - **Release vs. Dev:** Release retuns the information of user with matching UUID as passed. Whereas dev is returning another user details. 

9. **API-1: Delete User by UUID**  
   - **Bug:** Inconsistent handling of non-existing UUID deletion requests.
   - **Release vs. Dev:** Release returns `404 Not Found` for non-existing UUIDs, while Dev returns `500 Internal Server Error`

---

### **Wishlist (3/3)** | [Tests](./03_wishlist_test.go)

1. **API-5: Add Item to User's Wishlist**  
   - **Bug:** Mismatch in wishlist item handling between Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Successfully adds an item to the wishlist (status `200 OK`).  
     - Dev: Returns `422 Unprocessable Entity` when trying to add an item to the wishlist before reaching the limit.  

2. **API-25: Add Item to User's Wishlist**  
   - **Bug:** Wishlist item handling inconsistencies between Release and Dev servers when removing items and adding new ones.  
   - **Release vs. Dev:**  
     - Release: Correctly adds an item to the wishlist (status `200 OK`) and persists it.  
     - Dev: Adds an item temporarily (status `200 OK`) but does not persist the change.  

3. **API-8: Remove Item from User's Wishlist**  
   - **Bug:** Mismatch in the response when removing an item from the wishlist on Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Successfully removes an item (status `200 OK`) and returns `404 Not Found` when attempting to remove the same item again.  
     - Dev: Returns `200 OK` even if the item was not removed (indicating no actual change in the wishlist).

---

### **Games (2/2)** | [Tests](./04_games_test.go) 

1. **API-2: Search Games**  
   - **Bug:** Query is not working as expected in Dev server.  
   - **Release vs. Dev:**  
     - Release: Returns game that matches the query.  
     - Dev: Returns every game indicates that the query is not working.  

2. **API-9: Get Game by UUID**  
   - **Bug:** Mismatch in the status code when fetching a game by UUID between Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Successfully returns the game details with status code `200 OK`.  
     - Dev: Returns `404 Not Found` indicating the game is not found on the Dev server.

---

### **Categories (3/3)** | [Tests](./05_categories_test.go)

1. **API-10: Get Games by Category**  
   - **Bug:** Mismatch in handling the category UUID in the response between Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Returns games with the correct category UUID, matching the requested category.  
     - Dev: Returns games where the category UUID does **not** match the requested category, indicating a category filtering issue on the Dev server.

---

### **Avatars (1/1)** | [Tests](./06_avatar_test.go)

1. **API-11: Update User's Avatar**  
   - **Bug:** Inconsistent handling of avatar updates between Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Successfully updates the user's avatar and correctly reflects the updated avatar URL during login and user details.  
     - Dev: Successfully updates the avatar but fails to persit the changes into database. Despite this, both servers return the same avatar URL during login (Release's updated avatar URL vs. Dev's non-updated one).

Here is the structured summary for **Cart (1/5) Category** based on the test cases you provided:

---

### **Cart (5/5)** | [Tests](./07_cart_test.go)

1. **API-12: Get a Cart**
   - **Bug:** The total price of the cart differs between Release and Dev servers.  
   - **Release vs. Dev:**  
     - Release: Returns the correct cart total price.  
     - Dev: Returns an incorrect cart total price, which is not synchronized with the Release server.

2. **API-13: Change an Item in User's Cart**
   - **Bug:** Item quantity updates on Dev are not properly reflected, and the cart items appears empty after modifications in dev response.  
   - **Release vs. Dev:**  
     - Release: Successfully updates the item quantity and reflects the correct price change.  
     - Dev: Item quantity change is successful but and the response shoes empty items list.

3. **API-14: Remove an Item from User's Cart**
   - **Bug:** Removal of items from the cart is clears all items on the Dev server.  
   - **Release vs. Dev:**  
     - Release: Correctly removes the item from the cart and adjusts the total price and items list.  
     - Dev: makes the cart empty after requesting to remove single item.

4. **API-15: Clear User's Cart**
   - **Bug:** Clearing the cart works on the Release server, but the Dev server does not reset the cart correctly.  
   - **Release vs. Dev:**  
     - Release: Clears the cart successfully, resetting the total price and items list to zero.  
     - Dev: Does not clear the cart, and the total price and items remain, indicating a bug in the Dev server.

---

### **Orders (3/3) Category** | [Tests](./08_orders_test.go)

1. **API-16: Create a New Order**
   - **Bug:** Duplicate items in the order are handled differently between the Release and Dev environments.  
   - **Release vs. Dev:**  
     - Release: Returns a `400 Bad Request` when duplicate items are added to the order.  
     - Dev: Does not return an error for duplicate items, and instead processes the order successfully (returns `200 OK`) and added to the orders.

2. **API-17: List All Orders for a User**
   - **Bug:** The `limit` and `offset` parameters are not properly enforced on the Dev server, returning more orders than expected.  
   - **Release vs. Dev:**  
     - Release: Correctly applies the `limit` and `offset` query parameters and returns only the specified number of orders.  
     - Dev: Ignores the `limit` and returns more orders than requested.

3. **API-18: Update an Order Status**
   - **Bug:** The order status update behavior differs between Release and Dev servers when updating the status of an order.  
   - **Release vs. Dev:**  
     - Release: Successfully updates the order status to `canceled` for orders with an `open` status.  
     - Dev: Returns a `422 Unprocessable Entity` when attempting to update an order status to `canceled` on orders also with an `open` status. Which is unexpected behaviour.

---

### **Payments (1/2) Category** | [Tests](./09_payment_test.go)

1. **API-19: Get a Payment**
   - **Bug:** The response from the Dev server is missing the `created_at` and `updated_at` fields, which are present in the Release server's response.  
   - **Release vs. Dev:**  
     - Release: Returns payment details including both `created_at` and `updated_at` timestamps, reflecting the accurate time of creation and last update.  
     - Dev: The response lacks both `created_at` and `updated_at` fields.


---

Done with reading? Clone and Run tests :)
