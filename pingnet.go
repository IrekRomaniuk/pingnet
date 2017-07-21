package main
//Based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
import (
	"os/exec"
	"fmt"
	"time"
	"flag"
	"os"	
	"github.com/IrekRomaniuk/pingnet/utils"
)

var (
	HOSTS = flag.String("a", "all", "destinations to ping, i.e. ./file.txt") // 'all', '/path/file' or i.e. '193'
	CONCURRENTMAX = flag.Int("r", 200, "max concurrent pings")
	PINGCOUNT = flag.String("c", "1", "ping count)")
	PINGTIMEOUT = flag.String("w", "1", "ping timout in s")
	version = flag.Bool("v", false, "Prints current version")
	PRINT = flag.String("p", "alive", "print metadata for 'alive' or dead 'targets'")
	SITE = flag.String("s", "DC1", "source location tag")
)
var (
	Version = "No Version Provided"
	BuildTime = ""
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Copyright 2017 @IrekRomaniuk. All rights reserved.\n")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *version {
		fmt.Printf("App Version: %s\nBuild Time : %s\n", Version, BuildTime)
		os.Exit(0)
	}
}
func ping(pingChan <-chan string, pongChan chan <- string) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c", *PINGCOUNT, "-W", *PINGTIMEOUT, ip).Output()  //Linux (timenout in s)
		//_, err := exec.Command("ping", "-n", *PINGCOUNT, "-w", *PINGTIMEOUT, ip).Output()  //Windows (timenout in ms)
		if err == nil {
			pongChan <- ip
		} else {
			pongChan <- ""
		}
	}
}

func receivePong(pongNum int, pongChan <-chan string, doneChan chan <- []string) {
	var alives []string
	for i := 0; i < pongNum; i++ {
		ip := <-pongChan
		//fmt.Println("received: ", ip)
		alives = append(alives, ip)
	}
	doneChan <- alives
}



func main() {
	//var hosts []string

	hosts, err := utils.Hosts(*HOSTS)

	if err != nil {
		fmt.Println(hosts)
		os.Exit(0)
	}

	start := time.Now()
	result := utils.Deletempty(Ping(*CONCURRENTMAX, hosts))
	
	if *PRINT  == "alive" {
		//fmt.Println(result)
		for _, ip := range result {
			fmt.Println(ip)
			}
		fmt.Printf("%.2fs alive/total: %d/%d cur: %d\n", time.Since(start).Seconds(),len(result),len(hosts),*CONCURRENTMAX)
	} else if *PRINT  == "dead" {
		dead := utils.Diff(hosts, result)
		for _, ip := range dead {
			fmt.Println(ip)
			}
		fmt.Printf("%.2fs dead/total: %d/%d cur: %d\n", time.Since(start).Seconds(),len(dead),len(hosts),*CONCURRENTMAX)
	}
	
	fmt.Printf("pingcount,site=%s,cur=%d total-up=%d\n", *SITE, *CONCURRENTMAX, len(result))

}
//Ping pings slice of targets with given concurrency and retunrs alives 
func Ping(conc int, hosts []string) []string {
	concurrentMax := conc
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan string, len(hosts))
	doneChan := make(chan []string)
	if *PRINT == "alive" {
		fmt.Printf("concurrentMax=%d hosts=%d -> %s...%s\n", concurrentMax, len(hosts), hosts[0], hosts[len(hosts) - 1])
	}
	//start := time.Now()
	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
		//fmt.Println("sent: ", ip)
	}

	alives := <-doneChan
	return alives
}
