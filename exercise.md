# Failed to Book Return Flight Ticket

In this exercise, it won't be possible to book **return flight** tickets (previously, it was inbound flight tickets).
`FlightBookingFailed_v1` should be already properly emitted (you did this in the previous exercise).
We will also refund tickets, but we are missing the command handler for handling the inbound flight ticket, which was already booked.

```mermaid
stateDiagram-v2
    [*] --> VipBundleInitialized
    VipBundleInitialized --> BookShowTickets : OnVipBundleInitialized
    BookShowTickets --> BookingMade : Success
    BookShowTickets --> BookingFailed : Failure
    BookingMade --> BookInboundFlight : OnBookingMade
    BookInboundFlight --> FlightBooked : Success
    BookInboundFlight --> FlightBookingFailed : Failure
    FlightBooked --> BookReturnFlight : InboundBooked & ReturnNotBooked
    FlightBooked --> BookTaxi : Both Flights Booked
    BookReturnFlight --> FlightBooked : Success
    BookReturnFlight --> FlightBookingFailed : Failure
    FlightBookingFailed --> Rollback : OnFlightBookingFailed
    BookTaxi --> TaxiBooked : Success
    BookTaxi --> TaxiBookingFailed : Failure
    TaxiBooked --> VipBundleFinalized : OnTaxiBooked
    TaxiBookingFailed --> Rollback : OnTaxiBookingFailed
    Rollback --> RefundTicket : Show Booking Exists
    Rollback --> CancelInboundFlight : Inbound Flight Exists
    Rollback --> CancelReturnFlight : Return Flight Exists
    Rollback --> FailedState
    VipBundleFinalized --> FinalizedState
    FailedState --> [*] : Failed Process
    FinalizedState --> [*] : Completed Process
```

## Exercise

File: `project/main.go`

Implement the `CancelFlightTickets` command handler.

It should call `DeleteFlightTicketsTicketIdWithResponse` for **each** ticket ID.

```go
resp, err := t.clients.Transportation.DeleteFlightTicketsTicketIdWithResponse(ctx, ticketID)
if err != nil {
	return fmt.Errorf("failed to cancel flight tickets: %w", err)
}
```

`CancelFlightTickets` should already be published for inbound tickets by your process manager.
**You should keep the original functionality, and show tickets should be canceled as well.**

