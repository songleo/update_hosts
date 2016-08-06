package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

func main() {

	var (
		resp  *http.Response
		re    *regexp.Regexp
		err   error
		hosts string
	)

	// check OS type
	if runtime.GOOS == "windows" {
		hosts = "C:/Windows/System32/drivers/etc/hosts"
	} else {
		hosts = "/etc/hosts"
	}

	urlList := []string{
		"http://googlehosts-hostsfiles.stor.sinaapp.com/hosts",
		"http://blog.my-eclipse.cn/hosts.txt",
		"https://raw.githubusercontent.com/racaljk/hosts/master/hosts",
		"http://gcat.gq/wp-content/uploads/2016/04/201604040806092.txt",
	}

	// search for available hosts url
	for _, url := range urlList {
		resp, err = http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Println("update hosts......")
			break
		} else {
			continue
		}
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("sorry :(\nfail to update hosts, exit.\n")
		os.Exit(-1)
	}

	fileBuf, err := ioutil.ReadFile(hosts)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	var pat = "(?s)#old hosts start.*#old hosts end"
	re, _ = regexp.Compile(pat)

	oldHosts := string(fileBuf)
	findResult := re.FindAllStringSubmatch(oldHosts, -1)
	if len(findResult) == 0 {
		oldHosts = "\n#old hosts start\n" + oldHosts + "\n#old hosts end\n"
	} else {
		oldHosts = findResult[0][0]
	}

	file, err := os.OpenFile(hosts, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer file.Close()

	file.WriteString(oldHosts)
	file.WriteString("\n \n")

	newHosts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	resp.Body.Close()

	file.WriteString(string(newHosts))

	fmt.Println("\nupdate the hosts success, press ENTER to exit.")
	fmt.Scanln()
}
