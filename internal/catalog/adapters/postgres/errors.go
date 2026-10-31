package postgres

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func uniqueConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return pgErr.ConstraintName
	}
	return ""
}

func constraintContains(err error, part string) bool {
	return strings.Contains(uniqueConstraint(err), part)
}
