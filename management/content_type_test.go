package management

import (
	"encoding/json"
	"testing"
)

func TestFlexString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    FlexString
		wantErr bool
	}{
		{
			name:  "string value",
			input: `"/:title"`,
			want:  "/:title",
		},
		{
			name:  "empty string",
			input: `""`,
			want:  "",
		},
		{
			name:  "boolean false treated as empty string",
			input: `false`,
			want:  "",
		},
		{
			name:  "boolean true treated as empty string",
			input: `true`,
			want:  "",
		},
		{
			name:    "number is invalid",
			input:   `42`,
			wantErr: true,
		},
		{
			// JSON null unmarshals into bool as false in Go, so null is treated
			// as empty string rather than an error.
			name:  "null treated as empty string",
			input: `null`,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FlexString
			err := json.Unmarshal([]byte(tt.input), &f)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && f != tt.want {
				t.Errorf("UnmarshalJSON() = %q, want %q", f, tt.want)
			}
		})
	}
}

func TestContentTypeOptions_UnmarshalJSON(t *testing.T) {
	t.Run("url_pattern and url_prefix as strings", func(t *testing.T) {
		input := `{
			"title": "My Type",
			"is_page": true,
			"singleton": false,
			"url_pattern": "/:title",
			"url_prefix": "/"
		}`
		var opts ContentTypeOptions
		if err := json.Unmarshal([]byte(input), &opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.UrlPattern != "/:title" {
			t.Errorf("UrlPattern = %q, want %q", opts.UrlPattern, "/:title")
		}
		if opts.UrlPrefix != "/" {
			t.Errorf("UrlPrefix = %q, want %q", opts.UrlPrefix, "/")
		}
	})

	t.Run("url_pattern and url_prefix as boolean false (unset)", func(t *testing.T) {
		input := `{
			"title": "My Type",
			"is_page": false,
			"singleton": false,
			"url_pattern": false,
			"url_prefix": false
		}`
		var opts ContentTypeOptions
		if err := json.Unmarshal([]byte(input), &opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.UrlPattern != "" {
			t.Errorf("UrlPattern = %q, want empty string", opts.UrlPattern)
		}
		if opts.UrlPrefix != "" {
			t.Errorf("UrlPrefix = %q, want empty string", opts.UrlPrefix)
		}
	})

	t.Run("singleton content type without url fields", func(t *testing.T) {
		input := `{
			"title": "Header",
			"is_page": false,
			"singleton": true,
			"url_pattern": false,
			"url_prefix": false
		}`
		var opts ContentTypeOptions
		if err := json.Unmarshal([]byte(input), &opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !opts.Singleton {
			t.Error("Singleton = false, want true")
		}
		if opts.UrlPattern != "" {
			t.Errorf("UrlPattern = %q, want empty string", opts.UrlPattern)
		}
	})
}
