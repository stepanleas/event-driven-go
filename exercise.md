# Project: Running the Service in Tests

You now understand how to mock external dependencies. Next, we'll focus on running our service during testing.

We'll create our component tests in the `tests/component_test.go` file.
(Note: It's crucial to use this directory to ensure that your solution is verifiable.)


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Unfortunately, we lack the technical capability to determine whether your tests assert the correct things.
However, we will confirm that your tests pass — we trust you! :-)

</span>
	</div>
	</div>

In the provided example, we created a struct `Service` in the `service` package.
This struct's constructor accepts all dependencies and includes a `Run` method that starts the application.

Note that `Service` is not defined in the `main` package, as doing so would prevent us from importing it during testing.

We only pass the dependencies that differ between production and component testing to the `Service` constructor, not all dependencies.

```go
package service


import (
    "context"
    "fmt"
    stdHTTP "net/http"
    "tickets/db"
    ticketsHttp "tickets/http"
    "tickets/message"
    "tickets/message/event"
    "tickets/message/outbox"
    
    "github.com/ThreeDotsLabs/go-event-driven/common/log"
    "github.com/ThreeDotsLabs/watermill/components/cqrs"
    watermillMessage "github.com/ThreeDotsLabs/watermill/message"
    "github.com/jmoiron/sqlx"
    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
    "github.com/redis/go-redis/v9"
    "github.com/sirupsen/logrus"
    "golang.org/x/sync/errgroup"
)

type Service struct {
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

func New(
	redisClient *redis.Client,
	spreadsheetsService event.SpreadsheetsAPI,
	receiptsService event.ReceiptsService,
) Service {
	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	
	watermillRouter := message.NewWatermillRouter(
		receiptsService,
		spreadsheetsService,
		redisClient,
		watermillLogger,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		redisPublisher,
		spreadsheetsService,
	)

	return Service{
		watermillRouter,
		echoRouter,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && err != stdHTTP.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
```

If you're unsure about how the service should look, refer to the solution for [the previous exercise](/trainings/go-event-driven/exercise/a28838fa-2f5b-4e65-8b8e-6b5c40888ca5).


### golang.org/x/sync/errgroup

Do you remember how to use `golang.org/x/sync/errgroup`? If not, check [the previous exercise](/trainings/go-event-driven/exercise/b8a0b000-9f74-4b86-a399-44ee342dd374).

## Writing Component Tests for Your Project

Next, let's write super simple component tests for our service.
**Don't be scared by the wall of text in this exercise — we'll guide you through it step by step.**

First of all, you need to run the service during the tests.
You just need to prepare the dependencies and mocks and then run the service in a goroutine.

```go
package tests_test

import ( 
	// ... 
	"tickets/api"
	"tickets/entities"
	"tickets/message"
	"tickets/service" 
	// ...
)

func TestComponent(t *testing.T) {
	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spreadsheetsService := &api.SpreadsheetsMock{}
	receiptsService := &api.ReceiptsMock{}

	go func() {
		svc := service.New(
			redisClient,
			spreadsheetsService,
			receiptsService,
		)
		assert.NoError(t, svc.Run(ctx))
	}()

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

Our examples use the `github.com/stretchr/testify` library.
It's not mandatory, but we highly recommend it.

Bear in mind that while some companies (like Google) don't advocate the use of assertion libraries, and some people follow their lead blindly, 
Google doesn't hold a monopoly on good practices — every company is different.
If you are not working at Google, you should find the good practices that work for you and your team.

</span>
	</div>
	</div>

You'll need to port the receipts service mock from the previous exercise to the project.
Also, create the `SpreadsheetService` and mock that to implement the following interface.

```go
type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}
```

For now, these can be dummy implementations: They can just return `nil` or empty values.
You'll extend them later.

The service starts in the goroutine, and you need to ensure that it's healthy.
Your service should already have a `/health` endpoint implemented and should expose an HTTP server on port 8080.

We can check the service's health with a helper function:

```go
func TestComponent(t *testing.T) {
	// ...

	go func() {
		svc := service.New(
			redisClient,
			spreadsheetsService,
			receiptsService,
		)
		assert.NoError(t, svc.Run(ctx))
	}()

	waitForHttpServer(t)

    // ...

    func waitForHttpServer(t *testing.T) {
        t.Helper()
    
        require.EventuallyWithT(
            t,
            func(t *assert.CollectT) {
                resp, err := http.Get("http://localhost:8080/health")
                if !assert.NoError(t, err) {
                    return
                }
                defer resp.Body.Close()
    
                if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
                    return
                }
            },
            time.Second*10,
            time.Millisecond*50,
        )
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

Avoid using `require.XXX` in goroutines, as `t.FailNow()` doesn't function properly in these contexts.

You can learn more [here](https://github.com/stretchr/testify/issues/772#issuecomment-945166599).

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

Avoid using asserts within `assert.Eventually`, as it won't work as expected.

```go
assert.Eventually(
      t,
      func() bool {
          err := someLogic()
          return assert.NoError(t, err)
      },
      10*time.Second,
      100*time.Millisecond,
  )
```

For example, `assert.NoError` won't fail the test, even if `Eventually` eventually succeeds.

The `EventuallyWithT` function was added in `v1.8.3` of `github.com/stretchr/testify`.
This function can "buffer" asserts and only fail the test if `Eventually` fails.

</span>
	</div>
	</div>

### Running Locally

When you run your solution with `tdl training run`, we'll set up all the necessary containers.

If you prefer to run locally, you can use `docker-compose`.
You should already have a `docker-compose.yml` file in your project.

To run tests locally, you'll need to:

- Run `docker-compose up --pull` in one terminal.
- After everything is up and running, run: 

```bash
REDIS_ADDR=localhost:6379 go test ./tests/ -v
```

## Exercise

File: `project/main.go`

For this exercise, run your service in the `TestComponent` function in the `tests/component_test.go` file.
To ensure that the service is running, you can use the provided `waitForHttpServer` function.

If you're uncertain about how to prepare the service for running in tests, refer to [the solution from the previous exercise](/trainings/go-event-driven/exercise/a28838fa-2f5b-4e65-8b8e-6b5c40888ca5).
