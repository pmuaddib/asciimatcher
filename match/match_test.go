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

func BenchmarkParseIncomingPattern(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, err := ParseIncomingPattern("b.txt")
        if err != nil {
            b.Error(err)
        }
    }
}

func BenchmarkFindByPattern(b *testing.B) {
    p, _ := ParseIncomingPattern("b.txt")
    for i := 0; i < b.N; i++ {
        _, err := FindByPattern(p, "source.txt")
        if err != nil {
            b.Error(err)
        }
    }
}
