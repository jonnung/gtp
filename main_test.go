package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestHelpText(t *testing.T) {
	result := helpText()

	// Check that help text contains all expected commands
	expectedCommands := []string{
		"usage: gtp",
		"{number}",
		"list",
		"add",
		"remove",
		"clear",
	}

	for _, cmd := range expectedCommands {
		if !strings.Contains(result, cmd) {
			t.Errorf("helpText() should contain '%s', got: %s", cmd, result)
		}
	}

	// Check that help text is not empty
	if len(result) == 0 {
		t.Error("helpText() should not return empty string")
	}
}

func TestErrorOutOfRangeText(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{
			name:     "positive index",
			index:    5,
			expected: "ERRORRRRR: Selected number 5 is out of range OTP list\n",
		},
		{
			name:     "zero index",
			index:    0,
			expected: "ERRORRRRR: Selected number 0 is out of range OTP list\n",
		},
		{
			name:     "negative index",
			index:    -1,
			expected: "ERRORRRRR: Selected number -1 is out of range OTP list\n",
		},
		{
			name:     "large index",
			index:    100,
			expected: "ERRORRRRR: Selected number 100 is out of range OTP list\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorOutOfRangeText(tt.index)
			if result != tt.expected {
				t.Errorf("errorOutOfRangeText(%d) = %q, want %q", tt.index, result, tt.expected)
			}
		})
	}
}

func TestGetTotalOtpsInfo(t *testing.T) {
	tests := []struct {
		name     string
		otpList  []otp
		expected []string
	}{
		{
			name:     "empty list",
			otpList:  []otp{},
			expected: []string{},
		},
		{
			name: "single otp",
			otpList: []otp{
				{Issuer: "GitHub", AccountName: "user@example.com", Secret: "SECRET123"},
			},
			expected: []string{"{1} GitHub:user@example.com:<secret>"},
		},
		{
			name: "multiple otps",
			otpList: []otp{
				{Issuer: "GitHub", AccountName: "user1@example.com", Secret: "SECRET1"},
				{Issuer: "Google", AccountName: "user2@example.com", Secret: "SECRET2"},
				{Issuer: "AWS", AccountName: "user3@example.com", Secret: "SECRET3"},
			},
			expected: []string{
				"{1} GitHub:user1@example.com:<secret>",
				"{2} Google:user2@example.com:<secret>",
				"{3} AWS:user3@example.com:<secret>",
			},
		},
		{
			name: "otp with special characters in account name",
			otpList: []otp{
				{Issuer: "Test-Service", AccountName: "user+test@example.com", Secret: "SECRET"},
			},
			expected: []string{"{1} Test-Service:user+test@example.com:<secret>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTotalOtpsInfo(tt.otpList)

			if len(result) != len(tt.expected) {
				t.Errorf("getTotalOtpsInfo() returned %d items, want %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("getTotalOtpsInfo()[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestRewriteConfFile(t *testing.T) {
	tests := []struct {
		name    string
		otpList []otp
	}{
		{
			name:    "empty list",
			otpList: []otp{},
		},
		{
			name: "single otp",
			otpList: []otp{
				{Issuer: "GitHub", AccountName: "user@example.com", Secret: "TESTSECRET"},
			},
		},
		{
			name: "multiple otps",
			otpList: []otp{
				{Issuer: "GitHub", AccountName: "user1@example.com", Secret: "SECRET1"},
				{Issuer: "Google", AccountName: "user2@example.com", Secret: "SECRET2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "gtptest_*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Call the function
			rewriteConfFile(tmpFile, tt.otpList)

			// Read the file contents
			content, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			// Verify content
			if len(tt.otpList) == 0 {
				// Empty list should write empty content
				if len(content) != 0 {
					t.Errorf("Expected empty content for empty list, got: %s", string(content))
				}
			} else {
				// Non-empty list should be valid JSON
				var readOtpList []otp
				if err := json.Unmarshal(content, &readOtpList); err != nil {
					t.Errorf("Failed to unmarshal written content: %v", err)
					return
				}

				if len(readOtpList) != len(tt.otpList) {
					t.Errorf("Read %d otps, want %d", len(readOtpList), len(tt.otpList))
					return
				}

				for i, v := range readOtpList {
					if v.Issuer != tt.otpList[i].Issuer {
						t.Errorf("Issuer[%d] = %q, want %q", i, v.Issuer, tt.otpList[i].Issuer)
					}
					if v.AccountName != tt.otpList[i].AccountName {
						t.Errorf("AccountName[%d] = %q, want %q", i, v.AccountName, tt.otpList[i].AccountName)
					}
					if v.Secret != tt.otpList[i].Secret {
						t.Errorf("Secret[%d] = %q, want %q", i, v.Secret, tt.otpList[i].Secret)
					}
				}
			}
		})
	}
}

func TestRewriteConfFileOverwrite(t *testing.T) {
	// Create a temporary file with initial content
	tmpFile, err := os.CreateTemp("", "gtptest_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write initial data
	initialData := []otp{
		{Issuer: "Initial", AccountName: "initial@test.com", Secret: "INITIAL"},
	}
	initialJSON, _ := json.Marshal(initialData)
	os.WriteFile(tmpFile.Name(), initialJSON, 0644)

	// Overwrite with new data
	newData := []otp{
		{Issuer: "New", AccountName: "new@test.com", Secret: "NEW"},
	}
	rewriteConfFile(tmpFile, newData)

	// Verify the file contains only new data
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	var readOtpList []otp
	if err := json.Unmarshal(content, &readOtpList); err != nil {
		t.Errorf("Failed to unmarshal content: %v", err)
		return
	}

	if len(readOtpList) != 1 {
		t.Errorf("Expected 1 otp after overwrite, got %d", len(readOtpList))
		return
	}

	if readOtpList[0].Issuer != "New" {
		t.Errorf("Expected Issuer 'New', got %q", readOtpList[0].Issuer)
	}
}

func TestOtpStructJSON(t *testing.T) {
	// Test that otp struct can be properly marshaled and unmarshaled
	original := otp{
		Issuer:      "TestIssuer",
		AccountName: "test@example.com",
		Secret:      "JBSWY3DPEHPK3PXP",
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal otp: %v", err)
	}

	// Unmarshal
	var result otp
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal otp: %v", err)
	}

	// Verify
	if result.Issuer != original.Issuer {
		t.Errorf("Issuer = %q, want %q", result.Issuer, original.Issuer)
	}
	if result.AccountName != original.AccountName {
		t.Errorf("AccountName = %q, want %q", result.AccountName, original.AccountName)
	}
	if result.Secret != original.Secret {
		t.Errorf("Secret = %q, want %q", result.Secret, original.Secret)
	}
}

func TestOtpListJSON(t *testing.T) {
	// Test marshaling and unmarshaling a list of otps
	original := []otp{
		{Issuer: "GitHub", AccountName: "user1@example.com", Secret: "SECRET1"},
		{Issuer: "Google", AccountName: "user2@example.com", Secret: "SECRET2"},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal otp list: %v", err)
	}

	var result []otp
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal otp list: %v", err)
	}

	if len(result) != len(original) {
		t.Errorf("Got %d items, want %d", len(result), len(original))
	}
}
