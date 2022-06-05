-- Local database definition.

DROP DATABASE IF EXISTS my_db;

CREATE DATABASE my_db;

USE my_db;

DROP TABLE IF EXISTS Course;
DROP TABLE IF EXISTS Module;
DROP TABLE IF EXISTS APIkey;

-- Create Course table
CREATE TABLE Course (
    id varchar (5) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(100)
);
-- Insert data to Course table
INSERT INTO Course (id, name, description) VALUES ("EEE", "Electrical and Electronic Engineering", "Teach about EEE");
INSERT INTO Course (id, name, description) VALUES ("ACC", "Accountancy", "Teach about ACC");
INSERT INTO Course (id, name, description) VALUES ("CS", "Computer Science", "Teach about CS");

-- Create Module table
CREATE TABLE Module (
    id varchar (10) NOT NULL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(50),
    course_id varchar (5) NOT NULL
);

-- Insert data to Module table
INSERT INTO Module (id, name, description, course_id) VALUES ("EEE2004", "Digital Electronics", "Teach about Digital Electronics.", "EEE");
INSERT INTO Module (id, name, description, course_id) VALUES ("EEE2002", "Analog Electronics", "Teach about Analog Electronics", "EEE");
INSERT INTO Module (id, name, description, course_id) VALUES ("CS1002", "Data Structure and Algorithms", "Teach about Data Structure and Algorithms.", "CS");
INSERT INTO Module (id, name, description, course_id) VALUES ("CS1008", "Golang", "Teach about Golang programming.", "CS");
INSERT INTO Module (id, name, description, course_id) VALUES ("ACC3010", "Accounting I", "Teach about Accounting I.", "ACC");
INSERT INTO Module (id, name, description, course_id) VALUES ("ACC3014", "Statistics and Analysis", "Teach about Statistics and Analysis.", "ACC");

-- Create APIkey table
CREATE TABLE APIkey (username varchar (10) NOT NULL PRIMARY KEY, apikey VARCHAR(100) NOT NULL);

-- Insert data to APIkey table
INSERT INTO APIkey (username, apikey) VALUES ("jpschew", "0fd5062eccc4b049f0ba75ca31db3a4cb12088bb9fe3addf33e9e2b481307aa2");

