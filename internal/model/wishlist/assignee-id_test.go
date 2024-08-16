package wishlist

import (
	"bytes"
	"github.com/google/uuid"
	"testing"
)

func TestAssigneeId(t *testing.T) {
	t.Run("should restore id", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sample := uid[:]

		assigneeId := RestoreAssigneeId(sample)

		if !bytes.Equal(assigneeId.Bytes(), sample) {
			t.Error("expected restored id to be equal to original id")
		}
	})

	t.Run("should parse id", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sample := uid.String()

		assigneeId, err := ParseAssigneeId(sample)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if assigneeId.String() != sample {
			t.Error("expected parsed id to be equal to original id")
		}
	})
}
