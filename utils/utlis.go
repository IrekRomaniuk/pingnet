package utils
import (
    "os"
    "bufio"
    "fmt"
    "errors"
    "strconv"
)
//Diff produces difference between two string slices
func Diff(slice1 []string, slice2 []string) ([]string){
    diffStr := []string{}
    m :=map [string]int{}

    for _, s1Val := range slice1 {
        m[s1Val] = 1
    }
    for _, s2Val := range slice2 {
        m[s2Val] = m[s2Val] + 1
    }

    for mKey, mVal := range m {
        if mVal==1 {
            diffStr = append(diffStr, mKey)
        }
    }

    return diffStr
}

//Hosts creates output slice of targets based on the HOSTS flag
func Hosts(flaghosts string) ([]string, error) {        
    if flaghosts == "all" {
		return Deletempty(list1s(208)), nil
		//fmt.Println(hosts, len(hosts))
	} else if num, err := strconv.Atoi(flaghosts); err == nil {
		if (192 < num) && (num <= 208) {
			return Deletempty(list1s(num)), nil
		} else {
			return Deletempty(list1s(208)), nil
		}
	} else if pathExists(flaghosts) {
		lines, err := readHosts(flaghosts)
		
		if err != nil {
            //fmt.Println("Error reading file", flag_hosts)
            return  []string{}, errors.New("Error reading file")
		} else {
			return Deletempty(lines), nil
		}
	} else {
        return []string{}, errors.New("Flag not recognized")		
	}
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
//Deletempty deletes empty slice members
func Deletempty(s []string) []string {
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
	defer file.Close()
	if err != nil {
		return nil, err
	}

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