# Printing Tickets - emit event


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

One team in our company wants to integrate with the printing system to automate the printing of tickets.
They want to integrate with our system by subscribing to the `TicketPrinted` event.

They need information about the ticket ID and file name.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Use the Event Bus to emit a `TicketPrinted` event after the ticket is printed.
You need to inject the Event Bus to your handler.

The emitted event should have the following format:

```go
type TicketPrinted struct {
	Header EventHeader `json:"header"`

	TicketID string `json:"ticket_id"`
	FileName string `json:"file_name"`
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

If you feel tempted to add the entire ticket model to the event, don't do it by default!

Remember that events become a contract between systems.
If you add an entire ticket model to the event, you will need to always keep adding this data to the event.

It's especially painful if you are refactoring in the future, and you want to split services or modules to smaller ones.
You may no longer have access to all the data that you emitted in the event in the past.

As an alternative, you can deprecate the old event and introduce a new one. 
However, it's always painful (as it may require a cross-team initiative).

[YAGNI!](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)

</span>
	</div>
	</div>
