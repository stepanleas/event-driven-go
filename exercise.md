# Tracing in your project

It's time to check what we just learned about tracing in our project!

You should apply everything you learned in this module to the project.
The bare minimum is adding tracing, which will be propagated via messages.

Using traces within a project adds some requirements.

## Exposing the trace ID

It's good to be able to know the trace ID of the request.
It's good enough to just log it somewhere, so you will be able to correlate logs with traces.
You can also add it to each log as a field (but this is not recommended locally — it will pollute your logs a lot).
You can extract the trace ID from the context by calling [`trace.SpanContextFromContext(ctx).TraceID().String()`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#SpanContextFromContext).

You can also use a trick here: Use the trace ID as a correlation ID.
This will let you simplify your service logic and remove the propagation of the correlation ID.
Of course, to do that, you need to adjust your logging logic.

## Configuring the tracer

In earlier exercises, we preconfigured a tracer for you. 
It's now your turn to configure it.

We prepared something to inspire you:

```go
func ConfigureTraceProvider() *tracesdk.TracerProvider {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = fmt.Sprintf("%s/jaeger-api/api/traces", os.Getenv("GATEWAY_ADDR"))
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerEndpoint),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		// WARNING: `tracesdk.WithSyncer` should be not used in production.
		// For production, you should use `tracesdk.WithBatcher`.
		tracesdk.WithSyncer(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("tickets"),
		)),
	)

	otel.SetTracerProvider(tp)
	
	// Don't forget this line! Omitting it will cause the trace to not be propagated via messages.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
```

You should call the function during startup of the service.
This is a configuration that is good for a local development; traces are not batched, so they will be instantly visible in the Jaeger UI.

In the environment we used to verify your solutions, we will route traces via the gateway. It's also possible to explicitly provide a Jaeger endpoint via an environment variable.

It's also good to implement graceful shutdown of the tracer:

```go
errgrp.Go(func() error {
    <-ctx.Done()
    return s.traceProvider.Shutdown(context.Background())	
})
```

Thanks to this, you will be sure that all traces are exported before the service is shut down when using `tracesdk.WithBatcher`.

We will use Jaeger to collect traces.

### Jaeger

Jaeger is an open-source solution for collecting traces.
It's currently an industry standard for open-source tracing.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you are deploying your services to any cloud and you don't want to self-host Jaeger,
you may also check the tracing solution provided by the cloud provider.
For example, AWS has X-Ray, GCP has Stackdriver, and Azure has Application Insights.

For official exporters, you should check out https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters.

We didn't find a list of other supported exporters, but googling _"opentelemetry go trace exporter [technology]"_ should give you some results.
Some proprietary solutions also use the [Open Telemetry Protocol](https://opentelemetry.io/docs/specs/otel/protocol/),
so you can use them with an Open Telemetry exporter.

</span>
	</div>
	</div>

The screenshots that you have seen earlier were from the Jaeger UI.

![Tickets trace](https://academy.threedots.tech/media/trainings/go-event-driven/tickets-trace.png)

This is the trace that we should achieve in the project after this exercise.

We have added Jaeger to your docker-compose file.

## Adding HTTP requests to traces

We will not check your solution for this, but we recommend adding HTTP requests to traces.
It will add a lot of useful context to your traces.

You can achieve this with the [`go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`](go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp) library.

You need to use a custom HTTP client for your API clients:

```diff
 import (
 	"context"
+	"fmt"
 	"net/http"
 	"os"
 	"os/signal"
 	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
 	"github.com/ThreeDotsLabs/go-event-driven/common/log"
 	"github.com/jmoiron/sqlx"
+	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
+	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
 )
 
 func main() {
 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
 	defer cancel()
 
-	apiClients, err := clients.NewClients(
+	traceHttpClient := &http.Client{Transport: otelhttp.NewTransport(
+		http.DefaultTransport,
+		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
+			return fmt.Sprintf("HTTP %s %s %s", r.Method, r.URL.String(), operation)
+		}),
+	)}
+
+	apiClients, err := clients.NewClientsWithHttpClient(
 		os.Getenv("GATEWAY_ADDR"),
 		func(ctx context.Context, req *http.Request) error {
 			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
 			return nil
 		},
+		traceHttpClient,
 	)
```

This will automagically add all outgoing HTTP requests to traces.
It will also add the `traceparent` header to all outgoing HTTP requests, but we are calling external APIs, so it won't be used.

However, if you would like to call your own services, and you add middleware that propagates the `traceparent` header,
you will be able to correlate traces between services.

## Adding SQL queries to traces

We will not check your solution for this as well, but it's also nice to have SQL queries in traces, 
especially when debugging some performance issues.

To do that, you need to wrap your SQL connection with [`github.com/uptrace/opentelemetry-go-extra/otelsql`](github.com/uptrace/opentelemetry-go-extra/otelsql).

```diff
 import (
 	"context"
 	"net/http"
 	"os"
 	"os/signal"
@@ -14,24 +15,43 @@ import (
 	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
 	"github.com/ThreeDotsLabs/go-event-driven/common/log"
 	"github.com/jmoiron/sqlx"
+	"github.com/uptrace/opentelemetry-go-extra/otelsql"
 )
 

-	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
+	traceDB, err := otelsql.Open("postgres", os.Getenv("POSTGRES_URL"),
+		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
+		otelsql.WithDBName("db"))
+	if err != nil {
+		panic(err)
+	}
+
+	db := sqlx.NewDb(traceDB, "postgres")
 	if err != nil {
 		panic(err)
 	}
```

As long as you pass context to all your SQL queries (`ExecContext`, `NamedExecContext`, etc.), they will be added to traces.

## Exercise

File: `project/main.go`

In this exercise, we will check if all events have a trace ID and span ID.
We won't check if all message traces are properly correlated, but at least some of them should be properly correlated.

To simplify debugging, we recommend that you run your service locally and debug traces there. 
This is not mandatory — we will try to guide you through the exercise without it — but it would be good to see the traces in the UI and play with Jaeger.

To run Jaeger and your service dependencies, you should run `docker-compose up --pull`.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you don't have Docker installed on your machine, you can download it [here](https://www.docker.com/products/docker-desktop/).

</span>
	</div>
	</div>

If you haven't been running your service locally yet, it's a good time to do that — it's a very important part of the development process.

After starting docker-compose, you can try to visit [http://localhost:16686/](http://localhost:16686/) to check if Jaeger is running.

After running docker-compose in one terminal, you can now run your application:

```bash
REDIS_ADDR=localhost:6379 GATEWAY_ADDR=http://localhost:8888 POSTGRES_URL=postgres://user:password@localhost:5432/db?sslmode=disable JAEGER_ENDPOINT=http://localhost:14268/api/traces go run .
```

After that, you can try to send some requests to your service. You need start by creating a show:

```bash
curl -v -X POST http://localhost:8080/shows \
-H "Content-Type: application/json" \
-d '{
  "dead_nation_id": "0fe9f3bf-160f-49be-9509-862e91ee8c33",
  "number_of_tickets": 10000,
  "start_time": "2024-02-04T19:00:00Z",
  "title": "Metallic Nëcrømancer: Raging Terror 2024 Tour",
  "venue": "Colosseum"
}'
```

Then you can try to book some tickets:

```bash
curl -v -X POST http://localhost:8080/book-tickets \
-H "Content-Type: application/json" \
-d '{
  "customer_email": "email@example.com",
  "number_of_tickets": 1,
  "show_id": "<SHOW ID FROM POST /shows>"
}'
```

If you configured everything correctly, your traces should be visible in the Jaeger UI.

**Please remember that running an exercise with `tdl tr run` will cause it to not be executed in your environment, so traces won't be exported.**

Nobody configures tracing that works right away — even we haven't.
It may require a couple of iterations to connect everything up properly.

You should check if context was passed to all functions, all publishers are decorated, all headers are set, etc.
It may be worth adding logging in some places, so you can have a better understanding where traces are lost.
