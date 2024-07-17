# Idempotency Key for issuing receipts

Time to add support for idempotency keys in our project.


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Our financial team is complaining that they are receiving multiple receipts for the same ticket.
And this is not a team you should mess with!
They are scary, and they are responsible for paying us... so we should fix this as soon as possible.

In the meantime, we asked the team responsible for sending `POST /tickets-status` 
to add `Idempotency-Key` in the incoming request.
We can use it to deduplicate the issuing receipts.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Use `Idempotency-Key` from incoming `POST /tickets-status` requests and propagate it through events.



<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Sending the `Idempotency-Key` header for this endpoint is now part of our contract.
It's a good idea to return `400 Bad Request` if it's not present.
In the case of regression, this error won't propagate.

Another strategy could be to generate an idempotency key on our side when it's missing.
We don't know your context better than you do, though — you should consider tradeoffs
and choose the best strategy for you.

</span>
	</div>
	</div>

Add it to the event body. In our example solution, all our events have the `EventHeader`.
This is a good place to add idempotency key, and you can do it in this way (but it's up to you):


```go
type EventHeader struct {
    ID             string    `json:"id"`
    PublishedAt    time.Time `json:"published_at"`
    IdempotencyKey string    `json:"idempotency_key"`
}

type TicketBookingConfirmed struct {
	Header EventHeader `json:"header"`

	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`

	BookingID string `json:"booking_id"`
}
```

In this example, the idempotency key is optional. If it's not present, it's generated.

```go
func NewEventHeader() EventHeader {
	return EventHeader{
		ID:             uuid.NewString(),
		PublishedAt:    time.Now().UTC(),
		IdempotencyKey: uuid.NewString(),
	}
}

func NewEventHeaderWithIdempotencyKey(idempotencyKey string) EventHeader {
	return EventHeader{
		ID:             uuid.NewString(),
		PublishedAt:    time.Now().UTC(),
		IdempotencyKey: idempotencyKey,
	}
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

Note that we are receiving multiple tickets at once in `POST /tickets-status`.
If we use the same idempotency key for all events generated from this request,
we will be able to deduplicate all receipts for this request.

Our API will return an error, when you will re-use the same idempotency key for a different ticket.
We implemented it in a [similar way to how Stripe](https://stripe.com/docs/api/idempotent_requests) did.

**To overcome this issue, you can concatenate the idempotency key with the ticket id.**

</span>
	</div>
	</div>

After adding the idempotency key, you should send it to the receipts service with `IdempotencyKey`:

```go
resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, receipts.CreateReceipt{
    IdempotencyKey: &request.IdempotencyKey + request.TicketID,

    Price: receipts.Money{
        MoneyAmount:   request.Price.Amount,
        MoneyCurrency: request.Price.Currency,
    },
    TicketId: request.TicketID,
})
```
