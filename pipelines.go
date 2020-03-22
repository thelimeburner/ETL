package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
)

//processLogFile takes in an uploaded logfile, stores the data, processes stats.
func processLogFile(rawLogFile []byte) bool {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(rawLogFile)))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	logFile := parseFile(lines)

	//Store parsed logs
	LogStore.StoreLogLine(logFile)

	countVisitors(LogStore)
	countBrowsers(LogStore)

	return true
}

//parseBrowser parses user agent
func parseBrowser(ua string) string {
	browsers := []string{"Firefox", "Chrome", "Opera", "Safari", "MSIE"}

	for _, b := range browsers {
		if strings.Contains(ua, b) {
			return b
		}
	}
	return "Other"
}

//format time appropriately
func parseTime(t string) time.Time {
	//January 2nd, 3:04:05 PM of 2006, UTC-0700.
	tm, err := time.Parse("02/Jan/2006:15:04:05 -0700", t)
	if err != nil {
		log.Fatal("Failed to parse time", err)
	}
	return tm
}

//countBrowsers counts browsers by day
func countBrowsers(LogStore Database) {
	visitTimes := []time.Time{}
	browsers := []string{}

	lf := LogStore.fetchData()
	for _, v := range lf.Logs {
		browsers = append(browsers, parseBrowser(v.HTTPUserAgent))
		v.TimeLocal = strings.Replace(v.TimeLocal, "[", "", -1)
		v.TimeLocal = strings.Replace(v.TimeLocal, "]", "", -1)
		visitTimes = append(visitTimes, parseTime(v.TimeLocal))
	}
	//place to store IP's
	browserCounts := make(map[string]int)

	//iterate and store IP addresses uniquely.
	for index, b := range browsers {
		key := b + "_" + visitTimes[index].Format("02-01-2006")
		if _, ok := browserCounts[key]; !ok {
			browserCounts[key] = 0
		}
		browserCounts[key]++
	}

	for b, v := range browserCounts {

		keyPieces := strings.Split(b, "_")
		if !LogStore.storeBrowserCount(b, keyPieces[1], keyPieces[0], v) {
			fmt.Println("Failed to store ", b)
		}
	}

}

//countVisitors counts visitors by day
func countVisitors(LogStore Database) {
	ips := []string{}
	visitTimes := []time.Time{}

	lf := LogStore.fetchData()
	for _, v := range lf.Logs {
		ips = append(ips, v.RemoteAddr)
		v.TimeLocal = strings.Replace(v.TimeLocal, "[", "", -1)
		v.TimeLocal = strings.Replace(v.TimeLocal, "]", "", -1)
		visitTimes = append(visitTimes, parseTime(v.TimeLocal))
	}
	//place to store IP's
	uniqueIPs := make(map[string][]string)

	//iterate and store IP addresses uniquely.
	for index, ipAddr := range ips {
		key := visitTimes[index].Format("02-01-2006")

		if _, ok := uniqueIPs[key]; !ok {
			uniqueIPs[key] = []string{}
		}
		if !checkExists(uniqueIPs[key], ipAddr) {
			uniqueIPs[key] = append(uniqueIPs[key], ipAddr)
		}
	}

	visitorCount := make(map[string]int)
	//count each of the days
	for k, v := range uniqueIPs {
		visitorCount[k] = len(v)
	}

	//print the counts
	for k, v := range visitorCount {
		//fmt.Println(k, ": ", v)
		if !LogStore.storeVisitorCount(k, v) {
			fmt.Println("Failed to store ", k)
		}
	}

	//write to database

}
