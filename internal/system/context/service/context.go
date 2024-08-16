package service

import (
	"api/internal/model/config"
	"api/internal/model/log"
	"api/internal/system/context/service/integration"
)

type Context interface {
	Config() *config.Config
	IntegrationContext() integration.Context
	Log(level log.Level, message string, labels ...*log.Label)
}
