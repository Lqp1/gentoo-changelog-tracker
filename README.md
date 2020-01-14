# gentoo-changelog-tracker
A script wrapping equery and Gentoo's gitweb to retrieve Changelogs

As it's not possible anymore to track Changelogs using `equery changes`, this
simple wrapper script tries to get info via Gentoo's official gitweb (RSS
feed thus it's limited to 10 entries).

Usage examples:

```
go get -d ./... && go build
./gentoo-changelog-tracker --limit 5 vim
./gentoo-changelog-tracker --limit 1 --full emacs
```

Parameters:
 - `--limit` : limit output to n entries. Can't show more than 10 entries anyway
 - `--full` : get patch for each entry and print diff instead of just the
   summary
