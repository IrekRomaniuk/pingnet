package main
//Based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
import (
	"os/exec"
	"fmt"
	"time"
	//"github.com/k0kubun/pp"
	"flag"
	"os"
	"bufio"
	"strconv"
	"gihub.com/IrekRomaniuk/pingnet/utils"
)

var (
	HOSTS = flag.String("a", "all", "destinations to ping, i.e. ./file.txt") // 'all', '/path/file' or i.e. '193'
	CONCURRENTMAX = flag.Int("r", 200, "max concurrent pings")
	PINGCOUNT = flag.String("c", "1", "ping count)")
	PINGTIMEOUT = flag.String("w", "1", "ping timout in s")
	version = flag.Bool("v", false, "Prints current version")
	PRINT = flag.Bool("p", true, "print metadata")
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

func list1s(limit2 int) []string {
	//Shield_Slice int
	res := make([]string, 256 * 64) //256*64
	for x := 192; x < limit2; x++ {
		//192-256
		for y := 0; y < 256; y++ {
			//0-256
			res = append(res, fmt.Sprintf("10.%d.%d.1", x, y))
			//fmt.Printf("10.%d.%d.1", x, y)
		}
	}
	return res //[:Shield_Slice]
}

func main() {
	var hosts []string
	if *HOSTS == "all" {
		hosts = delete_empty(list1s(208))
		//fmt.Println(hosts, len(hosts))
	} else if num, err := strconv.Atoi(*HOSTS); err == nil {
		if (192 < num) && (num <= 208) {
			hosts = delete_empty(list1s(num))
		} else {
			hosts = delete_empty(list1s(208))
		}
	} else if pathExists(*HOSTS) {
		lines, err := readHosts(*HOSTS)
		hosts = delete_empty(lines)
		if err != nil {
			fmt.Println("Error reading file", *HOSTS)
		}
	} else {
		fmt.Println(*HOSTS)
		os.Exit(0)
	}

	concurrentMax := *CONCURRENTMAX
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan string, len(hosts))
	doneChan := make(chan []string)
	if *PRINT {
		fmt.Printf("concurrentMax=%d hosts=%d -> %s...%s\n", concurrentMax, len(hosts), hosts[0], hosts[len(hosts) - 1])
	}
	start := time.Now()
	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
		//fmt.Println("sent: ", ip)
	}
	alives := <-doneChan
	result := delete_empty(alives)
	if *PRINT {
		//fmt.Println(result)
		for _, ip := range result {
			fmt.Println(ip)
			}
		fmt.Printf("%.2fs %d/%d %d\n", time.Since(start).Seconds(),len(result),len(hosts),concurrentMax)
	}
	//pp.Println(len(result))
	fmt.Printf("pingcount,site=%s,cur=%d total-up=%d\n", *SITE, concurrentMax, len(result))

}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func readHosts(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func pathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}
