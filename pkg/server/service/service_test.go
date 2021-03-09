package service

import (
	"20dojo-online/pkg/server/model/mock_model"

	"github.com/golang/mock/gomock"
)

type mockRepository struct {
	userRepository *mock_model.MockUserRepositoryInterface
}

func newMockRepository(ctrl *gomock.Controller) *mockRepository {
	return &mockRepository{
		userRepository: mock_model.NewMockUserRepositoryInterface(ctrl),
	}
}
