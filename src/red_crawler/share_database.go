
package red_crawler

import (
	"strings"
	"io/ioutil"
	"database/sql"
    _"github.com/go-sql-driver/mysql"
)


//===========================================================================


// File containing user and password for MySQL
const MYSQL_SETTINGS = "mysql_settings"
// Database used by MySQL
const DATABASE_NAME = "top_ten_shares"
// #REFACTOR: the table's name should be a parameter, currently it is 'shares'

// Global Database
// #REFACTOR: should be treated as an object and it's functions as methods 
var db *sql.DB = nil

//===========================================================================


// Reads user and password for MySQL environment from file and returns the login format
func mysql_login_info() string {
	f_mysql, err := ioutil.ReadFile(MYSQL_SETTINGS)
	if err != nil { panic("Invalid login and password formatting for MySQL") }
	login := strings.TrimSpace(string(f_mysql))
	
	return strings.Replace(login, "\n", ":", 1)
}


// Opens or creates the database if it doesnt exist
func open_or_create() {
	
	var err error

	// Login info formatted for openning the database
	login := mysql_login_info()

	// Open without the database
	db, err = sql.Open("mysql", login+"@tcp(127.0.0.1:3306)/")
	if err != nil { panic("Cannot open MySQL") }

	// If it doesnt exist, create, then close
	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DATABASE_NAME)
	if err != nil { panic(err) }
	db.Close()

	// Open again, now using the database
	db, err = sql.Open("mysql", login+"@tcp(127.0.0.1:3306)/" + DATABASE_NAME)
	if err != nil { panic(err) }
}


// Deletes the older table with outdated information, if it exists and creates a new empty table
func delete_old_and_create_new_table() {

	// Create only to delete
	// #REFACTOR: Must be a better way to check if it exists instead
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS shares(
	share_name VARCHAR(10) NOT NULL,
    company_name VARCHAR(100),
    market_value BIGINT,
    daily_fluctuation FLOAT,
    PRIMARY KEY (share_name));`)
	if err != nil { panic(err) }

	// Delete the old table for updating
	_, err = db.Exec("DROP TABLE shares;")
	if err != nil { panic(err) }

	// Create empty table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS shares(
	share_name VARCHAR(10) NOT NULL,
    company_name VARCHAR(100),
    market_value BIGINT,
    daily_fluctuation FLOAT,
    PRIMARY KEY (share_name));`)
	if err != nil { panic(err) }
}


// Set UTF8MB4 for company_name variable for compatibility with information read on Fundamentus
func configure_character_set() {

	_,err := db.Exec("ALTER DATABASE " + DATABASE_NAME + " CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;")
	if err != nil { panic(err) }
	_,err = db.Exec("ALTER TABLE shares CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	if err != nil { panic(err) }
	_,err = db.Exec("ALTER TABLE shares CHANGE company_name company_name VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	if err != nil { panic(err) }
}


// Populates the database
func insert_rows (top []Paper) {

	row, err := db.Prepare("INSERT INTO shares VALUES( ?, ?, ?, ? )")
	if err != nil { panic(err) }
	
	for i := range top {
		_,err = row.Exec(top[i].Share, top[i].Company, top[i].Value, top[i].Fluctuation)
		if err != nil { panic(err) }
	}
}


// Creates a database and inserts the itens
func DB_persist(top []Paper) {

	open_or_create()

	delete_old_and_create_new_table()

	configure_character_set()

	insert_rows(top)

	db.Close()
}


//===========================================================================
