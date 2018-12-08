package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/raff/godet"
	"github.com/urfave/cli"
)

var chromeapp string

func lowhangingfruits(username string) (yesno bool) {
	client := &http.Client{}
	fmt.Println(username)
	fmt.Printf("https://www.instagram.com/%s/", username)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.instagram.com/%s/", username), nil)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	gotcha, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(gotcha.StatusCode)
	fmt.Println("lol")
	if gotcha.StatusCode != 200 {
		return false
	} else {
		return true
	}
}

func createRequest() *http.Client {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://httpbin.org/user-agent", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
	fmt.Println(resp.StatusCode)
	return client
}

func loadJS(username string, remote *godet.RemoteDebugger) {
	remote.RuntimeEvents(true)
	remote.NetworkEvents(true)
	remote.PageEvents(true)
	remote.DOMEvents(true)
	remote.LogEvents(true)
	_, _ = remote.Navigate(fmt.Sprintf("https://www.instagram.com/%s/", username))
	_ = remote.SaveScreenshot("screenshot.png", 0644, 0, true)

}

func main() {
	switch runtime.GOOS {
	case "darwin":
		chromeapp = `open "/Applications/Google Chrome Canary.app" --args`
	case "linux":
		chromeapp = "chromium-browser"
	}
	if chromeapp != "" {
		chromeapp = " --headless --remote-debugging-port=9222 --hide-scrollbars"
	}
	exec.Command("open -a '/Applications/Google Chrome Canary.app' --args --headless --remote-debugging-port=9222 --hide-scrollbars").Run()
	port := "localhost:9222"
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := range nums {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}
		remote, err := godet.Connect(port, false)
		if err == nil { // connection succeeded
			break
		}
		log.Println("connect", err, remote)

		if err != nil {
			log.Println("cannot connect to browser")
		}
	}
	remote, err := godet.Connect(port, false)
	remote.RuntimeEvents(true)
	remote.NetworkEvents(true)
	remote.PageEvents(true)
	remote.DOMEvents(true)
	remote.LogEvents(true)

	app := cli.NewApp()
	app.Name = "go-spry"
	app.Usage = "check user accounts"
	app.Action = func(c *cli.Context) error {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username to check: ")
		text, _ := reader.ReadString('\n')
		fmt.Printf("Username: %s", text)
		fmt.Println(lowhangingfruits(text))
		createRequest()
		loadJS(text, remote)
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
