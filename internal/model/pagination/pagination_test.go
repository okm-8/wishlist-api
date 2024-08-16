package pagination

import "testing"

func TestPagination(t *testing.T) {
	t.Run("should return the correct offset", func(t *testing.T) {
		t.Parallel()

		pagination := New(1, 10)

		if pagination.Offset() != 0 {
			t.Errorf("expected offset to be 0")
		}

		pagination = New(2, 10)

		if pagination.Offset() != 10 {
			t.Errorf("expected offset to be 10")
		}

		pagination = New(3, 10)

		if pagination.Offset() != 20 {
			t.Errorf("expected offset to be 20")
		}

		pagination = New(10, 7)

		if pagination.Offset() != 63 {
			t.Errorf("expected offset to be 63")
		}
	})
}
