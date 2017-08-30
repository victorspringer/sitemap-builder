package smb

import (
	"reflect"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	type args struct {
		host     string
		path     string
		filename string
		compress bool
		verbose  bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Builder
		wantErr bool
	}{
		{
			"ShouldReturnValidBuilderClient",
			args{
				host:     "http://example.com",
				path:     "/sitemaps",
				filename: "example",
				compress: true,
				verbose:  false,
			},
			&Builder{
				host:       "http://example.com",
				publicPath: "/sitemaps",
				filename:   "example",
				compress:   true,
				verbose:    false,
			},
			false,
		},
		{
			"ShouldReturnError",
			args{
				host:     "http://example.com",
				path:     "/sitemaps",
				compress: true,
				verbose:  false,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuilder(tt.args.host, tt.args.path, tt.args.filename, tt.args.compress, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
