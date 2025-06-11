# `mail-tui`

A minimal terminal-based email client using [`bubbletea`](https://github.com/charmbracelet/bubbletea), loading emails over IMAP or via a fake backend with dummy data.

Run it with the following command, replacing the address and credentials as necessary.
Note that only insecure mode is supported currently.

```
$ go run ./cmd/tui --imap-address="localhost:1143" --username="user" --password="password"
```

It can also be run using a fake backend, which displays dummy data instead of connecting to an IMAP server.

```
$ go run ./cmd/tui --use-imap=false
```

## Structure

The main entrypoint of the application is in `cmd/tui/main.go`.
It wires things up and runs the rest of the [`bubbletea`](https://github.com/charmbracelet/bubbletea) app.

I've loosely split up the UI part of the app into components, which are in the `internal/ui` directory.
So far, this contains the following components:
- `app`: the root application component, responsible for switching between the other views
- `email_list`: renders a list of emails
- `email_viewer`: displays a single email
- `email_composer`: a form-esque component for composing a new email

The "domain model" is in `internal/core`.
In here we have some structs representing the email domain.
There is also the abstract `EmailBackend` interface, to allow the `internal/ui` components to remain decoupled from the underlying email backend implementation.
This has 2 implementations:
- `internal/backend/fake`: returns dummy data
- `internal/backend/imap`: connects to an IMAP server insecurely

## Testing

There are unit tests for some of the UI components, but I haven't added tests for all of them.
They verify that the [`bubbletea`](https://github.com/charmbracelet/bubbletea) components update their state models correctly in response to various messages.

None of the actual happy-path terminal output is currently tested, although [this experimental `teatest` package](https://github.com/charmbracelet/x/tree/main/exp/teatest/v2) could be used to help with this (it essentially does snapshot-based testing).

The tests can be run with:

```
$ go test ./...
```

## Future Enhancements

Of course it isn't really usable at this stage, so these are some things I could still add:

- An SMTP client for sending emails
- Add a SQLite database to cache the mailbox data and avoid needing to re-fetch the whole thing each time the app loads
  - Should also store the query state so we can efficiently request only what has changed since the last time the app ran
- A background goroutine to subscribe to changes and update the cache accordingly
- Allow the user to open their `$EDITOR` to compose an email instead of using the built-in form
- Add config file to `$XDG_CONFIG_HOME/mail-tui/` for storing email account settings
  - Also need to consider how to pass in secrets securely

## Non-goals

Things I'm not planning to support, in order to keep things simple:

- Multiple mailboxes or email accounts - you can't see your "Sent" mailbox or configure multiple email accounts
- Replies and threads - no responding to emails or viewing threads of email responses; each email is standalone
- Forwarding emails
- Attachments and HTML-formatted emails - plain-text only
- Cc/Bcc
- Drafts
- Contacts or address book to pre-populate email addresses
- Email search
- Real-time UI updates when changes occur
