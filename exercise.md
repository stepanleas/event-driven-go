# Using CQRS in the project

You should now know all the building blocks of Watermill CQRS. 
Time to use them in the project!

## Exercise

File: `project/main.go`

Update your project to use Watermill CQRS.

Don't forget to add the Correlation ID decorator to the publisher! 
Without it, you will have a hard time debugging your application if something goes wrong.

You can use the decorator you implemented yourself, or use the `log.CorrelationPublisherDecorator` from 
[`github.com/ThreeDotsLabs/go-event-driven/common`](https://github.com/ThreeDotsLabs/go-event-driven/tree/main/common/log).

Here are some tips on how to do this:

* Publish messages using the EventBus, not the Publisher directly.
* Replace the Router handlers with an EventProcessor and EventHandlers.
* You should not do any JSON marshaling yourself. Just use the `JSONMarshaler` from Watermill for both EventBus and EventProcessor.
* Similarly, you don't need to create messages manually to send them. Just pass the event struct to EventBus's `Publish`.
* After the changes, you should have no Subscriber created manually. Instead, create them within the `SubscriberConstructor` in the EventProcessor.
* You can use the handler name as the consumer group.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

The Watermill CQRS handler depends on some metadata fields set in the message by the JSON marshaler.
Using just CQRS event handlers without the CQRS event bus may not work (even if the payload matches).

We recommend changing both handlers and the Event Bus to CQRS at the same time.

</span>
	</div>
	</div>

Don't forget to use the JSON marshaler with the custom `GenerateName` option:

```go
var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
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

If your service is not working as expected, you should double-check if everything is using the right topics.
Logs should be useful for debugging.
You may wish to increase the log level to `debug` or `trace` to see more.   

</span>
	</div>
	</div>

After doing that change, you should have much less boilerplate code in your project.
It will be also much easier to add new handlers later.
