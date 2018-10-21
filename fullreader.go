package main

import (
	"bufio"
)

func readLine(io *bufio.Reader) string {
	lnbuf := make([]byte, 0)
	for true {
		cont, isprefix, err := io.ReadLine()
		if err != nil {
			if len(lnbuf) != 0 {
				lnbuf = append(lnbuf, cont...)
				return string(lnbuf)
			}
			return string(cont)
		}
		if isprefix == false {
			if len(lnbuf) != 0 {
				lnbuf = append(lnbuf, cont...)
				return string(lnbuf)
			}
			return string(cont)
		}
		lnbuf = append(lnbuf, cont...)
	}
	return string(lnbuf)
}

func readAllLines(textfile *bufio.Reader) []string {
	sarr := make([]string, 0)
	lnbuf := make([]byte, 0)
	for true {
		cont, isprefix, err := textfile.ReadLine()
		if err != nil {
			if len(lnbuf) != 0 {
				lnbuf = append(lnbuf, cont...)
				sarr = append(sarr, string(lnbuf))
			}
			break
		}
		if isprefix == false {
			if len(lnbuf) != 0 {
				lnbuf = append(lnbuf, cont...)
				sarr = append(sarr, string(lnbuf))
			} else {
				sarr = append(sarr, string(cont))
			}
		} else {
			lnbuf = append(lnbuf, cont...)
		}
	}
	return sarr
}
