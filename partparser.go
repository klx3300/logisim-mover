package main

import (
	"log"
	"os"
	"strings"
)

// TokenRanger !
type TokenRanger struct {
	Top    int
	Bottom int
	Type   string
	Params map[string]string
}

func parameterExtractor(lineno int, paramstr string) (map[string]string, bool) {
	pbarr := []byte(paramstr)
	namer := make([]byte, 0)
	valuer := make([]byte, 0)
	state := 0
	params := make(map[string]string)
	// 0 Naming
	// 1 Pre-value
	// 2 Valuing
	// 3 Post-value
	escaper := false
	for _, char := range pbarr {
		switch state {
		case 0:
			switch char {
			case '=':
				if escaper == true {
					namer = append(namer, char)
					escaper = false
				} else {
					state = 1
				}
			case '\\':
				if escaper == true {
					namer = append(namer, char)
					escaper = false
				} else {
					escaper = true
				}
			case '>': // terminage mark
				if escaper == true {
					namer = append(namer, char)
					escaper = false
				} else {
					return params, false
				}
			case ' ':
				continue
			default:
				namer = append(namer, char)
				escaper = false
			}
		case 1:
			switch char {
			case '"':
				state = 2
			case ' ':
				continue
			default:
				log.Fatalf("Token failure at %d, attempt to break pre-value state\n", lineno)
				os.Exit(-1)
			}
		case 2:
			switch char {
			case '"':
				if escaper == true {
					valuer = append(valuer, char)
					escaper = false
				} else {
					state = 3
				}
			case '\\':
				if escaper == true {
					valuer = append(valuer, char)
					escaper = false
				} else {
					escaper = true
				}
			default:
				valuer = append(valuer, char)
			}
		case 3:
			switch char {
			case ' ':
				// complete push
				params[string(namer)] = string(valuer)
				state = 0
			case '/': // terminage mark
				params[string(namer)] = string(valuer)
				return params, true
			case '>': // terminage mark
				params[string(namer)] = string(valuer)
				return params, false
			default:
				log.Fatalf("Token failure at %d, attempt to break post-value state\n", lineno)
				os.Exit(-1)
			}
		}
	}
	log.Fatalf("Token failure at %d, incomplete sentence\n", lineno)
	os.Exit(-1)
	return params, false
}

func getRange(in []string, tType string) []TokenRanger {
	ranges := make([]TokenRanger, 0)
	tkrecord := new(TokenStack)
	for ln, str := range in {
		trimmed := strings.TrimSpace(str)
		if strings.HasPrefix(trimmed, "<"+tType) {
			log.Printf("Found %s begin at %d\n", tType, ln)
			var tk TokenRanger
			tk.Top = ln
			tk.Type = tType
			isClosing := false
			tk.Params, isClosing = parameterExtractor(ln, strings.TrimPrefix(trimmed, "<"+tType))
			if isClosing == true {
				tk.Bottom = ln
				log.Printf("Found %s end immediate\n", tType, ln)
				ranges = append(ranges, tk)
			} else {
				tk.Bottom = -1
				tkrecord.TokenPush(tk)
			}
		} else if strings.HasPrefix(trimmed, "</"+tType) {
			log.Printf("Found %s end at %d\n", tType, ln)
			if tkrecord.TokenEmpty() == true {
				log.Fatalln("Pairing Failure: %s :: Extra Closer\n", tType)
				os.Exit(-1)
			}
			tk := tkrecord.TokenTop()
			tk.Bottom = ln
			ranges = append(ranges, tk)
			tkrecord.TokenPop()
		}
	}
	if tkrecord.TokenEmpty() == false {
		log.Fatalln("Pairing Failure: %s :: Insufficient Closer\n", tType)
		os.Exit(-1)
	}
	return ranges
}
