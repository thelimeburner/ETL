package main

import (
	"database/sql"
	"fmt"
)

type BrowserCountRow struct {
	Key     string
	Date    string
	Browser string
	Count   int
}

type VisitorCountRow struct {
	Key   string
	Count int
}

type User struct {
	ID       int
	User     string
	Password string
	Perm     string
}

//Database controls database functionality
type Database struct {
	db *sql.DB
}

func (d *Database) fetchUserAuth(perm string) []User {
	var sqlStmt string
	if perm == "read" {
		sqlStmt = "SELECT * from users where permission = 'all' or permission = 'read'"
	} else if perm == "write" {
		sqlStmt = "SELECT * from users where permission = 'all' or permission = 'write'"
	} else {
		return nil
	}

	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		fmt.Println("Failed to fetch browser data: ", err)
	}

	users := []User{}

	for rows.Next() {
		u := User{}
		rows.Scan(&u.ID, &u.User, &u.Password, &u.Perm)
		//fmt.Println(bcr.Id, bcr.Key, bcr.Date, bcr.Browser, bcr.Count)
		users = append(users, u)
	}
	return users
}

//storeBrowserCount stores the count for browsers in logs by date
func (d *Database) storeBrowserCount(key string, dt string, b string, c int) bool {
	sqlStmt := `
	INSERT INTO browsers (
		key,
		date,
		browser,
		count
		) VALUES (?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//add admin
	_, err = statement.Exec(key, dt, b, c)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//storeVisitorCount stores the count for visitors in logs by date
func (d *Database) storeVisitorCount(key string, vc int) bool {
	sqlStmt := `
	INSERT INTO visitors (
		key,
		visitor_count
		) VALUES (?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//add admin
	_, err = statement.Exec(key, vc)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//Store stores a logfile in a database
func (d *Database) StoreLogLine(lf LogFile) {

	sqlStmt := `
	INSERT INTO logs (raw_log, 
		remote_addr,
		time_local,
		request_type,
		request_path,
		status,
		body_bytes_sent,
		http_referer,
		http_user_agent,
		created
		) VALUES (?, ?,?,?,?,?,?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range lf.Logs {

		statement.Exec(v.RawLog, v.RemoteAddr, v.TimeLocal, v.RequestType, v.RequestPath, v.Status, v.BodyBytesSent, v.HTTPReferer, v.HTTPUserAgent, v.Created)
	}
}

func (d *Database) fetchBrowserData() []BrowserCountRow {

	sqlStmt := `
		SELECT * from browsers
	`
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		fmt.Println("Failed to fetch browser data: ", err)
	}

	browserStats := []BrowserCountRow{}

	for rows.Next() {
		bcr := BrowserCountRow{}
		rows.Scan(&bcr.Key, &bcr.Date, &bcr.Browser, &bcr.Count)
		//fmt.Println(bcr.Id, bcr.Key, bcr.Date, bcr.Browser, bcr.Count)
		browserStats = append(browserStats, bcr)
	}
	return browserStats
}

func (d *Database) fetchVisitorData() []VisitorCountRow {

	sqlStmt := `
		SELECT * from visitors
	`
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		fmt.Println("Failed to fetch browser data: ", err)
	}

	visitorStats := []VisitorCountRow{}

	for rows.Next() {
		vcr := VisitorCountRow{}
		rows.Scan(&vcr.Key, &vcr.Count)
		//fmt.Println(bcr.Id, bcr.Key, bcr.Date, bcr.Browser, bcr.Count)
		visitorStats = append(visitorStats, vcr)
	}
	return visitorStats
}

func (d *Database) fetchData() LogFile {
	rows, _ := d.db.Query("SELECT * FROM logs")
	lf := LogFile{}
	for rows.Next() {
		logLine := LogLine{}
		rows.Scan(&logLine.RawLog,
			&logLine.RemoteAddr,
			&logLine.TimeLocal,
			&logLine.RequestType,
			&logLine.RequestPath,
			&logLine.Status,
			&logLine.BodyBytesSent,
			&logLine.HTTPReferer,
			&logLine.HTTPUserAgent,
			&logLine.Created)
		lf.Logs = append(lf.Logs, logLine)
	}
	return lf
}

//dbInit sets up the db
func (d *Database) dbInit() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS logs (
		raw_log TEXT NOT NULL UNIQUE,
		remote_addr TEXT,
		time_local TEXT,
		request_type TEXT,
		request_path TEXT,
		status INTEGER,
		body_bytes_sent INTEGER,
		http_referer TEXT,
		http_user_agent TEXT,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()

	//create user table
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS users (
		id int NOT NULL UNIQUE,
		user TEXT,
		password TEXT,
		permission TEXT
		)
	`

	statement, _ = d.db.Prepare(sqlStmt)
	statement.Exec()

	sqlStmt = `
	INSERT INTO users (id, 
		user,
		password,
		permission
		) VALUES (?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		fmt.Println(err)
		return
	}

	//add admin
	statement.Exec(1, "admin", "password", "all")

	//add user
	statement.Exec(2, "user1", "1234", "read")

	//create user table
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS visitors (
		key TEXT,
		visitor_count int
		)
	`

	statement, _ = d.db.Prepare(sqlStmt)
	statement.Exec()

	//create browser table
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS browsers (
		key TEXT,
		date TEXT,
		browser TEXT,
		count int
		)
	`

	statement, _ = d.db.Prepare(sqlStmt)
	statement.Exec()

	//create browser table
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS lineCount (
		key TEXT,
		count int
		)
	`

	statement, _ = d.db.Prepare(sqlStmt)
	statement.Exec()

}
