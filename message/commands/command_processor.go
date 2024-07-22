package commands

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func AddCommandProcessorHandlers(cp *cqrs.CommandProcessor) {
	cp.AddHandlers()
}
