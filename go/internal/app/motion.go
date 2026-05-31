package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kanije-kalesi/kanije/internal/capture"
	"github.com/kanije-kalesi/kanije/internal/event"
)

// pollWhenMotionIdle is how often the loop re-checks the enable flag while the
// motion detector is off — so /hareket ac takes effect within a few seconds.
const pollWhenMotionIdle = 5 * time.Second

// motionWatch runs the camera motion detector. While enabled it grabs a frame
// every interval and, when the frame differs from the previous one beyond the
// threshold, raises a motion_detected event carrying the photo. The enable flag,
// interval and threshold are re-read every tick so /hareket applies live.
func (a *App) motionWatch(ctx context.Context) error {
	var prev []byte

	for {
		interval, threshold := a.cfg.MotionSettings()
		wait := interval

		if !a.cfg.MotionEnabled() {
			prev = nil // reset baseline so re-enabling never false-triggers
			wait = pollWhenMotionIdle
		} else if a.camera != nil {
			if frame, err := a.camera.Capture(ctx); err != nil {
				a.log.Debug("hareket: kare alınamadı", "err", err)
			} else {
				if prev != nil {
					if diff, derr := capture.FrameDiff(prev, frame); derr == nil && diff >= threshold {
						a.raiseMotion(ctx, frame, diff)
					}
				}
				prev = frame
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(wait):
		}
	}
}

// raiseMotion publishes a motion_detected event with a short photo BURST so a
// single motion captures the sequence of what happened, not just one instant.
func (a *App) raiseMotion(ctx context.Context, frame []byte, diff float64) {
	a.log.Info("hareket algılandı", "fark", diff)
	ev := event.New(event.TypeMotionDetected, "MotionWatch")
	ev.Hostname, _ = os.Hostname()
	ev.Extra = map[string]string{"🎥 Değişim": fmt.Sprintf("%.0f/255", diff)}
	ev.Attachments = a.motionBurst(ctx, frame)
	a.bus.Publish(ev)
}

// motionBurst returns the triggering frame plus up to BurstCount-1 follow-up
// frames captured ~1s apart, so the owner sees the motion unfold (a person
// walking in, picking something up) rather than a single frozen moment.
func (a *App) motionBurst(ctx context.Context, first []byte) []event.Attachment {
	count := a.cfg.MotionBurst()
	atts := make([]event.Attachment, 0, count)
	atts = append(atts, event.Attachment{
		Type:    event.AttachmentPhoto,
		Data:    first,
		Caption: "🎥 Hareket #1",
	})
	for i := 2; i <= count && a.camera != nil; i++ {
		select {
		case <-ctx.Done():
			return atts
		case <-time.After(time.Second):
		}
		data, err := a.camera.Capture(ctx)
		if err != nil {
			a.log.Debug("hareket burst karesi alınamadı", "kare", i, "err", err)
			break
		}
		atts = append(atts, event.Attachment{
			Type:    event.AttachmentPhoto,
			Data:    data,
			Caption: fmt.Sprintf("🎥 Hareket #%d", i),
		})
	}
	return atts
}
