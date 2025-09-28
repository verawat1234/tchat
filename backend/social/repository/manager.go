package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Manager implements RepositoryManager using GORM
type Manager struct {
	db            *gorm.DB
	userRepo      UserRepository
	postRepo      PostRepository
	commentRepo   CommentRepository
	reactionRepo  ReactionRepository
	communityRepo CommunityRepository
	shareRepo     ShareRepository
}

// NewManager creates a new repository manager
func NewManager(db *gorm.DB) RepositoryManager {
	return &Manager{
		db:            db,
		userRepo:      NewGormUserRepository(db),
		postRepo:      NewGormPostRepository(db),
		commentRepo:   NewGormCommentRepository(db),
		reactionRepo:  NewGormReactionRepository(db),
		communityRepo: NewGormCommunityRepository(db),
		shareRepo:     NewGormShareRepository(db),
	}
}

func (m *Manager) Users() UserRepository {
	return m.userRepo
}

func (m *Manager) Posts() PostRepository {
	return m.postRepo
}

func (m *Manager) Comments() CommentRepository {
	return m.commentRepo
}

func (m *Manager) Reactions() ReactionRepository {
	return m.reactionRepo
}

func (m *Manager) Communities() CommunityRepository {
	return m.communityRepo
}

func (m *Manager) Shares() ShareRepository {
	return m.shareRepo
}

func (m *Manager) WithTransaction(ctx context.Context, fn func(ctx context.Context, rm RepositoryManager) error) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	txRM := &Manager{
		db:            tx,
		userRepo:      NewGormUserRepository(tx),
		postRepo:      NewGormPostRepository(tx),
		commentRepo:   NewGormCommentRepository(tx),
		reactionRepo:  NewGormReactionRepository(tx),
		communityRepo: NewGormCommunityRepository(tx),
		shareRepo:     NewGormShareRepository(tx),
	}

	err := fn(ctx, txRM)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}

func (m *Manager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

