package weather

type WeatherNews struct {
	URL         string
	Observatory map[string]string

	Date struct {
		Last struct {
			Month  int `xml:"month"`
			Date   int `xml:"date"`
			Hour   int `xml:"hour"`
			Minute int `xml:"minute"`
		} `xml:"last"`
		Next struct {
			Month  int `xml:"month"`
			Date   int `xml:"date"`
			Hour   int `xml:"hour"`
			Minute int `xml:"minute"`
		} `xml:"next"`
		Week struct {
			Date    []int `xml:"date"`
			Day     []int `xml:"day"`
			Holiday []int `xml:"holiday"`
		} `xml:"week>day"`
	} `xml:"date"`
	Data struct {
		AreaNo int `xml:"areaNo,attr"`
		Day    struct {
			StartHour  int `xml:"startHour,attr"`
			StartYear  int `xml:"startYear,attr"`
			StartMonth int `xml:"startMonth,attr"`
			StartDate  int `xml:"startDate,attr"`
			StartDay   int `xml:"startDay,attr"`
			Weather    struct {
				Hour []int `xml:"hour"`
			} `xml:"weather"`
			Temperature struct {
				Unit string `xml:"unit,attr"`
				Hour []int  `xml:"hour"`
			} `xml:"temperature"`
			Precipitation struct {
				Unit string `xml:"unit,attr"`
				Hour []int  `xml:"hour"`
			} `xml:"precipitation"`
		} `xml:"day"`
		Week struct {
			Weather struct {
				Day []int `xml:"day"`
			} `xml:"weather"`
			Temperature struct {
				Unit string `xml:"unit,attr"`
				Day  []struct {
					Max int `xml:"max"`
					Min int `xml:"min"`
				} `xml:"day"`
			} `xml:"temperature"`
			ChanceOfRain struct {
				Unit string `xml:"unit,attr"`
				Day  []int  `xml:"day"`
			} `xml:"chance_of_rain"`
		} `xml:"week"`
	} `xml:"data"`
}
