package wishlist

import (
	"bytes"
	"testing"
)

func TestId(t *testing.T) {
	t.Run("should restore id", func(t *testing.T) {
		t.Parallel()

		sample := newId().Bytes()

		wishId := RestoreId(sample)

		if !bytes.Equal(wishId.Bytes(), sample) {
			t.Error("expected restored id to be equal to original id")
		}
	})

	t.Run("should parse id", func(t *testing.T) {
		t.Parallel()

		sample := newId().String()

		wishId, err := ParseId(sample)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wishId.String() != sample {
			t.Error("expected parsed id to be equal to original id")
		}
	})
}
