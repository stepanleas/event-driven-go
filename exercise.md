# Printing Tickets


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Our operations team is now generating and printing tickets by hand.
This was a good strategy to roll out the product quickly, but it's not a good long-term solution because
they already struggle with the number of tickets they need to print.

Let's help our ops team by generating tickets for them!

We will use the `files` service to store the tickets.
The `files` will be available via the `gateway`, like other services.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Implement an event handler that will be triggered by the `TicketBookingConfirmed` event.

It should store ticket content with the API client `*clients.Clients` from the `github.com/ThreeDotsLabs/go-event-driven/common/clients`
(`Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, fileContent)` method).

The file name should have the format `[Ticket ID]-ticket.html`.
The content doesn't matter, it's just important that it contain the ticket ID, price, and amount.

It's not necessary to do anything on `TicketBookingCanceled` — the volume of this is low, and it's not a problem for ops to handle it.

Do you remember the discussion of eventual consistency? The client will return 409 when the file already exists.
We will use a similar strategy as in the previous module. If this error happens, you need to handle it gracefully.
It's worth adding a log in that situation, so you will know what happened in case any issues arise.

```go
import (
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

if resp.StatusCode() == http.StatusConflict {
	log.FromContext(ctx).Infof("file %s already exists", fileID)
	return nil
}
```


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Note that we can add this functionality without changing any existing code.
In real life, it could be even implemented by a different team that has access to the events.

</span>
	</div>
	</div>

