package capture

import (
	"reflect"
	"testing"
)

// Modern ffmpeg (5.x/6.x/7.x — BtbN builds) tags each device line with its kind.
const modernList = `[dshow @ 000] "Integrated Camera" (video)
[dshow @ 000]   Alternative name "@device_pnp_\\?\usb#vid_0001"
[dshow @ 000] "HD WebCam" (video)
[dshow @ 000] "Microphone Array (Realtek)" (audio)
[dshow @ 000]   Alternative name "@device_cm_{guid}"
[dshow @ 000] "Stereo Mix" (audio)`

// Legacy ffmpeg (4.x) uses section headers.
const legacyList = `[dshow @ 000] DirectShow video devices (some may be both video and audio devices)
[dshow @ 000]  "Integrated Camera"
[dshow @ 000]     Alternative name "@device_pnp_\\?\usb#vid_0001"
[dshow @ 000] DirectShow audio devices
[dshow @ 000]  "Microphone Array (Realtek)"
[dshow @ 000]     Alternative name "@device_cm_{guid}"`

func TestParseDshowModernVideo(t *testing.T) {
	got := parseDshowDevices(modernList, "video")
	want := []string{"Integrated Camera", "HD WebCam"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("modern video = %v, beklenen %v", got, want)
	}
}

func TestParseDshowModernAudio(t *testing.T) {
	got := parseDshowDevices(modernList, "audio")
	want := []string{"Microphone Array (Realtek)", "Stereo Mix"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("modern audio = %v, beklenen %v", got, want)
	}
}

func TestParseDshowLegacyVideo(t *testing.T) {
	got := parseDshowDevices(legacyList, "video")
	want := []string{"Integrated Camera"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("legacy video = %v, beklenen %v", got, want)
	}
}

func TestParseDshowLegacyAudio(t *testing.T) {
	got := parseDshowDevices(legacyList, "audio")
	want := []string{"Microphone Array (Realtek)"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("legacy audio = %v, beklenen %v", got, want)
	}
}

func TestParseDshowEmpty(t *testing.T) {
	if got := parseDshowDevices("", "video"); len(got) != 0 {
		t.Errorf("boş girdi = %v, beklenen boş", got)
	}
}
