# Migrating the Read Model

We now have all the building blocks required to perform the migration of our read model.
It's time to put them together!
The migration is pretty simple once we have all the parts ready.

Usually, to migrate the read model, you need to follow these steps:

1. Query events from the data lake one by one, from oldest to newest.
2. If needed, do a mapping of versions (your read model may be built from newer versions of events, 
while in the data lake, you may have older versions).
3. Call your read model methods to build it. 
   Usually, it will be some form of a repository similar to what you implemented in [13-read-models/01-read-models](/trainings/go-event-driven/exercise/cc7047b9-4d4b-413e-abbc-cbe29a8cba9d).

In the real world, if you create a new version of a read model, you should probably write it to a new table.

You should call your read model repository methods directly, not via message handlers.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Usually, if the migration is long-running, you may want to do the migration in the background and have some resume mechanism. 
For example, you can store the last timestamp of the event that you processed and
start from that timestamp when you resume the migration.

In our case, it will be not needed.

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

What if you needed to build a new read model but didn't have a data lake?
You can always do some reverse-engineering and build a read model based on your [write model](/trainings/go-event-driven/exercise/cc7047b9-4d4b-413e-abbc-cbe29a8cba9d).

It will be harder to do than from events, but it's possible.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Please migrate your Ops read model (`read_model_ops_bookings`) based on events from the data lake (from the `events` table).

**You should do the migration after initializing the `events` table.
After the `events` table is created, we will insert some events there.**
We will insert all events at once — you can safely assume that if the `events` table is not empty, 
you can migrate your read model.

**Please do not start migration before the `events` table is populated by us.
You can run the migration in a goroutine in a function that runs your service or in the `main` function.**

Please ensure that the `events` table has the same schema that we used in [15-data-lake/03-project-store-events-to-data-lake](/trainings/go-event-driven/exercise/1d2da948-3b82-46bb-8cbe-bec57433798e).

We don't know the exact format of events that you are using, so we will publish events with version `v0`.

The events that we will add to the `events` table are:
- `BookingMade_v0`
- `TicketBookingConfirmed_v0`
- `TicketReceiptIssued_v0`
- `TicketRefunded_v0`

You should map all of these to your event's format, which will ensure that you can call your normal read model methods.
This is a common scenario when you are migrating a read model and doing a migration with pretty old data.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

In the real world, you may not need to build a read model from the oldest data.
You may decide on some cut-off date and build the read model only from that date.
It depends on the use case, but if you are going to build it from the newest data, mapping may be not needed.

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

We recommend adding a good amount of logs to your migration.
It's not unusual that such migrations take longer than expected or something goes wrong.
It's worth logging the progress, how many events were processed, how much time it took, etc.

</span>
	</div>
	</div>
