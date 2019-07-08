package gitlab

import (
	"testing"
)

func TestPage_FormatFooter(t *testing.T) {
	type fields struct {
		FooterTitle string
		FooterURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "title and url",
			fields: fields{
				FooterTitle: "GitHub",
				FooterURL:   "https://github.com/",
			},
			want: "<https://github.com/|GitHub>",
		},
		{
			name: "title only",
			fields: fields{
				FooterTitle: "GitHub",
				FooterURL:   "",
			},
			want: "GitHub",
		},
		{
			name: "url only",
			fields: fields{
				FooterTitle: "",
				FooterURL:   "https://github.com/",
			},
			want: "https://github.com/",
		},
		{
			name: "nothing",
			fields: fields{
				FooterTitle: "",
				FooterURL:   "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Page{
				FooterTitle: tt.fields.FooterTitle,
				FooterURL:   tt.fields.FooterURL,
			}
			if got := p.FormatFooter(); got != tt.want {
				t.Errorf("Page.formatFooter() = %v, want %+v", got, tt.want)
			}
		})
	}
}
