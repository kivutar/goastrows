package main

import (
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
