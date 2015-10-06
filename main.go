package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Device struct {
	XMLName xml.Name `xml:"device"`
	Id      string   `xml:"id,attr"`
	Make    string   `xml:"make"`
	Model   string   `xml:"model"`
}

type SummaryData struct {
	XMLName         xml.Name `xml:"summarydata"`
	Beginning       int      `xml:"beginning"`
	Duration        int      `xml:"duration"`
	DurationStopped int      `xml:"durationstopped"`
	Dist            float32  `xml:"dist"`
}

type Segment struct {
	XMLName     xml.Name `xml:"segment"`
	Name        string   `xml:"name"`
	SummaryData SummaryData
}

type Sample struct {
	XMLName    xml.Name `xml:"sample"`
	Timeoffset int      `xml:"timeoffset"`
	Hr         int      `xml:"hr"`
	Spd        float32  `xml:"spd"`
	Pwr        int      `xml:"pwr"`
	Cad        int      `xml:"cad"`
	Dist       float32  `xml:"dist"`
	Alt        float32  `xml:"alt"`
	Temp       int      `xml:"temp"`
}

type Workout struct {
	XMLName     xml.Name `xml:"workout"`
	AthleteName string   `xml:"athlete>name"`
	SportType   string   `xml:"sportType"`
	Title       string   `xml:"title"`
	Code        string   `xml:"code"`
	Devie       Device
	Time        string `xml:"time"`
	SummaryData SummaryData
	Segment     []Segment `xml:"segment"`
	Sample      []Sample  `xml:"sample"`
}

type Pwx struct {
	XMLName xml.Name `pwx"`
	Version string   `xml:"version,attr"`
	Creator string   `xml:"creator,attr"`
	Workout Workout
}

func main() {
	xmldoc, _ := ioutil.ReadAll(os.Stdin)
	pwx := Pwx{}
	xml.Unmarshal(xmldoc, &pwx)
	w := pwx.Workout
	for i := 0; i < len(w.Sample); i++ {
		if i < 2 {
			continue
		}
		if w.Sample[i].Pwr == w.Sample[i-1].Pwr && w.Sample[i].Pwr == w.Sample[i-2].Pwr && w.Sample[i].Cad == w.Sample[i-1].Cad && w.Sample[i].Cad == w.Sample[i-2].Cad {
			w.Sample[i].Pwr = 0
			w.Sample[i-1].Pwr = 0
			w.Sample[i].Cad = 0
			w.Sample[i-1].Cad = 0
		}
	}

	buf, _ := xml.MarshalIndent(pwx, "", "  ")
	fmt.Println(string(buf))
}
