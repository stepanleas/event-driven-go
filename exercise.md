# Testing idempotency

Being sure that our handlers and repositories are idempotent is crucial to making our system work correctly.

What's the best strategy to follow?

### Testing idempotency at the repository level

You could test our handlers' idempotency at the component test level, but this would be pretty time/resource-consuming.

It's much better to test it at the repository level. This gives you the same guarantee, and it's much faster to execute and easier to write.
It also creates a faster feedback loop when you're developing the functionality.

### Duplicator middleware

There is also another way to test idempotency on the level of the entire application: [duplicator middleware](https://watermill.io/docs/middlewares/#duplicator).

It works in a very simple way: it processes all your messages twice.
You can enable this middleware for your tests and check if your application still works correctly.

You can even go one step further and enable it in real environments (even production) to be sure that your application is working correctly.
Keep in mind that this may be expensive for bigger systems in terms of infrastructure costs.

Repository-level tests are definitely the cheaper and easier way to debug!
This is what we want for our application.

## Exercise

File: `project/main.go`

It's time to test if your repository is idempotent.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

If you just implemented the logic to add tickets to the database, in your handler, it may be good to refactor it to the repository pattern.
This will make your code more testable and easier to maintain.

If you need inspiration, you can check out [our article about the repository pattern in Go](https://threedots.tech/post/repository-pattern-in-go/).

</span>
	</div>
	</div>

Unfortunately, we are not able to guess where your repository code/tests are located.
**Because of that, for sake of this exercise, you need to put your repository tests in the `./db/` directory.**

Your code is a blackbox for us, so we can't really assert that you tested everything. 
In the end, we will just check that your tests pass.

Most of the work that you need to do is to set up the tests.
After that, **you need to call your repository function twice and assert that the second call succeeded and didn't change anything.**
The ticket should be stored only once. 

### Tips

Here are a couple tips that may help you with test implementation.

#### Running locally

If you want, you can run your tests locally. In one terminal you need to run:

```bash
docker-compose up --pull
```

And in another one:

```bash
POSTGRES_URL=postgres://user:password@localhost:5432/db?sslmode=disable go test ./db/ -v
```

#### It's good to write tests so that they are not dependent on cleanup

The most reliable integration tests should not depend on cleanup.
The cleanup may fail because the tests are killed or for some other reason, so 
it's good to write tests that are not dependent on them.

Some time ago, we wrote an article that covers this: [4 principles of high-quality database integration tests](https://threedots.tech/post/database-integration-testing/).

#### Write a helper to get a db singleton

It's good to write a helper to get a db singleton.
If you have multiple tests, you can be sure that you're not opening the connection for each test.

Feel free to use this one:

```go
var db *sqlx.DB
var getDbOnce sync.Once

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	})
	return db
}
```

#### Don't forget to call function that does schema initialization

You need your schema for tests to work properly. You can just call the needed function before each test or use [`TestMain`](https://pkg.go.dev/testing#hdr-Main).

