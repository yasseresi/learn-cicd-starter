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
			expectedError: errors.New("no authorization header included"),
		},
		{
			name: "Empty Authorization header",
			headers: http.Header{
				"Authorization": []string{""},
			},
			expectedKey:   "",
			expectedError: errors.New("no authorization header included"),
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
				"Authorization": []string{"ApiKey abc123"},
			},
			expectedKey:   "abc123",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey, err := GetAPIKey(tt.headers)
			if apiKey != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, apiKey)
			}
			if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if tt.expectedError != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError != nil && err != nil {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("expected error message %q, got %q", tt.expectedError.Error(), err.Error())
				}
			}
		})
	}
}
