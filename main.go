package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func loadPwxFile(path string) (*Pwx, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	xmldoc, err := ioutil.ReadAll(inFile)
	if err != nil {
		return nil, err
	}
	pwx := Pwx{}
	xml.Unmarshal(xmldoc, &pwx)
	return &pwx, nil
}

func removeCoastingError(pwx *Pwx) {
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
}

func createNewFileName(origPath string) (outPath string, err error) {
	outPath = strings.Replace(os.Args[1], ".pwx", "-1.pwx", 1)
	return outPath, nil
}

func savePwxFile(pwx *Pwx, outPath string) error {
	buf, err := xml.MarshalIndent(pwx, "", "    ")
	if err != nil {
		return err
	}
	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer outFile.Close()

	outFile.Write(buf)
	fmt.Println(outPath)
	return nil
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: srm-coasting-error-filter /path/to/TrainingPeaks-PWX-file.pwx\n")
		os.Exit(1)
	}
	pwx, err := loadPwxFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	removeCoastingError(pwx)

	outPath, err := createNewFileName(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = savePwxFile(pwx, outPath)
	if err != nil {
		panic(err)
	}
}
