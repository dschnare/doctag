PACKAGE DOCUMENTATION

package parse
    import "github.com/dschnare/doctag/parse"

    Package parse builds a slice of nodes from UTF-8 encoded text documents
    that have doctags.

    Doctags are simple tags that can be written in any text document used to
    indicate a named or tagged piece of content. An example document with a
    doctag could be:

	<{headline}>
	Today's News Stories

    When this document is parsed it would parse out a DoctagNode with the
    name "headline" and value "\nToday's News Stories". The parsing
    functions makes the assumption that leading and trailing whitespace
    characters are important. It's left for the consumer of the parse
    results to decide how to treat leading/trailing whtiespace.

    Doctag names can contain any unicode character so long as the tag suffix
    substring is not encountered. This example treats doctags like a file
    system hierarchy. Also notice that doctag names can contain leading and
    trailing whitespace characters for greater ledgability.

	<{ page/title }>
	Today's News Stories

	<{ page/content }>
	Blah ablah blab ablaha bal.

    Sometimes you may want to skip or ignore a doctag. This can be achieved
    by prefixing a doctag with the '!' character.

	<{ !page/title }>
	Yesterday's News Stories

	<{ page/title }>
	Today's News Stories

    You may also use this to facilitate a open/close pattern to your
    doctags, thereby giving more control over leading and trailing
    whitespace and/or improving ledgability.

	<{ page/title }>Today's News Stories<{!}>

	<{ page/content }>
	Blah ablah blab ablaha bal.
	<{!}>

CONSTANTS

const (
    DefaultTagPrefix = "<{"
    DefaultTagSuffix = "}>"
)
    The default tag prefix and suffix used by the Parse() function.

VARIABLES

var (
    Logger *log.Logger
)
    The optional Logger to have warnings logged to. The logger is useful in
    finding what might be typos when declaring a doctag in a document.

FUNCTIONS

func Parse(reader *bufio.Reader) ([]*DoctagNode, error)
    Parse parses a buffered reader for doctags using the default prefix and
    suffix substrings. The returned slice contains all parsed DoctagNodes in
    the order they appear in the document.

func ParseFile(fileName string) ([]*DoctagNode, error)
    Parse parses a text file for doctags using the default prefix and suffix
    substrings. The returned slice contains all parsed DoctagNodes in the
    order they appear in the document.

func ParseFileWithPrefixAndSuffix(fileName string, tagPrefix string, tagSuffix string) ([]*DoctagNode, error)
    ParseFileWithPrefixAndSuffix parses a text file for doctags using custom
    prefix and suffix substrings for doctags. The returned slice contains
    all parsed DoctagNodes in the order they appear in the document.

func ParseWithPrefixAndSuffix(reader *bufio.Reader, tagPrefix string, tagSuffix string) (doctags []*DoctagNode, err error)
    ParseWithPrefixAndSuffix parses a buffered reader for doctags using
    custom prefix and suffix substrings for doctags. The returned slice
    contains all parsed DoctagNodes in the order they appear in the
    document.

TYPES

type DoctagNode struct {
    Name   string
    Value  string
    Line   int
    Column int
}
    A DoctagNode represents a doctag parsed from a text document.

SUBDIRECTORIES

	fixtures

