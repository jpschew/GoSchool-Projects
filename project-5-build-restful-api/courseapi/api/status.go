package api

import (
	"fmt"
	"net/http"
)

func notFoundStatus(res http.ResponseWriter, isCourse bool) {
	// write status code to header
	res.WriteHeader(http.StatusNotFound)

	if isCourse {
		// write to the response
		// Write() takes in []byte as input
		// response code 404 means not found
		res.Write([]byte("404 - No course found"))
		fmt.Println(http.StatusNotFound, "- No course found.")
	} else {
		// write to the response
		// Write() takes in []byte as input
		// response code 404 means not found
		res.Write([]byte("404 - No module found"))
		fmt.Println(http.StatusNotFound, "- No module found.")
	}

}

func acceptedStatus(res http.ResponseWriter, urlParams map[string]string, isCourse bool, isUpdate bool) {
	// write status code to header
	res.WriteHeader(http.StatusAccepted)

	if isCourse {
		if isUpdate {
			// write to the response
			// Write() takes in []byte as input
			// response code 202 means accepted for processing
			res.Write([]byte("202 - Course updated:" + urlParams["courseid"]))
			fmt.Println(http.StatusAccepted, urlParams["courseid"], "updated.")
		} else {
			// write to the response
			// Write() takes in []byte as input
			// response code 202 means accepted for processing
			res.Write([]byte("202 -  Course deleted:" + urlParams["courseid"]))
			fmt.Println(http.StatusAccepted, urlParams["courseid"], "deleted.")
		}

	} else {
		if isUpdate {
			// write to the response
			// Write() takes in []byte as input
			// response code 202 means accepted
			res.Write([]byte("202 - Module updated:" + urlParams["moduleid"]))
			fmt.Println(http.StatusAccepted, urlParams["moduleid"], "updated.")
		} else {
			// write to the response
			// Write() takes in []byte as input
			// response code 202 means accepted
			res.Write([]byte("202 - Module deleted:" + urlParams["moduleid"]))
			fmt.Println(http.StatusAccepted, urlParams["moduleid"], "deleted.")
		}
	}
}

func unprocessableEntityStatus(res http.ResponseWriter, isCourse bool) {

	// write status code to header
	res.WriteHeader(http.StatusUnprocessableEntity)

	if isCourse {
		// write to the response
		// Write() takes in []byte as input
		// response code 422 means unprocessable entity
		res.Write([]byte("422 - Please supply course information with course code and course name in JSON format."))
		fmt.Println(http.StatusUnprocessableEntity, "Please supply course information with course code and course name in JSON format.")
	} else {
		// write to the response
		// Write() takes in []byte as input
		// response code 422 means unprocessable entity
		res.Write([]byte("422 - Please supply course information with module code and module name in JSON format."))
		fmt.Println(http.StatusUnprocessableEntity, "Please supply module information with module code and module name in JSON format.")
	}

}

func createdStatus(res http.ResponseWriter, urlParams map[string]string, isCourse bool) {

	// write status code to header
	res.WriteHeader(http.StatusCreated)

	if isCourse {
		// write to the response
		// Write() takes in []byte as input
		// response code 201 means created successfully
		res.Write([]byte("201 - Course added:" + urlParams["courseid"]))
		fmt.Println(http.StatusCreated, urlParams["courseid"], "added.")
	} else {
		// write to the response
		// Write() takes in []byte as input
		// response code 201 means created successfully
		res.Write([]byte("201 - Module added:" + urlParams["moduleid"]))
		fmt.Println(http.StatusCreated, urlParams["moduleid"], "added.")
	}

}

func conflictStatus(res http.ResponseWriter, isCourse bool) {
	// write status code to header
	res.WriteHeader(http.StatusConflict)

	if isCourse {
		// write to the response
		// Write() takes in []byte as input
		// response code 409 means conflict
		res.Write([]byte("409 - Duplicate course ID. Please use PUT method if you want to update the content."))
		fmt.Println(http.StatusConflict, "Duplicate course ID. Please use PUT method if you want to update the content.")
	} else {
		// write to the response
		// Write() takes in []byte as input
		// response code 409 means conflict
		res.Write([]byte("409 - Duplicate module ID. Please use PUT method if you want to update the content."))
		fmt.Println(http.StatusConflict, "Duplicate module ID. Please use PUT method if you want to update the content.")
	}
}

func userConflictStatus(res http.ResponseWriter) {
	// write status code to header
	res.WriteHeader(http.StatusConflict)

	// write to the response
	// Write() takes in []byte as input
	// response code 409 means conflict
	res.Write([]byte("409 - Duplicate user. You already have one api key."))
	fmt.Println(http.StatusConflict, "Duplicate user. You already have one api key.")
}

func userAPIProcessableEntityStatus(res http.ResponseWriter) {

	// write status code to header
	res.WriteHeader(http.StatusUnprocessableEntity)

	// write to the response
	// Write() takes in []byte as input
	// response code 422 means unprocessable entity
	res.Write([]byte("422 - Please supply user and key information in JSON format."))
	fmt.Println(http.StatusUnprocessableEntity, "Please supply user and key information in JSON format.")

}

func createdStatusKey(res http.ResponseWriter, user string) {

	// write status code to header
	res.WriteHeader(http.StatusCreated)

	res.Write([]byte("201 - Key added to " + user))
	fmt.Println(http.StatusCreated, "Key added to "+user)

}

func notFoundStatusUser(res http.ResponseWriter) {
	// write status code to header
	res.WriteHeader(http.StatusNotFound)

	// write to the response
	// Write() takes in []byte as input
	// response code 404 means not found
	res.Write([]byte("404 - User not found"))
	fmt.Println(http.StatusNotFound, "- User not found.")

}

func acceptedStatusKey(res http.ResponseWriter, user string, isUpdate bool) {
	// write status code to header
	res.WriteHeader(http.StatusAccepted)

	if isUpdate {
		// write to the response
		// Write() takes in []byte as input
		// response code 202 means accepted for processing
		res.Write([]byte("202 - Key updated to " + user))
		fmt.Println(http.StatusAccepted, "Key updated to ", user)
	} else {
		// write to the response
		// Write() takes in []byte as input
		// response code 202 means accepted for processing
		res.Write([]byte("202 -  Key deleted for " + user))
		fmt.Println(http.StatusAccepted, "Key deleted for ", user)
	}
}
