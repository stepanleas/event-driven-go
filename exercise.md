# Returning tickets from API

## Exercise

File: `project/main.go`

Implement a `GET /tickets` endpoint that returns all tickets in this format:

```json
[
    {
        "ticket_id": "<ticket_id>",
        "customer_email": "<customer_email>",
        "price": {
            "amount": "<amount>",
            "currency": "<currency>"
        }
    },
    {
        "ticket_id": "<ticket_id>",
        "customer_email": "<customer_email>",
        "price": {
            "amount": "<amount>",
            "currency": "<currency>"
        }
    }
]
```

In this exercise, we don't care about pagination, filtering, etc. It's good enough to return all tickets.
(Of course, if you want, you can add this functionality as well).
