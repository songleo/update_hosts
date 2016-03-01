package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

var hostsFile string

func init() {
	if runtime.GOOS == "windows" {
		hostsFile = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	} else {
		hostsFile = "/etc/hosts"
	}
}

func main() {

	var (
		resp *http.Response
		re   *regexp.Regexp
		err  error
	)

	urlList := []string{
		"https://github.com/racaljk/hosts/blob/master/hosts",
		"http://googleips-google.stor.sinaapp.com/hosts",
		"http://blog.my-eclipse.cn/hosts.txt",
	}

	for _, url := range urlList {

		resp, err = http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Update your hosts......")
			break
		} else {
			continue
		}
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Sorry :(\nUpdate your hosts fail, program will exit.\n")
		os.Exit(-1)
	}
	hosts, err := ioutil.ReadFile(hostsFile)
	hostsContext := string(hosts)
	var pat = "(?s)#old hosts start.*#old hosts end"

	re, _ = regexp.Compile(pat)
	oldStr := re.FindAllStringSubmatch(string(hostsContext), -1)
	if len(oldStr) == 0 {
		hostsContext = "\n#old hosts start\n" + hostsContext + "\n#old hosts end\n"
	} else {
		hostsContext = oldStr[0][0]
	}

	file, err := os.OpenFile(hostsFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(hostsContext)
	file.WriteString("\n \n")

	buf := make([]byte, 10240)
	for {
		numBytes, _ := resp.Body.Read(buf)
		if numBytes == 0 {
			break
		}
		file.WriteString(string(buf[:numBytes]))
	}
	fmt.Println("\nUpdate the hosts success, press ENTER to exit!")
	fmt.Scanln()
}
