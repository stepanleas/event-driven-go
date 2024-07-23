# Project: Returning the Read Model

We have our Read Model stored in the database. It's time expose it to our Operations Team!

## Exercise

File: `project/main.go`

Implement two endpoints:

- `/ops/bookings`, which will return all bookings,
- `/ops/bookings/:id`, which will return a single booking.

You should just return the payload as is from the database.
Pagination is not required.

For single booking, the response format should be:

```json
{
  "booking_id": "31bc30c1-115d-4272-8066-5351f8e17433",
  "booked_at": "2023-07-21T17:52:20.406021166Z",
  "tickets": {
    "751e957c-9745-4ae6-bdb7-b2e14befc729": {
      "price_amount": "20.00",
      "price_currency": "EUR",
      "customer_email": "email@example.com",
      "status": "refunded",
      "printed_at": "2023-07-21T17:52:20.580737625Z",
      "printed_file_name": "751e957c-9745-4ae6-bdb7-b2e14befc729-ticket.html",
      "receipt_issued_at": "2023-07-21T17:52:20.558529458Z",
      "receipt_number": "2023-004833"
    },
    "7d960a14-3f04-47e5-a238-3db9466a51b0": {
      "price_amount": "20.00",
      "price_currency": "EUR",
      "customer_email": "email@example.com",
      "status": "refunded",
      "printed_at": "2023-07-21T17:52:20.585196Z",
      "printed_file_name": "7d960a14-3f04-47e5-a238-3db9466a51b0-ticket.html",
      "receipt_issued_at": "2023-07-21T17:52:20.587919208Z",
      "receipt_number": "2023-004834"
    }
  },
  "last_update": "2023-07-21T17:52:21.050429667Z"
}
```

For all bookings, the response should be:

```json
[
  {
    "booking_id": "31bc30c1-115d-4272-8066-5351f8e17433",
    "booked_at": "2023-07-21T17:52:20.406021166Z",
    "tickets": {
      "751e957c-9745-4ae6-bdb7-b2e14befc729": {
        "price_amount": "20.00",
        "price_currency": "EUR",
        "customer_email": "email@example.com",
        "status": "refunded",
        "printed_at": "2023-07-21T17:52:20.580737625Z",
        "printed_file_name": "751e957c-9745-4ae6-bdb7-b2e14befc729-ticket.html",
        "receipt_issued_at": "2023-07-21T17:52:20.558529458Z",
        "receipt_number": "2023-004833"
      },
      "7d960a14-3f04-47e5-a238-3db9466a51b0": {
        "price_amount": "20.00",
        "price_currency": "EUR",
        "customer_email": "email@example.com",
        "status": "refunded",
        "printed_at": "2023-07-21T17:52:20.585196Z",
        "printed_file_name": "7d960a14-3f04-47e5-a238-3db9466a51b0-ticket.html",
        "receipt_issued_at": "2023-07-21T17:52:20.587919208Z",
        "receipt_number": "2023-004834"
      }
    },
    "last_update": "2023-07-21T17:52:21.050429667Z"
  }
]
```