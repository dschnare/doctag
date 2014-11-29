// Package identifier implements a converter for UTF-8 strings to UTF-8 identifiers.
package identifier

import (
  "unicode"
  "unicode/utf8"
)

// ValidRuneFunc is the callback used to test a rune for inclusion in an identifier.
type ValidRuneFunc func (r rune, identifierLength int) bool

// The ToGoIdentifier function converts strings into Go identifiers.
// This function removes all invalid characters from the string.
// May return an empty string.
func ToGoIdentifier(str string) string {
  return ToIdentifierFunc(str, func (r rune, identifierLength int) bool {
    if identifierLength == 0 {
      return r == '_' || unicode.IsLetter(r)
    } else if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
      return true
    }
    return false
  })
}

// The ToIdentifierFunc function converts strings into custom identifiers.
// This function removes all characters from the string where fn() returns false.
func ToIdentifierFunc(str string, fn ValidRuneFunc) string {
  maxRuneCount := utf8.RuneCount([]byte(str))
  id := make([]rune, 0, maxRuneCount)

  for _,r := range str {
    if fn(r, len(id)) {
      id = append(id, r)
    }
  }

  return string(id)
}