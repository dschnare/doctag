/*
Package parse builds a slice of nodes from UTF-8 encoded text documents that have doctags.

Doctags are simple tags that can be written in any text document used to indicate
a named or tagged piece of content. An example document with a doctag could be:

  <{headline}>
  Today's News Stories

When this document is parsed it would result in a DoctagNode with the name "headline" and
value "\nToday's News Stories". A doctag's value is all text between the end of the doctag
and the begining of the next doctag or the end of file, whichever occurs first.
The parsing functions makes the assumption that leading and trailing whitespace
characters are important. It's left for the consumer of the parse results to decide how to
treat leading/trailing whtiespace.

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
*/
package parse

import (
  "os"
  "io"
  "fmt"
  "bufio"
  "log"
  "strings"
  "errors"
  "unicode/utf8"
)

// The default tag prefix and suffix used by the Parse() function.
const (
  DefaultTagPrefix = "<{"
  DefaultTagSuffix = "}>"
)

// The optional Logger to have warnings logged to. The logger is useful
// in finding what might be typos when declaring a doctag in a document.
var (
  Logger *log.Logger
)

// A DoctagNode represents a doctag parsed from a text document.
type DoctagNode struct {
  Name string
  Value string
  Line int
  Column int
}

// Parse parses a text file for doctags using the default prefix and suffix substrings.
// The returned slice contains all parsed DoctagNodes in the order they appear in the document.
func ParseFile(fileName string) ([]*DoctagNode, error) {
  return ParseFileWithPrefixAndSuffix(fileName, DefaultTagPrefix, DefaultTagSuffix);
}

// ParseFileWithPrefixAndSuffix parses a text file for doctags using custom prefix and suffix substrings for doctags.
// The returned slice contains all parsed DoctagNodes in the order they appear in the document.
func ParseFileWithPrefixAndSuffix(fileName string, tagPrefix string, tagSuffix string) ([]*DoctagNode, error) {
  file,err := os.Open(fileName)
  defer file.Close()

  if err == nil {
    return ParseWithPrefixAndSuffix(bufio.NewReader(file), tagPrefix, tagSuffix)
  }

  return nil,err
}

// Parse parses a buffered reader for doctags using the default prefix and suffix substrings.
// The returned slice contains all parsed DoctagNodes in the order they appear in the document.
func Parse(reader *bufio.Reader) ([]*DoctagNode, error) {
  return ParseWithPrefixAndSuffix(reader, DefaultTagPrefix, DefaultTagSuffix);
}

// ParseWithPrefixAndSuffix parses a buffered reader for doctags using custom prefix and suffix substrings for doctags.
// The returned slice contains all parsed DoctagNodes in the order they appear in the document.
func ParseWithPrefixAndSuffix(reader *bufio.Reader, tagPrefix string, tagSuffix string) (doctags []*DoctagNode, err error) {
  if tagPrefix == tagSuffix {
    err = errors.New("Tag prefix and suffix cannot be the same.")
    return
  }
  if len(tagPrefix) == 0 {
    err = errors.New("Tag prefix cannot be the empty string.")
    return
  }
  if len(tagSuffix) == 0 {
    err = errors.New("Tag suffix cannot be the empty string.")
    return
  }

  // The capacity to create text buffers at (i.e. to capture text between doctags).
  const bufferSize = 512
  doctags = make([]*DoctagNode, 0, 50)
  buff := make([]byte, 0, bufferSize)
  line := 1
  column := 0
  var currTag *DoctagNode
  var b byte

  for b,err = reader.ReadByte(); err == nil || err == io.EOF; b,err = reader.ReadByte() {
    var ok bool

    if err == io.EOF {
      if currTag != nil && len(currTag.Name) > 0 {
        // buff is previous tag's value
        currTag.Value = string(buff)
        doctags = append(doctags, currTag)
        currTag = nil
      }
      err = nil
      break
    }

    if utf8.RuneStart(b) {
      column++
    }
    buff = append(buff, b)

    if b == '\n' {
      line++
      column = 0
    }

    if b == tagPrefix[0] {
      if ok,err = consume(reader, tagPrefix); ok {
        if currTag != nil && len(currTag.Name) > 0 {
          // buff is previous tag's value (we don't want the first byte of the prefix)
          currTag.Value = string(buff[:len(buff) - 1])
          doctags = append(doctags, currTag)
          currTag = nil
        } else if currTag != nil {
          warn(line, column, "doctag open encountered but the previous doctag was not closed properly or has no tag name.")
        }

        // Create an empty tag
        currTag = &DoctagNode{Line: line, Column: column}
        // Clear the buffer
        buff = make([]byte, 0, bufferSize)
        // Make sure we take into account the bytes we just consumed
        column += utf8.RuneCount([]byte(tagSuffix)) - 1
      }
    } else if b == tagSuffix[0] && currTag != nil && currTag.Line == line {
      if len(currTag.Name) == 0 {
        if ok,err = consume(reader, tagSuffix); ok {
          // buff is the tag name (we don't want the first byte of the suffix)
          currTag.Name = strings.TrimSpace(string(buff[:len(buff) - 1]))
          // Make sure we take into account the bytes we just consumed
          column += utf8.RuneCount([]byte(tagSuffix)) - 1

          if len(currTag.Name) == 0 {
            warn(line, column, "doctag close encountered but tag name not detected. Skipping doctag.")
          } else {
            // Check to see if we are to skip this tag
            if currTag.Name[0] == '!' {
              warn(line, column, fmt.Sprintf("skipping doctag '%v'", currTag.Name))
              currTag = nil
            }

            // Clear the buffer
            buff = make([]byte, 0, bufferSize)
          }
        }
      } else {
        warn(line, column, "doctag close encountered but the previous doctag was not closed properly or has no tag name.")
      }
    }
  }

  if err != nil {
    err = fmt.Errorf("Line: %v, Column: %v :: %v", line, column, err.Error())
  }

  return
}

// Attempts to consume token from reader.
// Expects the first rune to be already read from the reader.
// In other words the first rune of the token is not re-read or verified.
func consume(reader *bufio.Reader, token string) (ok bool, err error) {
  _,firstRuneSize := utf8.DecodeRuneInString(token)
  size := len(token) - firstRuneSize

  if size <= 0 {
    ok = true
    return
  }

  buff := make([]byte, size)

  if buff,err = reader.Peek(size); string(buff) == token[firstRuneSize:] {
    // Actually consume the bytes
    reader.Read(buff)
    ok = true
  }

  return
}

// Convenient wrapper function that will log a warning message.
func warn(line int, column int, message string) {
  if Logger != nil {
    Logger.Printf("\nLine: %v, Column: %v\n%v\n\n", line, column, message)
  }
}