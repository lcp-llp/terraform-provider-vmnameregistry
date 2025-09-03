package provider

import "strings"

// Helper function to extract location from VM name
func extractLocationFromVmName(vmName string) string {
	// VM names follow pattern: lcp{env}{location}-{number}
	// e.g., lcpdevuks-0001, lcpprduks-0002, etc.
	parts := strings.Split(vmName, "-")
	if len(parts) < 2 {
		return ""
	}

	// Remove the "lcp" prefix and environment code to get location
	prefix := parts[0]    // e.g., "lcpdevuks"
	if len(prefix) <= 3 { // Should be at least "lcp" + something
		return ""
	}

	// Remove "lcp" prefix
	withoutLcp := prefix[3:] // e.g., "devuks"

	// Common environment codes to remove
	envCodes := []string{"dev", "prd", "ppd", "dvt"}
	for _, envCode := range envCodes {
		if strings.HasPrefix(withoutLcp, envCode) {
			return withoutLcp[len(envCode):] // Return the location part
		}
	}

	return ""
}
