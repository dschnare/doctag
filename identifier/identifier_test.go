package identifier

import (
  "testing"
)

func TestToGoIdentifier(t *testing.T) {
  if id := ToGoIdentifier("45_aceIn45"); id != "_aceIn45" {
    t.Fatalf("expected %v : got %v", "_aceIn45", id)
  }
  if id := ToGoIdentifier("_ace3&$In45@"); id != "_ace3In45" {
    t.Fatalf("expected %v : got %v", "_ace3In45", id)
  }
  if id := ToGoIdentifier("_ace\u00ff3&$In45@"); id != "_aceÿ3In45" {
    t.Fatalf("expected %v : got %v", "_ace3In45", id)
  }
  if id := ToGoIdentifier("_aceÿ3&$In45@"); id != "_aceÿ3In45" {
    t.Fatalf("expected %v : got %v", "_ace3In45", id)
  }
}