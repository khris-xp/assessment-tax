package models

type KReceipt struct {
	ID       int     `gorm:"primaryKey" json:"id"`
	KReceipt float64 `json:"k_receipt"`
}

func (KReceipt) TableName() string {
	return "kreceipt"
}
