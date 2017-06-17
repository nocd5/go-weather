package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
	"github.com/nocd5/weather"
	"os"
	"strings"
)

const (
	GREY    = "\033[0;30;1m"
	RED     = "\033[0;31;1m"
	GREEN   = "\033[0;32;1m"
	YELLOW  = "\033[0;33;1m"
	BLUE    = "\033[0;34;1m"
	MAGENTA = "\033[0;35;1m"
	CYAN    = "\033[0;36;1m"
	WHITE   = "\033[0;37;1m"
	END     = "\033[0m"
)

var (
	dow = map[int]string{
		0: "日",
		1: "月",
		2: "火",
		3: "水",
		4: "木",
		5: "金",
		6: "土",
	}

	bar = map[int]string{
		0: "\u2595 ",
		1: "\u2595\u258f",
		2: "\u2595\u258e",
		3: "\u2595\u258d",
		4: "\u2595\u258c",
		5: "\u2595\u258b",
		6: "\u2595\u258a",
		7: "\u2595\u2589",
	}
)

var opts struct {
	List  ListCmd  `description:"List of AMEDAS observatory" command:"list"`
	Week  WeekCmd  `description:"Weather report of this week" command:"week"`
	Today TodayCmd `description:"Weather report of today" command:"today" subcommands-optional:"true"`
}

type ListCmd struct {
}

var (
	stdOut *bufio.Writer
)

func (x *ListCmd) Execute(args []string) error {
	obs := weather.Observatory()

	// maximum length of 都府県振興局、観測所名
	maxlen := map[string]int{
		"都府県振興局": 0,
		"観測所名":   0,
	}
	for _, v := range obs {
		for _, kw := range args {
			if strings.Contains(v["都府県振興局"], kw) || strings.Contains(v["観測所名"], kw) || strings.Contains(v["所在地"], kw) {
				maxlen["都府県振興局"] = max(maxlen["都府県振興局"], runewidth.StringWidth(v["都府県振興局"]))
				maxlen["観測所名"] = max(maxlen["観測所名"], runewidth.StringWidth(v["観測所名"]))
			}
		}
	}
	for k, v := range obs {
		for _, kw := range args {
			if strings.Contains(v["都府県振興局"], kw) || strings.Contains(v["観測所名"], kw) || strings.Contains(v["所在地"], kw) || strings.Index(fmt.Sprintf("%d", k), kw) == 0 {
				fmt.Fprintln(stdOut, strings.Replace(fmt.Sprintf("%d : %s %s %s",
					k,
					v["都府県振興局"]+strings.Repeat(" ", maxlen["都府県振興局"]-runewidth.StringWidth(v["都府県振興局"])),
					v["観測所名"]+strings.Repeat(" ", maxlen["観測所名"]-runewidth.StringWidth(v["観測所名"])),
					v["所在地"],
				), kw, RED+kw+END, -1))
			}
		}
	}
	return nil
}

type WeekCmd struct {
}

func getWeather(locs []int) map[int]weather.WeatherNews {
	wns := make(map[int]weather.WeatherNews, len(locs))
	for _, loc := range locs {
		wn := weather.Get(loc)
		if wn != nil {
			wns[loc] = *wn
		}
	}
	return wns
}

func (x *WeekCmd) Execute(args []string) error {
	locs := weather.FindLocId(args)
	wns := getWeather(locs)

	for _, v := range wns {
		fmt.Fprintf(stdOut, "%s %s (%s)\n",
			v.Observatory["都府県振興局"],
			v.Observatory["観測所名"],
			v.Observatory["所在地"])
		for i := range v.Data.Week.Weather.Day {
			dc := WHITE
			if v.Date.Week.Holiday[i] == 1 {
				dc = RED
			}
			fmt.Fprintf(stdOut, "%02d(%s%s"+END+") "+RED+"%3d"+END+" /"+BLUE+"%3d"+END+" %3d%s %s\n",
				v.Date.Week.Date[i],
				dc,
				dow[v.Date.Week.Day[i]],
				v.Data.Week.Temperature.Day[i].Max,
				v.Data.Week.Temperature.Day[i].Min,
				v.Data.Week.ChanceOfRain.Day[i],
				v.Data.Week.ChanceOfRain.Unit,
				weather.TelopList[v.Data.Week.Weather.Day[i]])
		}
		fmt.Fprintln(stdOut)
	}

	return nil
}

type TodayCmd struct {
	TimeSpan int `long:"span" short:"s" description:"Time span of today's weather report" default:"3"`
}

func (x *TodayCmd) Execute(args []string) error {
	locs := weather.FindLocId(args)
	wns := getWeather(locs)

	var minT int
	var maxT int

	for _, v := range wns {
		fmt.Fprintf(stdOut, "%s %s (%s)\n",
			v.Observatory["都府県振興局"],
			v.Observatory["観測所名"],
			v.Observatory["所在地"])
		minT = v.Data.Day.Temperature.Hour[0]
		maxT = v.Data.Day.Temperature.Hour[0]
		for _, t := range v.Data.Day.Temperature.Hour {
			minT = min(t, minT)
			maxT = max(t, maxT)
		}
		for d, w := range v.Data.Day.Weather.Hour {
			if d%x.TimeSpan == 0 {
				fmt.Fprintf(stdOut, "%2d時%s%3d%s %3d%s %s\n",
					(v.Data.Day.StartHour+d)%24,
					bar[int(float64(v.Data.Day.Temperature.Hour[d]-minT)/float64(maxT-minT)*float64(len(bar)-1)+0.5)],
					v.Data.Day.Temperature.Hour[d],
					v.Data.Day.Temperature.Unit,
					v.Data.Day.Precipitation.Hour[d],
					v.Data.Day.Precipitation.Unit,
					weather.TelopList[w])
			}
		}
		fmt.Fprintln(stdOut)
	}

	return nil
}

func main() {
	stdOut = bufio.NewWriter(colorable.NewColorableStdout())

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "weather"

	args, err := parser.Parse()
	if err != nil {
		return
	}
	if len(args) < 1 {
		parser.WriteHelp(os.Stdout)
		return
	}

	stdOut.Flush()
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
