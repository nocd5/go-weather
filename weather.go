package weather

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Get(loc int) *WeatherNews {
	proxy, err := url.Parse(os.Getenv("HTTPS_PROXY"))
	if err != nil {
		panic(err)
	}
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// ex) 東京:44132
	res, err := client.Get(fmt.Sprintf("https://weathernews.jp/pinpoint/xml/%d.xml", loc))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	wn := new(WeatherNews)
	wn.Observatory = Observatory()[loc]

	// minimum temperature has chance of `nodata`
	err = xml.Unmarshal([]byte(strings.Replace(string(b), "nodata", "99", -1)), &wn)
	if err != nil {
		return nil
	}

	return wn
}
