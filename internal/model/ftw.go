package model

import "time"

type FTWStatus string

const (
	FTWStatusFit    FTWStatus = "fit"
	FTWStatusSpare  FTWStatus = "spare"
	FTWStatusPulang FTWStatus = "pulang"
	FTWStatusBelum  FTWStatus = "belum"
)

const (
	SleepFitMin     = 330
	SleepSpare1HMin = 300
	SleepSpare2HMin = 240
)

type FTWEval struct {
	Status    FTWStatus `json:"status"`
	RestHours int       `json:"restHours"`
	CanWork   bool      `json:"canWork"`
}

func EvaluateFTW(sleepMin *int) FTWEval {
	if sleepMin == nil || *sleepMin <= 0 {
		return FTWEval{Status: FTWStatusBelum, RestHours: 0, CanWork: false}
	}
	m := *sleepMin
	if m >= SleepFitMin {
		return FTWEval{Status: FTWStatusFit, RestHours: 0, CanWork: true}
	}
	if m >= SleepSpare1HMin {
		return FTWEval{Status: FTWStatusSpare, RestHours: 1, CanWork: true}
	}
	if m >= SleepSpare2HMin {
		return FTWEval{Status: FTWStatusSpare, RestHours: 2, CanWork: true}
	}
	return FTWEval{Status: FTWStatusPulang, RestHours: 0, CanWork: false}
}

type FTWRecord struct {
	ID          uint      `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	NIK         string    `json:"nik" gorm:"column:employee_nik"`
	Name        string    `json:"name" gorm:"-"`
	Dept        string    `json:"dept" gorm:"-"`
	Shift       string    `json:"shift" gorm:"column:shift"`
	SleepMin    *int      `json:"sleepMin" gorm:"column:sleep_minutes"`
	Sleep       string    `json:"sleep" gorm:"column:sleep_formatted"`
	St          FTWStatus `json:"st" gorm:"column:status"`
	RestHours   int       `json:"restHours" gorm:"column:rest_hours"`
	Hist        []int     `json:"hist" gorm:"-"`
	CanWork     bool      `json:"canWork" gorm:"column:can_work"`
	SendTime    string    `json:"sendTime" gorm:"column:send_time"`
	Date        string    `json:"date" gorm:"column:log_date"`
	SubmittedAt time.Time `json:"submittedAt" gorm:"column:submitted_at"`
}

func (FTWRecord) TableName() string { return "ftw_logs" }

type FTWHistEntry struct {
	D         int       `json:"d"`
	ISO       string    `json:"iso"`
	Date      string    `json:"date"`
	St        int       `json:"st"`
	SleepMin  *int      `json:"sleepMin"`
	Sleep     string    `json:"sleep"`
	Status    FTWStatus `json:"status"`
	RestHours int       `json:"restHours"`
	SendTime  string    `json:"sendTime"`
}
