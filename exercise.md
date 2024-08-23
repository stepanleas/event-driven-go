# Adding a Process Manager to the Project

The core of the process manager is ready. It's time to connect it to the rest of our project.

In this exercise, you need to:

- Add an HTTP endpoint that generates the VipBundle ID and booking ID, stores `VipBundle` in the database, and emits `VipBundleInitialized_v1`.
  This should be done within the transaction (event published with [outbox](/trainings/go-event-driven/exercise/6eeac65c-ff1b-4956-9523-c34e5ccc59b5)).
- Implement a repository for `VipBundle` that will store it in the database and emit `VipBundleInitialized_v1`.
- Implement the `BookShowTickets` command handler; it should call the existing logic used by the `POST /book-tickets` endpoint.
- Connect the process manager to the command processor and event processor.

**In this exercise, we are just checking the happy path.
We are also not booking taxis or flights yet.**


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

This exercise can benefit from using [Clean Architecture](https://threedots.tech/post/introducing-clean-architecture/).
You can use the same function triggered by the HTTP handler and command handler.

Clean Architecture is beyond the scope of this training,
but we recommend that you try to refactor the project to Clean Architecture on your own.

</span>
	</div>
	</div>

In this exercise, we will use a lof of things that we learned earlier:

- [outbox](/trainings/go-event-driven/exercise/6eeac65c-ff1b-4956-9523-c34e5ccc59b5): We need to publish `VipBundleInitialized_v1` in a transaction.
- [commands](/trainings/go-event-driven/exercise/fddf318e-b851-4f75-8972-d3c5deea44a5): For booking show tickets, and then later, flights and taxis.
- [tracing](/trainings/go-event-driven/exercise/c61c6d59-84fa-46bc-bf87-f87fbcd4205d): Will help us debug and understand what is going on.
- [at-least-once-delivery and deduplication](/trainings/go-event-driven/exercise/8c31d18a-b5ae-4d6a-9d1b-a057be5e4b2c): All our command and event handlers within the process manager should be idempotent (we don't want to book flight tickets twice — they are expensive).

### HTTP endpoint specification

Vip Bundles can be booked by calling the `POST /book-vip-bundle` endpoint with the following body:

```go
type vipBundleRequest struct {
	CustomerEmail   string    `json:"customer_email"`
	InboundFlightId uuid.UUID `json:"inbound_flight_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	Passengers      []string  `json:"passengers"`
	ReturnFlightId  uuid.UUID `json:"return_flight_id"`
	ShowId          uuid.UUID `json:"show_id"`
}
```

This will result in this response:

```go
type vipBundleResponse struct {
	BookingId   uuid.UUID `json:"booking_id"`
	VipBundleId uuid.UUID `json:"vip_bundle_id"`
}
```

The endpoint should return `StatusCreated (201)` when the process manager has successfully started.

### De-duplicating the booking show ticket

Previously, tickets were bookied just by HTTP endpoint, so we didn't need to deduplicate them.

Now this will be also called by a command, which may be [redelivered](/trainings/go-event-driven/exercise/8c31d18a-b5ae-4d6a-9d1b-a057be5e4b2c).

You can use this helper to detect when the booking has already been done:

```go
const (
	postgresUniqueValueViolationErrorCode = "23505"
)

func isErrorUniqueViolation(err error) bool {
	var psqlErr *pq.Error
	return errors.As(err, &psqlErr) && psqlErr.Code == postgresUniqueValueViolationErrorCode
}
```

If this happens in your command handler, it should return the `nil` error (so it will acknowledge the command).
The HTTP handler should work as is — it should return an error.

### `VipBundleRepository`

It's up to you whether you want to keep the `VipBundleRepository` interface.
You can redesign it if you want.

The interface was inspired by our [article on the repository pattern](https://threedots.tech/post/repository-pattern-in-go/).

### Alerting

We won't check it in this exercise, but in a production system, you may want to add alerting for stuck process managers.

You can periodically query process manager instances in the database and check which are not completed within the provided threshold.
It should also be possible to check stuck process managers out of events in the data lake.

### Testing

We won't check this, but we recommend adding [component tests](/trainings/go-event-driven/exercise/502648e4-ff0d-41d8-82c1-6ce7f4479774) for this endpoint and for the process manager.

## Exercise

File: `project/main.go`

Don't worry if this doesn't work the first time.
Even when we were implementing this exercise, we made some mistakes, and it took us some time to make it work.
As we said earlier, process managers are complex things to implement.
If something doesn't work, please check the logs and try to understand what is going on.

We are implementing this incrementally, which is the approach that we recommend in real life.

At the end, our trace should look like this:

![Process manager trace](https://academy.threedots.tech/media/trainings/go-event-driven/trace-process-manager.png)

_(In this exercise, we are just booking show tickets; we'll take care of taxis and flights in the following exercises)_.


<div class="accordion" id="hints-accordion">

<div class="accordion-item">
	<h3 class="accordion-header" id="hints-accordion-header-1">
	<button class="accordion-button fs-4 fw-semibold collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#hints-accordion-body-1" aria-expanded="false" aria-controls="hints-accordion">
		Hint #1
	</button>
	</h3>
	<div id="hints-accordion-body-1" class="accordion-collapse collapse" aria-labelledby="hints-accordion-header-1" data-bs-parent="#hints-accordion">
	<div class="accordion-body">

Do you need inspiration for how to implement `VipBundleRepository`?

You can use this schema:

```sql
CREATE TABLE IF NOT EXISTS vip_bundles (
	vip_bundle_id UUID PRIMARY KEY,
	booking_id UUID NOT NULL UNIQUE,
	payload JSONB NOT NULL
); 
```

And use this code:

```go
type VipBundleRepository struct {
	db *sqlx.DB
}

func NewVipBundleRepository(db *sqlx.DB) *VipBundleRepository {
	if db == nil {
		panic("db must be set")
	}

	return &VipBundleRepository{db: db}
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (v VipBundleRepository) Add(ctx context.Context, vipBundle entities.VipBundle) error {
	payload, err := json.Marshal(vipBundle)
	if err != nil {
		return fmt.Errorf("could not marshal vip bundle: %w", err)
	}

	return updateInTx(
		ctx,
		v.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err = v.db.ExecContext(ctx, `
				INSERT INTO vip_bundles (vip_bundle_id, booking_id, payload)
				VALUES ($1, $2, $3)
			`, vipBundle.VipBundleID, vipBundle.BookingID, payload)

			if err != nil {
				return fmt.Errorf("could not insert vip bundle: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = event.NewBus(outboxPublisher).Publish(ctx, entities.VipBundleInitialized_v1{
				Header:      entities.NewEventHeader(),
				VipBundleID: vipBundle.VipBundleID,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}

func (v VipBundleRepository) Get(ctx context.Context, vipBundleID uuid.UUID) (entities.VipBundle, error) {
	return v.vipBundleByID(ctx, vipBundleID, v.db)
}

func (v VipBundleRepository) vipBundleByID(ctx context.Context, vipBundleID uuid.UUID, db Executor) (entities.VipBundle, error) {
	var payload []byte
	err := v.db.QueryRowContext(ctx, `
		SELECT payload FROM vip_bundles WHERE vip_bundle_id = $1
	`, vipBundleID).Scan(&payload)

	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not get vip bundle: %w", err)
	}

	var vipBundle entities.VipBundle
	err = json.Unmarshal(payload, &vipBundle)
	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not unmarshal vip bundle: %w", err)
	}

	return vipBundle, nil
}

func (v VipBundleRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (entities.VipBundle, error) {
	return v.getByBookingID(ctx, bookingID, v.db)
}

func (v VipBundleRepository) getByBookingID(ctx context.Context, bookingID uuid.UUID, db Executor) (entities.VipBundle, error) {
	var payload []byte
	err := db.QueryRowContext(ctx, `
		SELECT payload FROM vip_bundles WHERE booking_id = $1
	`, bookingID).Scan(&payload)

	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not get vip bundle: %w", err)
	}

	var vipBundle entities.VipBundle
	err = json.Unmarshal(payload, &vipBundle)
	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not unmarshal vip bundle: %w", err)
	}

	return vipBundle, nil
}

func (v VipBundleRepository) UpdateByID(ctx context.Context, bookingID uuid.UUID, updateFn func(vipBundle entities.VipBundle) (entities.VipBundle, error)) (entities.VipBundle, error) {
	var vb entities.VipBundle

	err := updateInTx(ctx, v.db, sql.LevelSerializable, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		vb, err = v.vipBundleByID(ctx, bookingID, tx)
		if err != nil {
			return err
		}

		vb, err = updateFn(vb)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(vb)
		if err != nil {
			return fmt.Errorf("could not marshal vip bundle: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE vip_bundles SET payload = $1 WHERE vip_bundle_id = $2
		`, payload, vb.VipBundleID)

		if err != nil {
			return fmt.Errorf("could not update vip bundle: %w", err)
		}

		return nil
	})
	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not update vip bundle: %w", err)
	}

	return vb, nil
}

func (v VipBundleRepository) UpdateByBookingID(ctx context.Context, bookingID uuid.UUID, updateFn func(vipBundle entities.VipBundle) (entities.VipBundle, error)) (entities.VipBundle, error) {
	var vb entities.VipBundle

	err := updateInTx(ctx, v.db, sql.LevelSerializable, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		vb, err = v.getByBookingID(ctx, bookingID, tx)
		if err != nil {
			return err
		}

		vb, err = updateFn(vb)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(vb)
		if err != nil {
			return fmt.Errorf("could not marshal vip bundle: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE vip_bundles SET payload = $1 WHERE booking_id = $2
		`, payload, vb.BookingID)

		if err != nil {
			return fmt.Errorf("could not update vip bundle: %w", err)
		}

		return nil
	})
	if err != nil {
		return entities.VipBundle{}, fmt.Errorf("could not update vip bundle: %w", err)
	}

	return vb, nil
}
```

</div>
	</div>
	</div>


</div>