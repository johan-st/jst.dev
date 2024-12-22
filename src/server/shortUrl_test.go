package server

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_validateAndParseUrl(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "valid url with scheme",
			args: args{u: "https://example.com/path"},
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/path",
			},
			wantErr: false,
		},
		{
			name: "valid url without scheme",
			args: args{u: "example.com"},
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name:    "empty url",
			args:    args{u: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing host",
			args:    args{u: "https:///path"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "url with query params",
			args: args{u: "example.com/path?key=value"},
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/path",
				RawQuery: "key=value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndParseUrl(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndParseUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateAndParseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
