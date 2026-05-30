// Package shell provides a stateful remote command runner for the /terminal
// bot command. Each chat keeps its own working directory (so `cd` persists
// between commands, like an interactive session). Commands run in the agent's
// privilege context — which is Administrator/SYSTEM when installed normally.
package shell

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

const (
	cmdTimeout = 60 * time.Second
	maxOutput  = 3500 // keep well under Telegram's 4096-char message limit
)

// Runner executes shell commands with a per-chat persistent working directory.
type Runner struct {
	mu   sync.Mutex
	cwds map[int64]string
}

// New creates a Runner.
func New() *Runner {
	return &Runner{cwds: make(map[int64]string)}
}

// Cwd returns the current working directory for a chat.
func (r *Runner) Cwd(chatID int64) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cwdLocked(chatID)
}

func (r *Runner) cwdLocked(chatID int64) string {
	if d, ok := r.cwds[chatID]; ok && d != "" {
		return d
	}
	d, err := os.UserHomeDir()
	if err != nil || d == "" {
		d, _ = os.Getwd()
	}
	r.cwds[chatID] = d
	return d
}

func (r *Runner) setCwd(chatID int64, dir string) {
	r.mu.Lock()
	r.cwds[chatID] = dir
	r.mu.Unlock()
}

// Run executes command for chatID and returns combined stdout+stderr (truncated
// for Telegram). A bare `cd <dir>` updates the session's working directory
// instead of spawning a shell, so subsequent commands run there.
func (r *Runner) Run(ctx context.Context, chatID int64, command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return "Boş komut."
	}
	cwd := r.Cwd(chatID)

	if target, isCd := parseCd(command); isCd {
		return r.changeDir(chatID, cwd, target)
	}

	runCtx, cancel := context.WithTimeout(ctx, cmdTimeout)
	defer cancel()

	cmd := shellCommand(runCtx, command)
	cmd.Dir = cwd
	sysproc.Hide(cmd)

	out, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(out), "\r\n")

	if runCtx.Err() == context.DeadlineExceeded {
		text += "\n⏱️ (60 sn zaman aşımı — komut iptal edildi)"
	} else if text == "" {
		if err != nil {
			text = "Hata: " + err.Error()
		} else {
			text = "(çıktı yok)"
		}
	}
	return truncateOutput(text)
}

// changeDir resolves target relative to cwd and persists it if it's a directory.
func (r *Runner) changeDir(chatID int64, cwd, target string) string {
	var dest string
	switch {
	case target == "" || target == "~":
		dest, _ = os.UserHomeDir()
	case filepath.IsAbs(target):
		dest = target
	default:
		dest = filepath.Join(cwd, target)
	}
	dest = filepath.Clean(dest)

	fi, err := os.Stat(dest)
	if err != nil || !fi.IsDir() {
		return "❌ Dizin bulunamadı: " + dest
	}
	r.setCwd(chatID, dest)
	return "📂 " + dest
}

// parseCd reports whether command is a bare `cd [dir]` and returns the target.
func parseCd(command string) (target string, isCd bool) {
	fields := strings.Fields(command)
	if len(fields) == 0 || !strings.EqualFold(fields[0], "cd") {
		return "", false
	}
	if len(fields) == 1 {
		return "", true
	}
	return strings.Join(fields[1:], " "), true
}

// shellCommand builds the platform shell invocation.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "powershell.exe",
			"-NoProfile", "-NonInteractive", "-Command", command)
	}
	return exec.CommandContext(ctx, "/bin/sh", "-c", command)
}

// truncateOutput trims output to a Telegram-safe length on a valid UTF-8
// boundary, noting how much was dropped.
func truncateOutput(s string) string {
	if len(s) <= maxOutput {
		return s
	}
	cut := strings.ToValidUTF8(s[:maxOutput], "")
	return cut + "\n…\n[çıktı kısaltıldı — " +
		itoa(len(s)-len(cut)) + " bayt daha]"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
