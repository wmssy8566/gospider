package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	CreateWrite()

	//NormalStart()
	ChannelStart(1, 101)

}

func Spider(url string, ch chan bool, i int) {
	file, err := os.OpenFile("test.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	url = url + strconv.Itoa(i) + "/comments"
	fmt.Println(url)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", "Windows")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	infoRe, _ := regexp.Compile("\\w+@\\w.*\\b")
	info5 := infoRe.FindAll(body, 5)

	for j := 0; j < 5; j++ {

		_, err = file.WriteString(string(info5[j]) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	if ch != nil {
		ch <- true
	}

}

func CreateWrite() {
	file, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func ChannelStart(stnum, endnum int) {
	start := time.Now()
	ch := make(chan bool)
	url := "https://jsonplaceholder.typicode.com/posts/"
	for i := stnum; i < endnum; i++ {
		go Spider(url, ch, i)
	}
	for i := stnum; i < endnum; i++ {
		<-ch
	}
	elapsed := time.Since(start)
	fmt.Printf("ChannelStart Time %s \n", elapsed)
}

func NormalStart(stnum, endnum int) {
	start := time.Now()
	url := "https://jsonplaceholder.typicode.com/posts/"
	for i := stnum; i < endnum; i++ {
		Spider(url, nil, i)
	}
	elapsed := time.Since(start)
	fmt.Printf("ChannelStart Time %s \n", elapsed)
}
