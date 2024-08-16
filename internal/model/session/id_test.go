package session

import "testing"

func TestId(t *testing.T) {
	t.Run("should parse id", func(t *testing.T) {
		t.Parallel()

		id := newId()

		parsedId, err := ParseId(id.String())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if parsedId.String() != id.String() {
			t.Errorf("expected parsed id to be equal to original id")
		}
	})

	t.Run("should restore id", func(t *testing.T) {
		t.Parallel()

		original := newId()

		restored := RestoreId(original.Bytes())

		if restored != original {
			t.Errorf("expected restored id to be equal to original id")
		}
	})
}
