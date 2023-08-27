package reporter

import (
	"log"

	appreporter "gotrack/pkg/app/reporter"
)

func NewReporter(logger *log.Logger) appreporter.Reporter {
	return &reporter{logger: logger}
}

type reporter struct {
	logger *log.Logger
}

func (r *reporter) Report(msg string) {
	r.logger.Println(msg)
}
