kažy-go
====

### What is it?
kažy is an application that highlights and filters output from STDIN

### How to use?
```
Highlights output from STDIN

Usage:
  kazy [<pattern>...] [flags]

Flags:
  -e, --exclude stringArray   exclude from output lines which match provided patterns
  -h, --help                  help for kazy
  -i, --include stringArray   only include lines which match provided patterns
  -l, --limit int             limit the length of the line, characters
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
