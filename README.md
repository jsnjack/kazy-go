kažy-go
====

[![CircleCI](https://circleci.com/gh/jsnjack/kazy-go.svg?style=svg)](https://circleci.com/gh/jsnjack/kazy-go)

### What is it?
kažy-go is an implementation of the https://github.com/jsnjack/kazy in golang

### How to use?
```
Highlights output from STDIN
kažy 1.0
Usage: kazy [--include INCLUDE] [--exclude EXCLUDE] [TAIL [TAIL ...]]

Positional arguments:
  TAIL                   highlight patters

Options:
  --include INCLUDE, -i INCLUDE
                         include lines which match patterns
  --exclude EXCLUDE, -e EXCLUDE
                         exclude lines which match patterns
  --help, -h             display this help and exit
  --version              display version and exit

```
kažy is extremely useful when piping a command:
```bash
./kazy -h | ./kazy include exclude lines "match patterns" -e version
```
![ScreenShot](https://raw.githubusercontent.com/jsnjack/kazy-go/master/screenshot.png)

### How to install

#### Debian
```bash
curl -s https://packagecloud.io/install/repositories/jsnjack/kazy-go/script.deb.sh | sudo bash
```

#### Fedora
```bash
curl -s https://packagecloud.io/install/repositories/jsnjack/kazy-go/script.rpm.sh | sudo bash
```
