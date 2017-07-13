# pingnet (under development)
concurrently ping multiple targets (continuation of  https://github.com/irom77/go-public/pingnet/pingnet.go)
### Flags and version

```
docker@ubuntu-DC1:~$ pingnet -v
App Version: No Version Provided
Build Time : 
docker@ubuntu-DC1:~$ pingnet -h
Copyright 2017 @IrekRomaniuk. All rights reserved.
Usage of pingnet:
  -a string
        destinations to ping, i.e. ./file.txt (default "all")
  -c string
        ping count) (default "1")
  -p string
        print metadata (default "alive")
  -r int
        max concurrent pings (default 200)
  -s string
        source location tag (default "DC1")
  -v    Prints current version
  -w string
        ping timout in s (default "1")
```

### To get list of dead targets:

```
docker@ubuntu-DC1:~$ go get -u github.com/IrekRomaniuk/pingnet
docker@ubuntu-DC1:~$ pingnet > pinglist.txt
docker@ubuntu-DC1:~$ pingnet -a="pinglist.txt" -p="dead"
10.197.78.1
10.195.244.1
10.198.17.1
10.196.45.1
1.73s dead/total: 4/1144 cur: 200
pingcount,site=DC1,cur=200 total-up=1140
```

