package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
	"github.com/organiza-aqui/backend/internal/util"
	"github.com/organiza-aqui/backend/pkg/logger"
	"go.uber.org/zap"
)

type BankService interface {
	GetAllBanks(ctx context.Context) ([]*dto.BankDTO, error)
	GetBankByID(ctx context.Context, id uuid.UUID) (*dto.BankDTO, error)
	SyncBanks(ctx context.Context) error
}

type bankService struct {
	bankRepo      repository.BankRepository
	brasilAPIClient *util.BrasilAPIClient
}

func NewBankService(bankRepo repository.BankRepository) BankService {
	return &bankService{
		bankRepo:        bankRepo,
		brasilAPIClient: util.NewBrasilAPIClient(),
	}
}

func (s *bankService) GetAllBanks(ctx context.Context) ([]*dto.BankDTO, error) {
	banks, err := s.bankRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar bancos: %w", err)
	}

	dtos := make([]*dto.BankDTO, len(banks))
	for i, bank := range banks {
		dtos[i] = s.modelToDTO(bank)
	}

	return dtos, nil
}

func (s *bankService) GetBankByID(ctx context.Context, id uuid.UUID) (*dto.BankDTO, error) {
	bank, err := s.bankRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar banco: %w", err)
	}
	if bank == nil {
		return nil, fmt.Errorf("banco não encontrado")
	}

	return s.modelToDTO(bank), nil
}

func (s *bankService) SyncBanks(ctx context.Context) error {
	logger.Info("Iniciando sincronização de bancos com BrasilAPI")

	banksFromAPI, err := s.brasilAPIClient.GetBanks()
	if err != nil {
		return fmt.Errorf("erro ao buscar bancos da API: %w", err)
	}

	logger.Info("Bancos recebidos da API",
		zap.Int("total", len(banksFromAPI)),
	)

	var created, updated, unchanged int

	for _, apiBank := range banksFromAPI {
		bank := &model.Bank{
			ISPB:     apiBank.ISPB,
			Code:     apiBank.Code,
			Name:     apiBank.Name,
			FullName: apiBank.FullName,
			UpdatedAt: time.Now(),
		}

		existing, err := s.bankRepo.FindByISPB(ctx, bank.ISPB)
		if err != nil {
			logger.Error("Erro ao buscar banco por ISPB", err,
				zap.String("ispb", bank.ISPB),
			)
			continue
		}

		if existing == nil {
			// Criar novo banco
			bank.ID = uuid.New()
			bank.CreatedAt = time.Now()
			if err := s.bankRepo.Create(ctx, bank); err != nil {
				logger.Error("Erro ao criar banco", err,
					zap.String("ispb", bank.ISPB),
					zap.String("name", bank.Name),
				)
				continue
			}
			created++
			logger.Debug("Banco criado",
				zap.String("ispb", bank.ISPB),
				zap.String("name", bank.Name),
			)
		} else {
			// Verificar se há mudanças
			if existing.Code != bank.Code || existing.Name != bank.Name || existing.FullName != bank.FullName {
				bank.ID = existing.ID
				bank.CreatedAt = existing.CreatedAt
				if err := s.bankRepo.Update(ctx, bank); err != nil {
					logger.Error("Erro ao atualizar banco", err,
						zap.String("ispb", bank.ISPB),
						zap.String("name", bank.Name),
					)
					continue
				}
				updated++
				logger.Debug("Banco atualizado",
					zap.String("ispb", bank.ISPB),
					zap.String("name", bank.Name),
				)
			} else {
				unchanged++
			}
		}
	}

	logger.Info("Sincronização de bancos concluída",
		zap.Int("total_consultados", len(banksFromAPI)),
		zap.Int("criados", created),
		zap.Int("atualizados", updated),
		zap.Int("inalterados", unchanged),
	)

	return nil
}

func (s *bankService) modelToDTO(bank *model.Bank) *dto.BankDTO {
	return &dto.BankDTO{
		ID:       bank.ID.String(),
		ISPB:     bank.ISPB,
		Code:     bank.Code,
		Name:     bank.Name,
		FullName: bank.FullName,
	}
}
