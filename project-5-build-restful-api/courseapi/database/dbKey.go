package database

import (
	"database/sql"
	"fmt"
)

// GetAllKeys retrieves data from APIkey table and assign to the map for usage.
func GetAllKeys(db *sql.DB) map[string]string {

	query := `
			  SELECT *
			  FROM APIkey
			 `

	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// map user to apikey
	// and pass to the api for usage
	keyData = make(map[string]string)
	// map the key type to the record in the table
	var userAPIKey Key

	//fmt.Println("Database start.")
	for results.Next() {

		// Scan() copy each row of data from db and assign to the address specified
		err = results.Scan(&userAPIKey.UserName, &userAPIKey.APIKey)
		if err != nil {
			//fmt.Println(course.CourseCode, course.CourseName, course.Description, module.ModuleCode, module.ModuleName, module.Description)
			panic(err.Error())
		}

		keyData[userAPIKey.UserName] = userAPIKey.APIKey
		//fmt.Println(userAPIKey.UserName, userAPIKey.APIKey)
	}
	//fmt.Println("Database end.")

	return keyData
}

// AddKey adds data to APIkey table in database.
func AddKey(db *sql.DB, userAPI Key) {

	query := fmt.Sprintf("INSERT INTO APIkey (username, apikey) VALUES ('%s', '%s')",
		userAPI.UserName, userAPI.APIKey)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// UpdateKey updates data to APIkey table in database based on username.
func UpdateKey(db *sql.DB, userAPI Key) {

	query := fmt.Sprintf("UPDATE APIkey SET apikey='%s' WHERE username='%s'",
		userAPI.APIKey, userAPI.UserName)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// DeleteKey deletes data to APIkey table in database.
func DeleteKey(db *sql.DB, userAPI Key) {

	query := fmt.Sprintf("DELETE FROM APIkey WHERE username='%s'",
		userAPI.UserName)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
