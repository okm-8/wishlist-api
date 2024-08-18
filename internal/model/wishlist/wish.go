package wishlist

import "errors"

var (
	ErrWishAlreadyFulfilled = errors.New("wish already fulfilled")
	ErrWishAlreadyPromised  = errors.New("wish already promised")
	ErrorWishNotPromised    = errors.New("wish not promised")
)

type Wish struct {
	id          WishId
	name        string
	description string
	fulfilled   bool
	hidden      bool
	wishlist    *Wishlist
	assignee    *Assignee
}

func NewWish(name, description string) *Wish {
	return &Wish{
		id:          newWishId(),
		name:        name,
		description: description,
		fulfilled:   false,
		hidden:      false,
		wishlist:    nil,
		assignee:    nil,
	}
}

func RestoreWish(
	id WishId,
	name string,
	description string,
	fulfilled bool,
	hidden bool,
	assignee *Assignee,
) *Wish {
	return &Wish{
		id:          id,
		name:        name,
		description: description,
		fulfilled:   fulfilled,
		hidden:      hidden,
		wishlist:    nil,
		assignee:    assignee,
	}
}

func (wish *Wish) Id() WishId {
	return wish.id
}

func (wish *Wish) Name() string {
	return wish.name
}

func (wish *Wish) Description() string {
	return wish.description
}

func (wish *Wish) Fulfilled() bool {
	return wish.fulfilled
}

func (wish *Wish) Hidden() bool {
	return wish.hidden
}

func (wish *Wish) Wishlist() *Wishlist {
	return wish.wishlist
}

func (wish *Wish) Assignee() *Assignee {
	return wish.assignee
}

func (wish *Wish) Promised() bool {
	return wish.assignee != nil
}

func (wish *Wish) Rename(name string) error {
	if wish.fulfilled {
		return ErrWishAlreadyFulfilled
	}

	if wish.Promised() {
		return ErrWishAlreadyPromised
	}

	wish.name = name

	return nil
}

func (wish *Wish) UpdateDescription(description string) error {
	if wish.Fulfilled() {
		return ErrWishAlreadyFulfilled
	}

	if wish.Promised() {
		return ErrWishAlreadyPromised
	}

	wish.description = description

	return nil
}

func (wish *Wish) Fulfill() error {
	if wish.Fulfilled() {
		return ErrWishAlreadyFulfilled
	}

	if !wish.Promised() {
		return ErrorWishNotPromised
	}

	wish.fulfilled = true

	return nil
}

func (wish *Wish) Hide() {
	wish.hidden = true
}

func (wish *Wish) Show() {
	wish.hidden = false
}

func (wish *Wish) Promise(assignee *Assignee) error {
	if wish.Promised() {
		return ErrWishAlreadyPromised
	}

	if wish.Fulfilled() {
		// actually, by logic, it is impossible to fulfill a wish that is not promised yet
		// also, it is impossible to dismiss a wish that is already fulfilled.
		// we keep this check for the sake of consistency
		return ErrWishAlreadyFulfilled
	}

	wish.assignee = assignee

	return nil
}

func (wish *Wish) Renege() error {
	if wish.Fulfilled() {
		return ErrWishAlreadyFulfilled
	}

	if !wish.Promised() {
		return ErrorWishNotPromised
	}

	wish.assignee = nil

	return nil
}

func (wish *Wish) Move(to *Wishlist) {
	currentWishlist := wish.wishlist
	if currentWishlist == to {
		return
	}

	for index, _wish := range currentWishlist.wishes {
		if _wish == wish {
			currentWishlist.wishes = append(currentWishlist.wishes[:index], currentWishlist.wishes[index+1:]...)
			break
		}
	}

	exists := false
	for _, _wish := range to.wishes {
		if _wish == wish {
			exists = true
			break
		}
	}

	if !exists {
		to.wishes = append(to.wishes, wish)
	}

	wish.wishlist = to
}
