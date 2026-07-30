package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var validTLDs = map[string]bool{
	"com": true,
	"net": true,
	"org": true,
	// tambahkan lainnya
}

func init() {
	validate.RegisterValidation("domain", validateDomain)
	validate.RegisterValidation("timestamp", validateTimestamp)
	validate.RegisterValidation("file", FileImagesValidation)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	validate.RegisterValidation("datetime", func(fl validator.FieldLevel) bool {
		_, err := time.Parse(time.RFC3339, fl.Field().String())
		return err == nil
	})
	validate.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(time.Time).After(time.Now())
	})
}
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = msgForTag(e.Tag(), e)
	}

	return errors

}
func CustomValidateBody(s interface{}) string {
	errors := ValidateStruct(s)
	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for field, message := range errors {
			return fmt.Sprintf("- %s: %s\n", field, message)
		}
	}
	return ""
}
func msgForTag(tag string, e validator.FieldError) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "eqfield":
		return fmt.Sprintf("Field %s must match %s", e.Field(), e.Param())
	case "domain":
		return "Invalid domain format"
	case "min":
		return fmt.Sprintf("Value is too short. Minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("Value is too long. Maximum length is %s", e.Param())
	case "number":
		return "Count must be a valid number"
	case "url":
		return "Invalid URL format"
	case "timestamp":
		return "Must be a future timestamp in format: YYYY-MM-DD HH:MM:SS[.SSSSSS] (min 1 hour from now)"
	case "datetime":
		return "Must be a valid datetime format"
	case "future":
		return "Must be a future timestamp in format: YYYY-MM-DD HH:MM:SS[.SSSSSS] (min 1 hour from now)"
	default:
		return "Invalid value"
	}
}
func validateDomain(fl validator.FieldLevel) bool {
	domain := fl.Field().String()

	// Regex untuk validasi domain dasar
	// Contoh: example.com atau sub.example.com
	domainRegex := `^([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]))+)$`

	if matched, _ := regexp.MatchString(domainRegex, domain); !matched {
		return false
	}

	// 2. Validasi panjang maksimal (253 karakter)
	if len(domain) > 253 {
		return false
	}

	// 3. Validasi bagian-bagian domain
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		// Setiap bagian antara titik maksimal 63 karakter
		if len(part) > 63 {
			return false
		}
	}

	// 4. Opsional: Cek DNS
	if !isDomainActive(domain) {
		return false
	}

	return true
}
func isDomainActive(domain string) bool {
	_, err := net.LookupHost(domain)
	return err == nil
}
func IsValidTLD(domain string) bool {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}
	tld := parts[len(parts)-1]
	return validTLDs[tld]
}
func FileImagesValidation(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return false
	}

	if file.Size > 5<<20 { // 5MB limit
		return false
	}

	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	buffer := make([]byte, 512)
	if _, err = f.Read(buffer); err != nil && err != io.EOF {
		return false
	}

	// Reset reader for actual upload
	if _, err = f.Seek(0, 0); err != nil {
		return false
	}

	allowedTypes := map[string]bool{
		"image/jpeg":               true,
		"image/jpg":                true,
		"image/png":                true,
		"image/webp":               true,
		"image/svg+xml":            true,  // Correct MIME for SVG
		"application/octet-stream": false, // Explicitly block
	}

	contentType := http.DetectContentType(buffer)

	// Check both detected type and extension
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	validExtension := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".svg":  true,
	}

	return allowedTypes[contentType] && validExtension[fileExt]
}
func IsValidEmail(email string) bool {
	// Trim whitespace
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Regular expression for basic email validation
	// More comprehensive than simple @ check but not overly strict
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if !match {
		return false
	}

	// Additional checks
	if len(email) > 254 {
		return false
	}

	// Check for common disposable email domains
	disposableDomains := []string{
		"tempmail.com", "mailinator.com", "10minutemail.com",
		"guerrillamail.com", "throwawaymail.com",
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])
	for _, d := range disposableDomains {
		if strings.Contains(domain, d) {
			return false
		}
	}

	return true
}
func ValidateImageFile(file *multipart.FileHeader) error {
	if file == nil {
		return errors.New("nil file provided")
	}
	// Check file size (max 5MB)
	const maxUploadSize = 5 << 20 // 5MB
	if file.Size == 0 {
		return errors.New("empty file")
	}
	if file.Size > maxUploadSize {
		return fmt.Errorf("file too large (%.2fMB, max 5MB allowed)", float64(file.Size)/float64(1<<20))
	}
	// Check filename security
	if strings.Contains(file.Filename, "..") || strings.Contains(file.Filename, "/") {
		return errors.New("invalid filename")
	}
	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return errors.New("missing file extension")
	}
	allowedExtensions := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
	expectedMime, validExt := allowedExtensions[ext]
	if !validExt {
		return fmt.Errorf("invalid file type %q, only JPG, JPEG, PNG, or WEBP allowed", ext)
	}
	// Open file to verify content type
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}
	defer f.Close()
	// Read first 512 bytes for content detection
	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("could not read file content: %w", err)
	}
	if n == 0 {
		return errors.New("empty file content")
	}
	// Reset read pointer for actual upload
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("could not reset file pointer: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buffer[:n])

	// Verify MIME type matches expected
	if !strings.HasPrefix(contentType, expectedMime) {
		return fmt.Errorf("invalid file content: expected %s but got %s", expectedMime, contentType)
	}

	// Additional magic number verification
	switch ext {
	case ".jpg", ".jpeg":
		if !bytes.HasPrefix(buffer, []byte{0xFF, 0xD8, 0xFF}) {
			return errors.New("invalid JPEG file signature")
		}
	case ".png":
		if !bytes.HasPrefix(buffer, []byte{0x89, 0x50, 0x4E, 0x47}) {
			return errors.New("invalid PNG file signature")
		}
	case ".webp":
		if !bytes.HasPrefix(buffer, []byte{0x52, 0x49, 0x46, 0x46}) ||
			!bytes.HasPrefix(buffer[8:], []byte{0x57, 0x45, 0x42, 0x50}) {
			return errors.New("invalid WebP file signature")
		}
	}

	return nil
}
func ValidateArrayImageFiles(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return errors.New("no files uploaded")
	}

	const maxTotalSize = 20 << 20 // 20MB total for all files
	totalSize := int64(0)

	for i, file := range files {
		// Check individual file
		if err := ValidateImageFile(file); err != nil {
			return fmt.Errorf("file %d: %w", i+1, err)
		}

		// Check cumulative size
		totalSize += file.Size
		if totalSize > maxTotalSize {
			return fmt.Errorf("total upload size exceeds 20MB limit")
		}
	}

	return nil
}
func validateTimestamp(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Match Zod's regex pattern
	matched, _ := regexp.MatchString(
		`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])(\.\d{1,6})?$`,
		value,
	)
	if !matched {
		return false
	}
	// Parse the time (similar to Zod's refine)
	t, err := time.Parse("2006-01-02 15:04:05.999999", value)
	if err != nil {
		return false
	}
	// Validate minimum time (1 hour from now)
	minTime := time.Now().Add(1 * time.Hour)
	return t.After(minTime)
}
