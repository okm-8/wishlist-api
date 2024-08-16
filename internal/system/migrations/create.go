package migrations

import (
	systemContext "api/internal/system/context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

type meta struct {
	createdAt   time.Time
	description string
}

func newMeta(description string) *meta {
	return &meta{
		createdAt:   time.Now(),
		description: description,
	}
}

const (
	TSFMT              = "20060102150405.000"
	TSLEN              = len(TSFMT)
	EXT                = ".sql"
	EXTLEN             = len(EXT)
	SEP                = "-"
	SPACER             = " "
	SEPLEN             = len(SEP)
	MAXLEN             = 255
	DESCRIPTION_MAXLEN = MAXLEN - TSLEN - EXTLEN - SEPLEN
)

func (meta *meta) Filename() string {
	description := meta.description

	if len(description) > DESCRIPTION_MAXLEN {
		description = description[:DESCRIPTION_MAXLEN]
	}

	description = strings.ReplaceAll(description, SPACER, SEP)
	date := meta.createdAt.Format(TSFMT)
	date = strings.ReplaceAll(date, ".", "") // remove the dot in the milliseconds

	return fmt.Sprintf("%s%s%s%s", date, SEP, description, EXT)
}

func (meta *meta) CreatedAt() time.Time {
	return meta.createdAt
}

func (meta *meta) Description() string {
	return meta.description
}

//go:embed migration.sql.tpl
var migrationTemplateContent string
var migrationTemplate = template.Must(template.New("migration").Parse(migrationTemplateContent))

func Create(ctx *systemContext.Context, description string) (string, error) {
	path := ctx.Config().Migrations().DirPath()
	_meta := newMeta(description)
	fileName := fmt.Sprintf("%s/%s", path, _meta.Filename())
	file, err := os.Create(fileName)

	if err != nil {
		return fileName, err
	}

	return fileName, migrationTemplate.Execute(file, _meta)
}
