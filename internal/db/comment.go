package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Roham-Ghotbi/go-api/internal/comment"
	uuid "github.com/nu7hatch/gouuid"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Body.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComment(ctx context.Context, uuid string) (comment.Comment, error) {
	var cmtRow CommentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, slug, body, author
		FROM comments
		WHERE id = $1`,
		uuid,
	)
	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching the comment by uuid: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}

func (d *Database) PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to generate uuid %w", err)
	}
	cmt.ID = id.String()
	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO comments
		(id, slug, Author, body)
		VALUES
		(:id, :slug, :author, :body)`,
		postRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}
	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments where id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment from database: %w", err)
	}
	return nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, cmt comment.Comment) (comment.Comment, error) {
	cmtRow := CommentRow{
		ID:     id,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE comments SET
		slug = :slug,
		author = :author,
		body = :body
		WHERE id = :id`,
		cmtRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to upfate comment: %w", err)
	}
	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}
