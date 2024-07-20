# Emitting BookingMade event

You should already know how to [publish an event in a transaction](/trainings/go-event-driven/exercise/5b87e341-31cc-4b44-bc8b-5da40c8be069)
and [use a forwarder](/trainings/go-event-driven/exercise/6494e13f-eb2e-4aa7-986f-001faf3afd2e).

It's time to use this knowledge in the project!

## Exercise

File: `project/main.go`

Publish the `BookingMade` event during a call to the `POST /book-tickets` endpoint.
It should be emitted in the same transaction as the booking is stored in the database.

The event should be emitted by the forwarder to the Redis Pub/Sub topic `BookingMade` (like other events).
In the payload, you should include the following:
- Booking ID
- Number of tickets
- Customer email
- Show ID

Use the `"events_to_forward"` topic for the forwarder and the default Postgres adapter for SQL Pub/Sub.
Thanks to that, messages will be stored to the `watermill_events_to_forward` table.
**We will use this table to verify your solution.**


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

You should use the publisher wrapped with the forwarder to create a new instance of the Event Bus per transaction.

Don't forget about passing the context — thanks to that you will not lose the correlation ID.

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

Many things need to be configured properly to make this work.

And if something doesn't work, logs are your friends!
If you still don't see what's wrong, please post your solution on Discord. 
Click the "Share your solution" button, choose the most recent one, and just copy the link to the training Discord channel.

</span>
	</div>
	</div>
