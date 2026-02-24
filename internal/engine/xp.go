package engine

import "math"

func xpToLevel(xp int) int {
	if xp <= 0 {
		return 0
	}
	return int(math.Sqrt(float64(xp) / 100.0))
}

// XPToNextLevel returns XP needed to reach the next level.
func XPToNextLevel(xp int) int {
	current := xpToLevel(xp)
	next := current + 1
	threshold := next * next * 100
	return threshold - xp
}
