package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unsafe"

	gokogirixml "github.com/jbowtie/gokogiri/xml"
	"github.com/jbowtie/ratago/xslt"
)

/*
#include "swephexp.h"
#cgo CFLAGS: -Iswe
#cgo LDFLAGS: -Lswe -lswe -lm -ldl
*/
import "C"

var bnames = []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter",
	"Saturn", "Uranus", "Neptune", "Pluto", "MeanNode", "TrueNode",
	"MeanApogee", "OscuApogee", "Earth", "Chiron", "Pholus", "Ceres", "Pallas",
	"Juno", "Vesta", "InterpretedApogee", "InterpretedPerigee", "MeanSouthNode",
	"TrueSouthNode"}

var hnames = []string{"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII",
	"IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX",
	"XX", "XXI", "XXII", "XXIII", "XXIV", "XXV", "XXVI", "XXVII", "XXVIII",
	"XXIX", "XXX", "XXXI", "XXXII", "XXXIII", "XXXIV", "XXXV", "XXXVI"}

var anames = []string{"Ascendant", "MC", "ARMC", "Vertex",
	"EquatorialAscendant", "Co-Ascendant1", "Co-Ascendant2", "PolarAscendant"}

var snames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo",
	"Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius",
	"Pisces"}

type ChartInfo struct {
	XMLName xml.Name `xml:"chartinfo"`
	Houses  []House  `xml:"houses>House"`
	Bodies  []Body   `xml:"bodies>Body"`
	AscMCs  []AscMC  `xml:"ascmcs>AscMC"`
	Aspects []Aspect `xml:"aspects>Aspect"`
	Year    int64    `xml:"year,attr"`
	Month   int64    `xml:"month,attr"`
	Day     int64    `xml:"day,attr"`
	Time    float64  `xml:"time,attr"`
	Lat     float64  `xml:"lat,attr"`
	Lon     float64  `xml:"lon,attr"`
	Name    string   `xml:"name,attr"`
	City    string   `xml:"city,attr"`
	Display string   `xml:"display,attr"`
}

type AscMC struct {
	XMLName  xml.Name
	ID       int     `xml:"id,attr"`
	Sign     int     `xml:"sign,attr"`
	SignName string  `xml:"sign_name,attr"`
	Degree   float64 `xml:"degree,attr"`
	DegreeUt float64 `xml:"degree_ut,attr"`
}

type House struct {
	SignName string  `xml:"sign_name,attr"`
	Degree   float64 `xml:"degree,attr"`
	Number   string  `xml:"number,attr"`
	Sign     int     `xml:"sign,attr"`
	House    int     `xml:"house,attr"`
	DegreeUt float64 `xml:"degree_ut,attr"`
}

// Body represent a planet, a fictional planet or an asteroid
type Body struct {
	XMLName    xml.Name
	Sign       int     `xml:"sign,attr"`
	SignName   string  `xml:"sign_name,attr"`
	Degree     float64 `xml:"degree,attr"`
	DegreeUt   float64 `xml:"degree_ut,attr"`
	Retrograde bool    `xml:"retrograde,attr"`
	ID         int     `xml:"id,attr"`
}

type Aspect struct {
	XMLName xml.Name
	Body1   string  `xml:"body1,attr"`
	Body2   string  `xml:"body2,attr"`
	Degree1 float64 `xml:"degree1,attr"`
	Degree2 float64 `xml:"degree2,attr"`
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}

func normalize(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

func testAspect(ci *ChartInfo, body1 Body, body2 Body, deg1 float64, deg2 float64, delta float64, orb float64, t string) {
	if (deg1 > (deg2+delta-orb) && deg1 < (deg2+delta+orb)) ||
		(deg1 > (deg2-delta-orb) && deg1 < (deg2-delta+orb)) ||
		(deg1 > (deg2+360+delta-orb) && deg1 < (deg2+360+delta+orb)) ||
		(deg1 > (deg2-360+delta-orb) && deg1 < (deg2-360+delta+orb)) ||
		(deg1 > (deg2+360-delta-orb) && deg1 < (deg2+360-delta+orb)) ||
		(deg1 > (deg2-360-delta-orb) && deg1 < (deg2-360-delta+orb)) {
		if deg1 > deg2 {
			aspect(ci, body1, body2, deg1, deg2, t)
		}
	}
}

func aspect(ci *ChartInfo, body1 Body, body2 Body, deg1 float64, deg2 float64, t string) {
	ci.Aspects = append(ci.Aspects,
		Aspect{
			XMLName: xml.Name{Local: t},
			Body1:   body1.XMLName.Local,
			Body2:   body2.XMLName.Local,
			Degree1: deg1,
			Degree2: deg2,
		},
	)
}

// ChartInfoHandler returns houses and planet positions for a location and time
func ChartInfoHandler(w http.ResponseWriter, r *http.Request) {
	var c = &ChartInfo{}

	var xx [6]C.double
	serr := make([]byte, 256)
	var julday C.double
	var cusp [37]C.double
	var ascmc [10]C.double
	var hsys int = 'E'
	c.Year = 1970
	c.Month = 1
	c.Day = 1
	c.Display = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23"
	display := make([]int, C.SE_NPLANETS)
	for i := 0; i < len(display); i++ {
		display[i] = i
	}

	if r.URL.Query().Get("hsys") != "" {
		hsys = int([]rune(r.URL.Query().Get("hsys"))[0])
	}

	if r.URL.Query().Get("year") != "" {
		i, err := strconv.ParseInt(r.URL.Query().Get("year"), 10, 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Year = i
	}

	if r.URL.Query().Get("month") != "" {
		i, err := strconv.ParseInt(r.URL.Query().Get("month"), 10, 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Month = i
	}

	if r.URL.Query().Get("day") != "" {
		i, err := strconv.ParseInt(r.URL.Query().Get("day"), 10, 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Day = i
	}

	if r.URL.Query().Get("time") != "" {
		i, err := strconv.ParseFloat(r.URL.Query().Get("time"), 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Time = i
	}

	if r.URL.Query().Get("lat") != "" {
		i, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Lat = i
	}

	if r.URL.Query().Get("lon") != "" {
		i, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		c.Lon = i
	}

	if r.URL.Query().Get("display") != "" {
		c.Display = r.URL.Query().Get("display")

		d, err := sliceAtoi(strings.Split(c.Display, ","))

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		display = d
	}

	c.Name = r.URL.Query().Get("name")
	c.City = r.URL.Query().Get("city")

	// The number of houses is 12 except when using Gauquelin sectors
	var numhouses = 12
	if hsys == 'G' {
		numhouses = 36
	}

	julday = C.swe_julday(C.int(c.Year), C.int(c.Month), C.int(c.Day), C.double(c.Time), C.SE_GREG_CAL)

	C.swe_set_topo(C.double(c.Lat), C.double(c.Lon), 0)

	C.swe_houses(julday, C.double(c.Lat), C.double(c.Lon), C.int(hsys), (*C.double)(&cusp[0]), (*C.double)(&ascmc[0]))

	// AscMC
	for index := 0; index < C.SE_NASCMC; index++ {
		degreeUt := float64(ascmc[index])

		for sign := 0; sign < 12; sign++ {
			degLow := float64(sign * 30)
			degHigh := float64((sign + 1) * 30)
			if degreeUt >= degLow && degreeUt <= degHigh {

				c.AscMCs = append(c.AscMCs,
					AscMC{
						XMLName:  xml.Name{Local: anames[index]},
						ID:       index + 1,
						Sign:     sign,
						SignName: snames[sign],
						Degree:   degreeUt - degLow,
						DegreeUt: degreeUt,
					},
				)
			}
		}
	}

	// Houses
	for house := 1; house <= numhouses; house++ {
		degreeUt := float64(cusp[house])

		for sign := 0; sign < 12; sign++ {
			degLow := float64(sign * 30)
			degHigh := float64((sign + 1) * 30)
			if degreeUt >= degLow && degreeUt <= degHigh {

				c.Houses = append(c.Houses,
					House{
						SignName: snames[sign],
						Degree:   degreeUt - degLow,
						Number:   hnames[house],
						Sign:     sign,
						House:    house,
						DegreeUt: degreeUt,
					},
				)
			}
		}
	}

	// Bodies
	for body := C.int32(0); body < C.SE_NPLANETS+2; body++ {

		if !contains(display[:], int(body)) {
			continue
		}

		var degreeUt float64
		if body == 23 {
			C.swe_calc_ut(julday, body, 10, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = normalize(float64(xx[0]) + 180)
		} else if body == 24 {
			C.swe_calc_ut(julday, 11, 0, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = normalize(float64(xx[0]) + 180)
		} else {
			C.swe_calc_ut(julday, body, 0, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = float64(xx[0])
		}

		retrograde := xx[3] < 0

		for sign := 0; sign < 12; sign++ {
			degLow := float64(sign * 30)
			degHigh := float64((sign + 1) * 30)
			if degreeUt >= degLow && degreeUt <= degHigh {

				c.Bodies = append(c.Bodies,
					Body{
						XMLName:    xml.Name{Local: bnames[body]},
						Sign:       sign,
						SignName:   snames[sign],
						Degree:     degreeUt - degLow,
						DegreeUt:   degreeUt,
						Retrograde: retrograde,
						ID:         int(body),
					},
				)
			}
		}
	}

	// Ascpects
	for _, body1 := range c.Bodies {
		deg1 := body1.DegreeUt - c.AscMCs[0].DegreeUt + 180

		for _, body2 := range c.Bodies {
			deg2 := body2.DegreeUt - c.AscMCs[0].DegreeUt + 180

			testAspect(c, body1, body2, deg1, deg2, 180, 10, "Opposition")
			testAspect(c, body1, body2, deg1, deg2, 150, 2, "Quincunx")
			testAspect(c, body1, body2, deg1, deg2, 120, 8, "Trine")
			testAspect(c, body1, body2, deg1, deg2, 90, 6, "Square")
			testAspect(c, body1, body2, deg1, deg2, 60, 4, "Sextile")
			testAspect(c, body1, body2, deg1, deg2, 30, 1, "Semi-sextile")
			testAspect(c, body1, body2, deg1, deg2, 0, 10, "Conjunction")
		}
	}

	out, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(out)
}

var globalStyle, globalDoc *gokogirixml.XmlDocument
var globalStylesheet *xslt.Stylesheet

// TransformHandler performs an XSLT transformation
func TransformHandler(w http.ResponseWriter, r *http.Request) {
	globalStyle, _ := gokogirixml.ReadFile("wheel.xsl", gokogirixml.StrictParseOption)
	globalStylesheet, _ := xslt.ParseStylesheet(globalStyle, "test.xsl")
	globalDoc, _ := gokogirixml.ReadFile("test.xml", gokogirixml.StrictParseOption)
	out, _ := globalStylesheet.Process(globalDoc, xslt.StylesheetOptions{true, nil})
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(out))
}

func main() {

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/chartinfo", ChartInfoHandler)
	http.HandleFunc("/transform", TransformHandler)

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Errorf("$PORT not set")
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
