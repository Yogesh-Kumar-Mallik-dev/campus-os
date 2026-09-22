package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"github.com/skip2/go-qrcode"
)

func TestExtractTokenFromQRPayload(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"claim_genesis_test_demo", "claim_genesis_test_demo"},
		{"  tok_trimmed_123  ", "tok_trimmed_123"},
		{"https://campus.bbdit.edu/activate?token=ABC123XYZ456", "ABC123XYZ456"},
		{"https://campus.bbdit.edu/onboard?claim=CLAIM_987654", "CLAIM_987654"},
		{"campus://activate?claim_token=APP_TOKEN_4321", "APP_TOKEN_4321"},
		{"https://bbdit.ac.in/verify?t=SHORT_TOK_99", "SHORT_TOK_99"},
		{"token=KEY_VAL_TOKEN_111", "KEY_VAL_TOKEN_111"},
		{"claim:COLON_CLAIM_222", "COLON_CLAIM_222"},
		{"https://campus.edu/activate/voucher_polytechnic_2024", "voucher_polytechnic_2024"},
		{"", ""},
		{"   ", ""},
	}

	for _, tc := range tests {
		got := ExtractTokenFromQRPayload(tc.input)
		if got != tc.expected {
			t.Errorf("ExtractTokenFromQRPayload(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestDecodeQRFromImage_InvalidAndBlank(t *testing.T) {
	// Test 1: Invalid image format
	invalidReader := strings.NewReader("not an image at all")
	_, err := DecodeQRFromImage(invalidReader)
	if err == nil {
		t.Errorf("expected error for invalid image stream, got nil")
	}

	// Test 2: Blank image with no QR code
	blankImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			blankImg.Set(x, y, color.White)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, blankImg); err != nil {
		t.Fatalf("failed to encode blank image: %v", err)
	}

	_, err = DecodeQRFromImage(&buf)
	if err == nil {
		t.Errorf("expected error for blank image without QR code, got nil")
	}
	if !strings.Contains(err.Error(), "no QR code detected") {
		t.Errorf("expected 'no QR code detected' in error, got: %v", err)
	}
}

func TestDecodeQRFromImage_RealQRCode(t *testing.T) {
	// Test with alphanumeric token including consecutive digits
	pngData, err := qrcode.Encode("CAMPUS_OS:CLAIM:v1:claim_genesis_valid_12345", qrcode.Medium, 256)
	if err != nil {
		t.Fatalf("failed to encode test QR: %v", err)
	}

	token, err := DecodeQRFromImage(bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("failed to decode real QR: %v", err)
	}

	t.Logf("Decoded token: %q", token)
	if token != "claim_genesis_valid_12345" {
		t.Errorf("expected 'claim_genesis_valid_12345', got '%s'", token)
	}
}


