package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	// "github.com/jhee086/learngo/accounts"
	"github.com/jhee086/learngo/mydict"
)

type requestResult struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request failed")

// 컴파일을 원하면 main
func main() {
	// nico := person("nico", 18,)
	// // P 가 대문자인 이유 :
	// fmt.Println("Hello World") // formatting
	// // 대문자로된 func만 다른 패키지로부터 export 된 것
	// something.SayHello() // private func이라 실행 안 됨 -something.sayBye()

	// type을 찾아 정해줌
	// name := "jhshim" // 같음 var name string = "jhshim"
	// name = "dfsaf"
	// fmt.Println(name)
	// fmt.Println(multiply(2, 2))
	// fmt.Println(lenAndUpper("jhshim"))

	// _, upperName := lenAndUpper("jhsim")
	// totalLenght, _ 하나만 반환하고 싶으면 _ 처리
	// fmt.Println(upperName)

	// repeatMe("dkj", "sdf", "sdfagw")

	// result := supperAdd(1, 2, 3, 4, 5)
	// fmt.Println(result)

	// fmt.Println(canIDrink(18))

	// a := 2
	// b := &a // a가 바뀌면 b도 계속 바뀜
	// *b = 20
	// // a = 10
	// fmt.Println(a)
	/*
		Go에서 low level programming
		1. &a : a의 메모리 주소 address
		2. b := &a : 주소를 저장
		3. *b : 메모리 주소를 통해 값을 살펴보기 see through (pointer로 부름)
		4. *b = 20 -> fmt.Println(a) : 주소에 담긴 값을 변경
	*/

	// slice : length가 없는 array []string
	// names := []string{"jh", "nico", "lynn"}
	// names[3] = "lalala"
	// names[4] = "lalala"

	// append 하나의 slice와 값을 받아서 새로운 slice return
	// names = append(names, "kla")
	// fmt.Println(names)

	// 데이터 구조 map: map[key type]value type
	// nico_map := map[string]string{"name": "nico", "age": "12"}
	// for key, _ := range nico_map {
	// 	fmt.Println(key)
	// }
	// fmt.Println(nico)

	// // struct : 구조체 (class+object)
	// favFood := []string{"kimchi", "ramen"}
	// // nico_struct := person{"nico", 18, favFood} // 보기 불편
	// // field: value 사용시 value만 사용 불가능
	// nico_struct := person{name: "noico", age: 18, favFood: favFood}
	// fmt.Println(nico_struct)

	/* Go struct에는 constructor가 없다.
	Python: __init__
	JS: constructor() */

	/* 2. Bank Account - public method vs private method */
	// account := banking.Account{Owner: "jhshim", Balance: 0}
	// account := accounts.NewAccount("jhshim")
	// fmt.Println(account) // &{jhshim 0} : addresss - object
	// account.Deposit(10)
	// fmt.Println(account.Balance()) // 0 출력
	// err := account.Withdraw(20)
	// if err != nil {
	// 	// log.Fatalln(err) // Println 실행하고 프로그램 종료시킴
	// 	fmt.Println(err)
	// }
	// account.ChangeOwner("jhs")
	// fmt.Println(account.Balance(), account.Owner())

	/*
		struct가 가지고 있는 method
		struct에서 (자동으로) 호출하는 method
	*/
	// fmt.Println(account)

	/* 1. dictionary (map을 이용한) 시물레이션 - method를 type에도 추가가능*/
	dictionary := mydict.Dictionary{} // 빈 dictionary
	// dictionary["hello"] = "hello" fmt.Println(dictionary)
	// fmt.Println(dictionary["first"]) // dictionary := mydict.Dictionary{"first": "First word"}
	// definition, err := dictionary.Search("second") //("first")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(definition)
	// }
	baseword := "hello"
	// definition := "Greeting"
	// err := dictionary.Add(word, definition)
	dictionary.Add(baseword, "First")
	dictionary.Search(baseword)
	dictionary.Delete(baseword)
	// word, err := dictionary.Search(baseword)
	// err := dictionary.Update("asggsdg", "Second")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(word, dictionary) // Not Found : 삭제 후 찾는 것 불가능
	// }
	// helllo, _ := dictionary.Search(word)
	// fmt.Println("found", word, "difinition:", helllo)
	// err2 := dictionary.Add(word, definition)
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }

	/* 3.URL checker -> scraper */
	// url과 결과값 사이의 map
	// 해결방안1 (empty map 초기화)
	// var results = map[string]string{}
	// 해결방안2 (empty map 초기화)
	results := make(map[string]string)
	c := make(chan requestResult) // struct로 보내기
	// 배열형태의 URL 받아와 체크 -> 정상적으로 나오면 UP / 성공 혹은 실패 표시
	// 이걸이제 최적화하자 -> 순서대로 말고 동시에 처리하자
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://academy.nomadcoders.co/",
	}
	// results["hello"] = "HELLO"
	/* panic: assignment to entry in nil map
	   -> 컴파일러가 못 찾아냄, 컴파일러는 모르는 에러
	   1. results라는 map(type)정의
	   2. 값을 넣으려고 write
	!! 원인: 초기화되지 않은 map에 값을 넣을 수 없음
	=> (empty map 초기화) 해야함
	*/
	for _, url := range urls {
		go hitURL(url, c)
	}
	// 메시지 기다리기
	for i := 0; i < len(urls); i++ {
		// fmt.Println(<-c) // concurrency 동시성 때문에 반환 순서 다름
		result := <-c
		results[result.url] = result.status
	}
	for url, status := range results {
		fmt.Println(url, status)
	}

	// 중요!! channels & go routines
	/* Goroutines: 다른 함수와 동시에 실행시키는 함수
	   1. 프로그램이 작동하는 동안만 유효 -> 메인 함수가 실행되는 동안만!
	   => 아래의 경우는 메인 함수가 sexyCount("flynn")
	   => 둘 다 go를 붙이면 메인 함수에서 실행하는 것이 없어 종료됨
	   2. Main 함수와 goroutines 사이 정보를 주고 받는 방법
	   - 일반 Top-down 방식 경우, 순서대로 실행 후 ok or failed 결과를 Main에 보내줌
	*/

	/*channel: pipe*/
	// c := make(chan bool)
	// c := make(chan string)
	// // sexyCount("nico") // 일반 Top-down 방식
	// people := [3]string{"nico", "flynn", "jhshim"}
	// for _, person := range people {
	// 	// 이 함수 작업이 끝나면 true 값을 channel을 통해서 전송
	// 	// goroutine으로부터 return받는 것 대신 channel 통해 메시지 전송
	// 	go isSexy(person, c) // 두 개의 인수를 받음
	// }
	// for i := 0; i < len(people); i++ {
	// 	fmt.Print("Waiting for ", i, "\n")
	// 	fmt.Println(<-c) // concurrency 동시성 때문에 반환 순서 다름
	// }
	// resultOne := <-c
	// resultTwo := <-c
	// fmt.Println(resultOne) // 메시지를 받는 것 blocking operation
	// fmt.Println(resultTwo)
	// fmt.Println(<-c) // 배열보다 반환 많이하면 deadlock!
	// fmt.Println(result) // 이걸로 두개하면 하나의 반환값만 두번 출력됨
	// go sexyCount("nico") // go 추가 시 동시다발
	// go sexyCount("flynn")
	// 5초간 sleep: goroutines는 5초간 살아있음 그 이후 메인 함수 종료됨
	// time.Sleep(time.Second * 5)

} // main 끝

func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 5)
	// fmt.Println(person)
	c <- person + " is sexy"
}

func sexyCount(person string) {
	for i := 0; i < 10; i++ {
		fmt.Println(person, "is sexy", i)
		time.Sleep(time.Second) // 1초 동안 sleep - time은 GO 패키지
	}
}

// function
// 웹 사이트로 접속(hit: 인터넷 웹 서버의 파일 1개에 접속)하고 그 결과를 알려줌

// hitURL로 채널 보내기
// chan<- : this channel send only (direction 설정)
func hitURL(url string, c chan<- requestResult) {
	// fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	status := "OK"
	// 조건: err가 존재할 때
	if err != nil || resp.StatusCode >= 400 {
		status = "Failed"
	}
	c <- requestResult{url: url, status: status}
}

type person struct {
	name    string
	age     int
	favFood []string
}

func canIDrink(age int) bool {

	switch koreanAge := age + 2; koreanAge {
	//case age < 18:
	//	return false
	case 18:
		return true
	case 50:
		return false
	}
	return false
	//// if 문
	// if koreanAge := age - 2; koreanAge < 18 {
	// 	return false
	// } else {
	// 	return true
	// }
}

// loop 문: for 사용으로 가능
// range를 사용해서 loop 만들기
func supperAdd(numbers ...int) int {
	// for index, number := range numbers {
	// 	fmt.Println(index, number)
	// }
	// for i := 0; i < len(numbers); i++ {
	// 	fmt.Println(numbers[i])
	// }
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

// naked return : return 할 variable 명시하지 않아도 됨
// defer: func이 끝났을 때 추가적으로 하는 것을 정함 -> return 후 실행
func lenAndUpper(name string) (lenght int, uppercase string) {
	//
	defer fmt.Println("I'm done")
	lenght = len(name) // 선언된 걸 적용할 때는 = , 다시 생성할 때는 :=
	uppercase = strings.ToUpper(name)
	return
}

// func repeatMe(words ...string) {
// 	fmt.Println(words) // array 형태
// }

// 여러개 리턴 값 가질 수 있음
// func lenAndUpper(name string) (int, string) {
// 	return len(name), strings.ToUpper(name)
// }

// func multiply(a, b int) int {
// 	return a * b
// }
