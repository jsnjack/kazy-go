kažy-go
====

### What is it?
kažy is an application that highlights, filters and extracts string patterns from STDIN

### How to use?
```
Highlights, filters and extracts string patterns from STDIN

Usage:
  kazy [<pattern>...] [flags]

Flags:
  -b, --buffer int            buffer size in KB (default 64)
  -e, --exclude stringArray   exclude from output lines which match provided patterns
  -x, --extract               extract matched strings (leftmost) instead of highlighting them
  -h, --help                  help for kazy
  -i, --include stringArray   only include lines which match provided patterns
  -l, --limit int             limit the length of the line, characters
  -r, --regexp                use RegExp patterns instead of string patterns
      --version               print version and exit

```
kažy is extremely useful when piping a command:
```bash
./kazy -h | ./kazy include exclude lines "match patterns" -e version
```
![ScreenShot](https://raw.githubusercontent.com/jsnjack/kazy-go/master/screenshot.png)

### How to install

```
grm install jsnjack/kazy-go
```
