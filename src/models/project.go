package models

import (
	"github.com/Dataman-Cloud/swan/src/types"
)

type Project struct {
	Name         string           `json:"name"`
	CreateTime   string           `json:"createtime"`
	Applications []AppUpdateStage `json:"applications"`
	Status       int              `json:"-"`
}

type AppUpdateStage struct {
	AppId               string            `json:"app_id"`
	App                 types.Version     `json:"orchestration"`
	RollingUpdatePolicy []AppUpdatePolicy `json:"rolling_update_policy"`
	NextStage           int64             `json:"next_stage"`
	Status              string            `json:"status"`
}

type AppUpdatePolicy struct {
	InstancesToUpdate int64  `json:"instances_to_update"`
	Trigger           string `json:"trigger"`
	//RollbackPolicy    AppRollbackPolicy `json:"rollback_policy"`
}

type AppRollbackPolicy struct {
	AutoRollback      bool  `json:"auto_rollback"`
	RollbackCondition int64 `json:"rollback_condition"`
}

type Hamal struct {
	Name string `json:"name"`
}

type RollPolicy struct {
	AppId string `json:"app_id"`
}
