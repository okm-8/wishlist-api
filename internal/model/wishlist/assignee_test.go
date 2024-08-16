package wishlist

import (
	"github.com/google/uuid"
	"testing"
)

func TestAssignee(t *testing.T) {
	t.Run("should restore assignee", func(t *testing.T) {
		t.Parallel()

		uid := uuid.New()
		sampleId := RestoreAssigneeId(uid[:])
		sampleAssigneeName := "assignee name"
		sampleAssigneeEmail := "assignee@test.com"

		assignee := RestoreAssignee(sampleId, sampleAssigneeName, sampleAssigneeEmail)

		if assignee.Id() != sampleId {
			t.Error("expected restored assignee id to be equal to original assignee id")
		}

		if assignee.Name() != sampleAssigneeName {
			t.Error("expected restored assignee name to be equal to original assignee name")
		}

		if assignee.Email() != sampleAssigneeEmail {
			t.Error("expected restored assignee email to be equal to original assignee email")
		}
	})
}
