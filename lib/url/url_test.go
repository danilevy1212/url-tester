package url

import (
	"testing"
)

func TestSanitizeURL(t *testing.T) {
	testCases := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{"Valid HTTP URL", "http://example.com", "http://example.com", false},
		{"Valid HTTPS URL", "https://example.com", "https://example.com", false},
		{"Valid HTTP URL with encoded path", "http://example.com/100% unencoded path", "http://example.com/100%25%20unencoded%20path", false},
		{"Unsupported scheme (FTP)", "ftp://example.com", "", true},
		{"No scheme", "example.com", "", true},
		{"Invalid URL", "http://%gh&%$", "", true},
		{"Empty string", "", "", true},
		{"URL with port", "http://example.com:8080/path", "http://example.com:8080/path", false},
		{"HTTPS URL with query", "https://example.com/path?query=param", "https://example.com/path?query=param", false},
		{"URL with fragment", "http://example.com/path#section", "http://example.com/path#section", false},
		{"URL with username and password", "http://user:pass@example.com", "http://user:pass@example.com", false},
		{"URL with special characters in path", "http://example.com/path/äöü", "http://example.com/path/%C3%A4%C3%B6%C3%BC", false},
		{"URL with multiple path segments", "http://example.com/path/to/resource", "http://example.com/path/to/resource", false},
		{"URL with encoded query parameters", "http://example.com?param1=value1&param2=value%202", "http://example.com?param1=value1&param2=value%202", false},
		{"URL with encoded query parameters and fragment", "http://example.com?param1=value1&param2=value%202#fragment", "http://example.com?param1=value1&param2=value%202#fragment", false},
		{"URL with just a scheme", "http://", "", true},
		{"URL with malformed scheme", "ht!tp://example.com", "", true},
		{"URL with space in hostname", "http://exam ple.com", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SanitizeURL(tc.rawURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("SanitizeURL(%q) error = %v, wantErr %v", tc.rawURL, err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("SanitizeURL(%q) = %q, want %q", tc.rawURL, got, tc.want)
			}
		})
	}
}
