package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	gk "github.com/jbowtie/gokogiri/xml"
	"github.com/kivutar/ratago/xslt"
)

/*
#include "swephexp.h"
#cgo CFLAGS: -Iswe
#cgo LDFLAGS: -Lswe -lswe -lm -ldl
*/
import "C"

// Celestial bodies names
var bnames = []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter",
	"Saturn", "Uranus", "Neptune", "Pluto", "MeanNode", "TrueNode",
	"MeanApogee", "OscuApogee", "Earth", "Chiron", "Pholus", "Ceres", "Pallas",
	"Juno", "Vesta", "InterpretedApogee", "InterpretedPerigee", "MeanSouthNode",
	"TrueSouthNode"}

// Houses names
var hnames = []string{"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII",
	"IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX",
	"XX", "XXI", "XXII", "XXIII", "XXIV", "XXV", "XXVI", "XXVII", "XXVIII",
	"XXIX", "XXX", "XXXI", "XXXII", "XXXIII", "XXXIV", "XXXV", "XXXVI"}

// Names for the ascendant and other marks
var anames = []string{"Ascendant", "MC", "ARMC", "Vertex",
	"EquatorialAscendant", "Co-Ascendant1", "Co-Ascendant2", "PolarAscendant"}

// Sign names
var snames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo",
	"Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius",
	"Pisces"}

type aspectsetting struct {
	delta float64
	orb   float64
	title string
}

var aspectsettings = []aspectsetting{
	{180, 10, "Opposition"},
	{150, 2, "Quincunx"},
	{120, 8, "Trine"},
	{90, 6, "Square"},
	{60, 4, "Sextile"},
	{30, 1, "Semi-sextile"},
	{0, 10, "Conjunction"},
}

// ChartInfo is the root node of our xml output
type ChartInfo struct {
	XMLName xml.Name `xml:"chartinfo"`
	AscMCs  []AscMC  `xml:"ascmcs>AscMC"`
	Houses  []House  `xml:"houses>House"`
	Aspects []Aspect `xml:"aspects>Aspect"`
	Bodies  []Body   `xml:"bodies>Body"`
	Display string   `xml:"display,attr,omitempty"`
	Year    int64    `xml:"year,attr,omitempty"`
	Month   int64    `xml:"month,attr,omitempty"`
	Day     int64    `xml:"day,attr,omitempty"`
	Time    float64  `xml:"time,attr,omitempty"`
	Lat     float64  `xml:"lat,attr,omitempty"`
	Lon     float64  `xml:"lon,attr,omitempty"`
	Name    string   `xml:"name,attr,omitempty"`
	City    string   `xml:"city,attr,omitempty"`
	Hsys    string   `xml:"hsys,attr,omitempty"`
}

// AscMC represents special marks like the ascendants
type AscMC struct {
	XMLName  xml.Name
	SignName string  `xml:"sign_name,attr"`
	DegreeUt float64 `xml:"degree_ut,attr"`
	Degree   float64 `xml:"degree,attr"`
	Sign     int     `xml:"sign,attr"`
	ID       int     `xml:"id,attr"`
}

// House represents an astrological house cuspid
type House struct {
	SignName string  `xml:"sign_name,attr"`
	Degree   float64 `xml:"degree,attr"`
	Number   string  `xml:"number,attr"`
	Sign     int     `xml:"sign,attr"`
	ID       int     `xml:"id,attr"`
	DegreeUt float64 `xml:"degree_ut,attr"`
}

// Body represents a planet, a fictional planet or an asteroid
type Body struct {
	XMLName    xml.Name
	SignName   string  `xml:"sign_name,attr"`
	Dist       int     `xml:"dist,attr"`
	DegreeUt   float64 `xml:"degree_ut,attr"`
	Degree     float64 `xml:"degree,attr"`
	Sign       int     `xml:"sign,attr"`
	Retrograde bool    `xml:"retrograde,attr"`
	ID         int     `xml:"id,attr"`
}

// Aspect represents a astrological aspect like a Conjunction or a Sextile
type Aspect struct {
	XMLName xml.Name
	Body1   string  `xml:"body1,attr"`
	Body2   string  `xml:"body2,attr"`
	Degree1 float64 `xml:"degree1,attr"`
	Degree2 float64 `xml:"degree2,attr"`
}

// Checks if an int is contained in an int array
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Converts a slice of strings to a slice of integers
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

// Make sure angle values are in within the 0 to 360 range
func normalize(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

// makeAspect returns an Aspect for a given orb and two celectial bodies
func makeAspect(body1 Body, body2 Body, ascendant float64, delta float64, orb float64, t string) (aspect Aspect) {
	deg1 := normalize(body1.DegreeUt - ascendant + 180)
	deg2 := normalize(body2.DegreeUt - ascendant + 180)

	if (deg1 > (deg2+delta-orb) && deg1 < (deg2+delta+orb)) ||
		(deg1 > (deg2-delta-orb) && deg1 < (deg2-delta+orb)) ||
		(deg1 > (deg2+360+delta-orb) && deg1 < (deg2+360+delta+orb)) ||
		(deg1 > (deg2-360+delta-orb) && deg1 < (deg2-360+delta+orb)) ||
		(deg1 > (deg2+360-delta-orb) && deg1 < (deg2+360-delta+orb)) ||
		(deg1 > (deg2-360-delta-orb) && deg1 < (deg2-360-delta+orb)) {
		if deg1 > deg2 {
			aspect = Aspect{
				XMLName: xml.Name{Local: t},
				Body1:   body1.XMLName.Local,
				Body2:   body2.XMLName.Local,
				Degree1: deg1,
				Degree2: deg2,
			}
			return
		}
	}

	return
}

// ChartInfoHandler returns houses and planet positions for a location and time
func ChartInfoHandler(w http.ResponseWriter, r *http.Request) {
	var c = &ChartInfo{}

	var xx [6]C.double
	serr := make([]byte, 256)
	var julday C.double
	var cusp [37]C.double
	var ascmc [10]C.double
	c.Hsys = "E"
	c.Year = 1970
	c.Month = 1
	c.Day = 1
	c.Display = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23"
	display := make([]int, C.SE_NPLANETS)
	for i := 0; i < len(display); i++ {
		display[i] = i
	}

	if r.URL.Query().Get("hsys") != "" {
		c.Hsys = r.URL.Query().Get("hsys")
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
	if c.Hsys == "G" {
		numhouses = 36
	}

	julday = C.swe_julday(C.int(c.Year), C.int(c.Month), C.int(c.Day), C.double(c.Time), C.SE_GREG_CAL)

	C.swe_set_topo(C.double(c.Lat), C.double(c.Lon), 0)

	C.swe_houses(julday, C.double(c.Lat), C.double(c.Lon), C.int(rune(c.Hsys[0])), (*C.double)(&cusp[0]), (*C.double)(&ascmc[0]))

	// Add ascendant and other marks to the chart
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

	// Add house cuspids to the chart
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
						ID:       house,
						DegreeUt: degreeUt,
					},
				)
			}
		}
	}

	// Add celestial bodies to the chart
	for body := C.int32(0); body < C.SE_NPLANETS+2; body++ {

		if !contains(display[:], int(body)) {
			continue
		}

		var degreeUt float64
		var ret C.int32
		if body == 23 {
			ret = C.swe_calc_ut(julday, 10, 0, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = normalize(float64(xx[0]) + 180)
		} else if body == 24 {
			ret = C.swe_calc_ut(julday, 11, 0, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = normalize(float64(xx[0]) + 180)
		} else {
			ret = C.swe_calc_ut(julday, body, 0, &xx[0], (*C.char)(unsafe.Pointer(&serr[0])))
			degreeUt = float64(xx[0])
		}

		if ret < 0 {
			log.Fatal(string(serr))
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
		for _, body2 := range c.Bodies {
			ascendant := c.AscMCs[0].DegreeUt
			for _, s := range aspectsettings {
				aspect := makeAspect(body1, body2, ascendant, s.delta, s.orb, s.title)
				if aspect != (Aspect{}) {
					c.Aspects = append(c.Aspects, aspect)
				}
			}
		}
	}

	// Sort bodies on DegreeUt
	sort.Slice(c.Bodies, func(i, j int) bool {
		return c.Bodies[i].DegreeUt < c.Bodies[j].DegreeUt
	})

	// Bodies distance
	oldDeg := -1000.
	dist := 0
	for i, body := range c.Bodies {
		deg := body.DegreeUt - c.AscMCs[0].DegreeUt + 180
		if math.Abs(oldDeg-deg) < 5 {
			dist++
		} else {
			dist = 0
		}
		c.Bodies[i].Dist = dist
		oldDeg = deg
	}

	out, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	out = []byte("<?xml version='1.0' encoding='UTF-8'?>" + string(out))
	w.Write(out)
}

// TransformHandler performs an XSLT transformation
func TransformHandler(w http.ResponseWriter, r *http.Request) {

	XMLURI := r.URL.Query().Get("xml")

	XMLResponse, err := http.Get(XMLURI)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer XMLResponse.Body.Close()

	XMLContent, err := ioutil.ReadAll(XMLResponse.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	XSLURI := r.URL.Query().Get("xsl")

	XSLContent, err := gk.ReadFile(XSLURI, gk.StrictParseOption)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer XSLContent.Free()

	parsedXSL, err := xslt.ParseStylesheet(XSLContent, XSLURI)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	parsedXML, err := gk.Parse(XMLContent, gk.DefaultEncodingBytes, nil, gk.DefaultParseOption, gk.DefaultEncodingBytes)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer parsedXML.Free()

	out, err := parsedXSL.Process(parsedXML, xslt.StylesheetOptions{
		IndentOutput: true,
		Parameters:   nil,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(out))
}

func sweSetEphePath(path string) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.swe_set_ephe_path(cpath)
}

func sweClose() {
	C.swe_close()
}

func main() {
	sweSetEphePath("swe")
	defer sweClose()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/chartinfo.py", ChartInfoHandler)
	http.HandleFunc("/chartinfo", ChartInfoHandler)
	http.HandleFunc("/transform.py", TransformHandler)
	http.HandleFunc("/transform", TransformHandler)

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("$PORT not set")
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
