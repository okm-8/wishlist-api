package wishlist

type Wishlist struct {
	id          Id
	wisher      WisherId
	wishes      []*Wish
	name        string
	description string
	hidden      bool
}

func NewWishlist(wisher WisherId, name, description string) *Wishlist {
	return &Wishlist{
		id:          newId(),
		name:        name,
		description: description,
		wisher:      wisher,
		wishes:      make([]*Wish, 0),
		hidden:      false,
	}
}

func RestoreWishlist(id Id, wisher WisherId, name, description string, wishes []*Wish, hidden bool) *Wishlist {
	return &Wishlist{
		id:          id,
		wisher:      wisher,
		wishes:      wishes,
		name:        name,
		description: description,
		hidden:      hidden,
	}
}

func (wishlist *Wishlist) Id() Id {
	return wishlist.id
}

func (wishlist *Wishlist) Name() string {
	return wishlist.name
}

func (wishlist *Wishlist) Description() string {
	return wishlist.description
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

func (wishlist *Wishlist) Show() {
	wishlist.hidden = false
}

func (wishlist *Wishlist) AddWish(name, description string) *Wish {
	wish := NewWish(name, description)
	wishlist.wishes = append(wishlist.wishes, wish)

	return wish
}

func (wishlist *Wishlist) MoveWish(wishId WishId, destination *Wishlist) {
	for index, wish := range wishlist.wishes {
		if wish.Id() == wishId {
			wishlist.wishes = append(wishlist.wishes[:index], wishlist.wishes[index+1:]...)
			destination.wishes = append(destination.wishes, wish)

			return
		}
	}
}
