package models

import (
	"github.com/Dataman-Cloud/swan/src/types"
)

type Project struct {
	Name          string         `json:"name"`
	CreateTime    string         `json:"createtime"`
	Applications  AppUpdateStage `json:"applications"`
	UpdateHistory []ExecHistory  `json:"update_history"`
}

type AppUpdateStage struct {
	App                 types.Version
	RollingUpdatePolicy []AppUpdatePolicy
}

type AppUpdatePolicy struct {
	InstancesToUpdate int64             `json:"instances_to_update"`
	Trigger           string            `json:"trigger"`
	RollbackPolicy    AppRollbackPolicy `json:"rollback_policy"`
}

type AppRollbackPolicy struct {
	AutoRollback      bool  `json:"auto_rollback"`
	RollbackCondition int64 `json:"rollback_condition"`
}

type ExecHistory struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}
