# Add ticket limits

The final functionality that we need to implement before using our new endpoint in production is imposing limits on how many tickets we can sell.
We have this information in the `shows` table, so now we just need to enforce it.


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

If shows are more popular, our API may receive up to 20 concurrent requests to book tickets for the same show.
We will verify that our solution can handle that without overbooking.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Enforce the limit of available tickets in the `POST /book-tickets` endpoint.
Our endpoint should return `http.StatusBadRequest` if there are not enough tickets available.

It's fully up to you how you implement this logic.
The simplest approach may be doing it inside the repository method used to store the booking:
You can just simply get the number of available tickets and already booked tickets and compare them.

Make sure that you are doing this within the same transaction as the booking is stored in the database.
You should also use the `sql.LevelSerializable` isolation level to make sure that you are not overbooking.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

You can read more about PostgreSQL serializable transactions [in official documentation](https://www.postgresql.org/docs/13/transaction-iso.html#XACT-SERIALIZABLE)
and [this article](https://mkdev.me/posts/transaction-isolation-levels-with-postgresql-as-an-example).

</span>
	</div>
	</div>

We will not check your component tests, but we recommend implementing them. 
This it's a critical functionality, so it's good to have some tests to make sure that it works as expected.

Please try to implement tests yourself. In example solution you can see how we made it.

You can check the [component testing](/trainings/go-event-driven/exercise/56368c8e-1998-4e02-9a7d-b30b8a4af14f) exercise to remind yourself of how to run tests locally.
