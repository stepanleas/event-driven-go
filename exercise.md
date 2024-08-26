# Book Transportation for the VIP Bundle

The hardest part is done — our process manager works.
Now we can add the next part of the flow: booking transportation.

We need to book flight tickets and taxis.
We have prepared API clients for the `transportation` API in `github.com/ThreeDotsLabs/go-event-driven/common/clients`.

Flights can be booked with:

```go
resp, err := t.clients.Transportation.PutFlightTicketsWithResponse(ctx, transportation.BookFlightTicketRequest{
	CustomerEmail:  request.CustomerEmail,
	FlightId:       request.FlightID,
	PassengerNames: request.PassengerNames,
	ReferenceId:    request.ReferenceId,
	IdempotencyKey: request.IdempotencyKey,
})
```

Taxis can be booked with:

```go
resp, err := t.clients.Transportation.PutTaxiBookingWithResponse(ctx, transportation.TaxiBookingRequest{
	CustomerEmail:      request.CustomerEmail,
	NumberOfPassengers: request.NumberOfPassengers,
	PassengerName:      request.PassengerName, // this should be name of the first passenger in Vip Bundle
	ReferenceId:        request.ReferenceId,
	IdempotencyKey:     request.IdempotencyKey,
})
```
Use the name of the first passenger in the VIP Bundle as `PassengerName` for taxi booking.

The taxi provider will provide the right number of cars based on the number of passengers (so we don't need to think about that).

If booking was successful, the API will return Status Created (201) and a response body with:

- The taxi booking ID for a taxi
- Flight ticket IDs for a flight

If your request is  invalid, the API will return Status Bad Request (400) and an error message.


Taxis should be booked via the `BookTaxi` command and flights with `BookFlight`.

## Exercise

File: `project/main.go`

Add support for booking taxi and flight tickets to our VIP bundle when the `BookFlight` and `BookTaxi` 
commands are received from our process manager.


<div class="accordion" id="hints-accordion">

<div class="accordion-item">
	<h3 class="accordion-header" id="hints-accordion-header-1">
	<button class="accordion-button fs-4 fw-semibold collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#hints-accordion-body-1" aria-expanded="false" aria-controls="hints-accordion">
		Hint #1
	</button>
	</h3>
	<div id="hints-accordion-body-1" class="accordion-collapse collapse" aria-labelledby="hints-accordion-header-1" data-bs-parent="#hints-accordion">
	<div class="accordion-body">

Make sure to pick a different idempotency key for each request you make to the transportation API.

Transportation provider will ignore any extra requests that have the same idempotency key, no matter what information is inside them. 
This means if you use the same key for two different rides or flights, one might not get counted. So, always use a unique key for each request to avoid missing out on any services or losing data.

You can read more about that in [Idempotency Key for issuing receipts](/trainings/go-event-driven/exercise/92ec4eb7-2507-4ad0-850d-28089a587d3e).

</div>
	</div>
	</div>

</div>
