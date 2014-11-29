PACKAGE DOCUMENTATION

package identifier
    import "github.com/dschnare/doctag/identifier"

    Package identifier implements a converter for UTF-8 strings to UTF-8
    identifiers.

FUNCTIONS

func ToGoIdentifier(str string) string
    The ToGoIdentifier function converts strings into Go identifiers. This
    function removes all invalid characters from the string. May return an
    empty string.

func ToIdentifierFunc(str string, fn ValidRuneFunc) string
    The ToIdentifierFunc function converts strings into custom identifiers.
    This function removes all characters from the string where fn() returns
    false.

TYPES

type ValidRuneFunc func(r rune, identifierLength int) bool
    ValidRuneFunc is the callback used to test a rune for inclusion in an
    identifier.


