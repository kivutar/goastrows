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
	req, err := http.NewRequest("GET", "/chartinfo.py?name,&city,(null)&country,(null)&lat,0.000000&lon,0.000000&year,2019&month,2&day,18&time,16.083334&hsys,E&display,0,1,2,3,4,5,6,7,8,9,10,12,23&tz,Asia/Saigon", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChartInfoHandler)

	handler.ServeHTTP(rr, req)

	want := `<?xml version='1.0' encoding='UTF-8'?><chartinfo display="1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23" year="1970" month="1" day="1" hsys="E">
  <ascmcs>
    <Ascendant sign_name="Libra" degree_ut="191.12942456262076" degree="11.129424562620756" sign="6" id="1"></Ascendant>
    <MC sign_name="Cancer" degree_ut="99.40188603109995" degree="9.401886031099949" sign="3" id="2"></MC>
    <ARMC sign_name="Cancer" degree_ut="100.23081494650327" degree="10.230814946503273" sign="3" id="3"></ARMC>
    <Vertex sign_name="Aries" degree_ut="0" degree="0" sign="0" id="4"></Vertex>
    <EquatorialAscendant sign_name="Libra" degree_ut="191.12942456262076" degree="11.129424562620756" sign="6" id="5"></EquatorialAscendant>
    <Co-Ascendant1 sign_name="Libra" degree_ut="191.12942456262076" degree="11.129424562620756" sign="6" id="6"></Co-Ascendant1>
    <Co-Ascendant2 sign_name="Virgo" degree_ut="180" degree="30" sign="5" id="7"></Co-Ascendant2>
    <Co-Ascendant2 sign_name="Libra" degree_ut="180" degree="0" sign="6" id="7"></Co-Ascendant2>
    <PolarAscendant sign_name="Aries" degree_ut="11.129424562620759" degree="11.129424562620759" sign="0" id="8"></PolarAscendant>
  </ascmcs>
  <houses>
    <House sign_name="Libra" degree="11.129424562620756" number="I" sign="6" id="1" degree_ut="191.12942456262076"></House>
    <House sign_name="Scorpio" degree="11.129424562620756" number="II" sign="7" id="2" degree_ut="221.12942456262076"></House>
    <House sign_name="Sagittarius" degree="11.129424562620756" number="III" sign="8" id="3" degree_ut="251.12942456262076"></House>
    <House sign_name="Capricorn" degree="11.129424562620784" number="IV" sign="9" id="4" degree_ut="281.1294245626208"></House>
    <House sign_name="Aquarius" degree="11.129424562620784" number="V" sign="10" id="5" degree_ut="311.1294245626208"></House>
    <House sign_name="Pisces" degree="11.129424562620784" number="VI" sign="11" id="6" degree_ut="341.1294245626208"></House>
    <House sign_name="Aries" degree="11.129424562620784" number="VII" sign="0" id="7" degree_ut="11.129424562620784"></House>
    <House sign_name="Taurus" degree="11.129424562620784" number="VIII" sign="1" id="8" degree_ut="41.129424562620784"></House>
    <House sign_name="Gemini" degree="11.129424562620784" number="IX" sign="2" id="9" degree_ut="71.12942456262078"></House>
    <House sign_name="Cancer" degree="11.129424562620784" number="X" sign="3" id="10" degree_ut="101.12942456262078"></House>
    <House sign_name="Leo" degree="11.129424562620784" number="XI" sign="4" id="11" degree_ut="131.12942456262078"></House>
    <House sign_name="Virgo" degree="11.129424562620784" number="XII" sign="5" id="12" degree_ut="161.12942456262078"></House>
  </houses>
  <aspects>
    <Square body1="Sun" body2="Moon" degree1="269.0268594396571" degree2="179.5696251252196"></Square>
    <Conjunction body1="Sun" body2="Venus" degree1="269.0268594396571" degree2="263.3238424164573"></Conjunction>
    <Square body1="Sun" body2="Uranus" degree1="269.0268594396571" degree2="177.59120413790973"></Square>
    <Quincunx body1="Sun" body2="OscuApogee" degree1="269.0268594396571" degree2="120.97145833629688"></Quincunx>
    <Conjunction body1="Moon" body2="Uranus" degree1="179.5696251252196" degree2="177.59120413790973"></Conjunction>
    <Sextile body1="Moon" body2="OscuApogee" degree1="179.5696251252196" degree2="120.97145833629688"></Sextile>
    <Sextile body1="Moon" body2="InterpretedApogee" degree1="179.5696251252196" degree2="116.09777594684972"></Sextile>
    <Square body1="Mercury" body2="Jupiter" degree1="287.8920436548366" degree2="201.19559555304576"></Square>
    <Square body1="Mercury" body2="Saturn" degree1="287.8920436548366" degree2="20.9328864223898"></Square>
    <Sextile body1="Mercury" body2="Neptune" degree1="287.8920436548366" degree2="228.7564433086492"></Sextile>
    <Trine body1="Mercury" body2="Pluto" degree1="287.8920436548366" degree2="166.26273992846998"></Trine>
    <Opposition body1="Mercury" body2="MeanApogee" degree1="287.8920436548366" degree2="111.63538553494561"></Opposition>
    <Conjunction body1="Mercury" body2="Pallas" degree1="287.8920436548366" degree2="286.9402433480859"></Conjunction>
    <Opposition body1="Mercury" body2="InterpretedApogee" degree1="287.8920436548366" degree2="116.09777594684972"></Opposition>
    <Conjunction body1="Mercury" body2="InterpretedPerigee" degree1="287.8920436548366" degree2="280.37635039643567"></Conjunction>
    <Sextile body1="Venus" body2="Jupiter" degree1="263.3238424164573" degree2="201.19559555304576"></Sextile>
    <Trine body1="Venus" body2="Saturn" degree1="263.3238424164573" degree2="20.9328864223898"></Trine>
    <Square body1="Venus" body2="Uranus" degree1="263.3238424164573" degree2="177.59120413790973"></Square>
    <Quincunx body1="Venus" body2="MeanApogee" degree1="263.3238424164573" degree2="111.63538553494561"></Quincunx>
    <Sextile body1="Mars" body2="Sun" degree1="331.106825503049" degree2="269.0268594396571"></Sextile>
    <Quincunx body1="Mars" body2="Moon" degree1="331.106825503049" degree2="179.5696251252196"></Quincunx>
    <Quincunx body1="Mars" body2="OscuApogee" degree1="331.106825503049" degree2="120.97145833629688"></Quincunx>
    <Opposition body1="Jupiter" body2="Saturn" degree1="201.19559555304576" degree2="20.9328864223898"></Opposition>
    <Square body1="Jupiter" body2="MeanApogee" degree1="201.19559555304576" degree2="111.63538553494561"></Square>
    <Square body1="Jupiter" body2="InterpretedApogee" degree1="201.19559555304576" degree2="116.09777594684972"></Square>
    <Sextile body1="Uranus" body2="OscuApogee" degree1="177.59120413790973" degree2="120.97145833629688"></Sextile>
    <Sextile body1="Uranus" body2="InterpretedApogee" degree1="177.59120413790973" degree2="116.09777594684972"></Sextile>
    <Sextile body1="Neptune" body2="Pluto" degree1="228.7564433086492" degree2="166.26273992846998"></Sextile>
    <Trine body1="Neptune" body2="MeanApogee" degree1="228.7564433086492" degree2="111.63538553494561"></Trine>
    <Square body1="Neptune" body2="Vesta" degree1="228.7564433086492" degree2="134.99098809324295"></Square>
    <Trine body1="Neptune" body2="InterpretedApogee" degree1="228.7564433086492" degree2="116.09777594684972"></Trine>
    <Conjunction body1="MeanNode" body2="Mars" degree1="334.15744153189837" degree2="331.106825503049"></Conjunction>
    <Conjunction body1="MeanNode" body2="TrueNode" degree1="334.15744153189837" degree2="332.9305252185683"></Conjunction>
    <Sextile body1="TrueNode" body2="Sun" degree1="332.9305252185683" degree2="269.0268594396571"></Sextile>
    <Conjunction body1="TrueNode" body2="Mars" degree1="332.9305252185683" degree2="331.106825503049"></Conjunction>
    <Quincunx body1="TrueNode" body2="OscuApogee" degree1="332.9305252185683" degree2="120.97145833629688"></Quincunx>
    <Square body1="MeanApogee" body2="Saturn" degree1="111.63538553494561" degree2="20.9328864223898"></Square>
    <Conjunction body1="OscuApogee" body2="MeanApogee" degree1="120.97145833629688" degree2="111.63538553494561"></Conjunction>
    <Conjunction body1="OscuApogee" body2="InterpretedApogee" degree1="120.97145833629688" degree2="116.09777594684972"></Conjunction>
    <Sextile body1="Earth" body2="Mercury" degree1="348.8705754373792" degree2="287.8920436548366"></Sextile>
    <Square body1="Earth" body2="Venus" degree1="348.8705754373792" degree2="263.3238424164573"></Square>
    <Opposition body1="Earth" body2="Uranus" degree1="348.8705754373792" degree2="177.59120413790973"></Opposition>
    <Trine body1="Earth" body2="Neptune" degree1="348.8705754373792" degree2="228.7564433086492"></Trine>
    <Opposition body1="Earth" body2="Pluto" degree1="348.8705754373792" degree2="166.26273992846998"></Opposition>
    <Trine body1="Earth" body2="MeanApogee" degree1="348.8705754373792" degree2="111.63538553494561"></Trine>
    <Sextile body1="Earth" body2="Pallas" degree1="348.8705754373792" degree2="286.9402433480859"></Sextile>
    <Trine body1="Earth" body2="InterpretedApogee" degree1="348.8705754373792" degree2="116.09777594684972"></Trine>
    <Opposition body1="Chiron" body2="Moon" degree1="351.3909796929443" degree2="179.5696251252196"></Opposition>
    <Sextile body1="Chiron" body2="Mercury" degree1="351.3909796929443" degree2="287.8920436548366"></Sextile>
    <Square body1="Chiron" body2="Venus" degree1="351.3909796929443" degree2="263.3238424164573"></Square>
    <Quincunx body1="Chiron" body2="Jupiter" degree1="351.3909796929443" degree2="201.19559555304576"></Quincunx>
    <Semi-sextile body1="Chiron" body2="Saturn" degree1="351.3909796929443" degree2="20.9328864223898"></Semi-sextile>
    <Opposition body1="Chiron" body2="Uranus" degree1="351.3909796929443" degree2="177.59120413790973"></Opposition>
    <Trine body1="Chiron" body2="Neptune" degree1="351.3909796929443" degree2="228.7564433086492"></Trine>
    <Opposition body1="Chiron" body2="Pluto" degree1="351.3909796929443" degree2="166.26273992846998"></Opposition>
    <Trine body1="Chiron" body2="MeanApogee" degree1="351.3909796929443" degree2="111.63538553494561"></Trine>
    <Conjunction body1="Chiron" body2="Earth" degree1="351.3909796929443" degree2="348.8705754373792"></Conjunction>
    <Sextile body1="Chiron" body2="Juno" degree1="351.3909796929443" degree2="293.5476873633056"></Sextile>
    <Trine body1="Chiron" body2="InterpretedApogee" degree1="351.3909796929443" degree2="116.09777594684972"></Trine>
    <Trine body1="Pholus" body2="Jupiter" degree1="314.9121756198093" degree2="201.19559555304576"></Trine>
    <Square body1="Pholus" body2="Neptune" degree1="314.9121756198093" degree2="228.7564433086492"></Square>
    <Quincunx body1="Pholus" body2="Pluto" degree1="314.9121756198093" degree2="166.26273992846998"></Quincunx>
    <Conjunction body1="Pholus" body2="Ceres" degree1="314.9121756198093" degree2="309.99002380150625"></Conjunction>
    <Opposition body1="Pholus" body2="Vesta" degree1="314.9121756198093" degree2="134.99098809324295"></Opposition>
    <Opposition body1="Ceres" body2="OscuApogee" degree1="309.99002380150625" degree2="120.97145833629688"></Opposition>
    <Opposition body1="Ceres" body2="Vesta" degree1="309.99002380150625" degree2="134.99098809324295"></Opposition>
    <Semi-sextile body1="Ceres" body2="InterpretedPerigee" degree1="309.99002380150625" degree2="280.37635039643567"></Semi-sextile>
    <Square body1="Pallas" body2="Jupiter" degree1="286.9402433480859" degree2="201.19559555304576"></Square>
    <Square body1="Pallas" body2="Saturn" degree1="286.9402433480859" degree2="20.9328864223898"></Square>
    <Sextile body1="Pallas" body2="Neptune" degree1="286.9402433480859" degree2="228.7564433086492"></Sextile>
    <Trine body1="Pallas" body2="Pluto" degree1="286.9402433480859" degree2="166.26273992846998"></Trine>
    <Opposition body1="Pallas" body2="MeanApogee" degree1="286.9402433480859" degree2="111.63538553494561"></Opposition>
    <Quincunx body1="Pallas" body2="Vesta" degree1="286.9402433480859" degree2="134.99098809324295"></Quincunx>
    <Opposition body1="Pallas" body2="InterpretedApogee" degree1="286.9402433480859" degree2="116.09777594684972"></Opposition>
    <Conjunction body1="Pallas" body2="InterpretedPerigee" degree1="286.9402433480859" degree2="280.37635039643567"></Conjunction>
    <Trine body1="Juno" body2="Moon" degree1="293.5476873633056" degree2="179.5696251252196"></Trine>
    <Conjunction body1="Juno" body2="Mercury" degree1="293.5476873633056" degree2="287.8920436548366"></Conjunction>
    <Semi-sextile body1="Juno" body2="Venus" degree1="293.5476873633056" degree2="263.3238424164573"></Semi-sextile>
    <Square body1="Juno" body2="Jupiter" degree1="293.5476873633056" degree2="201.19559555304576"></Square>
    <Square body1="Juno" body2="Saturn" degree1="293.5476873633056" degree2="20.9328864223898"></Square>
    <Trine body1="Juno" body2="Uranus" degree1="293.5476873633056" degree2="177.59120413790973"></Trine>
    <Trine body1="Juno" body2="Pluto" degree1="293.5476873633056" degree2="166.26273992846998"></Trine>
    <Opposition body1="Juno" body2="MeanApogee" degree1="293.5476873633056" degree2="111.63538553494561"></Opposition>
    <Opposition body1="Juno" body2="OscuApogee" degree1="293.5476873633056" degree2="120.97145833629688"></Opposition>
    <Conjunction body1="Juno" body2="Pallas" degree1="293.5476873633056" degree2="286.9402433480859"></Conjunction>
    <Opposition body1="Juno" body2="InterpretedApogee" degree1="293.5476873633056" degree2="116.09777594684972"></Opposition>
    <Trine body1="Vesta" body2="Saturn" degree1="134.99098809324295" degree2="20.9328864223898"></Trine>
    <Square body1="InterpretedApogee" body2="Saturn" degree1="116.09777594684972" degree2="20.9328864223898"></Square>
    <Conjunction body1="InterpretedApogee" body2="MeanApogee" degree1="116.09777594684972" degree2="111.63538553494561"></Conjunction>
    <Trine body1="InterpretedPerigee" body2="Pluto" degree1="280.37635039643567" degree2="166.26273992846998"></Trine>
  </aspects>
  <bodies>
    <Earth sign_name="Aries" dist="0" degree_ut="0" degree="0" sign="0" retrograde="false" id="14"></Earth>
    <Chiron sign_name="Aries" dist="1" degree_ut="2.5204042555650723" degree="2.5204042555650723" sign="0" retrograde="false" id="15"></Chiron>
    <Saturn sign_name="Taurus" dist="0" degree_ut="32.06231098501057" degree="2.062310985010569" sign="1" retrograde="false" id="6"></Saturn>
    <MeanApogee sign_name="Leo" dist="0" degree_ut="122.76481009756637" degree="2.764810097566368" sign="4" retrograde="false" id="12"></MeanApogee>
    <InterpretedApogee sign_name="Leo" dist="1" degree_ut="127.22720050947048" degree="7.227200509470478" sign="4" retrograde="false" id="21"></InterpretedApogee>
    <OscuApogee sign_name="Leo" dist="2" degree_ut="132.10088289891763" degree="12.100882898917632" sign="4" retrograde="false" id="13"></OscuApogee>
    <Vesta sign_name="Leo" dist="0" degree_ut="146.1204126558637" degree="26.120412655863703" sign="4" retrograde="false" id="20"></Vesta>
    <Pluto sign_name="Virgo" dist="0" degree_ut="177.39216449109074" degree="27.39216449109074" sign="5" retrograde="false" id="9"></Pluto>
    <Uranus sign_name="Libra" dist="0" degree_ut="188.72062870053048" degree="8.720628700530483" sign="6" retrograde="false" id="7"></Uranus>
    <Moon sign_name="Libra" dist="1" degree_ut="190.69904968784036" degree="10.699049687840358" sign="6" retrograde="false" id="1"></Moon>
    <Jupiter sign_name="Scorpio" dist="0" degree_ut="212.32502011566652" degree="2.3250201156665184" sign="7" retrograde="false" id="5"></Jupiter>
    <Neptune sign_name="Scorpio" dist="0" degree_ut="239.88586787126997" degree="29.885867871269966" sign="7" retrograde="false" id="8"></Neptune>
    <Venus sign_name="Capricorn" dist="0" degree_ut="274.4532669790781" degree="4.45326697907808" sign="9" retrograde="false" id="3"></Venus>
    <Sun sign_name="Capricorn" dist="0" degree_ut="280.1562840022778" degree="10.1562840022778" sign="9" retrograde="false" id="0"></Sun>
    <InterpretedPerigee sign_name="Capricorn" dist="0" degree_ut="291.50577495905645" degree="21.505774959056453" sign="9" retrograde="false" id="22"></InterpretedPerigee>
    <Pallas sign_name="Capricorn" dist="0" degree_ut="298.0696679107067" degree="28.0696679107067" sign="9" retrograde="false" id="18"></Pallas>
    <Mercury sign_name="Capricorn" dist="1" degree_ut="299.02146821745737" degree="29.021468217457368" sign="9" retrograde="false" id="2"></Mercury>
    <Juno sign_name="Aquarius" dist="0" degree_ut="304.67711192592634" degree="4.67711192592634" sign="10" retrograde="false" id="19"></Juno>
    <Ceres sign_name="Aquarius" dist="0" degree_ut="321.11944836412704" degree="21.119448364127038" sign="10" retrograde="false" id="17"></Ceres>
    <Pholus sign_name="Aquarius" dist="1" degree_ut="326.0416001824301" degree="26.041600182430102" sign="10" retrograde="false" id="16"></Pholus>
    <Mars sign_name="Pisces" dist="0" degree_ut="342.23625006566976" degree="12.236250065669765" sign="11" retrograde="false" id="4"></Mars>
    <TrueNode sign_name="Pisces" dist="1" degree_ut="344.05994978118906" degree="14.059949781189061" sign="11" retrograde="false" id="11"></TrueNode>
    <MeanNode sign_name="Pisces" dist="2" degree_ut="345.28686609451916" degree="15.286866094519155" sign="11" retrograde="false" id="10"></MeanNode>
  </bodies>
</chartinfo>`

	got := rr.Body.String()
	if got != want {
		t.Errorf("handler returned wrong xml: got %v want %v",
			got, want)
	}
}
