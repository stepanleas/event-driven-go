# Calling Dead Nation

Remember that our goal is to proxy requests to Dead Nation?
To replace the endpoint in production, we need to call the Dead Nation API.

## Exercise

File: `project/main.go`

Use the `BookingMade` event that you emitted in the [previous](/trainings/go-event-driven/exercise/148ed9cd-29d4-4132-a57d-d1499e298897) exercise.

Clients from `github.com/ThreeDotsLabs/go-event-driven/common/clients` support calling the Dead Nation API.

```go
resp, err := h.deadNationClient.PostTicketBookingWithResponse(
    ctx,
    dead_nation.PostTicketBookingRequest{
        CustomerAddress: booking.CustomerEmail,
        EventId:         booking.DeadNationEventID,
        NumberOfTickets: booking.NumberOfTickets,
        BookingId:       booking.BookingID,
    },
)
```


As usually occurs, names from external APIs do not usually correspond 1:1 to our codebase.
For example: `CustomerAddress` is `CustomerEmail` in our codebase.

**Warning:** `EventId` should be the `dead_nation_id` from the store show request [previous exercise](/trainings/go-event-driven/exercise/d5291467-1e0b-442b-ac2e-e79141e96ff9).
**You should get this value from the database.**
This is intentionally a different name: EventID is a term used by Dead Nation, but we prefer name `ShowID` (so it's not confusing with our events in Pub/Sub).


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Repositories or clients are usually good places for making the translation from external language to internal.
Thanks to that, we can keep language inside our application free from external influences.

</span>
	</div>
	</div>

If everythng went fine, Dead Nation should call your `POST /ticket-status` endpoint.


