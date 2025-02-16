package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"avito-shop/internal/models"
	"avito-shop/internal/services"
)

func TestItemService_GetItemByName(t *testing.T) {
	t.Run("Успешное получение товара", func(t *testing.T) {
		mockItemRepo := new(MockItemRepository)

		expectedItem := &models.Item{
			Name:  "t-shirt",
			Price: 80,
		}

		mockItemRepo.On("GetByName", "t-shirt").Return(expectedItem, nil)

		itemService := services.NewItemService(mockItemRepo)

		item, err := itemService.GetItemByName("t-shirt")

		assert.NoError(t, err)
		assert.Equal(t, expectedItem, item)
	})
}
