package proxy

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	total, err := strconv.Atoi(doc.Find("#div_result > span").Text())
	if err != nil {
		return nil, err
	}
	pageNum := total / 500
	if total%500 > 0 {
		pageNum += 1
	}
	proxys := extractProxys(doc.Find("#content > script").Text())
	for i := 2; i <= pageNum; i++ {
		time.Sleep(1000 * time.Millisecond)
		resp, err := http.Get("http://www.cybersyndrome.net/search.cgi?q=CN&a=ABC&f=s&s=new&n=500&p=" + strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return nil, err
		}
		proxys = append(proxys, extractProxys(doc.Find("#content > script").Text())...)
	}
	return proxys, nil
}

func extractProxys(str string) (proxys []*Proxy) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "may be blocked by source site, sleep 60s")
			fmt.Println(err, "may be blocked by source site, sleep 60s")
			time.Sleep(60 * time.Second)
		}
	}()
	exp := regexp.MustCompile(`\[(.+?)\]`)
	expN1 := regexp.MustCompile(`\((.+)\)%(\d+)`)
	expN2 := regexp.MustCompile(`ps\[(\d+)\]`)
	expN3 := regexp.MustCompile(`\d+\*\d+`)
	matched := exp.FindAllStringSubmatch(str, 2)
	as := strings.Split(matched[0][1], ",")
	ps := strings.Split(matched[1][1], ",")
	s1 := expN1.FindStringSubmatch(str)[1]
	i, err := strconv.Atoi(expN1.FindStringSubmatch(str)[2])
	if err != nil {
		panic(err)

	}
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
	n = n % int64(i)

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
