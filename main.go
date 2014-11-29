/*
Package main implements a command line utility that exposes a
doctag parser and hierarchy transformer. Use `doctag --help` for
details about each command line argument.

This utilty supports standard in piping of a doctag file contents.
If the `--output` argument is not specified then output will be 
piped to standard out.
*/
package main

import (
  "flag"
  "fmt"
  "os"
  "log"
  "bufio"
  "bytes"
  "encoding/json"
  "strings"
  "unicode/utf8"
  "github.com/dschnare/doctag/parse"
  "github.com/dschnare/doctag/identifier"
  "github.com/dschnare/doctag/hierarchy"
)

var (
  fileName string
  tagSeparator rune
  // Flags
  tagPrefix string
  tagSuffix string
  tagSeparatorStr string
  output string
  help bool
  warn bool
  prettyPrint bool
  hierarchical bool
  trim bool
)

func usage() {
  fmt.Fprintf(os.Stderr, "Usage: %s {file path} | %s [help|/?]\n", os.Args[0], os.Args[0])
  flag.PrintDefaults()
}

func init() {
  const (
    helpDefault = false
    helpUsage = "Show the help message."
    prettyPrintDefault = false
    prettyPrintUsage = "Print JSON result with indentation."
    warnDefault = false
    warnUsage = "Print warning messages."
    hierarchicalDefault = false
    hierarchicalUsage = "Converts the flat doctag tree into a nested JSON object."
    trimDefault = false
    trimUsage = "Trim the leading and trailing whitespace from all doctag values."
    tagPrefixDefault = parse.DefaultTagPrefix
    tagPrefixUsage = "The prefix to use for doc tags."
    tagSuffixDefault = parse.DefaultTagSuffix
    tagSuffixUsage = "The suffix to use for doc tags."
    tagSeparatorDefault = string(hierarchy.DefaultSeparator)
    tagSeparatorUsage = "The separator character to use for hierarchical doc tags."
    outputDefault = ""
    outputUsage = "The output file to write to."
  )

  flag.Usage = usage
  
  flag.BoolVar(&help, "help", helpDefault, helpUsage)

  flag.BoolVar(&prettyPrint, "pretty-print", prettyPrintDefault, prettyPrintUsage)
  flag.BoolVar(&prettyPrint, "pretty", prettyPrintDefault, prettyPrintUsage + " (shorthand)")

  flag.BoolVar(&warn, "warn", warnDefault, warnUsage)

  flag.BoolVar(&hierarchical, "hierarchical", hierarchicalDefault, hierarchicalUsage)
  flag.BoolVar(&hierarchical, "hierarchy", hierarchicalDefault, hierarchicalUsage + " (shorthand)")

  flag.BoolVar(&trim, "trim", trimDefault, trimUsage)

  flag.StringVar(&tagPrefix, "tag-prefix", tagPrefixDefault, tagPrefixUsage)

  flag.StringVar(&tagSuffix, "tag-suffix", tagSuffixDefault, tagSuffixUsage)

  flag.StringVar(&tagSeparatorStr, "tag-separator", tagSeparatorDefault, tagSeparatorUsage)

  flag.StringVar(&output, "output", outputDefault, outputUsage)

  flag.Parse()



  if warn {
    parse.Logger = log.New(os.Stderr, "doctag warning: ", log.Lshortfile)
  }

  if len(tagSeparatorStr) == 0 {
    tagSeparator = hierarchy.DefaultSeparator
  } else {
    tagSeparator,_ = utf8.DecodeRuneInString(tagSeparatorStr)
  }

  if help {
    flag.Usage()
    os.Exit(0)
  } else if len(flag.Args()) == 1 && (flag.Arg(0) == "/?" || flag.Arg(0) == "help") {
    flag.Usage()
    os.Exit(0)
  } else if len(flag.Args()) == 1 {
    fileName = flag.Arg(0)
  } else {
    flag.Usage()
    os.Exit(1)
  }
}

func main() {
  if doctags,err := doParse(); err == nil {
    if writer,err := createWriter(); err == nil {
      if err := doWrite(writer, doctags); err != nil {
        panic(err)
      }
    } else {
      panic(err)
    }
  } else {
    panic(err)
  }
}

func doParse() (doctags []*parse.DoctagNode, err error) {
  if isPiped(os.Stdin) {
    doctags,err = parse.ParseWithPrefixAndSuffix(bufio.NewReader(os.Stdin), tagPrefix, tagSuffix)
  } else {
    doctags,err = parse.ParseFileWithPrefixAndSuffix(fileName, tagPrefix, tagSuffix)
  }

  return
}

func isPiped(file *os.File) bool {
  if info,err := file.Stat(); err == nil {
    return info.Mode() == os.ModeNamedPipe
  }
  return false
}

func createWriter() (*bufio.Writer, error) {
  var writer *bufio.Writer

  if len(output) == 0 || isPiped(os.Stdout) {
    writer = bufio.NewWriter(os.Stdout)
  } else if file,err := os.Create(output); err == nil {
    writer = bufio.NewWriter(file)
  } else {
    return nil,err
  }

  return writer,nil
}

func doWrite(writer *bufio.Writer, doctags []*parse.DoctagNode) (err error) {
  var value interface{}

  for _,doctag := range doctags {
    if !hierarchical {
      // This will remove the separator characters and convert JSON keys to identifiers.
      doctag.Name = identifier.ToGoIdentifier(strings.Replace(doctag.Name, string(tagSeparator), "_", -1))
    }
    if trim {
      doctag.Value = strings.TrimSpace(doctag.Value)
    }
  }

  if value,err = hierarchy.TransformWithSeparator(doctags, hierarchical, tagSeparator); err != nil {
    return
  }

  if prettyPrint {
    if b,err := json.Marshal(value); err == nil {
      var out bytes.Buffer
      if err = json.Indent(&out, b, "", "  "); err == nil {
        out.WriteTo(writer)
        writer.Flush()
      } else {
        return err
      }
    } else {
      return err
    }
  } else {
    jsonEncoder := json.NewEncoder(writer)
    if err := jsonEncoder.Encode(value); err == nil {
      writer.Flush()
    } else {
      return err
    }
  }

  return
}