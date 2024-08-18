package wishlist

import (
	"errors"
	"github.com/google/uuid"
	"testing"
)

func makeAssignee() *Assignee {
	assigneeUid := uuid.New()
	return RestoreAssignee(
		RestoreAssigneeId(assigneeUid[:]),
		"assignee name",
		"assignee email",
	)
}

func TestWish(t *testing.T) {
	t.Run("should restore wish", func(t *testing.T) {
		t.Parallel()

		sampleId := newWishId()
		sampleName := "wish name"
		sampleDescription := "wish description"

		sampleAssignee := makeAssignee()

		wish := RestoreWish(
			sampleId,
			sampleName,
			sampleDescription,
			true,
			true,
			sampleAssignee,
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

		if wish.Assignee() != sampleAssignee {
			t.Error("expected restored wish assignee to be equal to original wish assignee")
		}
	})

	t.Run("should wish be fulfilled", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())

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
		_ = wish.Promise(makeAssignee())
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

		err := wish.Promise(makeAssignee())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !wish.Promised() {
			t.Error("expected wish to be promised")
		}
	})

	t.Run("should fail to promise promised wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())

		err := wish.Promise(makeAssignee())

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
		_ = wish.Promise(makeAssignee())

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
		_ = wish.Promise(makeAssignee())
		_ = wish.Fulfill()

		err := wish.Renege()

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyFulfilled) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyFulfilled, err)
		}
	})

	t.Run("should move wish to another wishlist", func(t *testing.T) {
		t.Parallel()

		sampleWisher := makeWisher()

		wishlist := New(sampleWisher, "wishlist name", "wishlist description")
		wish := wishlist.AddWish("wish name", "wish description")

		otherWishList := New(sampleWisher, "wishlist name", "wishlist description")

		wish.Move(otherWishList)

		if len(wishlist.Wishes()) != 0 {
			t.Error("expected wish to be removed from wishlist")
		}

		if len(otherWishList.Wishes()) != 1 {
			t.Error("expected wish to be added to other wishlist")
		}
	})

	t.Run("should rename wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")

		newName := "new wish name"
		err := wish.Rename(newName)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wish.Name() != newName {
			t.Error("expected wish name to be changed")
		}
	})

	t.Run("should fail to rename fulfilled wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())
		_ = wish.Fulfill()

		err := wish.Rename("new wish name")

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyFulfilled) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyFulfilled, err)
		}
	})

	t.Run("should fail to rename promised wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())

		err := wish.Rename("new wish name")

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyPromised) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyPromised, err)
		}
	})

	t.Run("should update wish description", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")

		newDescription := "new wish description"
		err := wish.UpdateDescription(newDescription)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wish.Description() != newDescription {
			t.Error("expected wish description to be changed")
		}
	})

	t.Run("should fail to update description of fulfilled wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())
		_ = wish.Fulfill()

		err := wish.UpdateDescription("new wish description")

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyFulfilled) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyFulfilled, err)
		}
	})

	t.Run("should fail to update description of promised wish", func(t *testing.T) {
		t.Parallel()

		wish := NewWish("wish name", "wish description")
		_ = wish.Promise(makeAssignee())

		err := wish.UpdateDescription("new wish description")

		if err == nil {
			t.Error("expected error to be returned")
		}

		if !errors.Is(err, ErrWishAlreadyPromised) {
			t.Errorf("expected error to be %v, got %v", ErrWishAlreadyPromised, err)
		}
	})
}
