package wishlist

import "github.com/google/uuid"

type WisherId uuid.UUID

var NilWisherId = WisherId(uuid.Nil)

func RestoreWisherId(id []byte) WisherId {
	return WisherId(id)
}

func ParseWisherId(id string) (WisherId, error) {
	uid, err := uuid.Parse(id)

	if err != nil {
		return NilWisherId, err
	}

	return WisherId(uid), nil
}

func (id WisherId) String() string {
	return uuid.UUID(id).String()
}

func (id WisherId) Bytes() []byte {
	return id[:]
}
