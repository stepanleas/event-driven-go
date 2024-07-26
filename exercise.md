# Handling out-of-order `TicketBookingConfirmed` and `TicketBookingCanceled` events

We can't guess how you implemented removing tickets from the `tickets` table.
However, there is a chance that your code doesn't assume that 
`TicketBookingConfirmed` and `TicketBookingCanceled` events can be out of order.

In our previous example solution, we were using a query like this one:

```sql
DELETE FROM tickets WHERE ticket_id = $1
```

It will execute successfully even if there is no row for the ticket in the table yet.
**However, some tickets may be not removed if the `TicketBookingCanceled` event is 
processed before `TicketBookingConfirmed` because there is nothing to remove.**

In theory, we could check if there is a row for the ticket in the table before removing it, and return an error if there isn't.
However there is one downside: Our repository will be no longer be [idempotent](/trainings/go-event-driven/exercise/8c31d18a-b5ae-4d6a-9d1b-a057be5e4b2c).
**When we receive the `TicketBookingCanceled` event again, it will fail because there will be no row for the ticket in the table.**

One of the solutions here is to use soft delete.
We need to add a `deleted_at` column to the `tickets` table and exclude those rows from the result when querying for tickets.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

A repository pattern is pretty useful in such cases — you have one central place where you can change your querying logic.
When you are not querying your data in multiple places in your codebase,
you don't need to be afraid that you will forget to update your logic somewhere.

Bonus points for writing [integration tests](/trainings/go-event-driven/exercise/462e1ede-56d0-4aa0-ae2c-f51493606bcc) for that logic!

</span>
	</div>
	</div>

When `TicketBookingCanceled` arrives before `TicketBookingConfirmed` (in other words, when no ticket was stored in `tickets` yet), 
we will return an error. `TicketBookingCanceled` will be redelivered after a while, when the ticket should already exist.
Upon redelivery of `TicketBookingCanceled`, we will just ignore the update.
In other words, our repository will be idempotent and resilient to out-of-order events.

## Exercise

File: `project/main.go`

Add a `deleted_at` column to the `tickets` table and set it to a current timestamp when the `TicketBookingCanceled` event is processed.

**You also need to update your querying logic, so tickets with `deleted_at` set are not returned.**


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

To verify your solution, we will simulate sending `POST /tickets-status` with an out-of-order status update.
You will receive information about cancellation of a ticket that was not booked yet.
You should return an error in your event handler, so the event will be redelivered after a while, 
when `TicketBookingConfirmed` has been processed.

Please ensure that [your retry middleware](/trainings/go-event-driven/exercise/4620c713-50f9-4926-a602-7df110944cd0) doesn't redeliver messages too slowly.
If the message will be redelivered after too long time, it may be already after test timeout.

In the end, canceled tickets should not be returned from your `GET /tickets` endpoint.

</span>
	</div>
	</div>
