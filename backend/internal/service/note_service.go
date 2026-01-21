package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type NoteService interface {
	CreateNote(ctx context.Context, userID uuid.UUID, req *dto.CreateNoteRequest) (*dto.NoteDTO, error)
	GetNoteByID(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) (*dto.NoteDTO, error)
	GetNotesByUserID(ctx context.Context, userID uuid.UUID, filters *dto.NoteFilters) ([]*dto.NoteDTO, error)
	UpdateNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID, req *dto.UpdateNoteRequest) (*dto.NoteDTO, error)
	DeleteNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) error
}

type noteService struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{
		noteRepo: noteRepo,
	}
}

func (s *noteService) CreateNote(ctx context.Context, userID uuid.UUID, req *dto.CreateNoteRequest) (*dto.NoteDTO, error) {
	now := time.Now()
	note := &model.Note{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      pq.StringArray(req.Tags),
		IsPinned:  req.IsPinned,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("erro ao criar anotação: %w", err)
	}

	return s.modelToDTO(note), nil
}

func (s *noteService) GetNoteByID(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) (*dto.NoteDTO, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar anotação: %w", err)
	}
	if note == nil {
		return nil, fmt.Errorf("anotação %s: %w", noteID, appError.ErrNoteNotFound)
	}
	if note.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(note), nil
}

func (s *noteService) GetNotesByUserID(ctx context.Context, userID uuid.UUID, filters *dto.NoteFilters) ([]*dto.NoteDTO, error) {
	notes, err := s.noteRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar anotações: %w", err)
	}

	dtos := make([]*dto.NoteDTO, len(notes))
	for i, note := range notes {
		dtos[i] = s.modelToDTO(note)
	}

	return dtos, nil
}

func (s *noteService) UpdateNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID, req *dto.UpdateNoteRequest) (*dto.NoteDTO, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar anotação: %w", err)
	}
	if note == nil {
		return nil, fmt.Errorf("anotação %s: %w", noteID, appError.ErrNoteNotFound)
	}
	if note.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	note.Title = req.Title
	note.Content = req.Content
	note.Tags = pq.StringArray(req.Tags)
	note.IsPinned = req.IsPinned
	note.UpdatedAt = time.Now()

	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, fmt.Errorf("erro ao atualizar anotação: %w", err)
	}

	return s.modelToDTO(note), nil
}

func (s *noteService) DeleteNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) error {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("erro ao buscar anotação: %w", err)
	}
	if note == nil {
		return fmt.Errorf("anotação %s: %w", noteID, appError.ErrNoteNotFound)
	}
	if note.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.noteRepo.Delete(ctx, noteID); err != nil {
		return fmt.Errorf("erro ao deletar anotação: %w", err)
	}

	return nil
}

func (s *noteService) modelToDTO(note *model.Note) *dto.NoteDTO {
	return &dto.NoteDTO{
		ID:        note.ID.String(),
		UserID:    note.UserID.String(),
		Title:     note.Title,
		Content:   note.Content,
		Tags:      []string(note.Tags),
		IsPinned:  note.IsPinned,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}
}
