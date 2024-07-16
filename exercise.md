# Save tickets in the database - initialize schema


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Our application already supports multiple scenarios related to the tickets, 
but we are not storing any information about tickets in our system.

That's a problem because one of the other teams asked us to provide them an endpoint that will provide them information about tickets.
They said that data can be _eventually consistent_ for them.

It's time to store tickets in the database!

</span>
	</div>
	</div>

We will use a PostgreSQL database to store tickets. 
A PostgreSQL container will be added to each execution of your application.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>
The containers are recreated after each execution.
You don't need to worry about data cleanups.
</span>
	</div>
	</div>

The database connection string will be available in the `POSTGRES_URL` environment variable.

```go
db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
```

In example solutions, we'll use [`github.com/jmoiron/sqlx`](https://github.com/jmoiron/sqlx) and [`github.com/lib/pq`](https://github.com/lib/pq). 
This is not obligatory — you can use any other libraries.

### Creating a schema

In this exercise, we want to just create a database schema. We want to store all tickets in the database.

Normally, we would use more advanced tools for handling SQL migrations.
For the sake of this training, however, it's good enough to just create a schema during the startup of the application.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you want to check out the tools that we recommend for migration, check our [22 Go Libraries you need to know article](https://threedots.tech/post/list-of-recommended-libraries/#migrations).

</span>
	</div>
	</div>

To work properly, our application should be able to handle the situation when the schema already exists.

You can create your own schema, or you can use the database schema prepared by us.

We will need the following data:

- Ticket ID
- Price amount and currency
- Customer email

```sql
CREATE TABLE IF NOT EXISTS tickets (
	ticket_id UUID PRIMARY KEY,
	price_amount NUMERIC(10, 2) NOT NULL,
	price_currency CHAR(3) NOT NULL,
	customer_email VARCHAR(255) NOT NULL
);
```

You can also write your own. **However, it's important that the table be called `tickets`.**

## Exercise

File: `project/main.go`

Connect to the database, and create a `tickets` table. We will use this table in the following exercises to store tickets.

If you want to use [`github.com/lib/pq`](https://github.com/lib/pq), you can connect to the database with:

```go
db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
if err != nil {
	panic(err)
}
defer db.Close()
```

Then call `db.Exec` to execute a query to create the schema.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Note that this function can't just be part of your `main` function.
If your component tests were executed on a fresh database, they would fail.
You should set up the schema inside the function that is called both by the `main` package and the component tests.

</span>
	</div>
	</div>

If you use `github.com/lib/pq`, don't forget to add the empty import:

```go
import ( 
	// ... 
	_ "github.com/lib/pq"
	// ...
)
```
