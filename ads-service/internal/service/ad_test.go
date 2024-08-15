package services

import (
	"ads-service/internal/lib/types"
	"ads-service/internal/models"
	"ads-service/internal/repository"
	"ads-service/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetOne(t *testing.T) {

	testAd := &models.Ad{
		Id:          1,
		Title:       "title",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}

	testCases := []struct {
		name       string
		mockExpect func(userRepo *mocks.AdRepoInterface)
		testAd     *models.Ad
		err        error
	}{
		{
			name: "объявление найдено",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("GetOne", testAd.Id).Return(testAd, nil)
			},
			testAd: testAd,
			err:    nil,
		},
		{
			name: "объявление не найдено",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("GetOne", testAd.Id).Return(nil, pgx.ErrNoRows)
			},
			testAd: testAd,
			err:    types.ErrAdNotFound,
		},
		{name: "Получена любая другая ошибка",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("GetOne", testAd.Id).Return(nil, errors.New("some error"))
			},
			testAd: testAd,
			err:    errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockAdRepoInterface := &mocks.AdRepoInterface{}
			repository := &repository.Repository{Ad: mockAdRepoInterface}
			service := NewService(repository)

			testCase.mockExpect(mockAdRepoInterface)
			_, err := service.Ad.GetOne(testAd.Id)

			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestCreate(t *testing.T) {
	testAd := &models.Ad{
		Title:       "title",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
		Description: "description",
	}
	testId := 2

	testCases := []struct {
		name       string
		mockExpect func(userRepo *mocks.AdRepoInterface)
		testAd     *models.Ad
		returnedId int
		err        error
	}{
		{
			name: "Успешный кейс, валидный ввод",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("Create", testAd).Return(testId, nil)
			},
			testAd:     testAd,
			returnedId: testId,
			err:        nil,
		},
		{
			name: "Любая ошибка из базы",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("Create", testAd).Return(0, errors.New("some error"))
			},
			testAd:     testAd,
			returnedId: -1,
			err:        errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockAdRepoInterface := &mocks.AdRepoInterface{}
			repository := &repository.Repository{Ad: mockAdRepoInterface}
			service := NewService(repository)

			testCase.mockExpect(mockAdRepoInterface)
			id, err := service.Ad.Create(testAd)

			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.returnedId, id)
		})
	}
}

func TestGetAll(t *testing.T) {
	var testAds []*models.Ad = []*models.Ad{
		{
			Id:          1,
			Title:       "title",
			Price:       100,
			Photos:      []string{"photo1", "photo2"},
			Description: "description",
		},
		{
			Id:          2,
			Title:       "title",
			Price:       100,
			Photos:      []string{"photo1", "photo2"},
			Description: "description",
		},
	}
	testPriceSort := "price asc"
	testDateSort := "date desc"
	testUserId := 6
	testPageNumber := 1

	testCases := []struct {
		name        string
		mockExpect  func(userRepo *mocks.AdRepoInterface)
		expectedAds []*models.Ad
		err         error
	}{
		{
			name: "Успешный кейс, валидный ввод",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("GetAll", testPriceSort, testDateSort, testPageNumber).Return(testAds, nil)
			},
			expectedAds: testAds,
			err:         nil,
		},
		{
			name: "Любая ошибка из базы",
			mockExpect: func(userRepo *mocks.AdRepoInterface) {
				userRepo.On("GetAll", testPriceSort, testDateSort, testPageNumber).
					Return(nil, errors.New("some error"))
			},
			expectedAds: nil,
			err:         errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockAdRepoInterface := &mocks.AdRepoInterface{}
			repository := &repository.Repository{Ad: mockAdRepoInterface}
			service := NewService(repository)

			testCase.mockExpect(mockAdRepoInterface)
			ads, err := service.Ad.GetAll(testPriceSort, testDateSort, testPageNumber, testUserId)

			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.expectedAds, ads)
			// проверка, что все ожидания мока (mockExpect) были верно выполнены
			mockAdRepoInterface.AssertExpectations(t)
		})
	}
}
