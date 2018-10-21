package main

import (
	"fmt"
)

// TokenStack !
type TokenStack []TokenRanger

// TokenPush !
func (arr *TokenStack) TokenPush(tk TokenRanger) {
	*arr = append(*arr, tk)
}

// TokenPop !
func (arr *TokenStack) TokenPop() error {
	if len(*arr) == 0 {
		return fmt.Errorf("Stack Empty")
	}
	if len(*arr) == 1 {
		*arr = make([]TokenRanger, 0)
	} else {
		*arr = append(make([]TokenRanger, 0), (*arr)[0:len(*arr)-1]...)
	}
	return nil
}

// TokenTop !
func (arr *TokenStack) TokenTop() TokenRanger {
	return (*arr)[len(*arr)-1]
}

// TokenEmpty !
func (arr *TokenStack) TokenEmpty() bool {
	return len(*arr) == 0
}

// ToXML !
func (tk TokenRanger) ToXML() string {
	builder := "<"
	builder += tk.Type
	builder += " "
	for k, v := range tk.Params {
		builder += k
		builder += "=\""
		builder += v
		builder += "\" "
	}
	if tk.Top == tk.Bottom {
		builder += "/>"
	} else {
		builder += ">"
	}
	return builder
}
