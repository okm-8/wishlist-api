package wishlist

import "github.com/google/uuid"

type WishId uuid.UUID

var NilWishId = WishId(uuid.Nil)

func newWishId() WishId {
	return WishId(uuid.New())
}

func RestoreWishId(id []byte) WishId {
	return WishId(id)
}

func ParseWishId(id string) (WishId, error) {
	uid, err := uuid.Parse(id)

	if err != nil {
		return NilWishId, err
	}

	return WishId(uid), nil
}

func (id WishId) String() string {
	return uuid.UUID(id).String()
}

func (id WishId) Bytes() []byte {
	return id[:]
}
