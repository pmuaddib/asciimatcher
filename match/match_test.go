package match

import "testing"

func TestParseIncomingPattern(t *testing.T) {
    _, err := ParseIncomingPattern("b.txt")
    if err != nil {
        t.Error(err)
    }
}

func TestFindByPattern(t *testing.T) {
    p, err := ParseIncomingPattern("b.txt")
    n, err :=FindByPattern(p, "source.txt")
    if err != nil {
        t.Error(err)
    }
    expected := 3
    if n != expected {
        t.Errorf("Expected: %d, actual: %d\n", expected, n)
    }
}
