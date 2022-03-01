package main

import (
	"os"
	"strings"

	"github.com/jhee086/learngo/scrapper"
	"github.com/labstack/echo" // install: go get github.com/labstack/echo/v4
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, World!")
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	// 사용자가 파일을 다운로드하면 서버에서 파일 삭제
	defer os.Remove(fileName)

	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	// 첨부파일을 리턴하는 기능, 사용자에게 전달할 파일 이름은 job.csv
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323")) // http://localhost:1323/
	// scrapper.Scrape("term")

}
