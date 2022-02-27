package main

/* Job Scrapper
1. getPages 함수 실행 (총 몇 페이지인지 알 수 있음)
2. 각 페이지 별로 getPage 실행
   => 모든 페이지를 동시에 요청하도록 goroutine 생성
3. getPage 실행되는 중에 extractJob도 실행됨 (각 페이지에는 50개의 일자리 정보 존재)
   => getPage안에 실행되는 extractJob 50개를 동시에 실행하도록 goroutine 생성
   => 요청된 정보는 channel 기능 이용하여 정보를 주고 받음 (main함수와)
4. getPage 실행 종료 후 main함수로 채널 전송

Challenge: 파일 쓰기 부분 goroutine으로
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

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	// 많은 배열들의 조합 [] + [] + [] ...
	var jobs []extractedjob
	// channel
	c := make(chan []extractedjob)

	/* scrapper project
	: indeed site -> /job/ search -> how many pages
	-> changes start time -> get first page  */
	// 1. go get github.com/PuerkitoBio/goquery
	// 2. get all indeed page
	totalPages := getPages()
	// fmt.Println(totalPages)
	// 3. go query : getPages()
	// 4. hit URL
	// 총 페이지 숫자를 확인한 후 각 페이지별로 getPage함수 호출
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
		// extractedJobs := getPage(i)
		// // ... : 하나의 배열로 합치는 것
		// jobs = append(jobs, extractedJobs...)
	}
	// 몇개의 goroutine을 기다려야 할까?
	// : 총 n개의 페이지 조회되면 n개의 goroutine
	for i := 0; i < totalPages; i++ {
		// 메시지 기다리는 중
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))

} // main 끝

// csv로 저장
func writeJobs(jobs []extractedjob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	// Flush 함수: 작성된 모든 것들을 파일에 입력
	// 일자리 정보가 writer로 전달되어서 저장하는 함수
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	// 모든 페이지에서 가져온 jobs가 입력, for가 끝나고 defer w.Flush() 실행
	// Challenge : goroutine 만들어보기
	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func cleanString(str string) string {
	// Fields: space를 기준으로 배열안에 텍스트를 넣음 -> 반환값: []string
	// Join: 배열을 가져와서 합치는 역할
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 각 페이지에 있는 일자리를 모두 반환하는 함수
// getPage는 중간다리 - goroutine을 생성해서 일자리를 전달하고 main함수의 채널로 전송
func getPage(page int, mainC chan<- []extractedjob) {
	var jobs []extractedjob
	// channel
	c := make(chan extractedjob)
	// 필요한 주소 만들기
	pageUrl := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageUrl)
	// 정보가져오는 요청
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	// res.Body - byte , input & output IO
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// job card에서 일자리 정보 가져오기
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		// id, _ := card.Attr("data-jk") // id 찾기
		// // fmt.Println(id)
		// title := cleanString(card.Find("h2>span").Text())
		// // fmt.Println(title)
		// location := cleanString(card.Find(".companyLocation").Text())
		// fmt.Println(id, title, location)
		// extractedJob struct를 반환 후 job변수에 저장
		go extractedJob(card, c) // 각각의 카드에 대해서 job을 추출
		// jobs = append(jobs, job)     // jobs에 저장
	})
	// extractJob goroutine에서 정보가 전달되면 jobs를 mainC로 전달
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	// mainC(채널)에는 []extractedjob가 입력됨
	mainC <- jobs
}

func extractedJob(card *goquery.Selection, c chan<- extractedjob) {
	id, _ := card.Attr("data-jk") // id 찾기
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

// 페이지의 총 숫자 찾는 함수
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	// res.Body - byte , input & output IO
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// 클래스 찾기 fmt.Println(doc) // &{0xc0004604e0 <nil> 0xc000206000}
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
