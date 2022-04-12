package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	// "github.com/valyala/fasthttp/fasthttpproxy"
)

var (
	threads      int
	client       fasthttp.Client
	target       string
	sessionSlice []string
	running      bool = true
	mutex        sync.Mutex
	ratelimit    int
	attempts     int
	per_sec      int
)

func main() {
	fmt.Printf("[\x1B[1;36m!\x1B[0m] Instagram Turbo\n") // Originally was "Divine Instagram Turbo" as this build was for Divine.
	fmt.Printf("[\x1B[1;36m!\x1B[0m] Routines: ")
	fmt.Scanln(&threads)
	client = fasthttp.Client{
		//Dial: fasthttpproxy.FasthttpHTTPDialerTimeout("infproxy_kneb:EqGXrKkivHzLv0Uz@resi.infiniteproxies.com:1111", time.Duration(threads+20)), add your own proxies in here if you would like to use them
	}
	fmt.Printf("[\x1B[1;36m!\x1B[0m] Username: ")
	fmt.Scanln(&target)
	fmt.Printf("\n[\x1B[1;36m!\x1B[0m] Press Enter to start")
	fmt.Scanln()
	channel := make(chan string)
	readFile()
	for i := 0; i < threads; i++ {
		go sendReqs(channel)
	}
	go ThreadPrint()
	go updateRS()
	fillUserChan(channel)
	fmt.Scanln()

}

func fillUserChan(sender chan string) {
	for {
		for _, session := range sessionSlice {
			sender <- session
		}
	}
}

func readFile() {
	file, _ := os.Open("freshies.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sessionSlice = append(sessionSlice, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func sendReqs(receiver chan string) {
	for running {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		req.SetBodyString("username=" + target)
		req.SetRequestURI("https://b.secure.instagram.com/api/v1/accounts/set_username/")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.SetMethod("POST")
		req.Header.Set("Accept-language", "en-us")
		req.Header.Set("User-Agent", "Instagram 113.0.0.39.122 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
		session := <-receiver
		req.Header.Set("Cookie", "sessionid="+session)
		client.Do(req, resp)

		if strings.Contains(string(resp.Body()), "email") {
			running = false
			mutex.Lock()
			fmt.Printf("[\x1B[1;32m+\x1B[0m] Successfully Claimed @%s\n", target)
			fmt.Printf("[\x1B[1;32m+\x1B[0m] SessionID: %s\n", session)
			removalLine(session)
			mutex.Unlock()
			os.Exit(0)

		} else if strings.Contains(string(resp.Body()), "Please wait") {
			ratelimit++
		} else {
			attempts++
		}

	}

}

// do later
func removalLine(session string) {
	content, err := os.ReadFile("freshies.txt")
	if err != nil {
		log.Fatal(err)
	}
	filelines := strings.Split(string(content), "\n")
	for index, line := range filelines {
		if line == session {
			filelines = append(filelines[:index], filelines[index+1:]...)
			break
		}
	}
	new_content := []byte(strings.Join(filelines, "\n"))
	err = os.WriteFile("freshies.txt", new_content, 0644) // file permission to read and write
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[\x1B[1;32m+\x1B[0m] Removed Session: %s\n", session)
}

func ThreadPrint() {
	for running {
		fmt.Printf("\r[\x1B[1;34m!\x1B[0m] Requests: %d | R/S: %d | RL: %d\r", attempts, per_sec, ratelimit)
		time.Sleep(150 * time.Millisecond)
	}
}

func updateRS() {
	for {
		oldAttempts := attempts
		time.Sleep(time.Second)
		per_sec = attempts - oldAttempts
	}
}

//Honorable function:
/*

	func send_hooks(){} <- Uses third party package for webhook

*/
