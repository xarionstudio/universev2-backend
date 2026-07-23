package model

type PrestasiScore struct {
	ID                 uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeNIK        string  `json:"nik" gorm:"column:employee_nik"`
	PeriodDays         int     `json:"-" gorm:"column:period_days"`
	TotalPoints        int     `json:"ptsTotal" gorm:"column:total_points"`
	Rank               int     `json:"rank" gorm:"column:rank"`
	StreakDays         int     `json:"streak" gorm:"column:streak_days"`
	CurrentStreak      int     `json:"currentStreak" gorm:"column:current_streak"`
	AttCount           int     `json:"attCount" gorm:"column:att_count"`
	SleepPct           float64 `json:"sleepPct" gorm:"column:sleep_pct"`
	LateCount          int     `json:"lateCount" gorm:"column:late_count"`
	PenaltyCount       int     `json:"penaltyCount" gorm:"column:penalty_count"`
	SleepOkCount       int     `json:"sleepOkCount" gorm:"column:sleep_ok_count"`
	TotalScheduledDays int     `json:"totalScheduledDays" gorm:"column:total_scheduled_days"`
	QualifiedDays      int     `json:"qualifiedDays" gorm:"column:qualified_days"`
	CoverDays          int     `json:"coverDays" gorm:"column:cover_days"`
	AvgSleepMin        int     `json:"avgSleepMin" gorm:"column:avg_sleep_min"`
}

func (PrestasiScore) TableName() string { return "prestasi_scores" }

type PrestasiBadge struct {
	ID          uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeNIK string `json:"-" gorm:"column:employee_nik"`
	BadgeKey    string `json:"badge" gorm:"column:badge_key"`
}

func (PrestasiBadge) TableName() string { return "prestasi_badges" }

// PrestasiRecord — API response matching FE PrestasiEntry
type PrestasiRecord struct {
	NIK           string        `json:"nik"`
	Name          string        `json:"name"`
	Dept          string        `json:"dept"`
	Pos           string        `json:"pos"`
	Foto          string        `json:"foto,omitempty"`
	Rank          int           `json:"rank"`
	Points        int           `json:"points"`
	BestStreak    int           `json:"bestStreak"`
	CurrentStreak int           `json:"currentStreak"`
	AttCount      int           `json:"attCount"`
	SleepPct      float64       `json:"sleepPct"`
	AttRate       float64       `json:"attRate"`
	SleepRate     float64       `json:"sleepRate"`
	AvgSleepMin   int           `json:"avgSleepMin"`
	Badges        []string      `json:"badges"`
	LateCount     int           `json:"lateCount"`
	PenaltyDays   int           `json:"penaltyDays"`
	QualifiedDays int           `json:"qualifiedDays"`
	ScheduledDays int           `json:"scheduledDays"`
	CoverDays     int           `json:"coverDays"`
	Days          []PrestasiDay `json:"days,omitempty"`
}

// PrestasiDay — daily audit trail entry matching FE PrestasiDay
type PrestasiDay struct {
	ISO             string `json:"iso"`
	Code            string `json:"code"`
	UnitCode        string `json:"unitCode,omitempty"`
	Att             string `json:"att"`
	ClockIn         string `json:"clockIn"`
	Late            bool   `json:"late"`
	SleepMin        int    `json:"sleepMin"`
	AttOk           bool   `json:"attOk"`
	SleepOk         bool   `json:"sleepOk"`
	FtwStatus       string `json:"ftwStatus"`
	RestHours       int    `json:"restHours"`
	Outcome         string `json:"outcome"`
	CounterpartNik  string `json:"counterpartNik,omitempty"`
	CounterpartName string `json:"counterpartName,omitempty"`
	Points          int    `json:"points"`
}

// PrestasiHistoryEntry — daily audit trail for an operator
type PrestasiHistoryEntry struct {
	ID              uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeNIK     string `json:"-" gorm:"column:employee_nik"`
	RecordDate      string `json:"iso" gorm:"column:record_date"`
	PeriodDays      int    `json:"-" gorm:"column:period_days"`
	ShiftCode       string `json:"code" gorm:"column:shift_code"`
	UnitCode        string `json:"unitCode" gorm:"column:unit_code"`
	AttStatus       string `json:"att" gorm:"column:att_status"`
	ClockIn         string `json:"clockIn" gorm:"column:clock_in"`
	Late            bool   `json:"late" gorm:"column:is_late"`
	SleepMin        int    `json:"sleepMin" gorm:"column:sleep_min"`
	AttOk           bool   `json:"attOk" gorm:"column:att_ok"`
	SleepOk         bool   `json:"sleepOk" gorm:"column:sleep_ok"`
	FtwStatus       string `json:"ftwStatus" gorm:"column:ftw_status"`
	RestHours       int    `json:"restHours" gorm:"column:rest_hours"`
	Outcome         string `json:"outcome" gorm:"column:outcome"`
	CounterpartNik  string `json:"counterpartNik,omitempty" gorm:"column:counterpart_nik"`
	CounterpartName string `json:"counterpartName,omitempty" gorm:"column:counterpart_name"`
	Points          int    `json:"points" gorm:"column:points"`
}

func (PrestasiHistoryEntry) TableName() string { return "prestasi_history" }

const (
	PtsBase      = 10
	PtsSleep     = 3
	PtsOnTime    = 2
	PtsCover     = 5
	PtsStreak    = 2
	PtsStreakCap = 10
	PtsPenalty   = -15
)
