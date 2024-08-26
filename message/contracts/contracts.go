package contracts

import (
	"context"
	"tickets/entities"

	"github.com/google/uuid"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type TicketRepository interface {
	FindAll(ctx context.Context) ([]entities.Ticket, error)
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type ShowRepository interface {
	Add(ctx context.Context, show entities.Show) error
	FindAll(ctx context.Context) ([]entities.Show, error)
	FindByID(ctx context.Context, showID uuid.UUID) (entities.Show, error)
}

type BookingRepository interface {
	Add(ctx context.Context, booking entities.Booking) error
}

type VipBundleRepository interface {
	Add(ctx context.Context, vipBundle entities.VipBundle) error
	Get(ctx context.Context, vipBundleID uuid.UUID) (entities.VipBundle, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (entities.VipBundle, error)

	UpdateByID(
		ctx context.Context,
		bookingID uuid.UUID,
		updateFn func(vipBundle entities.VipBundle) (entities.VipBundle, error),
	) (entities.VipBundle, error)

	UpdateByBookingID(
		ctx context.Context,
		bookingID uuid.UUID,
		updateFn func(vipBundle entities.VipBundle) (entities.VipBundle, error),
	) (entities.VipBundle, error)
}

type DataLake interface {
	FindAll(ctx context.Context) ([]entities.DataLakeEvent, error)
	Store(ctx context.Context, event entities.DataLakeEvent) error
}

type FilesAPI interface {
	UploadFile(ctx context.Context, fileID string, fileContent string) error
}

type DeadNationApi interface {
	BookInDeadNation(ctx context.Context, booking entities.DeadNationBooking) error
}

type TransportationService interface {
	BookFlight(ctx context.Context, request entities.BookFlightTicketRequest) (entities.BookFlightTicketResponse, error)
	BookTaxi(ctx context.Context, request entities.BookTaxiRequest) (entities.BookTaxiResponse, error)
}
