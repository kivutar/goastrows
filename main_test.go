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
				body1:     Body{XMLName: xml.Name{Local: "Sun"}, DegreeUt: 180},
				body2:     Body{XMLName: xml.Name{Local: "Moon"}, DegreeUt: 0},
				ascendant: 180,
				delta:     10,
				orb:       180,
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

	req, err := http.NewRequest("GET", "/chartinfo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChartInfoHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
