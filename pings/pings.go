//Package pings based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
package pings
import (
	"os/exec"
)
//Ping pings slice of targets with given concurrency and retunrs alives 
func Ping(conc int, count, timeout string, hosts []string) []string {
	concurrentMax := conc
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan string, len(hosts))
	doneChan := make(chan []string)
	
	
	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan, count, timeout)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
		//fmt.Println("sent: ", ip)
	}

	alives := <-doneChan
	return alives
}

func ping(pingChan <-chan string, pongChan chan <- string, count, timeout string) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c", count, "-W", timeout, ip).Output()  //Linux (timenout in s)
		//_, err := exec.Command("ping", "-n", count, "-w", timeout, ip).Output()  //Windows (timenout in ms)
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

