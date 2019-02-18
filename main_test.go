package main

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_sliceAtoi(t *testing.T) {
	type args struct {
		sa []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{name: "Simple slice", args: args{sa: []string{"1", "2", "3"}}, want: []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sliceAtoi(tt.args.sa)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceAtoi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceAtoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalize(t *testing.T) {
	type args struct {
		angle float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "Angle greater than 360", args: args{angle: 362}, want: 2},
		{name: "Angle smaller than 0", args: args{angle: -2}, want: 358},
		{name: "Angle greater than many times 360", args: args{angle: 2000}, want: 200},
		{name: "Angle smaller than many times -360", args: args{angle: -2000}, want: 160},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.angle); got != tt.want {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Presence", args: args{s: []int{0, 1, 2, 3}, e: 2}, want: true},
		{name: "Absence", args: args{s: []int{0, 1, 2, 3}, e: 4}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeAspect(t *testing.T) {
	type args struct {
		body1     Body
		body2     Body
		ascendant float64
		delta     float64
		orb       float64
		t         string
	}
	tests := []struct {
		name       string
		args       args
		wantAspect Aspect
	}{
		{
			name: "Simple opposition",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 0},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 180},
				ascendant: 0,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 180,
				Degree2: 0,
			},
		},
		{
			name: "Opposition with negative angle",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 175},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: -5},
				ascendant: 0,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 355,
				Degree2: 175,
			},
		},
		{
			name: "Opposition with angle greater than 360",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 365},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 185},
				ascendant: 0,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 185,
				Degree2: 5,
			},
		},
		{
			name: "Opposition with an ascendant",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 130},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: -50},
				ascendant: 50,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 260,
				Degree2: 80,
			},
		},
		{
			name: "Opposition with an negative ascendant",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 260},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 80},
				ascendant: -100,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 180,
				Degree2: 0,
			},
		},
		{
			name: "Simple opposition within orb",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 0},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 189},
				ascendant: 0,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Opposition"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 180,
				Degree2: 9,
			},
		},
		{
			name: "Simple opposition out of orb",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: -15},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 180},
				ascendant: 0,
				delta:     180,
				orb:       10,
				t:         "Opposition",
			},
			wantAspect: Aspect{},
		},
		{
			name: "Complex square",
			args: args{
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 75},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 340},
				ascendant: 0,
				delta:     90,
				orb:       6,
				t:         "Square",
			},
			wantAspect: Aspect{
				XMLName: xml.Name{Local: "Square"},
				Body1:   "Sun",
				Body2:   "Moon",
				Degree1: 255,
				Degree2: 160,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAspect := makeAspect(tt.args.body1, tt.args.body2, tt.args.ascendant, tt.args.delta, tt.args.orb, tt.args.t); !reflect.DeepEqual(gotAspect, tt.wantAspect) {
				t.Errorf("makeAspect() = %v, want %v", gotAspect, tt.wantAspect)
			}
		})
	}
}

func TestChartInfoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/chartinfo.py?name=&city=(null)&country=(null)&lat=0.000000&lon=0.000000&year=2019&month=2&day=18&time=16.083334&hsys=E&display,0,1,2,3,4,5,6,7,8,9,10,12,23&tz=Asia/Saigon", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChartInfoHandler)

	handler.ServeHTTP(rr, req)

	want := `<?xml version='1.0' encoding='UTF-8'?><chartinfo display="1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23" year="2019" month="2" day="18" time="16.083334" city="(null)" hsys="E">
  <ascmcs>
    <Ascendant sign_name="Cancer" degree_ut="117.50872494334897" degree="27.508724943348966" sign="3" id="1"></Ascendant>
    <MC sign_name="Taurus" degree_ut="31.741532888977215" degree="1.7415328889772148" sign="1" id="2"></MC>
    <ARMC sign_name="Aries" degree_ut="29.57849084041179" degree="29.57849084041179" sign="0" id="3"></ARMC>
    <Vertex sign_name="Aries" degree_ut="0" degree="0" sign="0" id="4"></Vertex>
    <EquatorialAscendant sign_name="Cancer" degree_ut="117.50872494334894" degree="27.508724943348938" sign="3" id="5"></EquatorialAscendant>
    <Co-Ascendant1 sign_name="Cancer" degree_ut="117.50872494334897" degree="27.508724943348966" sign="3" id="6"></Co-Ascendant1>
    <Co-Ascendant2 sign_name="Virgo" degree_ut="180" degree="30" sign="5" id="7"></Co-Ascendant2>
    <Co-Ascendant2 sign_name="Libra" degree_ut="180" degree="0" sign="6" id="7"></Co-Ascendant2>
    <PolarAscendant sign_name="Capricorn" degree_ut="297.50872494334897" degree="27.508724943348966" sign="9" id="8"></PolarAscendant>
  </ascmcs>
  <houses>
    <House sign_name="Cancer" degree="27.508724943348966" number="I" sign="3" id="1" degree_ut="117.50872494334897"></House>
    <House sign_name="Leo" degree="27.508724943348966" number="II" sign="4" id="2" degree_ut="147.50872494334897"></House>
    <House sign_name="Virgo" degree="27.508724943348966" number="III" sign="5" id="3" degree_ut="177.50872494334897"></House>
    <House sign_name="Libra" degree="27.508724943348966" number="IV" sign="6" id="4" degree_ut="207.50872494334897"></House>
    <House sign_name="Scorpio" degree="27.508724943348966" number="V" sign="7" id="5" degree_ut="237.50872494334897"></House>
    <House sign_name="Sagittarius" degree="27.508724943348966" number="VI" sign="8" id="6" degree_ut="267.50872494334897"></House>
    <House sign_name="Capricorn" degree="27.508724943348966" number="VII" sign="9" id="7" degree_ut="297.50872494334897"></House>
    <House sign_name="Aquarius" degree="27.508724943348966" number="VIII" sign="10" id="8" degree_ut="327.50872494334897"></House>
    <House sign_name="Pisces" degree="27.508724943348966" number="IX" sign="11" id="9" degree_ut="357.50872494334897"></House>
    <House sign_name="Aries" degree="27.508724943348966" number="X" sign="0" id="10" degree_ut="27.508724943348966"></House>
    <House sign_name="Taurus" degree="27.508724943348966" number="XI" sign="1" id="11" degree_ut="57.508724943348966"></House>
    <House sign_name="Gemini" degree="27.508724943348966" number="XII" sign="2" id="12" degree_ut="87.50872494334897"></House>
  </houses>
  <aspects>
    <Conjunction body1="Sun" body2="MeanApogee" degree1="32.19789409202457" degree2="24.253872912601196"></Conjunction>
    <Conjunction body1="Sun" body2="OscuApogee" degree1="32.19789409202457" degree2="30.155917318864"></Conjunction>
    <Conjunction body1="Sun" body2="InterpretedApogee" degree1="32.19789409202457" degree2="22.776064798135508"></Conjunction>
    <Quincunx body1="Moon" body2="Mercury" degree1="198.05918019893105" degree2="47.05335285723277"></Quincunx>
    <Quincunx body1="Moon" body2="Neptune" degree1="198.05918019893105" degree2="48.037820595324035"></Quincunx>
    <Opposition body1="Moon" body2="MeanApogee" degree1="198.05918019893105" degree2="24.253872912601196"></Opposition>
    <Opposition body1="Moon" body2="InterpretedApogee" degree1="198.05918019893105" degree2="22.776064798135508"></Opposition>
    <Conjunction body1="Mercury" body2="Vesta" degree1="47.05335285723277" degree2="41.029149088999645"></Conjunction>
    <Quincunx body1="Venus" body2="Moon" degree1="349.51497520957923" degree2="198.05918019893105"></Quincunx>
    <Sextile body1="Venus" body2="Mercury" degree1="349.51497520957923" degree2="47.05335285723277"></Sextile>
    <Conjunction body1="Venus" body2="Saturn" degree1="349.51497520957923" degree2="349.28198299152405"></Conjunction>
    <Sextile body1="Venus" body2="Neptune" degree1="349.51497520957923" degree2="48.037820595324035"></Sextile>
    <Opposition body1="Venus" body2="MeanNode" degree1="349.51497520957923" degree2="177.48197354436465"></Opposition>
    <Opposition body1="Venus" body2="TrueNode" degree1="349.51497520957923" degree2="178.96236317353817"></Opposition>
    <Sextile body1="Mars" body2="Sun" degree1="95.34012060660855" degree2="32.19789409202457"></Sextile>
    <Conjunction body1="Mars" body2="Uranus" degree1="95.34012060660855" degree2="91.86644519578465"></Conjunction>
    <Trine body1="Jupiter" body2="Moon" degree1="323.03629424557374" degree2="198.05918019893105"></Trine>
    <Square body1="Jupiter" body2="Mercury" degree1="323.03629424557374" degree2="47.05335285723277"></Square>
    <Square body1="Jupiter" body2="Neptune" degree1="323.03629424557374" degree2="48.037820595324035"></Square>
    <Sextile body1="Jupiter" body2="MeanApogee" degree1="323.03629424557374" degree2="24.253872912601196"></Sextile>
    <Sextile body1="Jupiter" body2="InterpretedApogee" degree1="323.03629424557374" degree2="22.776064798135508"></Sextile>
    <Trine body1="Jupiter" body2="InterpretedPerigee" degree1="323.03629424557374" degree2="208.44234428598787"></Trine>
    <Quincunx body1="Saturn" body2="Moon" degree1="349.28198299152405" degree2="198.05918019893105"></Quincunx>
    <Sextile body1="Saturn" body2="Mercury" degree1="349.28198299152405" degree2="47.05335285723277"></Sextile>
    <Sextile body1="Saturn" body2="Neptune" degree1="349.28198299152405" degree2="48.037820595324035"></Sextile>
    <Opposition body1="Saturn" body2="MeanNode" degree1="349.28198299152405" degree2="177.48197354436465"></Opposition>
    <Opposition body1="Saturn" body2="TrueNode" degree1="349.28198299152405" degree2="178.96236317353817"></Opposition>
    <Sextile body1="Uranus" body2="Sun" degree1="91.86644519578465" degree2="32.19789409202457"></Sextile>
    <Sextile body1="Uranus" body2="OscuApogee" degree1="91.86644519578465" degree2="30.155917318864"></Sextile>
    <Semi-sextile body1="Uranus" body2="Earth" degree1="91.86644519578465" degree2="62.491275056651034"></Semi-sextile>
    <Semi-sextile body1="Uranus" body2="Chiron" degree1="91.86644519578465" degree2="62.50655412902415"></Semi-sextile>
    <Conjunction body1="Neptune" body2="Mercury" degree1="48.037820595324035" degree2="47.05335285723277"></Conjunction>
    <Conjunction body1="Neptune" body2="Vesta" degree1="48.037820595324035" degree2="41.029149088999645"></Conjunction>
    <Conjunction body1="Pluto" body2="Venus" degree1="354.65623396399155" degree2="349.51497520957923"></Conjunction>
    <Conjunction body1="Pluto" body2="Saturn" degree1="354.65623396399155" degree2="349.28198299152405"></Conjunction>
    <Opposition body1="Pluto" body2="MeanNode" degree1="354.65623396399155" degree2="177.48197354436465"></Opposition>
    <Opposition body1="Pluto" body2="TrueNode" degree1="354.65623396399155" degree2="178.96236317353817"></Opposition>
    <Semi-sextile body1="Pluto" body2="MeanApogee" degree1="354.65623396399155" degree2="24.253872912601196"></Semi-sextile>
    <Square body1="MeanNode" body2="Uranus" degree1="177.48197354436465" degree2="91.86644519578465"></Square>
    <Trine body1="MeanNode" body2="Earth" degree1="177.48197354436465" degree2="62.491275056651034"></Trine>
    <Trine body1="MeanNode" body2="Chiron" degree1="177.48197354436465" degree2="62.50655412902415"></Trine>
    <Square body1="TrueNode" body2="Uranus" degree1="178.96236317353817" degree2="91.86644519578465"></Square>
    <Conjunction body1="TrueNode" body2="MeanNode" degree1="178.96236317353817" degree2="177.48197354436465"></Conjunction>
    <Quincunx body1="TrueNode" body2="OscuApogee" degree1="178.96236317353817" degree2="30.155917318864"></Quincunx>
    <Trine body1="TrueNode" body2="Earth" degree1="178.96236317353817" degree2="62.491275056651034"></Trine>
    <Trine body1="TrueNode" body2="Chiron" degree1="178.96236317353817" degree2="62.50655412902415"></Trine>
    <Conjunction body1="MeanApogee" body2="InterpretedApogee" degree1="24.253872912601196" degree2="22.776064798135508"></Conjunction>
    <Conjunction body1="OscuApogee" body2="MeanApogee" degree1="30.155917318864" degree2="24.253872912601196"></Conjunction>
    <Conjunction body1="OscuApogee" body2="InterpretedApogee" degree1="30.155917318864" degree2="22.776064798135508"></Conjunction>
    <Semi-sextile body1="Earth" body2="Sun" degree1="62.491275056651034" degree2="32.19789409202457"></Semi-sextile>
    <Semi-sextile body1="Chiron" body2="Sun" degree1="62.50655412902415" degree2="32.19789409202457"></Semi-sextile>
    <Conjunction body1="Chiron" body2="Earth" degree1="62.50655412902415" degree2="62.491275056651034"></Conjunction>
    <Sextile body1="Pholus" body2="Sun" degree1="334.93606064706603" degree2="32.19789409202457"></Sextile>
    <Trine body1="Pholus" body2="Mars" degree1="334.93606064706603" degree2="95.34012060660855"></Trine>
    <Trine body1="Pholus" body2="Uranus" degree1="334.93606064706603" degree2="91.86644519578465"></Trine>
    <Square body1="Pholus" body2="Earth" degree1="334.93606064706603" degree2="62.491275056651034"></Square>
    <Square body1="Pholus" body2="Chiron" degree1="334.93606064706603" degree2="62.50655412902415"></Square>
    <Sextile body1="Pholus" body2="Pallas" degree1="334.93606064706603" degree2="271.9800016978568"></Sextile>
    <Quincunx body1="Pholus" body2="Juno" degree1="334.93606064706603" degree2="125.22217151509598"></Quincunx>
    <Trine body1="Pholus" body2="InterpretedPerigee" degree1="334.93606064706603" degree2="208.44234428598787"></Trine>
    <Trine body1="Ceres" body2="Earth" degree1="309.6874361812895" degree2="62.491275056651034"></Trine>
    <Trine body1="Ceres" body2="Chiron" degree1="309.6874361812895" degree2="62.50655412902415"></Trine>
    <Opposition body1="Ceres" body2="Juno" degree1="309.6874361812895" degree2="125.22217151509598"></Opposition>
    <Square body1="Ceres" body2="Vesta" degree1="309.6874361812895" degree2="41.029149088999645"></Square>
    <Trine body1="Pallas" body2="Sun" degree1="271.9800016978568" degree2="32.19789409202457"></Trine>
    <Opposition body1="Pallas" body2="Mars" degree1="271.9800016978568" degree2="95.34012060660855"></Opposition>
    <Opposition body1="Pallas" body2="Uranus" degree1="271.9800016978568" degree2="91.86644519578465"></Opposition>
    <Square body1="Pallas" body2="MeanNode" degree1="271.9800016978568" degree2="177.48197354436465"></Square>
    <Square body1="Pallas" body2="TrueNode" degree1="271.9800016978568" degree2="178.96236317353817"></Square>
    <Trine body1="Pallas" body2="MeanApogee" degree1="271.9800016978568" degree2="24.253872912601196"></Trine>
    <Trine body1="Pallas" body2="OscuApogee" degree1="271.9800016978568" degree2="30.155917318864"></Trine>
    <Quincunx body1="Pallas" body2="Earth" degree1="271.9800016978568" degree2="62.491275056651034"></Quincunx>
    <Quincunx body1="Pallas" body2="Chiron" degree1="271.9800016978568" degree2="62.50655412902415"></Quincunx>
    <Sextile body1="Pallas" body2="InterpretedPerigee" degree1="271.9800016978568" degree2="208.44234428598787"></Sextile>
    <Square body1="Juno" body2="Sun" degree1="125.22217151509598" degree2="32.19789409202457"></Square>
    <Semi-sextile body1="Juno" body2="Mars" degree1="125.22217151509598" degree2="95.34012060660855"></Semi-sextile>
    <Square body1="Juno" body2="OscuApogee" degree1="125.22217151509598" degree2="30.155917318864"></Square>
    <Sextile body1="Juno" body2="Earth" degree1="125.22217151509598" degree2="62.491275056651034"></Sextile>
    <Sextile body1="Juno" body2="Chiron" degree1="125.22217151509598" degree2="62.50655412902415"></Sextile>
    <Square body1="Juno" body2="Vesta" degree1="125.22217151509598" degree2="41.029149088999645"></Square>
    <Conjunction body1="Vesta" body2="Sun" degree1="41.029149088999645" degree2="32.19789409202457"></Conjunction>
    <Opposition body1="InterpretedPerigee" body2="Sun" degree1="208.44234428598787" degree2="32.19789409202457"></Opposition>
    <Trine body1="InterpretedPerigee" body2="Mars" degree1="208.44234428598787" degree2="95.34012060660855"></Trine>
    <Trine body1="InterpretedPerigee" body2="Uranus" degree1="208.44234428598787" degree2="91.86644519578465"></Trine>
    <Semi-sextile body1="InterpretedPerigee" body2="MeanNode" degree1="208.44234428598787" degree2="177.48197354436465"></Semi-sextile>
    <Semi-sextile body1="InterpretedPerigee" body2="TrueNode" degree1="208.44234428598787" degree2="178.96236317353817"></Semi-sextile>
    <Opposition body1="InterpretedPerigee" body2="MeanApogee" degree1="208.44234428598787" degree2="24.253872912601196"></Opposition>
    <Opposition body1="InterpretedPerigee" body2="OscuApogee" degree1="208.44234428598787" degree2="30.155917318864"></Opposition>
    <Opposition body1="InterpretedPerigee" body2="InterpretedApogee" degree1="208.44234428598787" degree2="22.776064798135508"></Opposition>
  </aspects>
  <bodies>
    <Earth sign_name="Aries" dist="0" degree_ut="0" degree="0" sign="0" retrograde="false" id="14"></Earth>
    <Chiron sign_name="Aries" dist="1" degree_ut="0.015279072373118058" degree="0.015279072373118058" sign="0" retrograde="false" id="15"></Chiron>
    <Uranus sign_name="Aries" dist="0" degree_ut="29.375170139133623" degree="29.375170139133623" sign="0" retrograde="false" id="7"></Uranus>
    <Mars sign_name="Taurus" dist="1" degree_ut="32.84884554995752" degree="2.8488455499575167" sign="1" retrograde="false" id="4"></Mars>
    <Juno sign_name="Gemini" dist="0" degree_ut="62.730896458444946" degree="2.7308964584449456" sign="2" retrograde="false" id="19"></Juno>
    <MeanNode sign_name="Cancer" dist="0" degree_ut="114.99069848771362" degree="24.990698487713615" sign="3" retrograde="false" id="10"></MeanNode>
    <TrueNode sign_name="Cancer" dist="1" degree_ut="116.47108811688712" degree="26.471088116887117" sign="3" retrograde="false" id="11"></TrueNode>
    <Moon sign_name="Leo" dist="0" degree_ut="135.56790514228" degree="15.567905142280011" sign="4" retrograde="false" id="1"></Moon>
    <InterpretedPerigee sign_name="Leo" dist="0" degree_ut="145.95106922933684" degree="25.951069229336838" sign="4" retrograde="false" id="22"></InterpretedPerigee>
    <Pallas sign_name="Libra" dist="0" degree_ut="209.48872664120572" degree="29.488726641205716" sign="6" retrograde="false" id="18"></Pallas>
    <Ceres sign_name="Sagittarius" dist="0" degree_ut="247.1961611246385" degree="7.1961611246385075" sign="8" retrograde="false" id="17"></Ceres>
    <Jupiter sign_name="Sagittarius" dist="0" degree_ut="260.5450191889227" degree="20.54501918892271" sign="8" retrograde="false" id="5"></Jupiter>
    <Pholus sign_name="Capricorn" dist="0" degree_ut="272.444785590415" degree="2.4447855904149947" sign="9" retrograde="false" id="16"></Pholus>
    <Saturn sign_name="Capricorn" dist="0" degree_ut="286.790707934873" degree="16.790707934873012" sign="9" retrograde="false" id="6"></Saturn>
    <Venus sign_name="Capricorn" dist="1" degree_ut="287.0237001529282" degree="17.0237001529282" sign="9" retrograde="false" id="3"></Venus>
    <Pluto sign_name="Capricorn" dist="0" degree_ut="292.1649589073405" degree="22.16495890734052" sign="9" retrograde="false" id="9"></Pluto>
    <InterpretedApogee sign_name="Aquarius" dist="0" degree_ut="320.2847897414845" degree="20.284789741484474" sign="10" retrograde="false" id="21"></InterpretedApogee>
    <MeanApogee sign_name="Aquarius" dist="1" degree_ut="321.76259785595016" degree="21.762597855950162" sign="10" retrograde="false" id="12"></MeanApogee>
    <OscuApogee sign_name="Aquarius" dist="0" degree_ut="327.66464226221296" degree="27.664642262212965" sign="10" retrograde="false" id="13"></OscuApogee>
    <Sun sign_name="Aquarius" dist="1" degree_ut="329.70661903537354" degree="29.706619035373535" sign="10" retrograde="false" id="0"></Sun>
    <Vesta sign_name="Pisces" dist="0" degree_ut="338.5378740323486" degree="8.537874032348611" sign="11" retrograde="false" id="20"></Vesta>
    <Mercury sign_name="Pisces" dist="0" degree_ut="344.56207780058173" degree="14.562077800581733" sign="11" retrograde="false" id="2"></Mercury>
    <Neptune sign_name="Pisces" dist="1" degree_ut="345.546545538673" degree="15.546545538673001" sign="11" retrograde="false" id="8"></Neptune>
  </bodies>
</chartinfo>`

	got := rr.Body.String()
	if got != want {
		t.Errorf("handler returned wrong xml: got %v want %v",
			got, want)
	}
}
