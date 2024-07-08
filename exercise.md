# Project: Payloads


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Our integrations work correctly, but adding just the ticket IDs isn't that helpful because the operations team still has to manually look up the ticket details.
We can help them out by adding more details to the external systems.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

First, we're going to extend the webhook endpoint to accept more than ticket IDs.
We're dropping support for `POST /tickets-confirmation`.
Instead, we're now going to expose `POST /tickets-status`.

Previously, the endpoint accepted a list of ticket IDs; now it's going to accept a list of tickets.
Each ticket has a ticket ID, a status, a customer email, and a price.

Here's an example incoming HTTP request:

```json
{
  "tickets": [
    {
      "ticket_id": "ticket-1",
      "status": "confirmed",
      "customer_email": "user@example.com",
      "price": {
        "amount": "50.00",
        "currency": "EUR"
      }
    }
  ]
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

Note that we use `string` types for the price amount.
This is on purpose.
Don't use `float64` for money, or you may lose precision.

</span>
	</div>
	</div>

The endpoint should still publish two messages per ticket: one for issuing a receipt and one for appending it to the spreadsheet.
However, we want to send more details to both external systems.

For issuing the receipts, we need to include the price:

```go
body := receipts.PutReceiptsJSONRequestBody{
	TicketId: request.TicketID, 
	Price: receipts.Money{
		MoneyAmount:   request.Price.Amount, 
		MoneyCurrency: request.Price.Currency,
	},
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

Note that there are two `Money` types: one comes from the HTTP request, and the other from the `receipts` package.
This is correct.
They come from different contexts, and you need to map one to the other.

</span>
	</div>
	</div>

And for each row appended to the sheet, we're going to add the customer's email and the price:

```go
spreadsheetsClient.AppendRow(
	msg.Context(), 
	"tickets-to-print", 
	[]string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency},
)
```

To make this work, we need to publish proper JSON message payloads to the Pub/Sub topics.

The `IssueReceiptPayload` looks like this:

```json
{
  "ticket_id": "ticket-1",
  "price": {
    "amount": "50.00",
    "currency": "EUR"
  }
}
```

And the `AppendToTrackerPayload` like this:

```json
{
  "ticket_id": "ticket-1",
  "customer_email": "user@example.com",
  "price": {
    "amount": "50.00",
    "currency": "EUR"
  }
}
```

Create structs for these payloads. 
In the HTTP handler, marshal them to JSON and publish messages with them on the topics.

Change the Router handlers so that they unmarshal the payload on the structs.
Then use the values from the structs to call the external systems.

Modify the receipts client so it accepts a struct as an argument:

```go
type IssueReceiptRequest struct {
	TicketID string
	Price    Money
}

func (c ReceiptsClient) IssueReceipt(ctx context.Context, request IssueReceiptRequest) error {
```
