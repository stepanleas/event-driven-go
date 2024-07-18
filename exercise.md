# Catching the Book Ticket request

Now we need to implement the `POST /book-tickets` endpoint.

This endpoint should accept requests in this format:

```json
{
  "show_id": "0299a177-a68a-47fb-a9fb-7362a36efa69",
  "number_of_tickets": 3,
  "customer_email": "email@example.com"
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

It's intended that booking be separate from tickets.
Earlier, we were operating on single tickets. 
Customers are able to book multiple tickets at once, so we need to introduce the concept of booking.

</span>
	</div>
	</div>

## Exercise

File: `project/main.go`

Implement the `POST /book-tickets` endpoint.
It should store bookings in the `bookings` table. 
It's important to use this table name: We'll use it to check your solution.

The booking ID is not sent in the request — you should generate it on your side.
The booking ID should be returned as the response from this endpoint.

```json
{
    "booking_id": "bde0bd8d-88df-4872-a099-d4cf5eb7b491"
}
```


