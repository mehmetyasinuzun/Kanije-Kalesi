package capture

import "strings"

// parseDshowDevices extracts DirectShow device friendly-names of the given kind
// ("video" or "audio") from `ffmpeg -list_devices` output. It handles BOTH:
//
//	modern ffmpeg (5.x+):  [dshow @ ..] "Integrated Camera" (video)
//	legacy ffmpeg (4.x):   [dshow @ ..] DirectShow video devices
//	                       [dshow @ ..]  "Integrated Camera"
//
// "Alternative name "@device_..."" lines are skipped (not selectable). This is the
// reason camera/audio "couldn't find device": the auto-downloaded build uses the
// modern format, which the old header-only parser missed.
func parseDshowDevices(output, kind string) []string {
	var devices []string
	section := "" // legacy: which section we're currently in ("video"/"audio")

	for _, line := range strings.Split(output, "\n") {
		low := strings.ToLower(line)

		if strings.Contains(low, "alternative name") {
			continue
		}

		// Modern format: the line itself is tagged "(video)" / "(audio)".
		switch {
		case strings.Contains(low, "(video)"):
			if kind == "video" {
				if n := firstQuoted(line); n != "" {
					devices = append(devices, n)
				}
			}
			continue
		case strings.Contains(low, "(audio)"):
			if kind == "audio" {
				if n := firstQuoted(line); n != "" {
					devices = append(devices, n)
				}
			}
			continue
		}

		// Legacy format: section headers switch context. Check "video" first —
		// the video header text also contains the word "audio".
		if strings.Contains(low, "video devices") {
			section = "video"
			continue
		}
		if strings.Contains(low, "audio devices") {
			section = "audio"
			continue
		}

		if section == kind {
			if n := firstQuoted(line); n != "" {
				devices = append(devices, n)
			}
		}
	}
	return devices
}

// firstQuoted returns the first double-quoted substring of s, or "".
func firstQuoted(s string) string {
	start := strings.IndexByte(s, '"')
	if start < 0 {
		return ""
	}
	rest := s[start+1:]
	end := strings.IndexByte(rest, '"')
	if end <= 0 {
		return ""
	}
	return rest[:end]
}
