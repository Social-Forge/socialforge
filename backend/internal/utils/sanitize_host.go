package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/net/publicsuffix"
)

var (
	hostCache      = make(map[string]bool)
	hostCacheLock  sync.RWMutex
	ValidHostRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)
	ipv6Regex      = regexp.MustCompile(`^\[([a-fA-F0-9:]+)\]$`)
)

func GetOriginHost(c fiber.Ctx) string {
	if origin := c.Get("X-Forwarded-Origin"); origin != "" {
		if host, err := parseAndValidateHost(origin); err == nil {
			return host
		}
	}

	headersToCheck := []string{
		"Origin",
		"Referer",
		"X-Forwarded-Host",
		"X-Original-Host",
		"Host",
	}

	for _, header := range headersToCheck {
		if val := c.Get(header); val != "" {
			if host, err := parseAndValidateHost(val); err == nil {
				return host
			}
		}
	}

	return getFallbackHost(c)
}
func parseAndValidateHost(source string) (string, error) {
	source = strings.TrimSpace(source)

	if strings.HasSuffix(source, ":") && len(source) <= 8 {
		return "", fmt.Errorf("invalid URL: just scheme without domain")
	}

	if !strings.Contains(source, "://") && !strings.HasPrefix(source, "//") {
		if strings.Contains(source, ".") || strings.Contains(source, "localhost") {
			source = "http://" + source
		} else {
			return "", fmt.Errorf("does not appear to be a valid domain")
		}
	}

	u, err := url.Parse(source)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("empty hostname")
	}

	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	if !isValidHost(host) {
		return "", fmt.Errorf("invalid host format")
	}

	return strings.ToLower(host), nil
}
func isIPv6(host string) bool {
	// Check for [IPv6] format
	if matches := ipv6Regex.FindStringSubmatch(host); len(matches) > 1 {
		ip := net.ParseIP(matches[1])
		return ip != nil && ip.To4() == nil
	}
	return false
}
func isValidHost(host string) bool {
	if host == "" || len(host) > 253 || strings.ContainsAny(host, " \t\n\r") {
		return false
	}

	if host == "localhost" {
		return false
	}

	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) < 1 || len(label) > 63 {
			return false
		}

		if !isAlphanumeric(rune(label[0])) || !isAlphanumeric(rune(label[len(label)-1])) {
			return false
		}

		for _, r := range label {
			if !isAlphanumeric(r) && r != '-' {
				return false
			}
		}
	}

	return true
}
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
func IsValidPublicSuffix(host string) bool {
	hostCacheLock.RLock()
	valid, exists := hostCache[host]
	hostCacheLock.RUnlock()

	if exists {
		return valid
	}

	suffix, _ := publicsuffix.PublicSuffix(host)
	valid = suffix != ""

	hostCacheLock.Lock()
	hostCache[host] = valid
	hostCacheLock.Unlock()

	return valid
}
func IsValidDomainSuffix(domain string) bool {
	if domain == "" {
		return false
	}
	_, icann := publicsuffix.PublicSuffix(domain)
	if !icann {
		return false
	}
	if strings.Contains(domain, "..") ||
		strings.Contains(domain, ".-") ||
		strings.Contains(domain, "-.") {
		return false
	}

	return true
}
func sanitizeHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))

	if isIPv6(host) {
		return host
	}

	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	host = strings.TrimPrefix(host, "www.")

	return host
}
func getFallbackHost(c fiber.Ctx) string {
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return sanitizeHost(ip)
	}

	return sanitizeHost(c.IP())
}

// Contoh penggunaan:
// extractSubdomain("api.client.mydomain.com", 1) → "client"
// extractSubdomain("api.client.mydomain.com", 2) → "api"
func ExtractSubdomain(host string, level int) string {
	parts := strings.Split(host, ".")
	if len(parts) <= 2 || level <= 0 {
		return ""
	}
	if level >= len(parts)-1 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return parts[len(parts)-2-level]
}

// Mode 2 Satnitize Host
func ExtractCleanDomain(origin string) string {
	if origin == "" {
		return ""
	}
	if strings.Contains(origin, ",") {
		origins := strings.Split(origin, ",")
		for _, o := range origins {
			if domain := extractSingleDomain(o); domain != "" {
				return domain
			}
		}
		return ""
	}

	return extractSingleDomain(origin)
}
func extractSingleDomain(origin string) string {
	origin = strings.ToLower(strings.TrimSpace(origin))
	if !strings.Contains(origin, "://") && !strings.HasPrefix(origin, "//") {
		origin = "http://" + origin
	}
	u, err := url.Parse(origin)
	if err != nil {
		return cleanMalformedOrigin(origin)
	}

	hostname := u.Hostname()
	if hostname == "" {
		return ""
	}

	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}

	hostname = removeCommonPrefixes(hostname)

	if !isValidHostname(hostname) {
		return ""
	}

	return hostname
}

func cleanMalformedOrigin(origin string) string {
	origin = strings.TrimPrefix(strings.TrimPrefix(origin, "http://"), "https://")
	origin = strings.Split(origin, "/")[0]
	origin = strings.Split(origin, "?")[0]
	origin = strings.Split(origin, "#")[0]

	if strings.Contains(origin, ":") {
		origin = strings.Split(origin, ":")[0]
	}

	return removeCommonPrefixes(origin)
}
func removeCommonPrefixes(hostname string) string {
	// prefixes := []string{"www.", "web.", "app.", "api.", "mobile."}
	prefixes := []string{"www."}
	for _, prefix := range prefixes {
		if strings.HasPrefix(hostname, prefix) {
			hostname = hostname[len(prefix):]
			break // Only remove one prefix
		}
	}
	return hostname
}
func isValidHostname(hostname string) bool {
	if hostname == "" || strings.ContainsAny(hostname, " \t\n\r") {
		return false
	}

	if strings.ContainsAny(hostname, "#$%^&*()+=[]{}|;'\"<>?") {
		return false
	}

	if matched, _ := regexp.MatchString(`^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`, hostname); !matched {
		// Allow localhost for development
		if hostname != "localhost" {
			return false
		}
	}
	return true
}

func ExtractCleanDomainRegex(rawURL string) string {
	// Regex yang lebih komprehensif
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?([^\/:\?\#]+)`)
	matches := re.FindStringSubmatch(rawURL)

	if len(matches) < 2 {
		return rawURL // fallback jika tidak match
	}

	domain := matches[1]
	// Hapus port jika ada
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}

	return domain
}
func ValidateSiteUrl(origin string) string {
	if !strings.Contains(origin, "://") && !strings.HasPrefix(origin, "//") {
		origin = "https://" + origin
	}

	return origin
}

// test
func TestDajaIniMah() bool {
	labels := strings.Split("", ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if len(label) < 1 || len(label) > 63 {
			return false
		}

		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}

		if !isAlphanumeric(rune(label[0])) || !isAlphanumeric(rune(label[len(label)-1])) {
			return false
		}

		for _, r := range label {
			if !isAlphanumeric(r) && r != '-' {
				return false
			}
		}
	}
	return true
}
