package utils

import (
	"net/http"
	"strings"
)

type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
	DeviceUnknown DeviceType = "unknown"
)

func DetectDeviceType(r *http.Request) DeviceType {
	ua := strings.ToLower(r.UserAgent())

	// ===== Mobile =====
	if strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "ipod") {
		return DeviceMobile
	}

	// ===== Tablet =====
	if strings.Contains(ua, "ipad") ||
		strings.Contains(ua, "tablet") {
		return DeviceTablet
	}

	// ===== Desktop / Laptop =====
	if strings.Contains(ua, "windows") ||
		strings.Contains(ua, "macintosh") ||
		strings.Contains(ua, "linux") ||
		strings.Contains(ua, "x11") {
		return DeviceDesktop
	}

	return DeviceUnknown
}

func DetectDeviceTypeAdvanced(r *http.Request) DeviceType {
	// Client Hints (Chrome, Edge)
	if mobile := r.Header.Get("Sec-CH-UA-Mobile"); mobile != "" {
		if mobile == "?1" {
			return DeviceMobile
		}
		return DeviceDesktop
	}

	// Fallback ke User-Agent
	return DetectDeviceType(r)
}
