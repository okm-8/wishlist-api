package wishlist

import (
	"bytes"
	"testing"
)

func TestWishId(t *testing.T) {
	t.Run("should restore id", func(t *testing.T) {
		t.Parallel()

		sample := newWishId().Bytes()

		wishId := RestoreWishId(sample)

		if !bytes.Equal(wishId.Bytes(), sample) {
			t.Error("expected restored id to be equal to original id")
		}
	})

	t.Run("should parse id", func(t *testing.T) {
		t.Parallel()

		sample := newWishId().String()

		wishId, err := ParseWishId(sample)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wishId.String() != sample {
			t.Error("expected parsed id to be equal to original id")
		}
	})
}
