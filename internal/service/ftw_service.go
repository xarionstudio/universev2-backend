package service

import (
	"encoding/json"
	"fmt"

	"universev/internal/repository"
)

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

// FtwRules holds dynamic FTW thresholds from business_rules table.
type FtwRules struct {
	SleepFitMin  int `json:"sleep_fit_min"`
	SleepSpare1h int `json:"sleep_spare_1h_min"`
	SleepSpare2h int `json:"sleep_spare_2h_min"`
}

// DefaultFtwRules returns fallback values if business rules not found.
func DefaultFtwRules() FtwRules {
	return FtwRules{
		SleepFitMin:  SleepFitMin,
		SleepSpare1h: SleepSpare1h,
		SleepSpare2h: SleepSpare2h,
	}
}

// GetFtwRules fetches FTW thresholds from business_rules table.
func GetFtwRules(settingsRepo *repository.SettingsRepo) FtwRules {
	if settingsRepo == nil {
		return DefaultFtwRules()
	}

	rule, err := settingsRepo.GetBusinessRuleByCategory("ftw")
	if err != nil {
		return DefaultFtwRules()
	}

	var rules FtwRules
	if err := json.Unmarshal([]byte(rule.Rules), &rules); err != nil {
		return DefaultFtwRules()
	}

	// Ensure all fields have values (use defaults if missing/zero)
	defaults := DefaultFtwRules()
	if rules.SleepFitMin <= 0 {
		rules.SleepFitMin = defaults.SleepFitMin
	}
	if rules.SleepSpare1h <= 0 {
		rules.SleepSpare1h = defaults.SleepSpare1h
	}
	if rules.SleepSpare2h <= 0 {
		rules.SleepSpare2h = defaults.SleepSpare2h
	}

	return rules
}

func FtwEvaluate(sleepMin *int) FtwEval {
	return FtwEvaluateWithRules(sleepMin, DefaultFtwRules())
}

// FtwEvaluateWithRules evaluates fit-to-work status using dynamic thresholds.
func FtwEvaluateWithRules(sleepMin *int, rules FtwRules) FtwEval {
	if sleepMin == nil || *sleepMin <= 0 {
		return FtwEval{Status: "belum", RestHours: 0, CanWork: false}
	}
	if *sleepMin >= rules.SleepFitMin {
		return FtwEval{Status: "fit", RestHours: 0, CanWork: true}
	}
	if *sleepMin >= rules.SleepSpare1h {
		return FtwEval{Status: "spare", RestHours: 1, CanWork: true}
	}
	if *sleepMin >= rules.SleepSpare2h {
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
