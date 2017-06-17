package weather

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"strconv"
	"strings"
)

//go:generate go-assets-builder -p weather asset/ame_master.csv -o ame_master.csv.go
func Observatory() map[int]map[string]string {
	file, err := Assets.Open("/asset/ame_master.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	reader.LazyQuotes = true
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}

	key := 0
	for i, h := range header {
		key = i
		if h == "観測所番号" {
			break
		}
	}
	typ := 0
	for i, h := range header {
		typ = i
		if h == "種類" {
			break
		}
	}

	var obs = make(map[int]map[string]string)
	for {
		l, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if l[typ] == "雨" {
			continue
		}
		id, err := strconv.Atoi(l[key])
		if err != nil {
			panic(err)
		}
		obs[id] = make(map[string]string)
		for i, v := range l {
			if i == key {
				continue
			}
			obs[id][header[i]] = v
		}
	}
	return obs
}

func FindLocId(keyword []string) []int {
	var locs []int
	for k, v := range Observatory() {
		for _, kw := range keyword {
			if strings.Index(fmt.Sprintf("%d", k), kw) == 0 ||
				strings.Contains(v["都府県振興局"], kw) ||
				strings.Contains(v["観測所名"], kw) ||
				strings.Contains(v["所在地"], kw) {
				locs = append(locs, k)
			}
		}
	}
	return locs
}
