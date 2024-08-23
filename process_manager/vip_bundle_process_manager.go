package process_manager

import (
	"context"
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type VipBundleProcessManager struct {
	commandBus *cqrs.CommandBus
	eventBus   *cqrs.EventBus
	repository contracts.VipBundleRepository
}

func NewVipBundleProcessManager(
	commandBus *cqrs.CommandBus,
	eventBus *cqrs.EventBus,
	repository contracts.VipBundleRepository,
) *VipBundleProcessManager {
	return &VipBundleProcessManager{
		commandBus: commandBus,
		eventBus:   eventBus,
		repository: repository,
	}
}

func (v VipBundleProcessManager) OnVipBundleInitialized(ctx context.Context, event *entities.VipBundleInitialized_v1) error {
	vb, err := v.repository.Get(ctx, event.VipBundleID)
	if err != nil {
		return err
	}

	return v.commandBus.Send(ctx, entities.BookShowTickets{
		BookingID:       vb.BookingID,
		CustomerEmail:   vb.CustomerEmail,
		NumberOfTickets: vb.NumberOfTickets,
		ShowId:          vb.ShowId,
	})
}

func (v VipBundleProcessManager) OnBookingMade(ctx context.Context, event *entities.BookingMade_v1) error {
	vb, err := v.repository.UpdateByBookingID(
		ctx,
		event.BookingID,
		func(vipBundle entities.VipBundle) (entities.VipBundle, error) {
			vipBundle.BookingMadeAt = &event.Header.PublishedAt

			return vipBundle, nil
		},
	)
	if err != nil {
		return err
	}

	return v.commandBus.Send(ctx, entities.BookFlight{
		CustomerEmail:  vb.CustomerEmail,
		FlightID:       vb.InboundFlightID,
		Passengers:     vb.Passengers,
		ReferenceID:    vb.VipBundleID.String(),
		IdempotencyKey: uuid.NewString(),
	})
}

func (v VipBundleProcessManager) OnTicketBookingConfirmed(ctx context.Context, event *entities.TicketBookingConfirmed_v1) error {
	_, err := v.repository.UpdateByBookingID(
		ctx,
		uuid.MustParse(event.BookingID),
		func(vipBundle entities.VipBundle) (entities.VipBundle, error) {
			eventTicketID := uuid.MustParse(event.TicketID)

			for _, ticketID := range vipBundle.TicketIDs {
				if ticketID == eventTicketID {
					// re-delivery (already stored)
					continue
				}
			}

			vipBundle.TicketIDs = append(vipBundle.TicketIDs, eventTicketID)

			return vipBundle, nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (v VipBundleProcessManager) OnBookingFailed(ctx context.Context, event *entities.BookingFailed_v1) error {
	vb, err := v.repository.GetByBookingID(ctx, event.BookingID)
	if err != nil {
		return err
	}

	return v.rollbackProcess(ctx, vb.VipBundleID)
}

func (v VipBundleProcessManager) OnFlightBooked(ctx context.Context, event *entities.FlightBooked_v1) error {
	vb, err := v.repository.UpdateByID(
		ctx,
		uuid.MustParse(event.ReferenceID),
		func(vipBundle entities.VipBundle) (entities.VipBundle, error) {
			if vipBundle.InboundFlightID == event.FlightID {
				vipBundle.InboundFlightBookedAt = &event.Header.PublishedAt
				vipBundle.InboundFlightTicketsIDs = event.TicketIDs
			}
			if vipBundle.ReturnFlightID == event.FlightID {
				vipBundle.ReturnFlightBookedAt = &event.Header.PublishedAt
				vipBundle.ReturnFlightTicketsIDs = event.TicketIDs
			}

			return vipBundle, nil
		},
	)
	if err != nil {
		return err
	}

	switch {
	case vb.InboundFlightBookedAt != nil && vb.ReturnFlightBookedAt == nil:
		return v.commandBus.Send(ctx, entities.BookFlight{
			CustomerEmail:  vb.CustomerEmail,
			FlightID:       vb.ReturnFlightID,
			Passengers:     vb.Passengers,
			ReferenceID:    vb.VipBundleID.String(),
			IdempotencyKey: uuid.NewString(),
		})
	case vb.InboundFlightBookedAt != nil && vb.ReturnFlightBookedAt != nil:
		return v.commandBus.Send(ctx, entities.BookTaxi{
			CustomerEmail:      vb.CustomerEmail,
			CustomerName:       vb.Passengers[0],
			NumberOfPassengers: vb.NumberOfTickets,
			ReferenceID:        vb.VipBundleID.String(),
			IdempotencyKey:     uuid.NewString(),
		})
	default:
		return fmt.Errorf(
			"unsupported state: InboundFlightBookedAt: %v, ReturnFlightBookedAt: %v",
			vb.InboundFlightBookedAt,
			vb.ReturnFlightBookedAt,
		)
	}
}

func (v VipBundleProcessManager) OnFlightBookingFailed(ctx context.Context, event *entities.FlightBookingFailed_v1) error {
	return v.rollbackProcess(ctx, uuid.MustParse(event.ReferenceID))
}

func (v VipBundleProcessManager) OnTaxiBooked(ctx context.Context, event *entities.TaxiBooked_v1) error {
	vb, err := v.repository.UpdateByID(
		ctx,
		uuid.MustParse(event.ReferenceID),
		func(vb entities.VipBundle) (entities.VipBundle, error) {
			vb.TaxiBookedAt = &event.Header.PublishedAt
			vb.TaxiBookingID = &event.TaxiBookingID

			vb.IsFinalized = true

			return vb, nil
		},
	)
	if err != nil {
		return err
	}

	return v.eventBus.Publish(ctx, entities.VipBundleFinalized_v1{
		Header:      entities.NewEventHeader(),
		VipBundleID: vb.VipBundleID,
	})
}

func (v VipBundleProcessManager) OnTaxiBookingFailed(ctx context.Context, event *entities.TaxiBookingFailed_v1) error {
	return v.rollbackProcess(ctx, uuid.MustParse(event.ReferenceID))
}

func (v VipBundleProcessManager) rollbackProcess(ctx context.Context, vipBundleID uuid.UUID) error {
	vb, err := v.repository.Get(ctx, vipBundleID)
	if err != nil {
		return err
	}

	if vb.BookingMadeAt != nil {
		if err := v.rollbackTickets(ctx, vb); err != nil {
			return err
		}
	}
	if vb.InboundFlightBookedAt != nil {
		if err := v.commandBus.Send(ctx, entities.CancelFlightTickets{
			FlightTicketIDs: vb.InboundFlightTicketsIDs,
		}); err != nil {
			return err
		}
	}
	if vb.ReturnFlightBookedAt != nil {
		if err := v.commandBus.Send(ctx, entities.CancelFlightTickets{
			FlightTicketIDs: vb.ReturnFlightTicketsIDs,
		}); err != nil {
			return err
		}
	}

	_, err = v.repository.UpdateByID(
		ctx,
		vb.VipBundleID,
		func(vb entities.VipBundle) (entities.VipBundle, error) {
			vb.IsFinalized = true
			vb.Failed = true
			return vb, nil
		},
	)

	return err
}

func (v VipBundleProcessManager) rollbackTickets(ctx context.Context, vb entities.VipBundle) error {
	// TicketIDs is eventually consistent, we need to ensure that all tickets are stored
	// for alternative solutions please check "Message Ordering" module
	if len(vb.TicketIDs) != vb.NumberOfTickets {
		return fmt.Errorf(
			"invalid number of tickets, expected %d, has %d: not all of TicketBookingConfirmed_v1 events were processed",
			vb.NumberOfTickets,
			len(vb.TicketIDs),
		)
	}

	for _, ticketID := range vb.TicketIDs {
		if err := v.commandBus.Send(ctx, entities.RefundTicket{
			Header:   entities.NewEventHeader(),
			TicketID: ticketID.String(),
		}); err != nil {
			return err
		}
	}

	return nil
}
