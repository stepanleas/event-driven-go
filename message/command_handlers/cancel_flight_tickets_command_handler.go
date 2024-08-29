package command_handlers

import (
	"context"
	"tickets/entities"
	"tickets/message/contracts"
)

type CancelFlightTicketsCommandHandler struct {
	transportationClient contracts.TransportationService
}

func NewCancelFlightTicketsCommandHandler(transportationService contracts.TransportationService) CancelFlightTicketsCommandHandler {
	return CancelFlightTicketsCommandHandler{transportationClient: transportationService}
}

func (h CancelFlightTicketsCommandHandler) Handle(ctx context.Context, command *entities.CancelFlightTickets) error {
	return h.transportationClient.CancelFlightTickets(
		ctx,
		entities.CancelFlightTicketsRequest{
			TicketIds: command.FlightTicketIDs,
		},
	)
}
