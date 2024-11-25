package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedKey   string
		expectedError error
	}{
		{
			name:          "No Authorization header",
			headers:       http.Header{},
			expectedKey:   "",
			expectedError: ErrNoAuthHeaderIncluded,
		},
		{
			name: "Empty Authorization header",
			headers: http.Header{
				"Authorization": []string{""},
			},
			expectedKey:   "",
			expectedError: ErrNoAuthHeaderIncluded,
		},
		{
			name: "Malformed Authorization header - no space",
			headers: http.Header{
				"Authorization": []string{"ApiKey"},
			},
			expectedKey:   "",
			expectedError: errors.New("malformed authorization header"),
		},
		{
			name: "Malformed Authorization header - wrong scheme",
			headers: http.Header{
				"Authorization": []string{"Bearer token"},
			},
			expectedKey:   "",
			expectedError: errors.New("malformed authorization header"),
		},
		{
			name: "Valid Authorization header",
			headers: http.Header{
				"Authorization": []string{"ApiKey 12345"},
			},
			expectedKey:   "12345",
			expectedError: nil,
		},
		{
			name: "Valid Authorization header with extra spaces",
			headers: http.Header{
				"Authorization": []string{"ApiKey    12345"},
			},
			expectedKey:   "12345",
			expectedError: nil,
		},
		{
			name: "Authorization header with missing key",
			headers: http.Header{
				"Authorization": []string{"ApiKey "},
			},
			expectedKey:   "",
			expectedError: errors.New("malformed authorization header"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey, err := GetAPIKey(tt.headers)
			if apiKey != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, apiKey)
			}
			if tt.expectedError != nil && err != nil {
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %q, got %q", tt.expectedError.Error(), err.Error())
				}
			} else if tt.expectedError != err {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
