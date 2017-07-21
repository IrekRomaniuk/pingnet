package main

import (	
	"fmt"
	"time"
	"flag"
	"os"	
	"strings"
	"github.com/IrekRomaniuk/pingnet/utils"
	"github.com/IrekRomaniuk/pingnet/pings"
)

var (
	HOSTS = flag.String("a", "all", "destinations to ping, i.e. ./file.txt") // 'all', '/path/file' or i.e. '193'
	CONCURRENTMAX = flag.Int("r", 200, "max concurrent pings")
	PINGCOUNT = flag.String("c", "1", "ping count)")
	PINGTIMEOUT = flag.String("w", "1", "ping timout in s")
	version = flag.Bool("v", false, "Prints current version")
	PRINT = flag.String("p", "", "print or not metadata for 'alive' or dead 'targets'")
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

func main() {

	hosts, err := utils.Hosts(*HOSTS)

	if err != nil {
		fmt.Println(hosts)
		os.Exit(0)
	}
	if *PRINT != "" {
		fmt.Printf("concurrentMax=%d hosts=%d -> %s...%s\n", *CONCURRENTMAX, len(hosts), hosts[0], hosts[len(hosts) - 1])
	}
	start := time.Now()
	result := utils.Deletempty(pings.Ping(*CONCURRENTMAX, *PINGCOUNT, *PINGTIMEOUT, hosts))
	
	if strings.Contains(*PRINT, "alive") {
		for _, ip := range result {
			fmt.Println(ip)
			}
		fmt.Printf("%.2fs alive/total: %d/%d cur: %d\n", time.Since(start).Seconds(),len(result),len(hosts),*CONCURRENTMAX)
	} else if strings.Contains(*PRINT, "dead") {
		dead := utils.Diff(hosts, result)
		for _, ip := range dead {
			fmt.Println(ip)
			}
		fmt.Printf("%.2fs dead/total: %d/%d cur: %d\n", time.Since(start).Seconds(),len(dead),len(hosts),*CONCURRENTMAX)
	}
	//Telegraf compliant output be default
	fmt.Printf("pingcount,site=%s,cur=%d total-up=%d\n", *SITE, *CONCURRENTMAX, len(result))
	if strings.Contains(*PRINT, "firebase") {
		fmt.Printf("firebase output to be implemented\n")
	}
}

