# Project: Context logger

Since you now have the correlation ID propagated, it would be helpful to include it in all log messages.
To make it easier to use in any place in code, you can create a logger with a correlation ID field
and keep it in the request's context.

The common library provides the `log.ToContext` function that does this.

```go
ctx := log.ToContext(ctx, logrus.WithFields(logrus.Fields{"correlation_id": correlationID}))
```

There's also `log.FromContext`, which retrieves the logger:

```go
logger := log.FromContext(msg.Context())
logger.WithField("key", "value").Info("Log message")
```

## Exercise

File: `project/main.go`

Modify the correlation ID middleware you created before to store the logger in the context.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you prefer, you can also create a separate middleware function that gets the correlation ID out of the context (with `log.CorrelationIDFromContext`)
and stores the logger.

</span>
	</div>
	</div>

Modify the logger middleware to use the logger from the context to log messages.
This way, each `Handling a message` log should automatically include a `correlation_id` field.

Be careful about the order of the middleware functions!
You need the middleware storing the correlation ID and the logger to be before the logging middleware.
