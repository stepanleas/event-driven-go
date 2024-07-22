# Handling refunds

It's time to handle the refund command.
To handle refunds, we need to do two things:

1. Void the receipt.
2. Refund the payment.

Clients from `github.com/ThreeDotsLabs/go-event-driven/common/clients` support both operations:

```go
clients.Payments.PutRefundsWithResponse(ctx, payments.PaymentRefundRequest{
    // we are using TicketID as a payment reference
    PaymentReference: command.TicketID,
    Reason:           "customer requested refund",
    DeduplicationId:  &command.Header.IdempotencyKey,
})
```

and

```go
clients.Receipts.PutVoidReceiptWithResponse(ctx, receipts.VoidReceiptRequest{
    Reason:       "customer requested refund",
    TicketId:     command.TicketID,
    IdempotentId: &command.Header.IdempotencyKey,
})
```

They are both idempotent, so we can handle them in a single command handler.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Do you remember how idempotency works? If not, check the [idempotency key](/trainings/go-event-driven/exercise/d295e9e2-4cb4-49b2-bf73-28635208a78d) exercise.

This API is called from the browser, so the idempotency key will not be set by the client.
You should generate it by yourself in an HTTP handler.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Implement the refund command handler.

For now, let's assume that commands will be handled eventually and won't spin forever.
We will take care of that in the _observability and monitoring_ module.
