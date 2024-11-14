package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type BaseModel struct {
	Id        string         `gorm:"column:id;type:varchar(64);primaryKey;not null;unique;comment:id" json:"id,omitempty"`
	CreatedAt LocalTime      `gorm:"column:created_at;type:datetime;autoCreateTime:milli,comment:创建时间" json:"createdAt,omitempty"`
	UpdatedAt LocalTime      `gorm:"column:updated_at;type:datetime;autoUpdateTime:milli;comment:更新时间" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt,omitempty"`
	Remark    string         `gorm:"column:remark;type:varchar(128);comment:备注" json:"remark,omitempty"`
}

func (b *BaseModel) TableName() string {
	return "baseModel"
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = GetUuid()
	return
}
func GetUuid() string {
	str := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
	return str
}
