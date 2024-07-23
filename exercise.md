# Project: Emitting ReceiptIssued

In our read model, **we need to store when a receipt for the ticket was issued and the receipt number.**
We don't store this data anywhere yet.

Let's emit the `TicketReceiptIssued` event. 
You should emit an event that looks more or less like this (we will not assert the exact values):

```go
type TicketReceiptIssued struct {
	Header EventHeader `json:"header"`

	TicketID      string `json:"ticket_id"`
	ReceiptNumber string `json:"receipt_number"`

	IssuedAt time.Time `json:"issued_at"`
}
```

Please note that we have a separate `issued_at` field in this event, and we are not using the value from header.
The issue date is important in the accounting domain (in terms of taxes, etc.).
Time in the event header is informational, and it should be not treated as an accurate value.
Think of it like this: the time in the header is the time the event has been published, not the time the receipt has been issued.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

As an alternative to emitting an event, you could also store this data within a service that is responsible for issuing receipts.
You could return this data via API and call this API while building the read model.

This has a few drawbacks:
- you need to call another service - explicit coupling,
- you depend on another service's availability,
- you put higher pressure on another service.

But in some cases, if it's not that easy to emit a new event, it's a good alternative to consider.

</span>
	</div>
	</div>


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

What if publishing the event fails?
We can return an error, and retry. 
Thanks to the idempotency key, this operation will be safe.

</span>
	</div>
	</div>

## Exercise

Emit the `TicketReceiptIssued` event to the `events.TicketReceiptIssued` topic, after the receipt is issued.
The `receipts` service API client returns the receipt number and issue date.

File: `project/main.go`
