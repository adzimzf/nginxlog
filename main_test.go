package main

import (
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}
func Test_getTopTen(t *testing.T) {
	tests := []struct {
		name       string
		logs       []string
		wantResult []string
	}{
		{
			name:       "failed input is empty",
			logs:       []string{},
			wantResult: []string{},
		},
		{
			name:       "the only one input",
			logs:       []string{"/articles.html?id=4, 0.12s, Status Code: 200"},
			wantResult: []string{"/articles.html?id=4"},
		},
		{
			name: "there're 2 post and 2 get",
			logs: []string{
				"/articles.html?id=4, 0.12s, Status Code: 200",
				"/articles.php/login, 0.12s, Status Code: 200",
				"/articles/user?id=5, 1.12s, Status Code: 200",
				"/articles.html/user/update, 0.12s, Status Code: 200",
			},
			wantResult: []string{
				"/articles/user?id=5",
				"/articles.html?id=4",
			},
		},
		{
			name: "there're 2 POST and 2 GET, but GET url has only diff case",
			logs: []string{
				"/articles.html?id=4, 0.12s, Status Code: 200",
				"/articles.php/login, 0.12s, Status Code: 200",
				"/articles.html?ID=4, 1.12s, Status Code: 200",
				"/articles.html/user/update, 0.12s, Status Code: 200",
			},
			wantResult: []string{
				"/articles.html?id=4",
			},
		},
		{
			name: "there're .gif url and 3 different GET url",
			logs: []string{
				"/vendor/bootstrap.min.js, 0.15s, Status Code: 200",
				"/vendor/bootstrap.min.css, 0.12s, Status Code: 200",
				"/img/logo/vendor_4.gif, 1.12s, Status Code: 200",
				"/img/logo/vendor_3.png, 0.12s, Status Code: 200",
			},
			wantResult: []string{
				"/vendor/bootstrap.min.js",
				"/vendor/bootstrap.min.css",
				"/img/logo/vendor_3.png",
			},
		},
		{
			name: "there're 5 urls, but there are only 3 url get Status Code 200",
			logs: []string{
				"/admin/user?ID=1, 3.01s, Status Code: 200",
				"/admin/user?ID=2, 3.02s, Status Code: 404",
				"/admin/user?ID=3, 3.03s, Status Code: 200",
				"/admin/user?ID=4, 3.04s, Status Code: 404",
				"/admin/user?ID=5, 3.05s, Status Code: 200",
			},
			wantResult: []string{
				"/admin/user?id=5",
				"/admin/user?id=3",
				"/admin/user?id=1",
			},
		},
		{
			name: "there're 15 urls, but there are only 11 valid url",
			logs: []string{
				"/admin/user?ID=1, 3.01s, Status Code: 200",
				"/admin/user?ID=2, 3.02s, Status Code: 404",
				"/admin/user?ID=3, 3.03s, Status Code: 200",
				"/admin/user?ID=4, 3.04s, Status Code: 404",
				"/admin/user?ID=5, 3.05s, Status Code: 200",
				"/admin/user?ID=6, 3.06s, Status Code: 200",
				"/admin/user?ID=7, 3.07s, Status Code: 200",
				"/admin/user?ID=8, 3.08s, Status Code: 200",
				"/admin/user?ID=9, 3.09s, Status Code: 200",
				"/admin/user?ID=10, 3.10s, Status Code: 200",
				"/admin/user?ID=11, 3.11s, Status Codee: 200",
				"/admin/user?ID=12, 3.12s, Status Code: 200",
				"/admin/user?ID=13, 3.13s, Status Code: 200",
				"/admin/user/update, 3.14s, Status Code: 200",
				"/admin/user?ID=15, 3.15s, Status Code: 200",
			},
			wantResult: []string{
				"/admin/user?id=15",
				"/admin/user?id=13",
				"/admin/user?id=12",
				"/admin/user?id=10",
				"/admin/user?id=9",
				"/admin/user?id=8",
				"/admin/user?id=7",
				"/admin/user?id=6",
				"/admin/user?id=5",
				"/admin/user?id=3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := getTopTen(tt.logs); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("getTopTen() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_getURLResTimeStaCode(t *testing.T) {
	tests := []struct {
		name         string
		log          string
		wantUrl      string
		wantRespTime float64
		wantStatCode string
		wantErr      bool
	}{
		{
			name:    "failed this is not nginx log",
			log:     "this not log",
			wantErr: true,
		},
		{
			name:    "failed empty nginx log",
			log:     "",
			wantErr: true,
		},
		{
			name:    "status code is not 200",
			log:     "/articles.html?id=1, 0.33s, Status Code: 401",
			wantErr: true,
		},
		{
			name:    ".gif file",
			log:     "/articles/politics.gif, 0.33s, Status Code: 401",
			wantErr: true,
		},
		{
			name:    "POST URL",
			log:     "/articles/politics/update, 0.33s, Status Code: 200",
			wantErr: true,
		},
		{
			name:         "success GET url",
			log:          "/articles/politics/update?id=4, 0.33s, Status Code: 200",
			wantErr:      false,
			wantUrl:      "/articles/politics/update?id=4",
			wantRespTime: 0.33,
			wantStatCode: "200",
		},
		{
			name:         "success GET url",
			log:          "/articles/jQuery.css, 0.33s, Status Code: 200",
			wantErr:      false,
			wantUrl:      "/articles/jquery.css",
			wantRespTime: 0.33,
			wantStatCode: "200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrl, gotRespTime, gotStatCode, err := getURLResTimeStaCode(tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("getURLResTimeStaCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if gotUrl != tt.wantUrl {
				t.Errorf("getURLResTimeStaCode() gotUrl = %v, want %v", gotUrl, tt.wantUrl)
			}
			if gotRespTime != tt.wantRespTime {
				t.Errorf("getURLResTimeStaCode() gotRespTime = %v, want %v", gotRespTime, tt.wantRespTime)
			}
			if gotStatCode != tt.wantStatCode {
				t.Errorf("getURLResTimeStaCode() gotStatCode = %v, want %v", gotStatCode, tt.wantStatCode)
			}
		})
	}
}
