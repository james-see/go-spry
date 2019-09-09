package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/raff/godet"
	"github.com/sajari/docconv"
	"github.com/urfave/cli"
)

var chromeapp string

func lowhangingfruits(username string) (yesno bool) {
	client := &http.Client{}
	// fmt.Println(username)
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
	// fmt.Println(gotcha.StatusCode)
	// fmt.Println("lol")
	if gotcha.StatusCode != 200 {
		return false
	}
	return true
}

// createRequest creates the http request object identified as client
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
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

// loadJS navigates to the requested webpage TODO pass in url as variable to load
func loadJS(username string, remote *godet.RemoteDebugger) {
	remote.RuntimeEvents(true)
	remote.NetworkEvents(true)
	remote.PageEvents(true)
	remote.DOMEvents(true)
	remote.LogEvents(true)
	_, _ = remote.Navigate(fmt.Sprintf("https://www.instagram.com/%s/", username))
	_ = remote.SaveScreenshot("screenshot.png", 0644, 0, true)
	_ = remote.SavePDF("page.pdf", 0644)

}

// This function extracts text from the saved PDF for each check
func getText(pdfname string) string {
	res, err := docconv.ConvertPath(pdfname)
	if err != nil {
		log.Fatal(err)
	}
	return string(res.Body)
}

// WriteToFile writes text string to a text file
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func main() {
	cmd := exec.Command("/Applications/canary.app/Contents/MacOS/canary", "--args", "--headless", "--remote-debugging-port=9222", "--hide-scrollbars")
	cmd.Start()

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
		text = strings.TrimSuffix(text, "\n")
		fmt.Printf("Username: %s\n", text)
		if lowhangingfruits(text) {
			createRequest()
			loadJS(text, remote)
			textfrompdf := getText("page.pdf")
			err := WriteToFile("result.txt", textfrompdf)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		}
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Process.Kill()
	if err != nil {
		panic(err) // panic as can't kill a process.
	}
}
