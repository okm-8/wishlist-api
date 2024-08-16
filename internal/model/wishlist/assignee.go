package wishlist

type Assignee struct {
	id    AssigneeId
	name  string
	email string
}

func RestoreAssignee(id AssigneeId, name string, email string) *Assignee {
	return &Assignee{id: id, name: name, email: email}
}

func (assignee *Assignee) Id() AssigneeId {
	return assignee.id
}

func (assignee *Assignee) Name() string {
	return assignee.name
}

func (assignee *Assignee) Email() string {
	return assignee.email
}
