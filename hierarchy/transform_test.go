package hierarchy

import (
  "testing"
  "github.com/dschnare/doctag/parse"
)

func TestTransform(t *testing.T) {
  if doctags,err := parse.ParseFile("./fixtures/nested.txt"); err == nil {
    expected := map[string]interface{}{
      "nums": &([]interface{}{"1\n", "2", "3\n", "4"}),
      "aa": map[string]interface{}{
        "b": (&[]interface{}{
          map[string]interface{}{"name": "Dave\n", "title": "Plumber\n", "age": "36\n"},
          map[string]interface{}{"name": "Max\n", "title": "3D Animator\n", "age": "24\n\n"},
        }),
      },
      "a": map[string]interface{}{
        "b": map[string]interface{}{
          "c": map[string]interface{}{
            "d": map[string]interface{}{
              "e": map[string]interface{}{
                "f": map[string]interface{}{
                  "test": map[string]interface{}{"title": "\nBoom\n\n"},
                },
              },
            },
          },
        },
      },
      "obj": map[string]interface{}{
        "urls": &([]interface{}{"http://google.com1\n", "http://google.com2\n", "http://google.com3\n", "http://google.com4\n\n"}),
      },
    }

    obj,err := Transform(doctags, true)
    if err != nil {
      t.Fatalf("unexpected error encountered : %v", err.Error())
    } else {
      // t.Fatalf("%v", obj)
      testValue(obj, expected, t)
    }
  } else {
    t.Fatalf("unexpected error encountered : %v", err.Error())
  }
}

func testValue(value interface{}, expected interface{}, t *testing.T) {
  switch expected.(type) {
  case string:
    str := expected.(string)
    if _str,ok := value.(string); ok {
      if str != _str {
        t.Fatalf("expected strings to be equal '%v' : got '%v'", str, _str)
      }
    } else {
      t.Fatalf("expected a string type")
    }
  case map[string]interface{}:
    obj := expected.(map[string]interface{})
    if _obj,ok := value.(map[string]interface{}); ok {
      testObj(_obj, obj, t)
    } else {
      t.Fatalf("expected a map type %v", value)
    }
  case *[]interface{}:
    slicePtr := expected.(*[]interface{})
    if _slicePtr,ok := value.(*[]interface{}); ok {
      testSlice(_slicePtr, slicePtr, t)
    } else {
      t.Fatalf("expected a pointer to a slice type")
    }
  }
}

func testObj(obj map[string]interface{}, expected map[string]interface{}, t *testing.T) {
  for k,v := range expected {
    if _v,ok := obj[k]; ok {
      testValue(_v, v, t)
    } else {
      t.Fatalf("expected map to have key '%v'", k)
    }
  }
}

func testSlice(slicePtr *[]interface{}, expectedPtr *[]interface{}, t *testing.T) {
  slice := *slicePtr
  expected := *expectedPtr

  if len(expected) != len(slice) {
    t.Fatalf("expected slices to be same length")
  }

  for k,v := range expected {
    testValue(slice[k], v, t)
  }
}