package service

import "fmt"

// FtwEvaluate evaluates fit-to-work status based on sleep minutes.
// Rules match frontend lib/data/ftw.ts:
//
//	>= 330 (5j30) → fit        : can work immediately
//	300-329 (5j-5j29) → spare  : rest 1 hour, then can work
//	240-299 (4j-4j59) → spare  : rest 2 hours, then can work
//	< 240 (<4j) → pulang       : cannot work on this shift
//	no log → belum              : no sleep log submitted
type FtwEval struct {
	Status    string `json:"status"`    // "fit" | "spare" | "pulang" | "belum"
	RestHours int    `json:"restHours"` // 0, 1, or 2
	CanWork   bool   `json:"canWork"`
}

const (
	SleepFitMin   = 330 // 5 hours 30 minutes
	SleepSpare1h  = 300 // 5 hours
	SleepSpare2h  = 240 // 4 hours
	SleepMinGreat = 420 // 7 hours — for bonus points
)

func FtwEvaluate(sleepMin *int) FtwEval {
	if sleepMin == nil || *sleepMin <= 0 {
		return FtwEval{Status: "belum", RestHours: 0, CanWork: false}
	}
	if *sleepMin >= SleepFitMin {
		return FtwEval{Status: "fit", RestHours: 0, CanWork: true}
	}
	if *sleepMin >= SleepSpare1h {
		return FtwEval{Status: "spare", RestHours: 1, CanWork: true}
	}
	if *sleepMin >= SleepSpare2h {
		return FtwEval{Status: "spare", RestHours: 2, CanWork: true}
	}
	return FtwEval{Status: "pulang", RestHours: 0, CanWork: false}
}

func FmtSleepMin(min *int, en bool) string {
	if min == nil || *min <= 0 {
		return "\u2014"
	}
	h := *min / 60
	m := *min % 60
	if en {
		return fmt.Sprintf("%d h %02d m", h, m)
	}
	return fmt.Sprintf("%d j %02d m", h, m)
}

func IsValidEmployeeStatus(status string) bool {
	return status == "aktif" || status == "cuti" || status == "nonaktif"
}
