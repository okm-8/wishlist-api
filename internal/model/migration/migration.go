package migration

import "time"

type Migration struct {
	id         *uint64
	filename   string
	executedAt *time.Time
}

func New(filename string) *Migration {
	return &Migration{
		filename: filename,
	}
}

func Restore(id uint64, filename string, executedAt time.Time) *Migration {
	return &Migration{
		id:         &id,
		filename:   filename,
		executedAt: &executedAt,
	}
}

func (migration *Migration) Id() *uint64 {
	return migration.id
}

func (migration *Migration) Filename() string {
	return migration.filename
}

func (migration *Migration) ExecutedAt() *time.Time {
	if migration.executedAt == nil {
		return nil
	}

	tsCopy := new(time.Time)
	*tsCopy = *migration.executedAt

	return migration.executedAt
}
