package wishlist

import (
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestWish(t *testing.T) {
	t.Run("should restore wish", func(t *testing.T) {
		t.Parallel()

		sampleId := newWishId()
		sampleName := "wish name"
		sampleDescription := "wish description"

		assigneeUid := uuid.New()
		sampleAssigneeId := RestoreAssigneeId(assigneeUid[:])

		wish := RestoreWish(
			sampleId,
			sampleName,
			sampleDescription,
			true,
			true,
			sampleAssigneeId,
		)

		if wish.Id() != sampleId {
			t.Error("expected restored wish id to be equal to original wish id")
		}

		if wish.Name() != sampleName {
			t.Error("expected restored wish name to be equal to original wish name")
		}

		if wish.Description() != sampleDescription {
			t.Error("expected restored wish description to be equal to original wish description")
		}

		if !wish.Fulfilled() {
			t.Error("expected restored wish fulfilled to be equal to original wish fulfilled")
		}

		if wish.Assignee() != sampleAssigneeId {
			t.Error("expected restored wish assignee to be equal to original wish assignee")
		}
	})

	t.Run("should wish be fulfilled", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		assigneeUid := uuid.New()
		_ = wish.Promise(RestoreAssigneeId(assigneeUid[:]))

		err := wish.Fulfill()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !wish.Fulfilled() {
			t.Error("expected wish to be fulfilled")
		}
	})

	t.Run("should fail to fulfill fulfilled wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		assigneeUid := uuid.New()
		_ = wish.Promise(RestoreAssigneeId(assigneeUid[:]))
		_ = wish.Fulfill()

		err := wish.Fulfill()

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyFulfilled) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyFulfilled, err)
		}
	})

	t.Run("should fail to fulfill un-promised wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")

		err := wish.Fulfill()

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrorWishNotPromised) {
			t.Errorf("expected error to be %v, got %v", ErrorWishNotPromised, err)
		}
	})

	t.Run("should wish be hidden", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")

		wish.Hide()

		if !wish.Hidden() {
			t.Error("expected wish to be archived")
		}
	})

	t.Run("should wish be shown", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		wish.Hide()
		wish.Show()

		if wish.Hidden() {
			t.Error("expected wish to be shown")
		}
	})

	t.Run("should wish be promised", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		assigneeUid := uuid.New()

		err := wish.Promise(RestoreAssigneeId(assigneeUid[:]))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !wish.Promised() {
			t.Error("expected wish to be promised")
		}
	})

	t.Run("should fail to promise promised wish", func(t *testing.T) {
		t.Parallel()

		assigneeUid := uuid.New()
		assigneeId := RestoreAssigneeId(assigneeUid[:])

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(assigneeId)

		err := wish.Promise(assigneeId)

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyPromised) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyPromised, err)
		}
	})

	t.Run("should wish be reneged", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		assigneeUid := uuid.New()
		_ = wish.Promise(RestoreAssigneeId(assigneeUid[:]))

		err := wish.Renege()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wish.Promised() {
			t.Error("expected wish to be reneged")
		}
	})

	t.Run("should fail to renege un-promised wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")

		err := wish.Renege()

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrorWishNotPromised) {
			t.Errorf("expected error to be %v, got %v", ErrorWishNotPromised, err)
		}
	})

	t.Run("should fail to renege fulfilled wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		assigneeUid := uuid.New()
		_ = wish.Promise(RestoreAssigneeId(assigneeUid[:]))
		_ = wish.Fulfill()

		err := wish.Renege()

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyFulfilled) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyFulfilled, err)
		}
	})
}
