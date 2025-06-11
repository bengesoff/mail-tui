# Running a dummy IMAP server to test against

There is a dummy IMAP server called [`tameimap`](https://github.com/bcampbell/tameimap), which uses the [`go-imap`](https://pkg.go.dev/github.com/emersion/go-imap/v2@v2.0.0-beta.5/imapclient) package to serve a mailbox from a set of static files.
This can be used to test the TUI application against, to avoid needing to set up a real IMAP server and add emails to it.

To do so, run the following commands:

```
$ cd imap_test_server
$ docker build --tag tameimap --file Dockerfile.tameimap .
$ docker run --rm --name imap-server-test --publish 1143:1143 tameimap
```

Then run the TUI application in another shell, using `localhost:1143` as the IMAP server address, `bob` as the username, and `pass` as the password.

> These are also the default values for the CLI flags, so you can omit them.
