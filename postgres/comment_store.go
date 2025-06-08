package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/titaniumcoder/golang-reddit-fake/goreddit"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var t goreddit.Comment
	if err := s.Get(&t, `SELECT * FROM comments WHERE id = $1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("error getting comment: %w", err)
	}
	return t, nil
}

func (s *CommentStore) Comments(postID uuid.UUID) ([]goreddit.Comment, error) {
	var tt []goreddit.Comment
	if err := s.Select(&tt, `SELECT * FROM comments where post_id=$1`, postID); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("error getting comments: %w", err)
	}
	return tt, nil
}

func (s *CommentStore) CreateComment(t *goreddit.Comment) error {
	if err := s.Get(t, `
		INSERT INTO comments(id, post_id, content, contes) 
		VALUES ($1, $2, $3, $4) 
		RETURNING *`, t.ID, t.PostID, t.Content, t.Votes); err != nil {
		return fmt.Errorf("error creating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(t *goreddit.Comment) error {
	if err := s.Get(t, `UPDATE comments SET content=$1, votes=$2, post_id=$3 WHERE id = $4`, t.Content, t.Votes, t.PostID, t.ID); err != nil {
		return fmt.Errorf("error updating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM comments WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}
