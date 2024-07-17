# Handling re-delivery in your project

## Exercise

File: `project/main.go`

There is a chance that you already implemented the storing of tickets in the `tickets` table in an idempotent way.
If you didn't, this is a good time to do that.

To support deduplication for adding tickets, you need to add `ON CONFLICT DO NOTHING` at the end of the `INSERT` query.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you defined the primary key properly in your database schema, your API will return just one ticket for the same `ticket_id`.
But without `ON CONFLICT DO NOTHING`, the message won't be acknowledged, and it will be redelivered forever.

You should not allow that to happen because it will make your system much harder to debug.
We will cover this in more depth in the observability module.

</span>
	</div>
	</div>
