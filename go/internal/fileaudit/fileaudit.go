// Package fileaudit answers "who opened/copied my files" using Windows native
// auditing: it turns on Object Access (File System) auditing, attaches a SACL
// audit rule to a folder, and reads the resulting Security-log 4663 events back
// as who/when/which-program/what-action records. Windows-only (admin/SYSTEM
// required to change audit policy); other platforms report it unsupported.
package fileaudit

// AccessEvent is one file-access record reconstructed from a Security 4663 event.
type AccessEvent struct {
	Time    string // when the access happened
	User    string // SubjectUserName — the account that accessed it
	Process string // the program that did it (full path)
	Object  string // the file/folder touched
	Access  string // human action: "okuma/kopya", "yazma", "silme"…
}
