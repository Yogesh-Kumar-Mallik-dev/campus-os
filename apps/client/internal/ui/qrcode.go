package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"strings"

	"github.com/liyue201/goqr"
	"github.com/makiuchi-d/gozxing"
	zxingqr "github.com/makiuchi-d/gozxing/qrcode"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// ExtractTokenFromQRPayload extracts the claim token from a decoded QR payload.
// It handles raw tokens, full activation URLs (query params or path segments), and key-value formats.
func ExtractTokenFromQRPayload(payload string) string {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return ""
	}

	// Case 0: Official Campus OS sealed format (e.g. CAMPUS_OS:CLAIM:v1:<token> or CAMPUS_OS:<token>)
	if strings.HasPrefix(strings.ToUpper(payload), "CAMPUS_OS:") {
		parts := strings.Split(payload, ":")
		if len(parts) > 0 {
			token := strings.TrimSpace(parts[len(parts)-1])
			if token != "" {
				return token
			}
		}
	}

	// Case 1: Full URL (e.g. https://campus.edu/activate?token=XYZ or campus://activate?claim=XYZ)
	if u, err := url.Parse(payload); err == nil && (u.Scheme != "" || strings.Contains(payload, "?")) {
		q := u.Query()
		if tok := q.Get("token"); tok != "" {
			return tok
		}
		if tok := q.Get("claim"); tok != "" {
			return tok
		}
		if tok := q.Get("claim_token"); tok != "" {
			return tok
		}
		if tok := q.Get("t"); tok != "" {
			return tok
		}
		// Check path segment if no query param
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			if len(last) >= 8 && !strings.Contains(last, ".") {
				return last
			}
		}
	}

	// Case 2: Key-value string like "token=XYZ" or "claim:XYZ"
	for _, sep := range []string{"=", ":"} {
		if strings.Contains(payload, sep) {
			parts := strings.SplitN(payload, sep, 2)
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			if key == "token" || key == "claim" || key == "claim_token" {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	// Case 3: Raw token
	return payload
}

// DecodeQRFromImage reads an image stream and decodes any embedded QR code.
func DecodeQRFromImage(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("failed to decode image format: %w", err)
	}

	return DecodeQRFromLoadedImage(img)
}

// DecodeQRFromLoadedImage decodes any embedded QR code from an in-memory image.Image.
// It prioritizes gozxing (full QR spec compliance, correct numeric byte ordering)
// with a fallback to goqr for edge-case recognition.
func DecodeQRFromLoadedImage(img image.Image) (string, error) {
	if img == nil {
		return "", fmt.Errorf("image is nil")
	}

	// 1. Primary: gozxing (production ZXing port)
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err == nil && bmp != nil {
		qrReader := zxingqr.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)
		if err == nil && result != nil && result.GetText() != "" {
			token := ExtractTokenFromQRPayload(result.GetText())
			if token != "" {
				return token, nil
			}
		}
	}

	// 2. Secondary fallback: goqr
	qrCodes, err := goqr.Recognize(img)
	if err == nil && len(qrCodes) > 0 {
		rawPayload := string(qrCodes[0].Payload)
		token := ExtractTokenFromQRPayload(rawPayload)
		if token != "" {
			return token, nil
		}
	}

	return "", fmt.Errorf("no QR code detected in image")
}

