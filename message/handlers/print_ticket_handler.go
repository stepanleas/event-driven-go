package handlers

import (
	"context"
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type PrintTicketHandler struct {
	filesAPI contracts.FilesAPI
}

func NewPrintTicketHandler(filesAPI contracts.FilesAPI) PrintTicketHandler {
	return PrintTicketHandler{filesAPI: filesAPI}
}

func (h PrintTicketHandler) Handle(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Printing ticket")

	ticketHTML := `
		<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket ` + event.TicketID + `</h1>
				<p>Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>	
			</body>
		</html>
`

	ticketFile := event.TicketID + "-ticket.html"

	err := h.filesAPI.UploadFile(ctx, ticketFile, ticketHTML)
	if err != nil {
		return fmt.Errorf("failed to upload ticket file: %w", err)
	}

	return nil
}
