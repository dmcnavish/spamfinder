# spamfinder

Spamfinder searches your Gmail account, finds spam emails, and writes their name and unsubscribe URL to a CSV file.

Before running, you will need to enable the Gmail API from the Google dev console, generate a client_secret.json file, and place it in the root directory.

Then, in Windows run:

```go get```

```go build && spamfinder.exe```
