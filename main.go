package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var LogStore Database

//LogLine represents fields in a given log line
type LogLine struct {
	RawLog        string
	RemoteAddr    string
	TimeLocal     string
	RequestType   string
	RequestPath   string
	Status        int
	BodyBytesSent int
	HTTPReferer   string
	HTTPUserAgent string
	Created       time.Time
}

//LogFile represents a logfile with multiple lines
type LogFile struct {
	Logs []LogLine
}

//Print prints the data for a log line
func (l *LogLine) Print() {
	fmt.Println(l.RemoteAddr)
	fmt.Println(l.TimeLocal)
	fmt.Println(l.RequestType)
	fmt.Println(l.RequestPath)
	fmt.Println(l.Status)
	fmt.Println(l.BodyBytesSent)
	fmt.Println(l.HTTPReferer)
	fmt.Println(l.HTTPUserAgent)
	fmt.Println(l.Created)
}

func main() {
	// lines := readFile("log_a.txt")
	lines := readFile("jumbo_log.txt")
	//fmt.Println(lines[0])
	logFile := parseFile(lines)

	//Start DB connection
	LogStore = Database{}
	var err error
	LogStore.db, err = sql.Open("sqlite3", "./ETL.db")
	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()

	//Create table if not found
	LogStore.dbInit()

	//Store parsed logs
	LogStore.StoreLogLine(logFile)

	countVisitors(LogStore)
	countBrowsers(LogStore)

	r := mux.NewRouter()
	r.HandleFunc("/browser/count", BasicAuth(handleBrowserCount, "harry", "potter", "Please enter your username and password for this site")).Methods("GET")
	r.HandleFunc("/visitor/count", BasicAuth(handleVisitorCount, "harry", "potter", "Please enter your username and password for this site")).Methods("GET")
	http.ListenAndServe(":8000", r)
}

//check if element found in golang string
func checkExists(list []string, v string) bool {
	for _, a := range list {
		if a == v {
			return true
		}
	}
	return false
}

func readFile(name string) []string {
	var lines []string

	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func parseFile(lines []string) LogFile {

	//list to store log lines
	//logLines := []LogLine{}
	lf := LogFile{}

	for _, line := range lines {

		lineSplit := strings.Split(line, " ")
		// fmt.Println(lineSplit[0])
		userAgent := strings.Join(lineSplit[11:], " ")
		status, _ := strconv.Atoi(lineSplit[8])
		totalBytes, _ := strconv.Atoi(lineSplit[9])
		tempLine := LogLine{
			line,
			lineSplit[0],
			lineSplit[3] + " " + lineSplit[4],
			lineSplit[5],
			lineSplit[6],
			status,
			totalBytes,
			lineSplit[10],
			userAgent,
			time.Now(),
		}

		// logLines = append(logLines, tempLine)
		lf.Logs = append(lf.Logs, tempLine)
	}
	return lf
}