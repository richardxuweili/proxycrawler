package proxy

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CyberSource() ([]*Proxy, error) {
	resp, err := http.Get("http://www.cybersyndrome.net/search.cgi?q=CN&a=ABC&f=s&s=new&n=500")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	proxys := extractProxys(doc.Find("#content > script").Text())
	return proxys, nil
}

func extractProxys(str string) (proxys []*Proxy) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	exp := regexp.MustCompile(`\[(.+?)\]`)
	expN1 := regexp.MustCompile(`\((.+)\)%2000`)
	expN2 := regexp.MustCompile(`ps\[(\d+)\]`)
	expN3 := regexp.MustCompile(`\d+\*\d+`)
	matched := exp.FindAllStringSubmatch(str, 2)
	as := strings.Split(matched[0][1], ",")
	ps := strings.Split(matched[1][1], ",")
	s1 := expN1.FindStringSubmatch(str)[1]
	s2 := expN2.ReplaceAllStringFunc(s1, func(str string) string {
		num, _ := strconv.Atoi(strings.Trim(str, "ps[]"))
		return ps[num]
	})
	s3 := expN3.ReplaceAllStringFunc(s2, func(str string) string {
		nums := strings.Split(str, "*")
		num1, _ := strconv.Atoi(nums[0])
		num2, _ := strconv.Atoi(nums[1])
		return strconv.Itoa(num1 * num2)
	})
	var n int64
	nums := strings.Split(s3, "+")
	for _, v := range nums {
		num, _ := strconv.Atoi(v)
		n += int64(num)
	}
	n = n % 2000

	as = append(as[n:], as[0:n]...)
	for i := range as {
		idx := i / 4
		if i%4 == 0 {
			proxys = append(proxys, &Proxy{"http", as[i] + ".", ps[idx], 0})
		} else if i%4 == 3 {
			proxys[idx].IP += as[i]
		} else {
			proxys[idx].IP += as[i] + "."
		}
	}
	return
}
