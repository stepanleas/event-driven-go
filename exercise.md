# Rename topics

Currently, topic names are equal to command/event names.

Even if the chance of conflict is not high, it would be good to separate them.
This will also make more explicit whether the topic is a command or event topic.

## Exercise

File: `project/main.go`

Change command and event topic names to the formats `events.{event_name}` and `commands.{command_name}`.

We will check your solution by checking if `TicketBookingConfirmed` was emitted on `events.TicketBookingConfirmed`
and `RefundTicket` on `commands.RefundTicket`.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

When you make such a change in production, you should be sure that you don't have any leftover messages on old topics.
If you want to ensure that this doesn't happen, you can make the change in three steps:

1. Add handlers for new topics while keeping topics for old names and then deploy.
2. Change the command and event bus to publish to new topics and deploy.
3. Remove the handlers for old topics and deploy.

Thanks to this, you will be sure that you don't have any leftover messages on old topics.

</span>
	</div>
	</div>
