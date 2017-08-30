package smb

import (
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

func TestRemoveFromURLSet(t *testing.T) {
	type args struct {
		urlset []*URL
		n      int
	}
	tests := []struct {
		name string
		args args
		want []*URL
	}{
		{
			"ShouldReturnValidSlice",
			args{
				[]*URL{
					&URL{
						Location:         "http://example.com/page1",
						Mobile:           true,
						ChangeFrequency:  "hourly",
						Priority:         1.0,
						LastModification: time.Now().Format(time.RFC3339),
					},
					&URL{
						Location:         "http://example.com/page2",
						Mobile:           true,
						ChangeFrequency:  "hourly",
						Priority:         0.8,
						LastModification: time.Now().Format(time.RFC3339),
					},
					&URL{
						Location:         "http://example.com/page3",
						Mobile:           true,
						ChangeFrequency:  "hourly",
						Priority:         0.7,
						LastModification: time.Now().Format(time.RFC3339),
					},
				},
				2,
			},
			[]*URL{
				&URL{
					Location:         "http://example.com/page3",
					Mobile:           true,
					ChangeFrequency:  "hourly",
					Priority:         0.7,
					LastModification: time.Now().Format(time.RFC3339),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeFromURLSet(tt.args.urlset, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeFromURLSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildURLMobile(t *testing.T) {
	type urlMobile struct {
		XMLName          xml.Name `xml:"url"`
		Location         string   `xml:"loc"`
		Mobile           string   `xml:"mobile:mobile,allowempty"`
		ChangeFrequency  string   `xml:"changefreq"`
		Priority         float32  `xml:"priority"`
		LastModification string   `xml:"lastmod"`
	}
	type args struct {
		url *URL
	}
	tests := []struct {
		name string
		args args
		want *urlMobile
	}{
		{
			"ShouldBuildValidURLMobileStruct",
			args{
				&URL{
					Location:         "http://example.com/page1",
					Mobile:           true,
					ChangeFrequency:  "hourly",
					Priority:         1.0,
					LastModification: time.Now().Format(time.RFC3339),
				},
			},
			&urlMobile{
				Location:         "http://example.com/page1",
				Mobile:           "",
				ChangeFrequency:  "hourly",
				Priority:         1.0,
				LastModification: time.Now().Format(time.RFC3339),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildURLMobile(tt.args.url); got == tt.want {
				t.Errorf("buildURLMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}
