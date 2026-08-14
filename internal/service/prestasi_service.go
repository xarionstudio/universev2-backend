package service

import (
	"encoding/json"
	"fmt"
	"time"

	"universev/internal/model"
	"universev/internal/repository"
)

type PrestasiService struct {
	repo         *repository.PrestasiRepo
	fleetRepo    *repository.FleetRepo
	settingsRepo *repository.SettingsRepo
}

func NewPrestasiService(repo *repository.PrestasiRepo, fleetRepo *repository.FleetRepo, settingsRepo *repository.SettingsRepo) *PrestasiService {
	return &PrestasiService{repo: repo, fleetRepo: fleetRepo, settingsRepo: settingsRepo}
}

// PrestasiRules holds dynamic prestasi point rules from business_rules table
type PrestasiRules struct {
	PtsBase                int
	PtsOntime              int
	PtsSleep               int
	PtsStreakStep          int
	PtsStreakCap           int
	PtsCover               int
	PtsPenalty             int
	SleepMinGreat          int
	BadgeStreak7Threshold  int
	BadgeStreak14Threshold int
}

// DefaultPrestasiRules returns fallback values if business rules not found
func DefaultPrestasiRules() PrestasiRules {
	return PrestasiRules{
		PtsBase:                10,
		PtsOntime:              2,
		PtsSleep:               3,
		PtsStreakStep:          2,
		PtsStreakCap:           10,
		PtsCover:               5,
		PtsPenalty:             -15,
		SleepMinGreat:          420,
		BadgeStreak7Threshold:  7,
		BadgeStreak14Threshold: 14,
	}
}

// GetPrestasiRules fetches prestasi rules from business_rules table
func (s *PrestasiService) GetPrestasiRules() (PrestasiRules, error) {
	if s.settingsRepo == nil {
		return DefaultPrestasiRules(), nil
	}

	rule, err := s.settingsRepo.GetBusinessRuleByCategory("prestasi")
	if err != nil {
		return DefaultPrestasiRules(), nil // Fallback to defaults
	}

	var rules PrestasiRules
	if err := json.Unmarshal([]byte(rule.Rules), &rules); err != nil {
		return DefaultPrestasiRules(), nil // Fallback to defaults
	}

	// Ensure all fields have values (use defaults if missing)
	defaults := DefaultPrestasiRules()
	if rules.PtsBase == 0 {
		rules.PtsBase = defaults.PtsBase
	}
	if rules.PtsOntime == 0 {
		rules.PtsOntime = defaults.PtsOntime
	}
	if rules.PtsSleep == 0 {
		rules.PtsSleep = defaults.PtsSleep
	}
	if rules.PtsStreakStep == 0 {
		rules.PtsStreakStep = defaults.PtsStreakStep
	}
	if rules.PtsStreakCap == 0 {
		rules.PtsStreakCap = defaults.PtsStreakCap
	}
	if rules.PtsCover == 0 {
		rules.PtsCover = defaults.PtsCover
	}
	if rules.PtsPenalty == 0 {
		rules.PtsPenalty = defaults.PtsPenalty
	}
	if rules.SleepMinGreat == 0 {
		rules.SleepMinGreat = defaults.SleepMinGreat
	}
	if rules.BadgeStreak7Threshold == 0 {
		rules.BadgeStreak7Threshold = defaults.BadgeStreak7Threshold
	}
	if rules.BadgeStreak14Threshold == 0 {
		rules.BadgeStreak14Threshold = defaults.BadgeStreak14Threshold
	}

	return rules, nil
}

// Outcome enumerations matching FE PrestasiDay.outcome
const (
	OutcomeNotScheduled   = "notScheduled"
	OutcomeQualified      = "qualified"
	OutcomeReplacedAbsent = "replacedAbsent"
	OutcomeReplacedSleep  = "replacedSleep"
	OutcomeReplacement    = "replacement"
)

// Badge keys matching FE PrestasiBadgeKey
const (
	BadgeStreak7      = "streak7"
	BadgeStreak14     = "streak14"
	BadgePerfectSleep = "perfectSleep"
	BadgeNeverLate    = "neverLate"
	BadgeNoPenalty    = "noPenalty"
)

func (s *PrestasiService) GetLeaderboard(periodDays int) ([]model.PrestasiRecord, error) {
	records, err := s.repo.GetLeaderboard(periodDays)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *PrestasiService) GetHistory(nik string, days int) ([]model.PrestasiHistoryEntry, error) {
	return s.repo.GetOperatorHistory(nik, days)
}

// Recalculate simulates ALL operators per date (like FE) because
// replacement events pair two operators on the same day.
func (s *PrestasiService) Recalculate(periodDays int) error {
	if periodDays <= 0 {
		periodDays = 30
	}

	// Load dynamic business rules
	rules, err := s.GetPrestasiRules()
	if err != nil {
		return err
	}
	ftwRules := GetFtwRules(s.settingsRepo)

	emps, err := s.repo.GetAllEmployeesWithCompetencies()
	if err != nil {
		return err
	}

	// Load real unit DB for fallback unit codes
	unitDB, _ := s.fleetRepo.GetUnitDB()
	unitCodes := make([]string, 0, len(unitDB))
	for _, u := range unitDB {
		if u.Active && !u.Breakdown {
			unitCodes = append(unitCodes, u.Code)
		}
	}

	allocByDate := s.loadAllocationsByDate(periodDays)

	// Simulate each date for ALL operators simultaneously
	dayResults := make(map[string]map[string]model.PrestasiHistoryEntry)
	for d := periodDays; d >= 1; d-- {
		dateStr := time.Now().AddDate(0, 0, -d).Format("2006-01-02")
		dayResults[dateStr] = s.simulateDay(dateStr, emps, allocByDate[dateStr], unitCodes, rules, ftwRules)
	}

	// Aggregate per operator
	var scores []model.PrestasiScore
	var history []model.PrestasiHistoryEntry
	var badges []model.PrestasiBadge

	for _, emp := range emps {
		totalPoints := 0
		currentStreak := 0
		bestStreak := 0
		attCount := 0
		sleepOkCount := 0
		totalScheduled := 0
		qualifiedDays := 0
		coverDays := 0
		penaltyCount := 0
		totalSleepMin := 0
		sleepEntriesCount := 0
		lateCount := 0

		var empDays []model.PrestasiHistoryEntry

		for d := periodDays; d >= 1; d-- {
			dateStr := time.Now().AddDate(0, 0, -d).Format("2006-01-02")
			entry, ok := dayResults[dateStr][emp.NIK]
			if !ok {
				continue
			}

			empDays = append(empDays, entry)
			totalPoints += entry.Points
			totalSleepMin += entry.SleepMin
			sleepEntriesCount++

			if entry.Outcome == OutcomeQualified || entry.Outcome == OutcomeReplacement {
				qualifiedDays++
				if entry.Outcome == OutcomeReplacement {
					coverDays++
				}
				currentStreak++
				if currentStreak > bestStreak {
					bestStreak = currentStreak
				}
			} else if entry.Outcome == OutcomeReplacedAbsent || entry.Outcome == OutcomeReplacedSleep {
				penaltyCount++
				currentStreak = 0
			}

			if entry.AttOk {
				attCount++
			}
			if entry.SleepOk {
				sleepOkCount++
			}
			if entry.Late {
				lateCount++
			}
			if entry.ShiftCode == "D" || entry.ShiftCode == "N" {
				totalScheduled++
			}
		}

		// Add streak bonus (capped)
		if bestStreak > 1 {
			bonus := (bestStreak - 1) * rules.PtsStreakStep
			if bonus > rules.PtsStreakCap {
				bonus = rules.PtsStreakCap
			}
			totalPoints += bonus
		}

		avgSleep := 0
		if sleepEntriesCount > 0 {
			avgSleep = totalSleepMin / sleepEntriesCount
		}

		sleepPct := 0.0
		if totalScheduled > 0 {
			sleepPct = float64(sleepOkCount) / float64(totalScheduled) * 100.0
		}

		sRow := model.PrestasiScore{
			EmployeeNIK:        emp.NIK,
			PeriodDays:         periodDays,
			TotalPoints:        totalPoints,
			StreakDays:         bestStreak,
			CurrentStreak:      currentStreak,
			AttCount:           attCount,
			SleepPct:           sleepPct,
			PenaltyCount:       penaltyCount,
			SleepOkCount:       sleepOkCount,
			TotalScheduledDays: totalScheduled,
			QualifiedDays:      qualifiedDays,
			CoverDays:          coverDays,
			AvgSleepMin:        avgSleep,
			LateCount:          lateCount,
		}
		scores = append(scores, sRow)

		// Assign Badges (matching FE keys) - using configurable thresholds
		badgeStreak7 := rules.BadgeStreak7Threshold
		badgeStreak14 := rules.BadgeStreak14Threshold
		if badgeStreak7 <= 0 {
			badgeStreak7 = 7
		}
		if badgeStreak14 <= 0 {
			badgeStreak14 = 14
		}
		if bestStreak >= badgeStreak14 {
			badges = append(badges, model.PrestasiBadge{EmployeeNIK: emp.NIK, BadgeKey: BadgeStreak14})
		} else if bestStreak >= badgeStreak7 {
			badges = append(badges, model.PrestasiBadge{EmployeeNIK: emp.NIK, BadgeKey: BadgeStreak7})
		}
		if totalScheduled > 0 && sleepOkCount == totalScheduled {
			badges = append(badges, model.PrestasiBadge{EmployeeNIK: emp.NIK, BadgeKey: BadgePerfectSleep})
		}
		if totalScheduled > 0 && lateCount == 0 {
			badges = append(badges, model.PrestasiBadge{EmployeeNIK: emp.NIK, BadgeKey: BadgeNeverLate})
		}
		if totalScheduled > 0 && penaltyCount == 0 {
			badges = append(badges, model.PrestasiBadge{EmployeeNIK: emp.NIK, BadgeKey: BadgeNoPenalty})
		}

		history = append(history, empDays...)
	}

	sortScoresByPoints(scores)
	for i := range scores {
		scores[i].Rank = i + 1
	}

	return s.repo.SavePrestasiData(periodDays, scores, history, badges)
}

// simulateDay simulates one date for ALL operators simultaneously.
// Handles replacement pairing: unfit operator gets penalty, replacement gets cover bonus.
func (s *PrestasiService) simulateDay(dateStr string, emps []model.Employee, alloc map[string]string, unitCodes []string, rules PrestasiRules, ftwRules FtwRules) map[string]model.PrestasiHistoryEntry {
	result := make(map[string]model.PrestasiHistoryEntry)

	type Row struct {
		NIK       string
		Name      string
		Code      string
		SleepMin  int
		SleepOk   bool
		FtwStatus string
		RestHours int
		Attended  bool
		Late      bool
		UnitCode  string
		Classes   []string
	}

	rows := make([]Row, 0, len(emps))
	for _, emp := range emps {
		hVal := hashStr(emp.NIK + dateStr)

		// Use real roster schedule from DB — fallback to OFF if no schedule exists
		code := "OFF"
		if schedCode, err := s.repo.GetRosterSchedule(emp.NIK, dateStr); err == nil && schedCode != "" {
			code = schedCode
		} else {
			// No roster data for this date — treat as OFF (not scheduled)
			code = "OFF"
		}

		realAtt, _ := s.repo.GetAttendanceRecord(emp.NIK, dateStr)
		realFTW, _ := s.repo.GetFTWRecord(emp.NIK, dateStr)

		isAttended := false
		isLate := false
		sleepMin := 0

		if realAtt != nil {
			isAttended = realAtt.St == "hadir" || realAtt.St == "terlambat"
			isLate = realAtt.St == "terlambat"
		} else {
			// No attendance record — only mark as attended if scheduled
			isAttended = false
			isLate = false
		}

		if realFTW != nil && realFTW.SleepMin != nil {
			sleepMin = *realFTW.SleepMin
		} else {
			// No FTW log — treat as no log (belum)
			sleepMin = 0
		}

		ftwEval := FtwEvaluateWithRules(&sleepMin, ftwRules)

		unitCode := ""
		if alloc != nil {
			unitCode = alloc[emp.NIK]
		}
		if unitCode == "" {
			unitCode = s.resolveUnitCodeFallback(emp.NIK, dateStr, hVal, unitCodes)
		}

		// Collect competency classes
		classes := make([]string, 0, len(emp.Komp))
		for _, k := range emp.Komp {
			classes = append(classes, k.Class)
		}

		rows = append(rows, Row{
			NIK:       emp.NIK,
			Name:      emp.Name,
			Code:      code,
			SleepMin:  sleepMin,
			SleepOk:   ftwEval.CanWork,
			FtwStatus: string(ftwEval.Status),
			RestHours: ftwEval.RestHours,
			Attended:  isAttended,
			Late:      isLate,
			UnitCode:  unitCode,
			Classes:   classes,
		})
	}

	// Separate scheduled vs not scheduled
	var scheduled, notScheduled []Row
	for _, r := range rows {
		if r.Code == "D" || r.Code == "N" {
			scheduled = append(scheduled, r)
		} else {
			notScheduled = append(notScheduled, r)
		}
	}

	// Fit = attended AND can work; unfit = scheduled but not both
	var fit, unfit []Row
	for _, r := range scheduled {
		if r.Attended && r.SleepOk {
			fit = append(fit, r)
		} else {
			unfit = append(unfit, r)
		}
	}

	// Replacement pool: not scheduled but can work
	var pool []Row
	for _, r := range notScheduled {
		if r.SleepOk {
			pool = append(pool, r)
		}
	}

	// 1. Fit operators get points
	for _, r := range fit {
		pts := rules.PtsBase
		if !r.Late {
			pts += rules.PtsOntime
		}
		if r.SleepMin >= rules.SleepMinGreat {
			pts += rules.PtsSleep
		}
		// Spare: base points only (no ontime/sleep bonus)
		if r.FtwStatus == "spare" {
			pts = rules.PtsBase
		}

		result[r.NIK] = model.PrestasiHistoryEntry{
			EmployeeNIK: r.NIK,
			RecordDate:  dateStr,
			ShiftCode:   r.Code,
			UnitCode:    r.UnitCode,
			AttStatus:   attStatus(r.Attended, r.Late),
			ClockIn:     clockInOf(r.Code, r.Late),
			Late:        r.Late,
			SleepMin:    r.SleepMin,
			AttOk:       true,
			SleepOk:     true,
			FtwStatus:   r.FtwStatus,
			RestHours:   r.RestHours,
			Outcome:     OutcomeQualified,
			Points:      pts,
		}
	}

	// 2. Unfit operators get penalty, and find replacement
	// Replacement MUST have competency matching the unit's class (like FE)
	usedCover := make(map[string]bool)
	for _, r := range unfit {
		// Determine required class from unit code
		needClass := s.classOfUnit(r.UnitCode)

		// Find replacement from pool with matching competency
		var cover *Row
		for i := range pool {
			if usedCover[pool[i].NIK] {
				continue
			}
			// If unit has a class, replacement must have matching competency
			if needClass != "" && !hasClass(pool[i].Classes, needClass) {
				continue
			}
			cover = &pool[i]
			usedCover[pool[i].NIK] = true
			break
		}

		outcome := OutcomeReplacedAbsent
		if r.Attended {
			outcome = OutcomeReplacedSleep
		}

		result[r.NIK] = model.PrestasiHistoryEntry{
			EmployeeNIK: r.NIK,
			RecordDate:  dateStr,
			ShiftCode:   r.Code,
			UnitCode:    r.UnitCode,
			AttStatus:   attStatus(r.Attended, r.Late),
			ClockIn:     clockInOf(r.Code, r.Late),
			Late:        r.Late,
			SleepMin:    r.SleepMin,
			AttOk:       r.Attended,
			SleepOk:     false,
			FtwStatus:   r.FtwStatus,
			RestHours:   r.RestHours,
			Outcome:     outcome,
			Points:      rules.PtsPenalty,
		}
		if cover != nil {
			entry := result[r.NIK]
			entry.CounterpartNik = cover.NIK
			entry.CounterpartName = cover.Name
			result[r.NIK] = entry
		}

		// 3. Replacement gets cover bonus
		if cover != nil {
			pts := rules.PtsBase + rules.PtsCover
			if !cover.Late {
				pts += rules.PtsOntime
			}
			if cover.SleepMin >= rules.SleepMinGreat {
				pts += rules.PtsSleep
			}
			if cover.FtwStatus == "spare" {
				pts = rules.PtsBase + rules.PtsCover
			}

			result[cover.NIK] = model.PrestasiHistoryEntry{
				EmployeeNIK:     cover.NIK,
				RecordDate:      dateStr,
				ShiftCode:       cover.Code,
				UnitCode:        r.UnitCode, // same unit as the replaced operator
				AttStatus:       attStatus(true, cover.Late),
				ClockIn:         clockInOf(r.Code, cover.Late),
				Late:            cover.Late,
				SleepMin:        cover.SleepMin,
				AttOk:           true,
				SleepOk:         true,
				FtwStatus:       cover.FtwStatus,
				RestHours:       cover.RestHours,
				Outcome:         OutcomeReplacement,
				CounterpartNik:  r.NIK,
				CounterpartName: r.Name,
				Points:          pts,
			}
		}
	}

	// 4. Not scheduled and not called as replacement → netral
	for _, r := range notScheduled {
		if _, ok := result[r.NIK]; ok {
			continue
		}
		result[r.NIK] = model.PrestasiHistoryEntry{
			EmployeeNIK: r.NIK,
			RecordDate:  dateStr,
			ShiftCode:   r.Code,
			UnitCode:    "",
			AttStatus:   "off",
			ClockIn:     "",
			Late:        false,
			SleepMin:    r.SleepMin,
			AttOk:       false,
			SleepOk:     r.SleepOk,
			FtwStatus:   r.FtwStatus,
			RestHours:   r.RestHours,
			Outcome:     OutcomeNotScheduled,
			Points:      0,
		}
	}

	return result
}

func attStatus(attended, late bool) string {
	if !attended {
		return "belum"
	}
	if late {
		return "terlambat"
	}
	return "hadir"
}

func clockInOf(code string, late bool) string {
	if code == "N" {
		if late {
			return "18:05"
		}
		return "18:00"
	}
	if late {
		return "06:05"
	}
	return "06:00"
}

// loadAllocationsByDate loads real unit allocations for the period
func (s *PrestasiService) loadAllocationsByDate(periodDays int) map[string]map[string]string {
	result := make(map[string]map[string]string)
	if s.fleetRepo == nil {
		return result
	}

	for d := periodDays; d >= 1; d-- {
		dateStr := time.Now().AddDate(0, 0, -d).Format("2006-01-02")
		for _, shift := range []string{"pagi", "malam"} {
			ops, err := s.fleetRepo.GetAllocationOperators(dateStr, shift)
			if err != nil {
				continue
			}
			if _, ok := result[dateStr]; !ok {
				result[dateStr] = make(map[string]string)
			}
			for unitCode, nik := range ops {
				if nik != "" {
					result[dateStr][nik] = unitCode
				}
			}
		}
	}
	return result
}

// resolveUnitCodeFallback returns a real unit code from units_db when no allocation exists
func (s *PrestasiService) resolveUnitCodeFallback(nik, dateStr string, hVal uint32, unitCodes []string) string {
	if len(unitCodes) == 0 {
		return fmt.Sprintf("EX-%d", 7000+int(hVal%13))
	}
	// Deterministic pick from real unit codes
	idx := int(hVal % uint32(len(unitCodes)))
	return unitCodes[idx]
}

// classOfUnit returns the eq class of a unit code (from units_db)
func (s *PrestasiService) classOfUnit(unitCode string) string {
	if unitCode == "" || s.fleetRepo == nil {
		return ""
	}
	units, err := s.fleetRepo.GetUnitDB()
	if err != nil {
		return ""
	}
	for _, u := range units {
		if u.Code == unitCode {
			return u.Cls
		}
	}
	return ""
}

// hasClass checks if a competency class list contains the required class
func hasClass(classes []string, need string) bool {
	for _, c := range classes {
		if c == need {
			return true
		}
	}
	return false
}

func hashStr(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func sortScoresByPoints(scores []model.PrestasiScore) {
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].TotalPoints > scores[i].TotalPoints {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
}
