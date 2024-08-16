package wishlist

type Wisher struct {
	id    WisherId
	name  string
	email string
}

func RestoreWisher(id WisherId, name string, email string) *Wisher {
	return &Wisher{id: id, name: name, email: email}
}

func (wisher *Wisher) Id() WisherId {
	return wisher.id
}

func (wisher *Wisher) Name() string {
	return wisher.name
}

func (wisher *Wisher) Email() string {
	return wisher.email
}
