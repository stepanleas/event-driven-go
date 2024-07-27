# Event Versioning

Emitting an event to a Pub/Sub, which anyone from the company can access, becomes a contract between you and other teams.
When the event is added to the data lake, it's even more important to not change it.
On our blog, we mentioned multiple times that it's important to build a system that is open to changes.
What should we do if we need to change the event format?

## Backward-compatible changes

The best strategy is to add new fields without removing the old ones.
In our experience, this is possible in most cases.
The tradeoff is that you need to keep the old fields forever (and keep the payload bigger than it could be).
However, this may be a good tradeoff.

## Non-backward-compatible changes

It would be wonderful if all changes were backward compatible,
but in the real world, that is not always possible.

Sometimes, adding a new field is not an option.
There may be multiple reasons for this:

- The event is too big, and adding a new field will make it even bigger.
- The event changed, so it was emitted in a different situation than before (due to domain logic changes).
- Some fields are no longer available.

## Migrating

1. Add a new event with a new name and new format while keeping the old one.
2. Emit new events.
3. Try to update all consumers to use the new events.
4. Remove the old events.

Sometimes, point number 3 never happens, which is fine.
It may sometimes be impossible to find all the use cases of the event.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

The process of event version migration is not simple, and sometimes it is not worth the effort.

With an API, it's easier to measure who is using an endpoint and how often.
Events are usually used in a more "indirect" way, so it's harder to find who is using them and with what frequency.

This shows how important it is to spend enough time on designing events: You may have no chance to change them later.

</span>
	</div>
	</div>

Even if we should avoid making non-backward-compatible changes, it's good to keep the door open for that.
To support this possibility, we will add versions to all of our events.

The simplest strategy to do that is to add a suffix for each event type.
For example, `TicketBookingConfirmed` will become `TicketBookingConfirmed_v1`.
This is practical because `TicketBookingConfirmed_v1` and `TicketBookingConfirmed_v2`
are treated as totally separate events, so you don't risk accidentally using the wrong version of the event.
Version 2 may have much different meaning than version 1 and may be triggered at a much different moment.

If you have an already existing system, and you need to introduce versioning, you can just assume that 
events without a version number are version 1.
It's not a big problem in our project, so we can change all events to contain an explicit version number.


## Exercise

File: `project/main.go`

Please add a `v1` suffix to  all events in your application.
For example, `TicketBookingConfirmed` should become `TicketBookingConfirmed_v1`, etc.

You don't need to change your event bus / event processor configuration.

Events will be published to new topics. 

For example: `TicketBookingConfirmed_v1` should be published to `events.TicketBookingConfirmed_v1`
Event processors will also listen to new topics.

For now, let's just introduce v1 for all our events. We will add new event with version 2 soon.
