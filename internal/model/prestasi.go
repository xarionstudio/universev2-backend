package model

type PrestasiScore struct {
	ID           uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeNIK  string  `json:"nik" gorm:"column:employee_nik"`
	PeriodDays   int     `json:"-" gorm:"column:period_days"`
	TotalPoints  int     `json:"ptsTotal" gorm:"column:total_points"`
	Rank         int     `json:"rank" gorm:"column:rank"`
	StreakDays   int     `json:"streak" gorm:"column:streak_days"`
	AttCount     int     `json:"attCount" gorm:"column:att_count"`
	SleepPct     float64 `json:"sleepPct" gorm:"column:sleep_pct"`
	LateCount    int     `json:"lateCount" gorm:"column:late_count"`
	PenaltyCount int     `json:"penaltyCount" gorm:"column:penalty_count"`
}

func (PrestasiScore) TableName() string { return "prestasi_scores" }

type PrestasiBadge struct {
	ID          uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeNIK string `json:"-" gorm:"column:employee_nik"`
	BadgeKey    string `json:"badge" gorm:"column:badge_key"`
}

func (PrestasiBadge) TableName() string { return "prestasi_badges" }

// PrestasiRecord — API response combining score + employee + badges
type PrestasiRecord struct {
	NIK          string   `json:"nik"`
	Name         string   `json:"name"`
	Dept         string   `json:"dept"`
	Rank         int      `json:"rank"`
	PtsTotal     int      `json:"ptsTotal"`
	Streak       int      `json:"streak"`
	AttCount     int      `json:"attCount"`
	SleepPct     float64  `json:"sleepPct"`
	Badges       []string `json:"badges"`
	LateCount    int      `json:"lateCount"`
	PenaltyCount int      `json:"penaltyCount"`
}

const (
	PtsBase    = 100
	PtsSleep   = 50
	PtsOnTime  = 25
	PtsCover   = 75
	PtsStreak  = 10
	PtsPenalty = -50
)
