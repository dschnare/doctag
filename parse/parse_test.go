package parse

import (
  "testing"
  "log"
  "os"
)

func TestParse_SamePrefixAndSuffix(t *testing.T) {
  doctags,err := ParseFileWithPrefixAndSuffix("./fixtures/empty.txt", "!", "!")

  if len(doctags) > 0 {
    t.Fatalf("expected no doc tags to be found: found %v doc tags", len(doctags))
  }

  if err == nil {
    t.Fatalf("expected error")
  }
}

func TestParse_Empty(t *testing.T) {
  doctags,err := ParseFile("./fixtures/empty.txt")

  if len(doctags) > 0 {
    t.Fatalf("expected no doc tags to be found: found %v doc tags", len(doctags))
  }

  if err != nil {
    t.Fatalf("expected no error: %v", err.Error())
  }
}

func TestParse_NoTags(t *testing.T) {
  doctags,err := ParseFile("./fixtures/no_tags.txt")
  
  if len(doctags) > 0 {
    t.Fatalf("expected no doc tags to be found: found %v doc tags", len(doctags))
  }

  if err != nil {
    t.Fatalf("expected no error: %v", err.Error())
  }
}

func TestParse_BeginingOfFile(t *testing.T) {
  var expected = []*DoctagNode {
    &DoctagNode{
      Name: "Headline",
      Value: "\nThis is a headline",
      Line: 1,
      Column: 1,
    },
  }

  doctags,err := ParseFile("./fixtures/begining_of_file.txt")

  if err != nil {
    t.Fatalf("expected no error: %v", err.Error())
  }

  testSlice(doctags, expected, t)
}

func TestParse_Complex(t *testing.T) {
  var expected = []*DoctagNode {
    &DoctagNode{
      Name: "Headline",
      Value: "\nThis is a headline\n\n",
      Line: 1,
      Column: 1,
    },
    &DoctagNode{
      Name: "Headline2",
      Value: "\nThis is a headline\n\n  ",
      Line: 4,
      Column: 13,
    },
    &DoctagNode{
      Name: "Headline3",
      Value: "\nThis is a headline\nThis is another line\n\n\n",
      Line: 7,
      Column: 3,
    },
    &DoctagNode{
      Name: "Headline4",
      Value: " Headline 4 ",
      Line: 12,
      Column: 1,
    },
    &DoctagNode{
      Name: "Headline4 / Link",
      Value: " Headline 4 Link\n",
      Line: 12,
      Column: 26,
    },
    &DoctagNode{
      Name: "Broke",
      Value: "}>Boom\n",
      Line: 13,
      Column: 1,
    },
    &DoctagNode{
      Name: "hello",
      Value: "hello",
      Line: 15,
      Column: 1,
    },
    &DoctagNode{
      Name: "hi",
      Value: "\nhi",
      Line: 16,
      Column: 1,
    },
  }

  Logger = log.New(os.Stderr, "doctag warning: ", log.Lshortfile)
  doctags,err := ParseFile("./fixtures/complex.txt")
  Logger = nil

  if err != nil {
    t.Fatalf("expected no error: %v", err.Error())
  }

  testSlice(doctags, expected, t)
}

func TestParse_ComplexWithPrefixAndSuffix(t *testing.T) {
  var expected = []*DoctagNode {
    &DoctagNode{
      Name: "Headline",
      Value: "\nThis is a headline\n\n",
      Line: 1,
      Column: 1,
    },
    &DoctagNode{
      Name: "Headline2",
      Value: "\nThis is a headline\n\n  ",
      Line: 4,
      Column: 13,
    },
    &DoctagNode{
      Name: "Headline3",
      Value: "\nThis is a headline\nThis is another line\n\n\n",
      Line: 7,
      Column: 3,
    },
    &DoctagNode{
      Name: "Headline4",
      Value: " Headline 4 ",
      Line: 12,
      Column: 1,
    },
    &DoctagNode{
      Name: "Headline4 / Link",
      Value: " Headline 4 Link\n",
      Line: 12,
      Column: 28,
    },
    &DoctagNode{
      Name: "Broke",
      Value: " --Boom",
      Line: 13,
      Column: 1,
    },
  }

  doctags,err := ParseFileWithPrefixAndSuffix("./fixtures/complex_with_prefix_and_suffix.txt", "-- ", " --")

  if err != nil {
    t.Fatalf("expected no error: %v", err.Error())
  }

  testSlice(doctags, expected, t)
}

func testSlice(doctags []*DoctagNode, expected []*DoctagNode, t *testing.T) {
  doctagsLen := len(doctags)

  if doctagsLen != len(expected) {
    t.Fatalf("expected the document to contain '%v' tags : only found %v tags", len(expected), doctagsLen)
  }

  for k,expectedDoctag := range expected {
    doctag := doctags[k]
    if doctag.Name != expectedDoctag.Name {
      t.Fatalf("expected the document to contain '%v' tag : got '%v'", expectedDoctag.Name, doctag.Name)
    }
    if doctag.Value != expectedDoctag.Value {
      t.Fatalf("expected the document to contain '%v:%v' tag : got '%v'", expectedDoctag.Name, []byte(expectedDoctag.Value), []byte(doctag.Value))
    }
    if doctag.Line != expectedDoctag.Line {
      t.Fatalf("expected tag '%v' to be on line %v : got '%v'", k, doctag.Line, doctag.Line)
    }
    if doctag.Column != expectedDoctag.Column {
      t.Fatalf("expected tag '%v' to be on column %v : got '%v'", k, doctag.Column, doctag.Column)
    }
  }
}