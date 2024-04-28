package models

type PersonalDeduction struct {
	ID                int     `gorm:"primaryKey" json:"id"`
	PersonalDeduction float64 `json:"personal_deduction"`
}


func (PersonalDeduction) TableName() string {
	return "personaldeduction"
}
