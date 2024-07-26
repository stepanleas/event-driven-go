# Store events in Data Lake

Now, when we have all of our events on a single topic, we can store them in the data lake.
For the purposes of the training, we will use PostgreSQL as our data lake.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Just a reminder: PostgreSQL is usually not a good choice for a data lake at scale.
If your dataset grows huge, PostgreSQL may turn out too expensive to store all the events.

</span>
	</div>
	</div>

You can think about a data lake as being like a big [read model](/trainings/go-event-driven/exercise/cc7047b9-4d4b-413e-abbc-cbe29a8cba9d) containing all events in raw form.

## Exercise

File: `project/main.go`

Please add a [message handler](/trainings/go-event-driven/exercise/6e0ddff2-aaf9-4188-aeea-9fc8eb9ac6ba) that will store all events in the data lake.
An event handler won't work here because we need to store our events in raw form.
It should listen to the `events` topic and store all events in the data lake.

**It's important to use exactly this schema for your `events` table:**

```sql
CREATE TABLE IF NOT EXISTS events (
    event_id UUID PRIMARY KEY,
    published_at TIMESTAMP NOT NULL,
    event_name VARCHAR(255) NOT NULL,
    event_payload JSONB NOT NULL
);
```

We will later depend on this exact schema.

Please don't forget about [at-least-once delivery](/trainings/go-event-driven/exercise/7c4d2754-3fec-44f9-9d46-63be65a76468)!
You should deduplicate potential redelivered events.


<div class="accordion" id="hints-accordion">

<div class="accordion-item">
	<h3 class="accordion-header" id="hints-accordion-header-1">
	<button class="accordion-button fs-4 fw-semibold collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#hints-accordion-body-1" aria-expanded="false" aria-controls="hints-accordion">
		Hint #1
	</button>
	</h3>
	<div id="hints-accordion-body-1" class="accordion-collapse collapse" aria-labelledby="hints-accordion-header-1" data-bs-parent="#hints-accordion">
	<div class="accordion-body">

To store events in a data lake, you need to extract the header from the event.
You can do this by unmarshaling the event to a struct that just has a header field:

```go
// We just need to unmarshal the event header; the rest is stored as is.
type Event struct {
	Header entities.EventHeader `json:"header"`
}

var event Event
if err := eventProcessorConfig.Marshaler.Unmarshal(msg, &event); err != nil {
	return fmt.Errorf("cannot unmarshal event: %w", err)
}
```

</div>
	</div>
	</div>

<div class="accordion-item">
	<h3 class="accordion-header" id="hints-accordion-header-2">
	<button class="accordion-button fs-4 fw-semibold collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#hints-accordion-body-2" aria-expanded="false" aria-controls="hints-accordion">
		Hint #2
	</button>
	</h3>
	<div id="hints-accordion-body-2" class="accordion-collapse collapse" aria-labelledby="hints-accordion-header-2" data-bs-parent="#hints-accordion">
	<div class="accordion-body">

As in previous exercises, you can extract the event name from a message by using the CQRS marshaler:

```go
eventName := eventProcessorConfig.Marshaler.NameFromMessage(msg)
if eventName == "" {
	return fmt.Errorf("cannot get event name from message")
}
```

</div>
	</div>
	</div>

</div>
