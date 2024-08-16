package wishlist

type Wishlist struct {
	id     Id
	wisher WisherId
	wishes []*Wish
	name   string
	hidden bool
}

func NewWishlist(wisher WisherId) *Wishlist {
	return &Wishlist{
		id:     newId(),
		wisher: wisher,
		wishes: make([]*Wish, 0),
	}
}

func RestoreWishlist(id Id, wisher WisherId, wishes []*Wish) *Wishlist {
	return &Wishlist{
		id:     id,
		wisher: wisher,
		wishes: wishes,
	}
}

func (wishlist *Wishlist) Id() Id {
	return wishlist.id
}

func (wishlist *Wishlist) Wisher() WisherId {
	return wishlist.wisher
}

func (wishlist *Wishlist) Wishes() []*Wish {
	return wishlist.wishes
}

func (wishlist *Wishlist) Hidden() bool {
	return wishlist.hidden
}

func (wishlist *Wishlist) Hide() {
	wishlist.hidden = true
}

func (wishlist *Wishlist) AddWish(name, description string) *Wish {
	wish := NewWish(name, description)
	wishlist.wishes = append(wishlist.wishes, wish)

	return wish
}

func (wishlist *Wishlist) MoveWish(wish *Wish, destination *Wishlist) {
	for i, w := range wishlist.wishes {
		if w == wish {
			wishlist.wishes = append(wishlist.wishes[:i], wishlist.wishes[i+1:]...)
			destination.wishes = append(destination.wishes, wish)

			return
		}
	}
}
