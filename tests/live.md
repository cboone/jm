# jm live integration tests

These tests require a valid JMAP token and explicit opt-in.
They are skipped automatically when preconditions are not met.

## Preconditions

```scrut
$ test "$JMAP_LIVE_TESTS" = "1" || exit 80
```

```scrut
$ test -n "$JMAP_TOKEN" || exit 80
```

## Session returns JSON with username

```scrut
$ $TESTDIR/../jm session 2>&1
{
  "username": *, (glob)
* (glob+)
}
```

## Session with text format

```scrut
$ $TESTDIR/../jm session --format text 2>&1
Username: * (glob)
* (glob+)
```

## Mailboxes returns JSON array

```scrut
$ $TESTDIR/../jm mailboxes 2>&1
* (glob+)
```

## Mailboxes roles-only returns entries with role field

```scrut
$ $TESTDIR/../jm mailboxes --roles-only 2>&1
* (glob+)
```

## List returns JSON with total and emails

```scrut
$ $TESTDIR/../jm list --limit 1 2>&1
{
  "total": *, (glob)
  "offset": 0,
* (glob+)
}
```

## List with text format

```scrut
$ $TESTDIR/../jm list --limit 1 --format text 2>&1
Total: * (glob)
* (glob+)
```

## Search with filter returns results

```scrut
$ $TESTDIR/../jm search --limit 1 2>&1
{
  "total": *, (glob)
  "offset": 0,
* (glob+)
}
```

## Search with text query returns snippets

```scrut
$ $TESTDIR/../jm search "the" --limit 1 2>&1
{
  "total": *, (glob)
  "offset": 0,
* (glob+)
}
```

## Mailboxes text format shows tabular output

```scrut
$ $TESTDIR/../jm mailboxes --format text 2>&1
* (glob+)
```

## Error in text format uses Error prefix

```scrut
$ env -u JMAP_TOKEN $TESTDIR/../jm session --format text 2>&1
Error [authentication_failed]: * (glob)
* (glob+)
[1]
```
