// Package schedule provides a tiny persistent task scheduler: named commands that
// run on a fixed interval. Tasks survive restarts (JSON file) and missed slots
// while the agent was down are skipped rather than fired in a burst.
package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// Task is one recurring scheduled command.
type Task struct {
	ID          int       `json:"id"`
	Command     string    `json:"command"`      // e.g. "/foto"
	IntervalSec int       `json:"interval_sec"` // run period in seconds
	NextRun     time.Time `json:"next_run"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store persists tasks to a JSON file. All methods are safe for concurrent use.
type Store struct {
	mu    sync.Mutex
	path  string
	tasks []Task
	seq   int
}

// NewStore opens (or creates) the store at path. A missing file is not an error.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.tasks); err != nil {
		return fmt.Errorf("zamanlama dosyası okunamadı: %w", err)
	}
	for _, t := range s.tasks {
		if t.ID > s.seq {
			s.seq = t.ID
		}
	}
	return nil
}

// save atomically writes the task list (caller holds the lock).
func (s *Store) save() error {
	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Add registers command to run every intervalSec, first firing one interval from
// now. Returns the created task.
func (s *Store) Add(command string, intervalSec int, now time.Time) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if intervalSec < 1 {
		return Task{}, fmt.Errorf("geçersiz aralık")
	}
	s.seq++
	t := Task{
		ID:          s.seq,
		Command:     command,
		IntervalSec: intervalSec,
		NextRun:     now.Add(time.Duration(intervalSec) * time.Second),
		CreatedAt:   now,
	}
	s.tasks = append(s.tasks, t)
	return t, s.save()
}

// Remove deletes the task with the given ID. Returns true if one was removed.
func (s *Store) Remove(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true, s.save()
		}
	}
	return false, nil
}

// List returns a copy of all tasks, sorted by ID.
func (s *Store) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Due returns the tasks whose NextRun has passed, advancing each one to its next
// future slot (skipping any missed while the app was down, so they never fire in
// a burst). The change is persisted.
func (s *Store) Due(now time.Time) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var due []Task
	changed := false
	for i := range s.tasks {
		if s.tasks[i].NextRun.After(now) {
			continue
		}
		due = append(due, s.tasks[i])
		step := time.Duration(s.tasks[i].IntervalSec) * time.Second
		for !s.tasks[i].NextRun.After(now) {
			s.tasks[i].NextRun = s.tasks[i].NextRun.Add(step)
		}
		changed = true
	}
	if changed {
		_ = s.save()
	}
	return due
}
