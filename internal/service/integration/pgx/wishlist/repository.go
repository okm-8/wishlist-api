package wishlist

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
	//go:embed query/select-by-wisher-id-active.sql
	selectByWisherIdActiveSQL string
)

var (
	ErrAssigneeIdNil         = errors.New("assignee id is nil")
	ErrorNilWishlistReturned = errors.New("nil wishlists returned")
	ErrorNilWisReturned      = errors.New("nil wish returned")
)

func upsert(ctx context.Context, executor driver.CommandExecutor, wishlist *wishlist.Wishlist) error {
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
	var assigneeId *string

	if wish.Promised() {
		assigneeIdStr := wish.Assignee().Id().String()
		assigneeId = &assigneeIdStr
	}

	_, err := executor.Exec(ctx, upsertWishSQL, pgx.NamedArgs{
		"id":          wish.Id().String(),
		"wishlistId":  wish.Wishlist().Id().String(),
		"name":        wish.Name(),
		"description": wish.Description(),
		"hidden":      wish.Hidden(),
		"fulfilled":   wish.Fulfilled(),
		"assigneeId":  assigneeId,
	})

	return err
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

	return scanWishesViewRows(rows)
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

	return scanWishesViewRows(rows)
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

	wishlists, err := scanWishesViewRows(rows)

	if err != nil {
		return nil, err
	}

	if len(wishlists) == 0 {
		return nil, nil
	}

	return wishlists[0], nil
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

	wishlists, err := scanWishesViewRows(rows)

	if err != nil {
		return nil, err
	}

	if len(wishlists) == 0 || len(wishlists[0].Wishes()) == 0 {
		return nil, nil
	}

	return wishlists[0].Wishes()[0], nil
}

func selectByWisherIdActive(
	ctx context.Context,
	executor driver.QueryExecutor,
	wisherId wishlist.WisherId,
) ([]*wishlist.Wishlist, error) {
	rows, err := executor.Query(ctx, selectByWisherIdActiveSQL, pgx.NamedArgs{
		"wisherId": wisherId.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanWishesViewRows(rows)
}

func Store(ctx Context, wishlist *wishlist.Wishlist) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		if err := upsert(ctx.RuntimeContext(), tx, wishlist); err != nil {
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

		if err = upsert(ctx.RuntimeContext(), tx, updatedWishlist); err != nil {
			return err
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

func GetByWisherIdActive(ctx Context, wisherId wishlist.WisherId) ([]*wishlist.Wishlist, error) {
	var wishlists []*wishlist.Wishlist
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		wishlists, err = selectByWisherIdActive(ctx.RuntimeContext(), conn, wisherId)
		return err
	})

	return wishlists, err
}
