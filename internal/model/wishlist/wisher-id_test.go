package wishlist

import (
	"bytes"
	"github.com/google/uuid"
	"testing"
)

func TestWisherId(t *testing.T) {
	t.Run("should restore id", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sample := uid[:]

		wisherId := RestoreWisherId(sample)

		if !bytes.Equal(wisherId.Bytes(), sample) {
			t.Error("expected restored id to be equal to original id")
		}
	})

	t.Run("should parse id", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sample := uid.String()

		wisherId, err := ParseWisherId(sample)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wisherId.String() != sample {
			t.Error("expected parsed id to be equal to original id")
		}
	})
}
