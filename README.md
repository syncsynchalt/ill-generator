# ill-generator

Reads a template on stdin and interprets it to stdout.

## Installation

```
go install github.com/syncsynchalt/ill-generator@latest
```

## Usage

```
ill-generator < input.template > output
```

## Template commands

Interprets the following commands in templates:

* `%%`:
   a literal '%'
* `%file filename`:
   loads filename
* `%next bytes`:
   takes the next {bytes} bytes from the file and loads them as current context
* `%empty`:
   assert that we've read to the end of the file with %next commands
* `%bytes`:
   prints all the current context as zero-padded hex
* `%{num}`:
   prints the numbered byte from current context as zero-padded hex
* `%-{num}`:
   prints the numbered byte from end current context as zero-padded hex (%-1 is last byte)
* `%x{num}`:
   prints the numbered byte from current context as capital hex with leading "0x"
* `%xx{num}`:
   prints the 16-bit number {num}..{num+1} as capital hex with leading "0x"
* `%xxx{num}`:
   prints the 24-bit number {num}..{num+2} as capital hex with leading "0x"
* `%xxxx{num}`:
   prints the 32-bit number {num}..{num+3} as capital hex with leading "0x"
* `%xxxxx{num}`:
   prints the 40-bit number {num}..{num+4} as capital hex with leading "0x"
* `%xxxxxx{num}`:
   prints the 48-bit number {num}..{num+5} as capital hex with leading "0x"
* `%d{num}`:
   prints the numbered byte from current context as decimal number
* `%dd{num}`:
   prints the 16-bit number {num}..{num+1} as decimal number
* `%ddd{num}`:
   prints the 24-bit number {num}..{num+2} as decimal number
* `%dddd{num}`:
   prints the 32-bit number {num}..{num+3} as decimal number
* `%ddddd{num}`:
   prints the 40-bit number {num}..{num+4} as decimal number
* `%dddddd{num}`:
   prints the 48-bit number {num}..{num+5} as decimal number
* `%n{num}`:
   prints the numbered byte from current context as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%nn{num}`:
   prints the 16-bit number {num}..{num+1} as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%nnn{num}`:
   prints the 24-bit number {num}..{num+2} as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%nnnn{num}`:
   prints the 32-bit number {num}..{num+3} as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%nnnnn{num}`:
   prints the 40-bit number {num}..{num+4} as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%nnnnnn{num}`:
   prints the 48-bit number {num}..{num+5} as capital hex with leading "0x"
   then in parenthetical decimal, unless the number is less than 0x0a in which case
   it only prints the decimal number (since it's the same in hex and dec).
* `%stop`:
   stop interpreting commands
