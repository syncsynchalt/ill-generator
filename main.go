package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
 * Reads a template on stdin and interprets it to stdout.
 *
 * Interprets the following commands:
 *
 * %%
 *   a literal '%'
 * %file filename
 *   loads filename
 * %next bytes
 *   takes the next {bytes} bytes from the file and loads them as current context
 * %empty
 *   assert that we've read to the end of the file with %next commands
 * %bytes
 *   prints all the current context as zero-padded hex
 * %{num}
 *   prints the numbered byte from current context as zero-padded hex
 * %-{num}
 *   prints the numbered byte from end current context as zero-padded hex (%-1 is last byte)
 * %x{num}
 *   prints the numbered byte from current context as capital hex with leading "0x"
 * %xx{num}
 *   prints the 16-bit number {num}..{num+1} as capital hex with leading "0x"
 * %xxx{num}
 *   prints the 24-bit number {num}..{num+2} as capital hex with leading "0x"
 * %xxxx{num}
 *   prints the 32-bit number {num}..{num+3} as capital hex with leading "0x"
 * %d{num}
 *   prints the numbered byte from current context as decimal number
 * %dd{num}
 *   prints the 16-bit number {num}..{num+1} as decimal number
 * %ddd{num}
 *   prints the 24-bit number {num}..{num+2} as decimal number
 * %dddd{num}
 *   prints the 32-bit number {num}..{num+3} as decimal number
 * %n{num}
 *   prints the numbered byte from current context as capital hex with leading "0x"
 *   then in parenthetical decimal, unless the number is less than 0x0a in which case
 *   it only prints the decimal number (since it's the same in hex and dec).
 * %nn{num}
 *   prints the 16-bit number {num}..{num+1} as capital hex with leading "0x"
 *   then in parenthetical decimal, unless the number is less than 0x0a in which case
 *   it only prints the decimal number (since it's the same in hex and dec).
 * %nnn{num}
 *   prints the 24-bit number {num}..{num+2} as capital hex with leading "0x"
 *   then in parenthetical decimal, unless the number is less than 0x0a in which case
 *   it only prints the decimal number (since it's the same in hex and dec).
 * %nnnn{num}
 *   prints the 32-bit number {num}..{num+3} as capital hex with leading "0x"
 *   then in parenthetical decimal, unless the number is less than 0x0a in which case
 *   it only prints the decimal number (since it's the same in hex and dec).
 * %stop
 *   stop interpreting commands
 */

func main() {
	btemplate, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	template := string(btemplate)

	data := make([]byte, 0)
	context := make([]byte, 0)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	i := 0
	for i < len(template) {
		if template[i] == '%' {
			// handle escape
			if len(template) > i && template[i+1] == '%' {
				writer.Write([]byte{'%'})
				i += 2
				continue
			}

			// read the full line for convenience
			end := strings.IndexByte(template[i:], '\n')
			if end == -1 {
				end = len(template) - i
			}
			line := template[i+1 : i+end]
			fmt.Fprintf(os.Stderr, "command %%%s\n", line)

			if len(line) > 5 && line[0:5] == "file " {
				// %file filename
				data, err = ioutil.ReadFile(line[5:])
				if err != nil {
					panic(err)
				}
				i += end + 1
			} else if len(line) > 5 && line[0:5] == "next " {
				// %next bytes
				toks := strings.Fields(line)
				nbytes, err := strconv.Atoi(toks[1])
				if err != nil {
					panic(err)
				}
				context = data[:nbytes]
				data = data[nbytes:]
				i += end + 1
			} else if len(line) >= 5 && line[0:5] == "empty" {
				if len(data) > 0 {
					panic(fmt.Sprintf("there are leftover bytes: %v", data))
				}
				i += end + 1
			} else if len(line) >= 5 && line[0:5] == "bytes" {
				// %bytes
				for i := range context {
					fmt.Fprintf(writer, " %02x", context[i])
					if i%32 == 31 {
						fmt.Fprintf(writer, "\n")
					}
				}
				if len(context)%32 != 0 {
					fmt.Fprintf(writer, "\n")
				}
				i += end + 1
			} else if len(line) >= 4 && line[0:4] == "stop" {
				// %stop
				i += end + 1
				writer.Write([]byte(template[i:]))
				i = len(template)
			} else if len(line) > 6 && line[0:6] == "dddddd" {
				// %dddddd{bytenum} (as decimal)
				num, span := parseNumber(line[6:])
				val := numFromContext(context, num, 6)
				fmt.Fprintf(writer, "%d", val)
				i += 1 + 6 + span
			} else if len(line) > 5 && line[0:5] == "ddddd" {
				// %ddddd{bytenum} (as decimal)
				num, span := parseNumber(line[5:])
				val := numFromContext(context, num, 5)
				fmt.Fprintf(writer, "%d", val)
				i += 1 + 5 + span
			} else if len(line) > 4 && line[0:4] == "dddd" {
				// %dddd{bytenum} (as decimal)
				num, span := parseNumber(line[4:])
				val := numFromContext(context, num, 4)
				fmt.Fprintf(writer, "%d", val)
				i += 1 + 4 + span
			} else if len(line) > 3 && line[0:3] == "ddd" {
				// %ddd{bytenum} (as decimal)
				num, span := parseNumber(line[3:])
				val := numFromContext(context, num, 3)
				fmt.Fprintf(writer, "%d", val)
				i += 1 + 3 + span
			} else if len(line) > 2 && line[0:2] == "dd" {
				// %dd{bytenum} (as decimal)
				num, span := parseNumber(line[2:])
				val := numFromContext(context, num, 2)
				fmt.Fprintf(writer, "%d", val)
				i += 1 + 2 + span
			} else if len(line) > 1 && line[0] == 'd' {
				// %d{bytenum} (as decimal)
				num, span := parseNumber(line[1:])
				fmt.Fprintf(writer, "%d", context[num])
				i += 1 + 1 + span
			} else if len(line) > 1 && line[0] == '-' {
				// %-{bytenum} (from end, as hex)
				num, span := parseNumber(line[1:])
				fmt.Fprintf(writer, "%02x", context[len(context)-num])
				i += 1 + 1 + span
			} else if len(line) > 6 && line[0:6] == "xxxxxx" {
				// %xxxxxx{bytenum} (as 0xHex)
				num, span := parseNumber(line[6:])
				val := numFromContext(context, num, 6)
				fmt.Fprintf(writer, "0x%X", val)
				i += 1 + 6 + span
			} else if len(line) > 5 && line[0:5] == "xxxxx" {
				// %xxxxx{bytenum} (as 0xHex)
				num, span := parseNumber(line[4:])
				val := numFromContext(context, num, 5)
				fmt.Fprintf(writer, "0x%X", val)
				i += 1 + 5 + span
			} else if len(line) > 4 && line[0:4] == "xxxx" {
				// %xxxx{bytenum} (as 0xHex)
				num, span := parseNumber(line[4:])
				val := numFromContext(context, num, 4)
				fmt.Fprintf(writer, "0x%X", val)
				i += 1 + 4 + span
			} else if len(line) > 3 && line[0:3] == "xxx" {
				// %xxx{bytenum} (as 0xHex)
				num, span := parseNumber(line[3:])
				val := numFromContext(context, num, 3)
				fmt.Fprintf(writer, "0x%X", val)
				i += 1 + 3 + span
			} else if len(line) > 2 && line[0:2] == "xx" {
				// %xx{bytenum} (as 0xHex)
				num, span := parseNumber(line[2:])
				val := numFromContext(context, num, 2)
				fmt.Fprintf(writer, "0x%X", val)
				i += 1 + 2 + span
			} else if len(line) > 1 && line[0] == 'x' {
				// %x{bytenum} (as 0xHex)
				num, span := parseNumber(line[1:])
				fmt.Fprintf(writer, "0x%X", context[num])
				i += 1 + 1 + span
			} else if len(line) > 4 && line[0:6] == "nnnnnn" {
				// %nnnnnn{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[6:])
				val := numFromContext(context, num, 6)
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 6 + span
			} else if len(line) > 5 && line[0:5] == "nnnnn" {
				// %nnnn{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[5:])
				val := numFromContext(context, num, 5)
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 5 + span
			} else if len(line) > 4 && line[0:4] == "nnnn" {
				// %nnnn{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[4:])
				val := numFromContext(context, num, 4)
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 4 + span
			} else if len(line) > 3 && line[0:3] == "nnn" {
				// %nnn{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[3:])
				val := numFromContext(context, num, 3)
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 3 + span
			} else if len(line) > 2 && line[0:2] == "nn" {
				// %nn{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[2:])
				val := numFromContext(context, num, 2)
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 2 + span
			} else if len(line) > 1 && line[0] == 'n' {
				// %n{bytenum} (as 0xHex (dec))
				num, span := parseNumber(line[1:])
				var val = context[num]
				if val >= 0 && val < 10 {
					fmt.Fprintf(writer, "%d", val)
				} else {
					fmt.Fprintf(writer, "0x%X (%d)", val, val)
				}
				i += 1 + 1 + span
			} else if len(line) >= 1 && isNum(line[0]) {
				// %{bytenum} (as hex)
				num, span := parseNumber(line)
				fmt.Fprintf(writer, "%02x", context[num])
				i += 1 + span
			} else {
				panic(fmt.Sprintf("Unhandled command %%%s", line))
			}
		} else {
			writer.Write([]byte{template[i]})
			i++
		}
	}
}

func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}

func numFromContext(context []byte, index int, byteCount int) int {
	val := 0
	for b := 0; b < byteCount; b++ {
		val |= int(context[index+b]) << (8 * b)
	}
	return val
}

func parseNumber(line string) (num int, span int) {
	nend := 0
	for nend < len(line) && isNum(line[nend]) {
		nend++
	}
	if nend == 0 {
		panic(fmt.Sprintf("bad number at [%s]", line))
	}
	num, err := strconv.Atoi(line[:nend])
	if err != nil {
		panic(err)
	}
	return num, nend
}
