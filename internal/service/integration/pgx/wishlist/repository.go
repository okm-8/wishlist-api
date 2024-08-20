package user

import (
	"api/internal/model/wishlist"
	"api/internal/service/integration/pgx/driver"
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgx/v5"
)

var (
	//go:embed query/insert.sql
	insertSQL string
	//go:embed query/upsert-wish.sql
	upsertWishSQL string
	//go:embed query/select-by-wisher-id.sql
	selectByWisherIdSQL string
	//go:embed query/select-by-assignee-id.sql
	selectByAssigneeIdSQL string
	//go:embed query/select-by-id.sql
	selectByIdSQL string
	//go:embed query/select-wish-by-id.sql
	selectWishByIdSQL string
)

var (
	ErrAssigneeIdNil         = errors.New("assignee id is nil")
	ErrorNilWishlistReturned = errors.New("nil wishlists returned")
	ErrorNilWisReturned      = errors.New("nil wish returned")
)

func insert(ctx context.Context, executor driver.CommandExecutor, wishlist *wishlist.Wishlist) error {
	_, err := executor.Exec(ctx, insertSQL, pgx.NamedArgs{
		"id":          wishlist.Id().String(),
		"wisherId":    wishlist.Wisher().Id().String(),
		"name":        wishlist.Name(),
		"description": wishlist.Description(),
		"hidden":      wishlist.Hidden(),
	})

	return err
}

func upsertWish(
	ctx context.Context,
	executor driver.CommandExecutor,
	wish *wishlist.Wish,
) error {
	assigneeId := wishlist.NilAssigneeId

	if wish.Promised() {
		assigneeId = wish.Assignee().Id()
	}

	_, err := executor.Exec(ctx, upsertWishSQL, pgx.NamedArgs{
		"id":          wish.Id().String(),
		"wishlistId":  wish.Wishlist().Id().String(),
		"name":        wish.Name(),
		"description": wish.Description(),
		"hidden":      wish.Hidden(),
		"fulfilled":   wish.Fulfilled(),
		"assigneeId":  assigneeId.String(),
	})

	return err
}

type wishlistRecord struct {
	Id          wishlist.Id
	Wisher      *wishlist.Wisher
	Name        string
	Description string
	Hidden      bool
	Wishes      []*wishlist.Wish
}

func (record *wishlistRecord) Wishlist() *wishlist.Wishlist {
	return wishlist.Restore(
		record.Id,
		record.Wisher,
		record.Name,
		record.Description,
		record.Hidden,
		record.Wishes,
	)
}

type wishlistCollection struct {
	recordMap map[wishlist.Id]*wishlistRecord
}

func (collection *wishlistCollection) add(record *wishlistRecord) {
	if collection.recordMap == nil {
		collection.recordMap = make(map[wishlist.Id]*wishlistRecord)
	}

	if _, ok := collection.recordMap[record.Id]; !ok {
		collection.recordMap[record.Id] = record
	}
}

func (collection *wishlistCollection) addWish(wishlistId wishlist.Id, wish *wishlist.Wish) {
	if record, ok := collection.recordMap[wishlistId]; ok {
		record.Wishes = append(record.Wishes, wish)
	}
}

func (collection *wishlistCollection) toWishlists() []*wishlist.Wishlist {
	wishlists := make([]*wishlist.Wishlist, 0, len(collection.recordMap))

	for _, record := range collection.recordMap {
		wishlists = append(wishlists, record.Wishlist())
	}

	return wishlists
}

func selectByWisherId(
	ctx context.Context,
	executor driver.QueryExecutor,
	wisherId wishlist.WisherId,
) ([]*wishlist.Wishlist, error) {
	rows, err := executor.Query(ctx, selectByWisherIdSQL, pgx.NamedArgs{
		"wisherId": wisherId.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wishlists := &wishlistCollection{}
	assignees := map[wishlist.AssigneeId]*wishlist.Assignee{
		wishlist.NilAssigneeId: nil,
	}
	var wisher *wishlist.Wisher

	for rows.Next() {
		var idStr *string
		var name *string
		var description *string
		var hidden *bool
		var fulfilled *bool
		var assigneeIdStr *string
		var assigneeEmail *string
		var assigneeName *string
		var wishlistIdStr string
		var wishlistName string
		var wishlistDescription string
		var wishlistHidden bool
		var wisherEmail string
		var wisherName string

		err = rows.Scan(
			&idStr,
			&name,
			&description,
			&hidden,
			&fulfilled,
			&assigneeIdStr,
			&assigneeEmail,
			&assigneeName,
			&wishlistIdStr,
			&wishlistName,
			&wishlistDescription,
			&wishlistHidden,
			&wisherEmail,
			&wisherName,
		)

		if err != nil {
			return nil, err
		}

		if wisher == nil {
			wisher = wishlist.RestoreWisher(wisherId, wisherName, wisherEmail)
		}

		var id wishlist.Id
		id, err = wishlist.ParseId(wishlistIdStr)

		if err != nil {
			return nil, err
		}

		wishlists.add(&wishlistRecord{
			Id:          id,
			Wisher:      wisher,
			Name:        wishlistName,
			Description: wishlistDescription,
			Hidden:      wishlistHidden,
		})

		var assignee *wishlist.Assignee
		var ok bool

		if assigneeIdStr != nil {
			assigneeId, err := wishlist.ParseAssigneeId(*assigneeIdStr)

			if err != nil {
				return nil, err
			}

			if assignee, ok = assignees[assigneeId]; !ok {
				assignee = wishlist.RestoreAssignee(assigneeId, *assigneeName, *assigneeEmail)
				assignees[assigneeId] = assignee
			}
		}

		if idStr != nil {
			var wishId wishlist.WishId
			wishId, err = wishlist.ParseWishId(*idStr)

			if err != nil {
				return nil, err
			}

			wishlists.addWish(id, wishlist.RestoreWish(
				wishId,
				*name,
				*description,
				*fulfilled,
				*hidden,
				assignee,
			))
		}
	}

	return wishlists.toWishlists(), nil
}

func selectByAssigneeId(
	ctx context.Context,
	executor driver.QueryExecutor,
	assigneeId wishlist.AssigneeId,
) ([]*wishlist.Wishlist, error) {
	if assigneeId == wishlist.NilAssigneeId {
		return nil, ErrAssigneeIdNil
	}

	rows, err := executor.Query(ctx, selectByAssigneeIdSQL, pgx.NamedArgs{
		"assigneeId": assigneeId.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wishlists := &wishlistCollection{}
	wishers := map[wishlist.WisherId]*wishlist.Wisher{
		wishlist.NilWisherId: nil,
	}
	var assignee *wishlist.Assignee

	for rows.Next() {
		var idStr string
		var name string
		var description string
		var hidden bool
		var fulfilled bool
		var assigneeEmail string
		var assigneeName string
		var wishlistIdStr string
		var wishlistName string
		var wishlistDescription string
		var wishlistHidden bool
		var wisherIdStr string
		var wisherEmail string
		var wisherName string

		err = rows.Scan(
			&idStr,
			&name,
			&description,
			&hidden,
			&fulfilled,
			&assigneeEmail,
			&assigneeName,
			&wishlistIdStr,
			&wishlistName,
			&wishlistDescription,
			&wishlistHidden,
			&wisherIdStr,
			&wisherEmail,
			&wisherName,
		)

		if err != nil {
			return nil, err
		}

		var id wishlist.Id
		id, err = wishlist.ParseId(wishlistIdStr)

		if err != nil {
			return nil, err
		}

		var wisher *wishlist.Wisher
		var ok bool

		if wisherIdStr != "" {
			var wisherId wishlist.WisherId
			wisherId, err = wishlist.ParseWisherId(wisherIdStr)

			if err != nil {
				return nil, err
			}

			if wisher, ok = wishers[wisherId]; !ok {
				wisher = wishlist.RestoreWisher(wisherId, wisherName, wisherEmail)
				wishers[wisherId] = wisher
			}
		}

		if assignee == nil {
			assignee = wishlist.RestoreAssignee(assigneeId, assigneeName, assigneeEmail)
		}

		wishlists.add(&wishlistRecord{
			Id:          id,
			Wisher:      wisher,
			Name:        wishlistName,
			Description: wishlistDescription,
			Hidden:      wishlistHidden,
		})

		var wishId wishlist.WishId
		wishId, err = wishlist.ParseWishId(idStr)

		if err != nil {
			return nil, err
		}

		wishlists.addWish(id, wishlist.RestoreWish(
			wishId,
			name,
			description,
			fulfilled,
			hidden,
			assignee,
		))
	}

	return wishlists.toWishlists(), nil
}

func selectById(
	ctx context.Context,
	executor driver.QueryExecutor,
	id wishlist.Id,
) (*wishlist.Wishlist, error) {
	rows, err := executor.Query(ctx, selectByIdSQL, pgx.NamedArgs{
		"wishlistId": id.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var wisher *wishlist.Wisher
	var record *wishlistRecord
	assignees := map[wishlist.AssigneeId]*wishlist.Assignee{
		wishlist.NilAssigneeId: nil,
	}

	for rows.Next() {
		var idStr *string
		var name *string
		var description *string
		var hidden *bool
		var fulfilled *bool
		var assigneeIdStr *string
		var assigneeEmail *string
		var assigneeName *string
		var wishlistName string
		var wishlistDescription string
		var wishlistHidden bool
		var wisherIdStr string
		var wisherEmail string
		var wisherName string

		err = rows.Scan(
			&idStr,
			&name,
			&description,
			&hidden,
			&fulfilled,
			&assigneeIdStr,
			&assigneeEmail,
			&assigneeName,
			&wishlistName,
			&wishlistDescription,
			&wishlistHidden,
			&wisherIdStr,
			&wisherEmail,
			&wisherName,
		)

		if err != nil {
			return nil, err
		}

		if wisher == nil {
			var wisherId wishlist.WisherId
			wisherId, err = wishlist.ParseWisherId(wisherIdStr)

			wisher = wishlist.RestoreWisher(
				wisherId,
				wisherName,
				wisherEmail,
			)
		}

		if record == nil {
			record = &wishlistRecord{
				Id:          id,
				Wisher:      wisher,
				Name:        wishlistName,
				Description: wishlistDescription,
				Hidden:      wishlistHidden,
			}
		}

		var assignee *wishlist.Assignee
		var ok bool

		if assigneeIdStr != nil {
			var assigneeId wishlist.AssigneeId
			assigneeId, err = wishlist.ParseAssigneeId(*assigneeIdStr)

			if err != nil {
				return nil, err
			}

			if assignee, ok = assignees[assigneeId]; !ok {
				assignee = wishlist.RestoreAssignee(assigneeId, *assigneeName, *assigneeEmail)
				assignees[assigneeId] = assignee
			}
		}

		if idStr != nil {
			var wishId wishlist.WishId
			wishId, err = wishlist.ParseWishId(*idStr)

			if err != nil {
				return nil, err
			}

			record.Wishes = append(record.Wishes, wishlist.RestoreWish(
				wishId,
				*name,
				*description,
				*fulfilled,
				*hidden,
				assignee,
			))
		}
	}

	if record == nil {
		return nil, nil
	}

	return record.Wishlist(), nil
}

func selectWishById(
	ctx context.Context,
	executor driver.QueryExecutor,
	id wishlist.WishId,
) (*wishlist.Wish, error) {
	rows, err := executor.Query(ctx, selectWishByIdSQL, pgx.NamedArgs{
		"id": id.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var name string
	var description string
	var hidden bool
	var fulfilled bool
	var assigneeIdStr *string
	var assigneeEmail *string
	var assigneeName *string
	var wishlistIdStr string
	var wishlistName string
	var wishlistDescription string
	var wisherIdStr string
	var wisherEmail string
	var wisherName string

	err = rows.Scan(
		&name,
		&description,
		&hidden,
		&fulfilled,
		&assigneeIdStr,
		&assigneeEmail,
		&assigneeName,
		&wishlistIdStr,
		&wishlistName,
		&wishlistDescription,
		&wisherIdStr,
		&wisherEmail,
		&wisherName,
	)

	if err != nil {
		return nil, err
	}

	wisherId, err := wishlist.ParseWisherId(wisherIdStr)

	if err != nil {
		return nil, err
	}

	wisher := wishlist.RestoreWisher(wisherId, wisherName, wisherEmail)

	var assignee *wishlist.Assignee

	if assigneeIdStr != nil {
		var assigneeId wishlist.AssigneeId
		assigneeId, err = wishlist.ParseAssigneeId(*assigneeIdStr)

		if err != nil {
			return nil, err
		}

		assignee = wishlist.RestoreAssignee(assigneeId, *assigneeName, *assigneeEmail)
	}

	wish := wishlist.RestoreWish(
		id,
		name,
		description,
		fulfilled,
		hidden,
		assignee,
	)

	wishlistId, err := wishlist.ParseId(wishlistIdStr)

	if err != nil {
		return nil, err
	}

	// set wishlist to wish
	_ = wishlist.Restore(
		wishlistId,
		wisher,
		wishlistName,
		wishlistDescription,
		hidden,
		[]*wishlist.Wish{wish},
	)

	return wish, nil
}

func Store(ctx Context, wishlist *wishlist.Wishlist) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		if err := insert(ctx.RuntimeContext(), tx, wishlist); err != nil {
			return err
		}

		for _, wish := range wishlist.Wishes() {
			if err := upsertWish(ctx.RuntimeContext(), tx, wish); err != nil {
				return err
			}
		}

		return nil
	}, driver.TxInsertDefaultOpts)
}

// Update updates a wishlist and its wishes.
// It passes the wishlist to the doUpdate function, which should return the updated wishlist.
// If wishlist is not found, it passes nil to the doUpdate function.
// If an error occurs during the update, it returns the error, transaction will be rolled back.
// If the updated wishlist is nil, it returns ErrorNilWishlistReturned.
func Update(
	ctx Context,
	wishlistId wishlist.Id,
	doUpdate func(*wishlist.Wishlist) (*wishlist.Wishlist, error),
) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		_wishlist, err := selectById(ctx.RuntimeContext(), tx, wishlistId)

		if err != nil {
			return err
		}

		var updatedWishlist *wishlist.Wishlist
		updatedWishlist, err = doUpdate(_wishlist)

		if err != nil {
			return err
		}

		if updatedWishlist == nil {
			return ErrorNilWishlistReturned
		}

		if _wishlist == nil {
			if err := insert(ctx.RuntimeContext(), tx, updatedWishlist); err != nil {
				return err
			}
		}

		for _, wish := range updatedWishlist.Wishes() {
			if err = upsertWish(ctx.RuntimeContext(), tx, wish); err != nil {
				return err
			}
		}

		return nil
	}, driver.TxUpdateDefaultOpts)
}

// UpdateWish updates a wish.
// It passes the wish to the doUpdate function, which should return the updated wish.
// If wish is not found, it passes nil to the doUpdate function.
// If an error occurs during the update, it returns the error, transaction will be rolled back.
// If the updated wish is nil, it returns ErrorNilWisReturned.
// Wish is passed to the doUpdate function with the wishlist field set,
// but it contains only the wish itself, not the entire wishlist.
// Also, it will not store any changes to the wishlist, only the wish will be updated.
func UpdateWish(
	ctx Context,
	wishId wishlist.WishId,
	doUpdate func(*wishlist.Wish) (*wishlist.Wish, error),
) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		_wish, err := selectWishById(ctx.RuntimeContext(), tx, wishId)

		if err != nil {
			return err
		}

		var updatedWish *wishlist.Wish
		updatedWish, err = doUpdate(_wish)

		if err != nil {
			return err
		}

		if updatedWish == nil {
			return ErrorNilWisReturned
		}

		return upsertWish(ctx.RuntimeContext(), tx, updatedWish)
	}, driver.TxUpdateDefaultOpts)
}

func GetByWisherId(ctx Context, wisherId wishlist.WisherId) ([]*wishlist.Wishlist, error) {
	var wishlists []*wishlist.Wishlist
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		wishlists, err = selectByWisherId(ctx.RuntimeContext(), conn, wisherId)
		return err
	})

	return wishlists, err
}

func GetByAssigneeId(ctx Context, assigneeId wishlist.AssigneeId) ([]*wishlist.Wishlist, error) {
	var wishlists []*wishlist.Wishlist
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		wishlists, err = selectByAssigneeId(ctx.RuntimeContext(), conn, assigneeId)
		return err
	})

	return wishlists, err
}
