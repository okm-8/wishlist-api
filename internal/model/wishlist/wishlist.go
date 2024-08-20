package wishlist

type Wishlist struct {
	id          Id
	wisher      *Wisher
	wishes      []*Wish
	name        string
	description string
	hidden      bool
}

func New(wisher *Wisher, name, description string) *Wishlist {
	return &Wishlist{
		id:          newId(),
		wisher:      wisher,
		wishes:      make([]*Wish, 0),
		name:        name,
		description: description,
		hidden:      false,
	}
}

func Restore(id Id, wisher *Wisher, name, description string, hidden bool, wishes []*Wish) *Wishlist {
	wishlist := &Wishlist{
		id:          id,
		wisher:      wisher,
		wishes:      wishes,
		name:        name,
		description: description,
		hidden:      hidden,
	}

	for _, wish := range wishes {
		wish.wishlist = wishlist
	}

	return wishlist
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

func (wishlist *Wishlist) Wisher() *Wisher {
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
	wish.wishlist = wishlist

	return wish
}

func (wishlist *Wishlist) Rename(name string) {
	wishlist.name = name
}

func (wishlist *Wishlist) UpdateDescription(description string) {
	wishlist.description = description
}
