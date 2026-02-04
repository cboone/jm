# jm error handling

Verify structured error output when things go wrong.

## Missing token produces structured JSON error

```scrut
$ $TESTDIR/../jm session 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Missing token with text format

```scrut
$ $TESTDIR/../jm session --format text 2>&1
Error [authentication_failed]: no token configured; set JMAP_TOKEN, --token, or token in config file
Hint: Check your token in JMAP_TOKEN or config file
[1]
```

## List without token

```scrut
$ $TESTDIR/../jm list 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Read without token

```scrut
$ $TESTDIR/../jm read some-email-id 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Search without token

```scrut
$ $TESTDIR/../jm search "test query" 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Archive without token

```scrut
$ $TESTDIR/../jm archive M123 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Spam without token

```scrut
$ $TESTDIR/../jm spam M123 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Move without token

```scrut
$ $TESTDIR/../jm move M123 --to Archive 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Mailboxes without token

```scrut
$ $TESTDIR/../jm mailboxes 2>&1
{
  "error": "authentication_failed",
  "message": "no token configured; set JMAP_TOKEN, --token, or token in config file",
  "hint": "Check your token in JMAP_TOKEN or config file"
}
[1]
```

## Invalid token produces connection error

```scrut
$ $TESTDIR/../jm session --token "invalid-token-xxx" 2>&1
{
  "error": "authentication_failed",
* (glob+)
[1]
```
