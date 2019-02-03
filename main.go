package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	tests := [][]string{
		{
			"/admin/user?ID=1, 3.01s, Status Code: 200",
			"/admin/user?ID=2, 3.02s, Status Code: 404",
			"/admin/user?ID=3, 3.03s, Status Code: 200",
			"/admin/user?ID=4, 3.04s, Status Code: 404",
			"/admin/user?ID=5, 3.05s, Status Code: 200",
			"/admin/user?ID=6, 3.06s, Status Code: 200",
			"/admin/user?ID=7, 3.07s, Status Code: 200",
			"/admin/user?ID=8, 3.08s, Status Code: 200",
			"/admin/user?ID=9, 3.09s, Status Code: 200",
			"/admin/user?ID=10, 3.10s, Status Code: 200",
			"/admin/user?ID=11, 3.11s, Status Codee: 200",
			"/admin/user?ID=12, 3.12s, Status Code: 200",
			"/admin/user?ID=13, 3.13s, Status Code: 200",
			"/admin/user/update, 3.14s, Status Code: 200",
			"/admin/user?ID=15, 3.15s, Status Code: 200",
		},
		{
			"/admin/user?ID=1, 3.01s, Status Code: 200",
			"/admin/user?ID=2, 3.02s, Status Code: 404",
			"/admin/user?ID=3, 3.03s, Status Code: 200",
			"/admin/user?ID=4, 3.04s, Status Code: 404",
			"/admin/user?ID=5, 3.05s, Status Code: 200",
		},
	}

	for _, test := range tests {
		fmt.Printf("TEST CASE :\n%v\nOUT PUT :%v\n\n", test, getTopTen(test))
	}

}

type tmp struct {
	url      string
	value    float64
	avg      float64
	quantity int
}

func getTopTen(logs []string) (result []string) {

	if len(logs) == 0 {
		return []string{}
	}

	logStore := map[string]tmp{}

	for _, log := range logs {
		url, respTime, _, err := getURLResTimeStaCode(log)
		if err != nil {
			fmt.Printf("%s does not pass with error %v\n", log, err)
			continue
		}
		logS, exist := logStore[url]
		if exist {
			logStore[url] = tmp{
				value:    respTime,
				avg:      respTime,
				quantity: 1,
			}
		} else {
			avg := (logS.value + respTime) / float64(logS.quantity+1)
			logStore[url] = tmp{
				value:    logS.value + respTime,
				avg:      avg,
				quantity: logS.quantity + 1,
			}
		}

	}

	n := map[float64][]string{}
	for k, v := range logStore {
		n[v.value] = append(n[v.value], k)
	}

	var a []float64
	for k := range n {
		a = append(a, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(a)))

	for _, v := range a {
		for _, value := range n[v] {
			result = append(result, value)
		}
	}

	if len(result) > 10 {
		return result[:10]
	}
	return result
}

func getURLResTimeStaCode(log string) (url string, respTime float64, statCode string, err error) {
	if log == "" {
		return url, respTime, statCode, fmt.Errorf("log is empty string")
	}
	datas := strings.Split(log, ",")
	if len(datas) < 3 {
		return url, respTime, statCode, fmt.Errorf("log is not nginx log")
	}

	url = strings.ToLower(strings.TrimSpace(datas[0]))
	statCode = strings.TrimSpace(strings.Replace(datas[2], "Status Code:", "", 1))

	respTime, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(datas[1], "s", "", 1)), 64)
	if err != nil {
		return url, respTime, statCode, errors.Wrapf(err, "failed to ParseFloat")
	}

	// only 200 status code
	if statCode != "200" {
		return url, respTime, statCode, fmt.Errorf("status code is not 200")
	}

	// exclude .gif
	if strings.HasSuffix(url, ".gif") {
		return url, respTime, statCode, fmt.Errorf("url has suffix .gif")
	}

	// since there's no method in prefix of log
	// I assume that if url is contain parameter or a file is GET method
	// otherwise is POST or another method
	if !regexp.MustCompile(`(?m).+\.\w+($|\n)|.+?\w+=\d+($|\n)`).Match([]byte(url)) {
		return url, respTime, statCode, fmt.Errorf("url Method is not GET")
	}

	return url, respTime, statCode, err
}
