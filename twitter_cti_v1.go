package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

var wg sync.WaitGroup

func outFile(str string, ioc_array []string) {
	t_time := time.Now()
	fo, e := os.OpenFile(str+t_time.String()+".txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if e != nil {
		panic(e)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	for _, i := range ioc_array {
		if i != "" {
			_, err2 := fo.WriteString(i + "\n")
			if err2 != nil {
				panic(err2)
			}
		}
	}
}

func main() {

	start := time.Now()

	no_tweets := flag.Int("n", 50, "Number of tweets to look up max value ~3200")
	path := flag.String("p", "", "Specify the path to create the output file")
	hashtag := flag.String("H", "#opendir c2 #malware", "Enter a hashtag or keyword to search eg. #emotet, c2 ...")
	flag.Parse()

	ntweets := *no_tweets
	p_path := *path
	Hashtag := *hashtag

	fmt.Println("\nSearching Twitter For", Hashtag, "You Can Specify The keywords Using -H Option")

	if p_path == "" {
		fmt.Println("\nWarning: File Path Is Not Specified,No Output File Will Be Created\n\nUse -h or --help \n")

	}

	ip_array := []string{}
	deip_array := []string{}
	hash_array := []string{}
	url_array := []string{}
	url_defang_array := []string{}

	ip_defang_re := regexp.MustCompile("(\\d.*)([[].][\\d+])([.\\d.+]{0,3}){0,3}")
	ip_re := regexp.MustCompile("(\\d+).([.]\\d{1,3}){3}")
	hash_re := regexp.MustCompile("([a-z]|[0-9]|[A-Z]){32,128}\\s")
	url_re := regexp.MustCompile("(h..ps?:)\\/\\/[-a-zA-Z0-9@:%._\\+~#=]{1,256}[^t.co]\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&\\/\\/=]*)?")
	url_defang_re := regexp.MustCompile("(h..ps?:)\\/\\/[-a-zA-Z0-9@:%._\\+~#=]{1,256}[^t.co]\\[[.]][a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&\\/\\/=]*)?")

	wg.Add(4)

	tweetStream := make(chan string)
	go func() {
		defer close(tweetStream)
		defer wg.Done()

		scraper := twitterscraper.New()
		for tweet := range scraper.SearchTweets(context.Background(),
			*hashtag+"-filter:retweets", ntweets) {
			if tweet.Error != nil {
				panic(tweet.Error)
			}
			tweetStream <- tweet.Text
		}

	}()
	go func() {
		defer wg.Done()
		for t := range tweetStream {

			ioc_ip_defang := ip_defang_re.FindString(t)
			ioc_ip := ip_re.FindString(t)
			if ioc_ip == "" && ioc_ip_defang == "" {
				continue
			} else {
				ip_array = append(ip_array, ioc_ip)
				deip_array = append(deip_array, ioc_ip_defang)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for t := range tweetStream {

			ioc_hash := hash_re.FindString(t)

			if ioc_hash == "" {
				continue
			} else {
				hash_array = append(hash_array, ioc_hash)

			}
		}
	}()
	go func() {
		defer wg.Done()
		for t := range tweetStream {

			ioc_url_defang := url_defang_re.FindString(t)
			ioc_url := url_re.FindString(t)
			if ioc_url == "" && ioc_url_defang == "" {
				continue
			} else {
				url_array = append(url_array, ioc_url)
				url_defang_array = append(url_defang_array, ioc_url_defang)

			}
		}
	}()

	wg.Wait()

	if p_path != "" {
		outFile(p_path, ip_array)
		time.Sleep(1 * time.Second)
		outFile(p_path, deip_array)
		time.Sleep(1 * time.Second)
		outFile(p_path, hash_array)
		time.Sleep(1 * time.Second)
		outFile(p_path, url_array)
		time.Sleep(1 * time.Second)
		outFile(p_path, url_defang_array)
		time.Sleep(1 * time.Second)
	} else {
		fmt.Println(ip_array)
		fmt.Println(deip_array)
		fmt.Println(hash_array)
		fmt.Println(url_array)
		fmt.Println(url_defang_array)
	}
	fmt.Println("This operation took,", time.Since(start))
}
