package token

import (
	"testing"
)

func TestGenerate_Unique(t *testing.T) {
	length := 10
	max := 100000

	exists := make(map[string]bool, max)

	generator := NewGenerator()

	for i := 0; i < max; i++ {
		token, err := generator.Generate(length)
		if err != nil {
			t.Errorf("unable to generate a token: %v", err)
		}

		if _, ok := exists[token]; ok {
			t.Errorf("found duplicated token: %s", token)
		}

		exists[token] = true
	}
}
