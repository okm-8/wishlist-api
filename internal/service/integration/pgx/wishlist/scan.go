package wishlist

import (
	"api/internal/model/wishlist"

	"github.com/jackc/pgx/v5"
)

type wishlistRecord struct {
	Id          wishlist.Id
	Wisher      *wishlist.Wisher
	Name        string
	Description string
	Hidden      bool
	Wishes      []*wishlist.Wish
}

func scanWishesViewRows(rows pgx.Rows) ([]*wishlist.Wishlist, error) {
	wishers := make(map[wishlist.WisherId]*wishlist.Wisher)
	assignees := make(map[wishlist.AssigneeId]*wishlist.Assignee)
	records := make([]*wishlistRecord, 0)               // This is a slice to keep the order of the rows
	recordsMap := make(map[wishlist.Id]*wishlistRecord) // This is a map to quickly find the record by id

	var err error
	var ok bool
	var wisherId wishlist.WisherId
	var wisher *wishlist.Wisher
	var wishlistId wishlist.Id
	var _wishlistRecord *wishlistRecord
	var assigneeId wishlist.AssigneeId
	var assignee *wishlist.Assignee
	var wishId wishlist.WishId

	for rows.Next() {
		var wishIdStr *string
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
		var wisherIdStr string
		var wisherEmail string
		var wisherName string

		rows.Scan(
			&wishIdStr,
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
			&wisherIdStr,
			&wisherEmail,
			&wisherName,
		)

		if wisherId, err = wishlist.ParseWisherId(wisherIdStr); err != nil {
			return nil, err
		}

		wisher, ok = wishers[wisherId]

		if !ok {
			wisher = wishlist.RestoreWisher(
				wisherId,
				wisherEmail,
				wisherName,
			)

			wishers[wisherId] = wisher
		}

		if wishlistId, err = wishlist.ParseId(wishlistIdStr); err != nil {
			return nil, err
		}

		_wishlistRecord, ok = recordsMap[wishlistId]

		if !ok {
			_wishlistRecord = &wishlistRecord{
				Id:          wishlistId,
				Wisher:      wisher,
				Name:        wishlistName,
				Description: wishlistDescription,
				Hidden:      wishlistHidden,
			}

			records = append(records, _wishlistRecord)
			recordsMap[wishlistId] = _wishlistRecord
		}

		if wishIdStr == nil {
			continue
		}

		if assigneeIdStr != nil {
			if assigneeId, err = wishlist.ParseAssigneeId(*assigneeIdStr); err != nil {
				return nil, err
			}

			assignee, ok = assignees[assigneeId]

			if !ok {
				assignee = wishlist.RestoreAssignee(
					assigneeId,
					*assigneeEmail,
					*assigneeName,
				)

				assignees[assigneeId] = assignee
			}
		}

		if wishId, err = wishlist.ParseWishId(*wishIdStr); err != nil {
			return nil, err
		}

		_wishlistRecord.Wishes = append(
			_wishlistRecord.Wishes,
			wishlist.RestoreWish(
				wishId,
				*name,
				*description,
				*fulfilled,
				*hidden,
				assignee,
			),
		)
	}

	wishlists := make([]*wishlist.Wishlist, len(records))

	for index, record := range records {
		wishlists[index] = wishlist.Restore(
			record.Id,
			record.Wisher,
			record.Name,
			record.Description,
			record.Hidden,
			record.Wishes,
		)
	}

	return wishlists, nil
}
