package wishlist

import "errors"

var (
	ErrWishAlreadyFulfilled = errors.New("wish already fulfilled")
	ErrWishAlreadyArchived  = errors.New("wish already archived")
	ErrWishAlreadyPromised  = errors.New("wish already promised")
	ErrorWishNotPromised    = errors.New("wish not promised")
)

type Wish struct {
	id          WishId
	name        string
	description string
	fulfilled   bool
	archived    bool
	assignee    AssigneeId
}

func NewWish(name, description string) *Wish {
	return &Wish{
		id:          newWishId(),
		name:        name,
		description: description,
		fulfilled:   false,
		archived:    false,
		assignee:    NilAssigneeId,
	}
}

func RestoreWish(
	id WishId,
	name string,
	description string,
	fulfilled bool,
	archived bool,
	assignee AssigneeId,
) *Wish {
	return &Wish{
		id:          id,
		name:        name,
		description: description,
		fulfilled:   fulfilled,
		archived:    archived,
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

func (wish *Wish) Archived() bool {
	return wish.archived
}

func (wish *Wish) Assignee() AssigneeId {
	return wish.assignee
}

func (wish *Wish) Fulfill() error {
	if wish.fulfilled {
		return ErrWishAlreadyFulfilled
	}

	if wish.archived {
		return ErrWishAlreadyArchived
	}

	if !wish.Promised() {
		return ErrorWishNotPromised
	}

	wish.fulfilled = true

	return nil
}

func (wish *Wish) Archive() error {
	if wish.archived {
		return ErrWishAlreadyArchived
	}

	wish.archived = true

	return nil
}

func (wish *Wish) Promise(assignee AssigneeId) error {
	if wish.Promised() {
		return ErrWishAlreadyPromised
	}

	if wish.archived {
		return ErrWishAlreadyArchived
	}

	if wish.fulfilled {
		// actually, by logic, it is impossible to fulfill a wish that is not promised yet
		// also, it is impossible to dismiss a wish that is already fulfilled.
		// we keep this check for the sake of consistency
		return ErrWishAlreadyFulfilled
	}

	wish.assignee = assignee

	return nil
}

func (wish *Wish) Renege() error {
	if wish.fulfilled {
		return ErrWishAlreadyFulfilled
	}

	if wish.archived {
		return ErrWishAlreadyArchived
	}

	if !wish.Promised() {
		return ErrorWishNotPromised
	}

	wish.assignee = NilAssigneeId

	return nil
}

func (wish *Wish) Promised() bool {
	return wish.assignee != NilAssigneeId
}
