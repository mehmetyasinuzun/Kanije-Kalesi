package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// pollWhenVOXIdle is how often the loop re-checks the enable flag while voice-
// activated recording is off, so /tetikses ac takes effect within a few seconds.
const pollWhenVOXIdle = 5 * time.Second

// voxWatch runs voice-activated recording (/tetikses). While enabled it samples
// the mic in short probes; when the level crosses the threshold it records — and
// KEEPS recording in parts while sound continues — then returns to silent
// listening. Nothing is recorded or sent while quiet. Settings are re-read each
// loop so /tetikses applies live.
func (a *App) voxWatch(ctx context.Context) error {
	for {
		if !a.cfg.VOXEnabled() || a.audioRec == nil {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(pollWhenVOXIdle):
			}
			continue
		}

		threshold, partSec, sampleSec, maxParts := a.cfg.VOXSettings()
		level, err := a.audioRec.MeasureLevel(ctx, sampleSec)
		if err != nil {
			a.log.Debug("VOX: seviye ölçülemedi", "err", err)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(10 * time.Second): // mic busy/unavailable — back off
			}
			continue
		}

		if level >= threshold {
			a.recordVOXSession(ctx, level, threshold, partSec, sampleSec, maxParts)
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

// recordVOXSession records the triggering sound in PartSec-long parts, sending
// each part separately, and keeps going as long as sound persists (re-probing
// the level after every part) up to maxParts. This splits a long disturbance
// into 1-2 min chunks instead of one huge file and stops when silence returns.
func (a *App) recordVOXSession(ctx context.Context, firstLevel, threshold float64, partSec, sampleSec, maxParts int) {
	a.log.Info("VOX tetiklendi — kayda başlanıyor", "seviye_dB", firstLevel, "eşik_dB", threshold)
	level := firstLevel

	for part := 1; part <= maxParts; part++ {
		if ctx.Err() != nil {
			return
		}
		data, err := a.audioRec.Record(ctx, partSec)
		if err != nil {
			a.log.Warn("VOX kaydı alınamadı", "parça", part, "err", err)
			return
		}
		a.publishVOX(data, level, part)

		// Still loud? Re-probe; if quiet, the session ends here.
		level, err = a.audioRec.MeasureLevel(ctx, sampleSec)
		if err != nil || level < threshold {
			a.log.Info("VOX oturumu bitti (sessizlik)", "parça_sayısı", part)
			return
		}
	}
	a.log.Info("VOX oturumu üst sınıra ulaştı", "max_parça", maxParts)
}

// publishVOX publishes one VOX audio part as a vox_triggered event.
func (a *App) publishVOX(data []byte, level float64, part int) {
	ev := event.New(event.TypeVOXTriggered, "VOXWatch")
	ev.Hostname, _ = os.Hostname()
	ev.Extra = map[string]string{
		"🔊 Seviye": fmt.Sprintf("%.0f dB", level),
		"🧩 Parça":  fmt.Sprintf("#%d", part),
	}
	ev.Attachments = []event.Attachment{{
		Type:     event.AttachmentAudio,
		Data:     data,
		Filename: fmt.Sprintf("vox_%d.mp3", part),
		Caption:  fmt.Sprintf("🔊 Ses algılandı — parça #%d (%.0f dB)", part, level),
	}}
	a.bus.Publish(ev)
}
