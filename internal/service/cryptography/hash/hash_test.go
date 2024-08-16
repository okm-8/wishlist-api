package hash

import (
	"testing"
)

func TestApiSignature(t *testing.T) {
	t.Run("should generate valid hash", func(t *testing.T) {
		ctx := newContextStub("secret")
		data := []byte("data")

		signature := Hash(ctx, data)

		if !Verify(ctx, data, signature) {
			t.Errorf("unexpected hash verification failure")
		}
	})
}
