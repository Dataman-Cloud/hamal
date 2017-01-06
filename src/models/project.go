package models

import (
	"github.com/Dataman-Cloud/swan/src/types"
)

type Project struct {
	Name         string         `json:"name"`
	Applications AppUpdateStage `json:"applications"`
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
