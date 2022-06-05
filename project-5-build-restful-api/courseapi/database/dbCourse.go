package database

import (
	"database/sql"
	"fmt"
	"log"

	// implementation of GO's database/sql/driver interface
	// only need to import and use the GO's database/sql API
	_ "github.com/go-sql-driver/mysql"
)

// Course added take in course data from db
type Course struct { // map this type to the record in the table
	CourseName    string            `json:"courseName"`
	Description   string            `json:"description"`
	CourseCode    string            `json:"courseCode"`
	CourseModules map[string]Module `json:"courseMod"`
}

// Module added take in module data from db
type Module struct { // map this type to the record in the table
	ModuleName  string `json:"moduleName"`
	ModuleCode  string `json:"moduleCode"`
	Description string `json:"description"`
}

// ModuleDB added take in module data from db to handle null string if any.
type ModuleDB struct { // map this type to the record in the table
	ModuleCode  sql.NullString
	ModuleName  sql.NullString
	Description sql.NullString
}

type DBRecord struct {
	//CourseCode  string
	CourseName   string
	CDescription string
	//ModuleCode  string
	ModuleName   string
	MDescription string
}

type Key struct {
	UserName string `json:"user"`
	APIKey   string `json:"apiKey"`
}

//var data map[string]map[string]DBRecord
var data map[string]Course
var keyData map[string]string

// CreateDBConn creates a connection to mysql database given the driver name, dsn and db name.
func CreateDBConn(driver string, dsn string, dbName string) *sql.DB {

	// Use mysql as driver Name and a valid DSN as data SourceName:
	source := dsn + "/" + dbName
	db, err := sql.Open(driver, source)

	// handle error
	if err != nil {
		//panic(err.Error())
		log.Panicln(err.Error())
	}

	return db
}

// GetAllRecords retrieves data from Course and Module table and assign to the map for usage.
func GetAllRecords(db *sql.DB) map[string]Course {
	query := `
			SELECT c.id, c.name, c.description, m.id, m.name, m.description
			FROM my_db.Course c
			LEFT JOIN my_db.Module m
			ON c.id = m.course_id
			`
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// map course code to course information
	// and pass to the api for usage
	data = make(map[string]Course)
	// map the course and module type to the record in the table
	var course Course
	var module ModuleDB
	//var moduleDB ModuleDB

	//fmt.Println("Database start.")
	// results.Next() means there is a next result
	for results.Next() {

		// Scan() copy each row of data from db and assign to the address specified
		err = results.Scan(&course.CourseCode, &course.CourseName, &course.Description, &module.ModuleCode, &module.ModuleName, &module.Description)
		if err != nil {
			//fmt.Println(course.CourseCode, course.CourseName, course.Description, module.ModuleCode, module.ModuleName, module.Description)
			panic(err.Error())

		}
		//fmt.Println(course.CourseCode, course.CourseName, course.Description, module.ModuleCode.String, module.ModuleName.String, module.Description.String)

		if _, ok := data[course.CourseCode]; !ok { // course code not in course map
			//var newMod map[string]Module
			newMod := make(map[string]Module)
			data[course.CourseCode] = Course{course.CourseCode, course.CourseName,
				course.Description, newMod}
			if module.ModuleName.String != "" && module.ModuleCode.String != "" {
				newMod[module.ModuleCode.String] = Module{module.ModuleCode.String,
					module.ModuleName.String, module.Description.String}
			}

		} else {
			editMod := data[course.CourseCode].CourseModules
			if module.ModuleName.String != "" && module.ModuleCode.String != "" {
				editMod[module.ModuleCode.String] = Module{module.ModuleCode.String,
					module.ModuleName.String, module.Description.String}
			}
			//editMod[module.ModuleCode] = Module{module.ModuleCode, module.ModuleName, module.Description}
		}
	}
	//fmt.Println("Database end.")
	return data
}

// DeleteCourse deletes data from Course table in database based on the course code.
func DeleteCourse(db *sql.DB, courseCode string) {

	query := fmt.Sprintf("DELETE FROM Course WHERE id='%s'", courseCode)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	DeleteModule(db, courseCode, true)

}

// DeleteModule deletes data from Module table in database based on the module code.
func DeleteModule(db *sql.DB, moduleCode string, delCourse bool) {

	var query string

	if delCourse {
		query = fmt.Sprintf("DELETE FROM Module WHERE course_id='%s'", moduleCode)
	} else {
		query = fmt.Sprintf("DELETE FROM Module WHERE id='%s'", moduleCode)
	}
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// InsertCourse adds data to Course table in database.
func InsertCourse(db *sql.DB, course Course) {
	query := fmt.Sprintf(`
						INSERT INTO Course (id, name, description) VALUES
						('%s', '%s', '%s')
						`, course.CourseCode, course.CourseName, course.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

}

// InsertModule adds data to Module table in database.
func InsertModule(db *sql.DB, module Module, courseCode string) {

	query := fmt.Sprintf(`
						INSERT INTO Module (id, name, description, course_id) VALUES 
						('%s', '%s', '%s', '%s')
						`, module.ModuleCode, module.ModuleName, module.Description, courseCode)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// EditCourse edits data to Course table in database.
func EditCourse(db *sql.DB, course Course) {

	query := fmt.Sprintf("UPDATE Course SET name='%s', description='%s' WHERE id='%s'",
		course.CourseName, course.Description, course.CourseCode)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// EditModule edits data to Module table in database.
func EditModule(db *sql.DB, module Module, courseCode string) {
	query := fmt.Sprintf("UPDATE Module SET name='%s', description='%s', course_id='%s' WHERE id='%s'",
		module.ModuleName, module.Description, courseCode, module.ModuleCode)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
