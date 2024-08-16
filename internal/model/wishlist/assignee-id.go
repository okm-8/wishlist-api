package wishlist

import "github.com/google/uuid"

type AssigneeId uuid.UUID

var NilAssigneeId = AssigneeId(uuid.Nil)

func RestoreAssigneeId(id []byte) AssigneeId {
	return AssigneeId(id)
}

func ParseAssigneeId(id string) (AssigneeId, error) {
	uid, err := uuid.Parse(id)

	if err != nil {
		return NilAssigneeId, err
	}

	return AssigneeId(uid), nil
}

func (id AssigneeId) String() string {
	return uuid.UUID(id).String()
}

func (id AssigneeId) Bytes() []byte {
	return id[:]
}
