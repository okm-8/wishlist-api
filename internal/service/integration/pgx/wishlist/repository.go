package user

import _ "embed"

var (
	//go:embed query/insert.sql
	insertSQL string
	//go:embed query/upsert-wish.sql
	upsertWishSQL string
	//go:embed query/select-wishlists-by-wisher-id.sql
	selectWishlistsByWisherIdSQL string
	//go:embed query/select-wishlists-by-assignee-id.sql
	selectWishlistsByAssigneeIdSQL string
)
