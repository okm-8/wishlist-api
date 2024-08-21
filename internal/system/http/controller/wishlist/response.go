package wishlist

import "api/internal/model/wishlist"

type wisher struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type assignee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type wish struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Fulfilled   bool      `json:"fulfilled"`
	Hidden      bool      `json:"hidden"`
	Assignee    *assignee `json:"assignee,omitempty"`
}

type wishList struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Wisher      wisher `json:"wisher"`
	Wishes      []wish `json:"wishes"`
}

func serializeWisher(_wisher *wishlist.Wisher) wisher {
	return wisher{
		Id:   _wisher.Id().String(),
		Name: _wisher.Name(),
	}
}

func serializeAssignee(_assignee *wishlist.Assignee) *assignee {
	if _assignee == nil {
		return nil
	}

	return &assignee{
		Id:   _assignee.Id().String(),
		Name: _assignee.Name(),
	}
}

func serializeWish(_wish *wishlist.Wish) wish {
	return wish{
		Id:          _wish.Id().String(),
		Name:        _wish.Name(),
		Description: _wish.Description(),
		Fulfilled:   _wish.Fulfilled(),
		Hidden:      _wish.Hidden(),
		Assignee:    serializeAssignee(_wish.Assignee()),
	}
}

func serializeWishlist(wishlist *wishlist.Wishlist) wishList {
	wishes := make([]wish, len(wishlist.Wishes()))

	for index, _wish := range wishlist.Wishes() {
		wishes[index] = serializeWish(_wish)
	}

	return wishList{
		Id:          wishlist.Id().String(),
		Name:        wishlist.Name(),
		Description: wishlist.Description(),
		Hidden:      wishlist.Hidden(),
		Wisher:      serializeWisher(wishlist.Wisher()),
		Wishes:      wishes,
	}
}

func serializeWishLists(wishlists []*wishlist.Wishlist) []wishList {
	result := make([]wishList, len(wishlists))

	for index, _wishlist := range wishlists {
		result[index] = serializeWishlist(_wishlist)
	}

	return result
}
