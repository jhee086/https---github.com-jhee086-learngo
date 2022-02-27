package main

/*
Job Scrapper
1. Run the getPages function (you can find out how many pages there are in total)
2. Execute getPage for each page
   => Create a goroutine to request all pages at the same time
3. While getPage is running, extractJob is also executed (each page has 50 job information)
    => Create goroutine to simultaneously execute 50 extractJobs executed in getPage
    => The requested information is exchanged using the channel function (with the main function)
4. After the execution of getPage is finished, the channel is sent to the main function.
*/

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedjob struct {
	id       string
	title    string
	location string
	// salary   string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=data&limit=50"

func main() {
	// combination of arrays [] + [] + [] ...
	var jobs []extractedjob
	// channel
	c := make(chan []extractedjob)

	/* scrapper project
	: indeed site -> /job/ search -> how many pages -> changes start time -> get first page */
	// 1. go get github.com/PuerkitoBio/goquery
	// 2. get all indeed page
	totalPages := getPages()
	// fmt.Println(totalPages)
	// 3. go query : getPages()
	// 4. hit URL
	// After checking the total number of pages, call the getPage function for each page
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
		// extractedJobs := getPage(i)
		// // ... : merging into one array
		// jobs = append(jobs, extractedJobs...)
	}
	// How many goroutines should we wait for? -> When a total of n pages are viewed, n goroutines
	for i := 0; i < totalPages; i++ {
		// waiting for message
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))

} // end of main

// // save as csv
// func writeJobs(jobs []extractedjob) {
// 	file, err := os.Create("jobs.csv")
// 	checkErr(err)
// 	w := csv.NewWriter(file)
// 	// Flush function: put everything written to a file
// 	// A function that stores job information by passing it to the writer
// 	defer w.Flush()

// 	headers := []string{"Link", "Title", "Location", "Summary"}
// 	wErr := w.Write(headers)
// 	checkErr(wErr)

// 	// All pages imported jobs are entered. Execute defer w.Flush() after for.
// 	for _, job := range jobs {
// 		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.summary}
// 		jwErr := w.Write(jobSlice)
// 		checkErr(jwErr)
// 	}
// }

// Challenge -> make a goroutine (Is it more fast?)
func writeJobs(jobs []extractedjob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	c := make(chan []string)
	for _, job := range jobs {
		go func(job extractedjob, c chan<- []string) {
			c <- []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.summary}
		}(job, c)
	}
	for i := 0; i < len(jobs); i++ {
		jobData := <-c
		writeErr := w.Write(jobData)
		checkErr(writeErr)
	}
}

func cleanString(str string) string {
	// Fields: Inserts text into an array based on space -> Return value: []string
	// Join: Taking an array and concatenating it
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// Function to return all jobs on each page
// getPage(an intermediate bridge) - make a goroutine, passes the jobs, and sends jobs to the main function's channel.
func getPage(page int, mainC chan<- []extractedjob) {
	var jobs []extractedjob
	// channel
	c := make(chan extractedjob)
	// create a pageUrl
	pageUrl := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageUrl)
	// Request to get information
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	// res.Body - byte , input & output IO
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// Get job information from job card
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		// id, _ := card.Attr("data-jk") // find a id
		// // fmt.Println(id)
		// title := cleanString(card.Find("h2>span").Text())
		// // fmt.Println(title)
		// location := cleanString(card.Find(".companyLocation").Text())
		// fmt.Println(id, title, location)
		// After returning the extractedJob struct, it is stored in the job variable.
		go extractedJob(card, c) // Extract a job for each card
		// jobs = append(jobs, job)     // save to jobs
	})
	// Passing jobs to mainC when information is passed from extractJob goroutine
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	// []extractedjob is entered in mainC (channel)
	mainC <- jobs
}

func extractedJob(card *goquery.Selection, c chan<- extractedjob) {
	id, _ := card.Attr("data-jk") // find a id
	title := cleanString(card.Find("h2>span").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	// salary :=
	summary := cleanString(card.Find(".job-snippet").Text())
	c <- extractedjob{
		id:       id,
		title:    title,
		location: location,
		// salary:   salary,
		summary: summary,
	}
}

// Function to find the total number of pages
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	// res.Body - byte , input & output IO
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// find a class,  fmt.Println(doc) // &{0xc0004604e0 <nil> 0xc000206000}
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Html())
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with statuscode")
	}
}
