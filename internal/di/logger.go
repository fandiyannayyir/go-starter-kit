package di

import "go-starter-kit/internal/shared/log"

func NewLogger() log.Logger {
	return log.GetLogger()
}
