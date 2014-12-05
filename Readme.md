[![GoDoc](https://godoc.org/github.com/dschnare/doctag?status.svg)](https://godoc.org/github.com/dschnare/doctag)

# Overview

Doctag exposes several Go packages that handle reading doctag files.

## What's a doctag?

Doctags are simple tags that can be written in any text document used to indicate
a named or tagged piece of content. An example document with a doctag could be:

    <{headline}>
    Today's News Stories

When this document is parsed it would parse out a DoctagNode with the name "headline" and
value "\nToday's News Stories". The parsing functions makes the assumption that leading
and trailing whitespace characters are important. It's left for the consumer of the parse
results to decide how to treat leading/trailing whtiespace.

Doctag names can contain any unicode character so long as the tag suffix substring is not encountered.
This example treats doctags like a file system hierarchy. Also notice that doctag names can contain
leading and trailing whitespace characters for greater ledgability.

    <{ page/title }>
    Today's News Stories

    <{ page/content }>
    Blah ablah blab ablaha bal.

Sometimes you may want to skip or ignore a doctag. This can be achieved by prefixing a doctag
with the '!' character.

    <{ !page/title }>
    Yesterday's News Stories

    <{ page/title }>
    Today's News Stories

You may also use this to facilitate a open/close pattern to your doctags, thereby giving more control
over leading and trailing whitespace and/or improving ledgability.

    <{ page/title }>Today's News Stories<{!}>

    <{ page/content }>
    Blah ablah blab ablaha bal.
    <{!}>


# Usage

    doctag {file path} | doctag [help|/?]
      -help=false: Show the help message.
      -hierarchical=false: Converts the flat doctag tree into a nested JSON object.
      -hierarchy=false: Converts the flat doctag tree into a nested JSON object. (shorthand)
      -output="": The output file to write to.
      -pretty=false: Print JSON result with indentation. (shorthand)
      -pretty-print=false: Print JSON result with indentation.
      -tag-prefix="<{": The prefix to use for doc tags.
      -tag-separator="/": The separator character to use for hierarchical doc tags.
      -tag-suffix="}>": The suffix to use for doc tags.
      -trim=false: Trim the leading and trailing whitespace from all doctag values.
      -warn=false: Print warning messages.

If no file path is specified as an argument then a file contents are expected to be piped into stdin.

If no output argument is specified then the out is piped to stdout.

# Packages

**[parse](http://godoc.org/github.com/dschnare/doctag/parse)** - Package parse builds a slice of nodes from UTF-8 encoded text documents that have doctags.

**[identifier](http://godoc.org/github.com/dschnare/doctag/identifier)** - Package identifier implements a converter for UTF-8 strings to UTF-8 identifiers.

**[hierarchy](http://godoc.org/github.com/dschnare/doctag/hierarchy)** - Package hierarchy implements a doctag transformer that transforms a list of doctags into a map hierarchy.

# Commands

**[doctag](http://godoc.org/github.com/dschnare/doctag)** - The doctag command exposes a doctag parser and hierarchy transformer.