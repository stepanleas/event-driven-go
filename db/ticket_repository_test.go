package db

import (
	"context"
	"os"
	"sync"
	"testing"
	"tickets/entities"

	_ "github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conn *sqlx.DB
var getDbOnce sync.Once

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		conn, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	})
	return conn
}

func TestTicketRepository_Add_idempotency(t *testing.T) {
	ctx := context.Background()

	dbConn := getDb()
	err := InitializeDatabaseSchema(dbConn)
	assert.NoError(t, err)

	repo := NewTicketRepository(dbConn)

	ticketToAdd := entities.Ticket{
		TicketID: uuid.NewString(),
		Price: entities.Money{
			Amount:   "30.00",
			Currency: "EUR",
		},
		CustomerEmail: "foo@bar.com",
	}

	for i := 0; i < 2; i++ {
		err = repo.Add(ctx, ticketToAdd)
		require.NoError(t, err)

		// probably it would be good to have a method to get ticket by ID
		tickets, err := repo.FindAll(ctx)
		require.NoError(t, err)

		foundTickets := lo.Filter(tickets, func(t entities.Ticket, _ int) bool {
			return t.TicketID == ticketToAdd.TicketID
		})
		// add should be idempotent, so the method should always return 1
		require.Len(t, foundTickets, 1)
	}
}
