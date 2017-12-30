# json + boltDB = joltDB
An api written in Go to simplify working with boltdb/bolt. Jolt saves structs as json and retrieves the stored data as json.

Please have a look at the example.

# Usage:

// 
package main

import (
	"fmt"
	"github.com/jasonmain/joltDB"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

func main() {

	// Open boltDB database in working dir.
	// If one isn't there it will be created.
	joltDB.Open("./data.db")

	// Add some users.
	users := []*User{
		{"100", "Ken", "39"},
		{"101", "Max", "25"},
		{"102", "Bill", "23"},
		{"103", "Martin", "54"},
		{"104", "Kross", "33"},
	}

	// Save users.
	for _, user := range users {
		// With simple key
		if err := joltDB.Save("users", user.ID, user); err != nil {
			fmt.Println(err)
		}
		// With a prefixed key.
		if err := joltDB.Save("users", user.Name+":"+user.ID, user); err != nil {
			fmt.Println(err)
		}
	}

	// Get one user.
	getUser, err := joltDB.GetOne("users", "103")
	if err != nil {
		fmt.Println(err)
	}

	// Get full user list.
	userList, err := joltDB.List("users")
	if err != nil {
		fmt.Println(err)
	}

	// Get user list by range.
	userListRange, err := joltDB.ListRange("users", "101", "103")
	if err != nil {
		fmt.Println(err)
	}

	// Get user list by Prefix.
	userListPrefix, err := joltDB.ListPrefix("users", "K")
	if err != nil {
		fmt.Println(err)
	}

	// Print the saved user lists.
	fmt.Printf("\n\rGet One User:\n\r%s\n\r\n\r", getUser)
	fmt.Printf("Full User List:\n\r%s\n\r\n\r", userList)
	fmt.Printf("User List By Range:\n\r%s\n\r\n\r", userListRange)
	fmt.Printf("User List By Prefix:\n\r%s", userListPrefix)
}
