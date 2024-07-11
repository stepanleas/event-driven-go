# Project: Logging middleware

Some middleware functions need a dependency or config. There are two approaches you can use.

First, simply pass the dependency as an argument to the middleware.
It looks quite complex, but it's basically a function returning a middleware function.

```go
func SaveMessageMiddleware(saver MessageSaver) func (h message.HandlerFunc) message.HandlerFunc {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			err := saver.Save(msg)
			if err != nil {
				return nil, err
			}
			
			return next(msg)
		}
	}
}
```

Use it like this:

```go
router.AddMiddleware(SaveMessageMiddleware(saver))
```

When you need more dependencies, you can consider making the function a struct.

```go
type RandomFail struct {
	Chance float64
	Error error
}

func (m RandomFail) Middleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		if rand.Float64() < m.Chance {
			return nil, m.Error
		}
		return next(msg)
	}
}
```

Use it like this (notice how the `Middleware` method is passed, not the struct itself):

```go
router.AddMiddleware(RandomFail{
	Chance: 0.1, 
	Error: errors.New("random error occurred"),
}.Middleware)
```

## Exercise

File: `project/main.go`

Add a middleware function to your project that logs incoming messages. The log message should be:

```text
Handling a message
```

It should include a log field called `message_uuid` with the message UUID.
Using logrus, you can add log fields like this:

```go
logrus.WithField("key", value).Info("Log message")
```
