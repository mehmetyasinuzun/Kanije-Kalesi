package telegram

import "strings"

// knownCommands is the union of all recognized command spellings — used to offer
// a "did you mean …?" hint when the user mistypes a command.
func knownCommands() []string {
	cmds := make([]string, 0, len(commandCaps)+12)
	for c := range commandCaps {
		cmds = append(cmds, c)
	}
	// Commands not in commandCaps (open to any authorized user, or aliases).
	cmds = append(cmds,
		"/yardim", "/help", "/start", "/menu", "/komutlar", "/ping", "/iptal",
		"/cihazlar", "/devices",
	)
	return cmds
}

// suggestCommand returns the closest known command to input within an edit
// distance of 2, or "" if nothing is close enough. This catches typos and
// wrong-alphabet look-alikes (e.g. a Cyrillic "о" in "/fоto").
func suggestCommand(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return ""
	}
	best, bestDist := "", 3 // suggest only when distance < 3
	for _, c := range knownCommands() {
		if d := levenshtein(input, c); d < bestDist {
			best, bestDist = c, d
		}
	}
	return best
}

// levenshtein computes the edit distance between two strings (rune-aware so
// multi-byte look-alikes are measured correctly).
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	switch {
	case la == 0:
		return lb
	case lb == 0:
		return la
	}
	prev := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		cur := make([]int, lb+1)
		cur[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[lb]
}
