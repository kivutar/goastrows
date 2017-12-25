package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

/*
#include "swephexp.h"
#cgo LDFLAGS: -lswe -lm -ldl
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

func main() {

	http.HandleFunc("/chartinfo", func(w http.ResponseWriter, r *http.Request) {

		var chartinfo = &ChartInfo{}

		var xx [6]C.double
		var serr string
		var serrC *C.char = C.CString(serr)
		var julday C.double
		var cusp [37]C.double
		var ascmc [10]C.double
		var hsys C.int = 'E'
		chartinfo.Year = 1970
		chartinfo.Month = 1
		display := []int{0, 1, 2, 3, 4}

		if r.URL.Query().Get("hsys") != "" {
			hsys = C.int([]rune(r.URL.Query().Get("hsys"))[0])
		}

		if r.URL.Query().Get("year") != "" {
			i, err := strconv.ParseInt(r.URL.Query().Get("year"), 10, 64)

			if err != nil {
				fmt.Printf("error: %v\n", err)
			}

			chartinfo.Year = i
		}

		if r.URL.Query().Get("month") != "" {
			i, err := strconv.ParseInt(r.URL.Query().Get("month"), 10, 64)

			if err != nil {
				fmt.Printf("error: %v\n", err)
			}

			chartinfo.Month = i
		}

		// The number of houses is 12 except when using Gauquelin sectors
		var numhouses = 12
		if hsys == 'G' {
			numhouses = 36
		}

		julday = C.swe_julday(C.int(chartinfo.Year), C.int(chartinfo.Month), 8, 13.25, C.SE_GREG_CAL)

		C.swe_set_topo(43, 5, 0)

		C.swe_houses(julday, 43.13517, 5.848, hsys, (*C.double)(&cusp[0]), (*C.double)(&ascmc[0]))

		// AscMC
		for index := 0; index < C.SE_NASCMC; index++ {
			degreeUt := float64(ascmc[index])

			for sign := 0; sign < 12; sign++ {
				degLow := float64(sign * 30)
				degHigh := float64((sign + 1) * 30)
				if degreeUt >= degLow && degreeUt <= degHigh {

					chartinfo.AscMCs = append(chartinfo.AscMCs,
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

					chartinfo.Houses = append(chartinfo.Houses,
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
				break
			}

			var degreeUt float64
			if body == 23 {
				C.swe_calc_ut(julday, body, 10, &xx[0], serrC)
				degreeUt = normalize(float64(xx[0]) + 180)
			} else if body == 24 {
				C.swe_calc_ut(julday, 11, 0, &xx[0], serrC)
				degreeUt = normalize(float64(xx[0]) + 180)
			} else {
				C.swe_calc_ut(julday, body, 0, &xx[0], serrC)
				degreeUt = float64(xx[0])
			}

			retrograde := xx[3] < 0

			for sign := 0; sign < 12; sign++ {
				degLow := float64(sign * 30)
				degHigh := float64((sign + 1) * 30)
				if degreeUt >= degLow && degreeUt <= degHigh {

					chartinfo.Bodies = append(chartinfo.Bodies,
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
		for _, body1 := range chartinfo.Bodies {
			deg1 := body1.DegreeUt - chartinfo.AscMCs[0].DegreeUt + 180

			for _, body2 := range chartinfo.Bodies {
				deg2 := body2.DegreeUt - chartinfo.AscMCs[0].DegreeUt + 180

				testAspect(chartinfo, body1, body2, deg1, deg2, 180, 10, "Opposition")
				testAspect(chartinfo, body1, body2, deg1, deg2, 150, 2, "Quincunx")
				testAspect(chartinfo, body1, body2, deg1, deg2, 120, 8, "Trine")
				testAspect(chartinfo, body1, body2, deg1, deg2, 90, 6, "Square")
				testAspect(chartinfo, body1, body2, deg1, deg2, 60, 4, "Sextile")
				testAspect(chartinfo, body1, body2, deg1, deg2, 30, 1, "Semi-sextile")
				testAspect(chartinfo, body1, body2, deg1, deg2, 0, 10, "Conjunction")
			}
		}

		out, err := xml.MarshalIndent(chartinfo, "", "  ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(out)

		C.swe_close()
	})

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Errorf("$PORT not set")
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
