package model

import (
	"time"
)

type Widget struct {
	Id               string    `db:"id" form:"id"  primary:"true"`
	WidgetIdentifier string    `db:"widget_identifier" form:"widget_identifier"`
	Property         string    `db:"property" form:"property"`
	Config           string    `db:"config" form:"config"`
	IsShare          int      `db:"is_share" form:"is_share"`
	UserId           string    `db:"user_id" form:"user_id"`
	LastUpdateAt     time.Time `db:"last_update_at" form:"last_update_at"`
	TenantId         int      `db:"tenant_id_" form:"tenant_id_"`
}

func (Widget) GetTableName() string { return "config_detail" }
