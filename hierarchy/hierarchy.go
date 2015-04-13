/*
Package hierarchy implements a doctag transformer
that transforms a list of doctags into a map hierarchy.
The name of each doctag must use a separator character
to describe a hierarchy. Doctags can also support adding
to a slice by prefixing the doctag name with "#".
If no separator character is found then the hierarchy will be flat.

Example:

Doctag document:
  <{ page/title }>This is the page title<{!}>
  <{ page/keywords }>awesome,stuff,aboutpeople<{!}>

Map hierarchy:

  map{
     "page": map{
        "title": "This is the page title",
        "keywords": "awesome,stuff,aboutpeople",
     },
  }

By prefixing doctags with the '#' character it indicates to the transformer
that the implicitly created map (or string value) will be appended to a slice.
Any doctags that add keys to the implicitly created map that has been added to
a slice will be set on the last map added to the slice. Use '#' in doctags to
explicitly append new objects to slices.

Example:

Doctag document:
  <{ page/title }>This is the page title<{!}>
  <{ ! These doctags will append each string value to a slice indexed by "keywords" }>
  <{ page/#keywords }>awesome<{!}>
  <{ page/#keywords }>stuff<{!}>
  <{ page/#keywords }>aboutpeople<{!}>

  <{ page/content }>
  Some stuff about people

  <{ ! These doctags will append the implicitly created map to a slice indexed by "links" }>
  <{ page/#links/rel }>alternate<{!}>
  <{ page/links/href }>http://my.domain.com/alternate.html<{!}>
  <{ page/#links/rel }>next<{!}>
  <{ page/links/href }>http://my.domain.com/next.html<{!}>
  <{ page/#links/rel }>prev<{!}>
  <{ page/links/href }>http://my.domain.com/prev.html<{!}>

Map hierarchy:

  map{
     "page": map{
        "title": "This is the page title",
        "keywords": ["awesome", "stuff", "aboutpeople"],
        "links": [
          map{
            "rel": "alternate",
            "href": "http://my.domain.com/alternate.html",
          },
          map{
            "rel": "next",
            "href": "http://my.domain.com/next.html",
          },
          map{
            "rel": "prev",
            "href": "http://my.domain.com/prev.html",
          },
        ],
     },
  }
*/
package hierarchy

import (
  "fmt"
  "strings"
  "unicode"
  "github.com/dschnare/doctag/parse"
  "github.com/dschnare/doctag/identifier"
)

// DefaultSeparator is a constant for the default character used to delimit separate doctag names.
const DefaultSeparator = '/'

// Transform transforms a slice of DoctagNodes into a hierarchical map that represents a JSON object.
// The default separater character will be used when parsing hierarchical doctags.
func Transform(doctags []*parse.DoctagNode, jsonKeysToIdentifiers bool) (map[string]interface{}, error) {
  return TransformWithSeparator(doctags, jsonKeysToIdentifiers, DefaultSeparator)
}

// TransformWithSeparator transforms a slice of DoctagNodes with a specific doctag separator character into a hierarchical map that represents a JSON object.
func TransformWithSeparator(doctags []*parse.DoctagNode, jsonKeysToIdentifiers bool, separator rune) (map[string]interface{}, error) {
  var err error
  object := make(map[string]interface{})

  for _,doctag := range doctags {
    pathNames := getPathNames(doctag.Name, separator)
    last := len(pathNames) - 1
    var o interface{} = object

    for g,pathName := range pathNames {
      if pathName == "#" {
        return nil,fmt.Errorf("Line: %v, Column: %v :: Path cannot equal '#'", doctag.Line, doctag.Column)
        break
      }
      if jsonKeysToIdentifiers {
        // When we convert to an identifier we prserve the "#" prefix.
        // The prefix is trimmed when actually saving to the map.
        pathName = identifier.ToIdentifierFunc(pathName, identifierValidRuneFunc)
        if len(pathName) == 0 {
          return nil,fmt.Errorf("Line: %v, Column: %v :: After converting to an identifier, path is empty", doctag.Line, doctag.Column)
          break
        }
      }
      if g == last {
        // setKey(o, pathName, doctag.Value)
        resolveWithValue(o, pathName, doctag.Value)
      } else {
        o = resolve(o, pathName)
      }
    }
  }

  return object,err
}

// Preseve the "#" prefix, otherwise same as ToGoIdentifier().
func identifierValidRuneFunc(r rune, idLen int) bool {
  if idLen == 0 {
    return r == '_' || r == '#' || unicode.IsLetter(r)
  } else if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
    return true
  }
  return false
}

// Takes a hierarchical doctag name and splits it into separate path names.
func getPathNames(tagName string, separator rune) []string {
  return strings.FieldsFunc(tagName, func (r rune) bool {
    return unicode.IsSpace(r) || r == separator
  })
}

// Resolve a key on the specified object o.
func resolve(o interface{}, key string) (interface{}) {
  value := make(map[string]interface{})
  return resolveWithValue(o, key, value)
}
func resolveWithValue(o interface{}, key string, value interface{}) (interface{}) {
  if value == nil {
    value = make(map[string]interface{})
  }

  switch o.(type) {
  case map[string]interface{}:
    m := o.(map[string]interface{})
    return resolveFromMapWithValue(m, key, value)
  case *[]interface{}:
    seqPtr := o.(*[]interface{})
    return resolveFromSliceWithValue(seqPtr, key, value)
  }

  return o
}

func resolveFromMapWithValue(m map[string]interface{}, key string, value interface{}) (interface{}) {
  // The key must refer refer to a slice with a new map appended.
  if strings.HasPrefix(key, "#") {
    key = key[1:]
    // If the key exists then we ensure it refers to a slice, replacing the key if required.
    if v,ok := m[key]; ok {
      if seqPtr,ok := v.(*[]interface{}); ok {
        seq := *seqPtr
        seq = append(seq, value)
        *seqPtr = seq
        return value
      } else {
        temp := make([]interface{}, 2, 50)
        temp[0] = v
        temp[1] = value
        m[key] = &temp
        return value
      }
    // If the key does not exist then we create a slice and set the key.
    } else {
      temp := make([]interface{}, 1, 50)
      temp[0] = value
      m[key] = &temp
      return value
    }
  // The key must refer to a map.
  } else {
    // If the key does not exist then we create a map and set the key.
    if v,ok := m[key]; !ok {
      m[key] = value
    // If the key exists and it's a string then we create a map and set the key.
    } else if _,ok := v.(string); ok {
      m[key] = value
    }
  }

  // Return the value referred to by the key.
  return m[key]
}

func resolveFromSliceWithValue(seqPtr *[]interface{}, key string, value interface{}) (interface{}) {
  seq := *seqPtr

  // Sequences (i.e. slices) are treated like the following:
  // - if the key is prefixed with "#" then the key must
  //    refer to another slice on the last map in the sequence.
  // - if the key is not prefixed with "#" then the key must
  //    refer to a map of on the last map in the sequence.

  // If the sequence has not items.
  if len(seq) == 0 {
    // The key must refer to a slice with a new map appended.
    if strings.HasPrefix(key, "#") {
      key = key[1:]
      // Create a map (i.e. the new last item in the sequence) and a slice,
      // set the key on the map to the slice and append the map to the sequence.
      obj := make(map[string]interface{})
      temp := make([]interface{}, 1, 50)
      temp[0] = value
      obj[key] = &temp
      seq = append(seq, obj)
      *seqPtr = seq
      return value
    // The key must refer to a map.
    } else {
      // Create a map (i.e. the new last item in the sequence) set the key
      // to another newly created map and append the first map to the sequence.
      obj := make(map[string]interface{})
      obj[key] = value
      seq = append(seq, obj)
      *seqPtr = seq
      // Return the second newly created map.
      return value
    }
  // If the sequence has items then we grab the last item and recursively resolve.
  // We can recursivly call resolve() because we'll never have slices of slices of slices ...
  } else {
    lastItem := seq[len(seq) - 1]
    return resolveWithValue(lastItem, key, value)
  }
}