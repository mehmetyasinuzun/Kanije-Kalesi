package capture

import (
	"log/slog"
	"strings"
	"testing"
)

func TestParseAudioDeviceList(t *testing.T) {
	// Mirrors real `ffmpeg -list_devices true -f dshow -i dummy` stderr: a video
	// section first, then audio, each device followed by an "Alternative name".
	out := strings.Join([]string{
		`[dshow @ 0x1] DirectShow video devices (some may be both video and audio devices)`,
		`[dshow @ 0x1]  "Integrated Camera"`,
		`[dshow @ 0x1]     Alternative name "@device_pnp_\\?\usb#vid_0001"`,
		`[dshow @ 0x1] DirectShow audio devices`,
		`[dshow @ 0x1]  "Microphone (Realtek(R) Audio)"`,
		`[dshow @ 0x1]     Alternative name "@device_cm_{ABC}"`,
		`[dshow @ 0x1]  "Stereo Mix (Realtek(R) Audio)"`,
		`[dshow @ 0x1]     Alternative name "@device_cm_{DEF}"`,
	}, "\n")

	got := parseDshowDevices(out, "audio")
	want := []string{"Microphone (Realtek(R) Audio)", "Stereo Mix (Realtek(R) Audio)"}

	if len(got) != len(want) {
		t.Fatalf("cihaz sayısı: got %d %v, want %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("cihaz[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseAudioDeviceListEmpty(t *testing.T) {
	if got := parseDshowDevices("", "audio"); got != nil {
		t.Errorf("boş çıktı için nil beklenir, got %v", got)
	}
	// Only a video device present → no audio devices.
	videoOnly := "[dshow] DirectShow video devices\n[dshow]  \"Cam\"\n"
	if got := parseDshowDevices(videoOnly, "audio"); len(got) != 0 {
		t.Errorf("yalnız video varken boş beklenir, got %v", got)
	}
}

func TestBuildArgsClampAndDuration(t *testing.T) {
	r := NewAudioRecorder(AudioConfig{Bitrate: "128k"}, slog.Default())
	args := r.buildArgs("default", 45)

	joined := strings.Join(args, " ")
	for _, want := range []string{"-t 45", "-b:a 128k", "-f mp3", "pipe:1", "-ac 1"} {
		if !strings.Contains(joined, want) {
			t.Errorf("args %q içinde %q bekleniyordu", joined, want)
		}
	}
}

func TestBuildArgsDefaultBitrate(t *testing.T) {
	r := NewAudioRecorder(AudioConfig{}, slog.Default())
	args := r.buildArgs("default", 30)
	if !strings.Contains(strings.Join(args, " "), "-b:a 96k") {
		t.Errorf("varsayılan bitrate 96k bekleniyordu: %v", args)
	}
}

func TestParseMaxVolume(t *testing.T) {
	out := strings.Join([]string{
		`[Parsed_volumedetect_0 @ 0x1] n_samples: 132300`,
		`[Parsed_volumedetect_0 @ 0x1] mean_volume: -41.2 dB`,
		`[Parsed_volumedetect_0 @ 0x1] max_volume: -23.4 dB`,
	}, "\n")
	v, err := parseMaxVolume(out)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if v != -23.4 {
		t.Errorf("max_volume = %v, beklenen -23.4", v)
	}
}

func TestParseMaxVolumeMissing(t *testing.T) {
	if _, err := parseMaxVolume("hiç ses bilgisi yok"); err == nil {
		t.Error("max_volume yokken hata beklenir")
	}
}
