package wishlist

import "github.com/google/uuid"

type Id uuid.UUID

var NilId = Id(uuid.Nil)

func newId() Id {
	return Id(uuid.New())
}

func RestoreId(id []byte) Id {
	return Id(id)
}

func ParseId(id string) (Id, error) {
	uid, err := uuid.Parse(id)

	if err != nil {
		return NilId, err
	}

	return Id(uid), nil
}

func (id Id) String() string {
	return uuid.UUID(id).String()
}

func (id Id) Bytes() []byte {
	return id[:]
}
