package wishlist

import (
	"github.com/google/uuid"
	"slices"
	"testing"
)

func TestWishlist(t *testing.T) {
	t.Run("should restore wishlist", func(t *testing.T) {
		t.Parallel()

		wisherUid := uuid.New()
		sampleWisherId := RestoreWisherId(wisherUid[:])
		sampleWishlistId := newId()
		sampleWishlistName := "wishlist name"
		sampleWishlistDescription := "wishlist description"
		sampleWishes := []*Wish{
			NewWish("wish name", "wish description"),
			NewWish("wish2 name", "wish2 description"),
		}

		wishlist := RestoreWishlist(
			sampleWishlistId,
			sampleWisherId,
			sampleWishlistName,
			sampleWishlistDescription,
			sampleWishes,
			true,
		)

		if wishlist.Id() != sampleWishlistId {
			t.Error("expected restored wishlist id to be equal to original wishlist id")
		}

		if wishlist.Name() != sampleWishlistName {
			t.Error("expected restored wishlist name to be equal to original wishlist name")
		}

		if wishlist.Description() != sampleWishlistDescription {
			t.Error("expected restored wishlist description to be equal to original wishlist description")
		}

		if wishlist.Wisher() != sampleWisherId {
			t.Error("expected restored wishlist wisher to be equal to original wishlist wisher")
		}

		if wishlist.Hidden() != true {
			t.Error("expected restored wishlist hidden to be equal to original wishlist hidden")
		}

		if len(wishlist.Wishes()) != len(sampleWishes) {
			t.Error("expected restored wishlist wishes to be equal to original wishlist wishes")
		}

		if !slices.EqualFunc(wishlist.Wishes(), sampleWishes, func(a, b *Wish) bool {
			return a.Id() == b.Id()
		}) {
			t.Error("expected restored wishlist wishes to be equal to original wishlist wishes")
		}
	})

	t.Run("should add wish", func(t *testing.T) {
		t.Parallel()

		wisherUid := uuid.New()
		sampleWisherId := RestoreWisherId(wisherUid[:])

		wishlist := NewWishlist(sampleWisherId, "wishlist name", "wishlist description")

		wish := wishlist.AddWish("wish name", "wish description")

		if len(wishlist.Wishes()) != 1 {
			t.Error("expected wishlist to have one wish")
		}

		if wishlist.Wishes()[0] != wish {
			t.Error("expected added wish to be in wishlist")
		}
	})

	t.Run("should hide wishlist", func(t *testing.T) {
		t.Parallel()

		wisherUid := uuid.New()
		sampleWisherId := RestoreWisherId(wisherUid[:])

		wishlist := NewWishlist(sampleWisherId, "wishlist name", "wishlist description")

		wishlist.Hide()

		if !wishlist.Hidden() {
			t.Error("expected wishlist to be hidden")
		}
	})

	t.Run("should show wishlist", func(t *testing.T) {
		t.Parallel()

		wisherUid := uuid.New()
		sampleWisherId := RestoreWisherId(wisherUid[:])

		wishlist := NewWishlist(sampleWisherId, "wishlist name", "wishlist description")

		wishlist.Show()

		if wishlist.Hidden() {
			t.Error("expected wishlist to be shown")
		}
	})

	t.Run("should move wish to wishlist", func(t *testing.T) {
		t.Parallel()

		wisherUid := uuid.New()
		sampleWisherId := RestoreWisherId(wisherUid[:])

		wishlist := NewWishlist(sampleWisherId, "wishlist name", "wishlist description")
		wish := wishlist.AddWish("wish name", "wish description")

		otherWishList := NewWishlist(sampleWisherId, "wishlist name", "wishlist description")

		wishlist.MoveWish(wish.Id(), otherWishList)

		if len(wishlist.Wishes()) != 0 {
			t.Error("expected wish to be removed from wishlist")
		}

		if len(otherWishList.Wishes()) != 1 {
			t.Error("expected wish to be added to other wishlist")
		}

		if otherWishList.Wishes()[0] != wish {
			t.Error("expected wish to be added to other wishlist")
		}
	})
}
