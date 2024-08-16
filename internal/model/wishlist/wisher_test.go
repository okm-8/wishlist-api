package wishlist

import (
	"github.com/google/uuid"
	"testing"
)

func TestWisher(t *testing.T) {
	t.Run("should restore wisher", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sampleId := RestoreWisherId(uid[:])
		sampleWisherName := "wisher name"
		sampleWisherEmail := "wisher@test.com"

		wisher := RestoreWisher(sampleId, sampleWisherName, sampleWisherEmail)

		if wisher.Id() != sampleId {
			t.Error("expected restored wisher id to be equal to original wisher id")
		}

		if wisher.Name() != sampleWisherName {
			t.Error("expected restored wisher name to be equal to original wisher name")
		}

		if wisher.Email() != sampleWisherEmail {
			t.Error("expected restored wisher email to be equal to original wisher email")
		}
	})
}
