package main

import (
	"fmt"
	"io/ioutil"
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
		"http://googlehosts-hostsfiles.stor.sinaapp.com/hosts",
		"http://blog.my-eclipse.cn/hosts.txt",
		"https://raw.githubusercontent.com/racaljk/hosts/master/hosts",
		"http://gcat.gq/wp-content/uploads/2016/04/201604040806092.txt",
	}

	for _, url := range urlList {

		resp, err = http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

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

	oldHosts, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	oldHostsContext := string(oldHosts)
	var pat = "(?s)#old hosts start.*#old hosts end"

	re, _ = regexp.Compile(pat)
	oldStr := re.FindAllStringSubmatch(string(oldHostsContext), -1)
	if len(oldStr) == 0 {
		oldHostsContext = "\n#old hosts start\n" + oldHostsContext + "\n#old hosts end\n"
	} else {
		oldHostsContext = oldStr[0][0]
	}

	file, err := os.OpenFile(hostsFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer file.Close()

	file.WriteString(oldHostsContext)
	file.WriteString("\n \n")

	newHosts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	file.WriteString(string(newHosts))
	fmt.Println("\nUpdate the hosts success, press ENTER to exit!")

	resp.Body.Close()
	fmt.Scanln()
}
