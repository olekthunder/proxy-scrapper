package main

import (
	"bufio"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var Sources = []string{
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=socks5",
	// "https://openproxy.space/list/socks5",
	"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
	"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks5.txt",
	"https://raw.githubusercontent.com/manuGMG/proxy-365/main/SOCKS5.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/socks5.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS5_RAW.txt",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",
	"https://raw.githubusercontent.com/UserR3X/proxy-list/main/socks5.txt",
	"https://www.proxy-list.download/api/v1/get?type=socks5",
	"https://www.proxyscan.io/download?type=socks5",
}

func scrape(client http.Client, url string, ch chan string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("ERR happened!")
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
}

func main() {
	client := http.Client{Timeout: 10 * time.Second}
	http.HandleFunc("/socks5.txt", func(rw http.ResponseWriter, r *http.Request) {
		ch := make(chan string)
        wg := sync.WaitGroup{}
        wg.Add(len(Sources))
		for _, source := range Sources {
            source := source
            go func() {
    			scrape(client, source, ch)
                wg.Done()
            }()
		}
        go func() {
            wg.Wait()
            close(ch)
        }()
		for proxy := range ch {
			rw.Write([]byte(proxy + "\n"))
		}
	})
	if err := http.ListenAndServe(":8000", nil); err != nil {
        fmt.Println(err)
    }
}
