package controllers

import (
	"ads-service/internal/lib/types"
	"ads-service/internal/models"
	"ads-service/internal/service/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	// добавить инициализацию логгера, чтобы логировать в файл
	testAd := &models.Ad{
		Title:       "title",
		Description: "description",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
	}

	tests := []struct {
		name         string
		mockExpect   func(adService *mocks.AdServiceInterface)
		ad           *models.Ad
		expectOutput string
		code         int
	}{
		{
			name: "Валидный кейс",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("Create", testAd).Return(3, nil)
			},
			ad:           testAd,
			expectOutput: `{"id": 3}`,
			code:         http.StatusCreated,
		},
		{
			name:       "Заголовок длиной более 200 символов",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       strings.Repeat("a", 201),
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: types.ErrInvalidTitle.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Заголовок пуст",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       "",
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: types.ErrEmptyTitle.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Длина описания больше 1000 символов",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      testAd.Photos,
				Title:       testAd.Title,
				Description: strings.Repeat("a", 1001),
				Price:       testAd.Price,
			},
			expectOutput: types.ErrInvalidDescription.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:       "Попытка загрузить более чем 3 ссылки на фото",
			mockExpect: func(adService *mocks.AdServiceInterface) {},
			ad: &models.Ad{
				Photos:      []string{"photo1", "photo2", "photo3", "photo4"},
				Title:       testAd.Title,
				Description: testAd.Description,
				Price:       testAd.Price,
			},
			expectOutput: types.ErrInvalidPhotos.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name: "Случайная ошибка от сервиса",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("Create", testAd).Return(0, errors.New("some error"))
			},
			ad:           testAd,
			expectOutput: errors.New("some error").Error() + "\n",
			code:         http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// инициализация контекста
			ctx := context.Background()
			// инициализация мок-сервиса
			mockService := &mocks.AdServiceInterface{}
			// инициализация тестируемого контроллера
			mux := http.NewServeMux()
			adController := InitAdController(ctx, mockService, mux)

			// превращаем тестовое объявление в JSON
			marshalled, _ := json.Marshal(test.ad)
			// создаем тестовый запрос, который будем пихать в тестируемый контроллер
			req, _ := http.NewRequest("POST", "/ad/create/", strings.NewReader(string(marshalled)))
			// создаем подставной обработчик, который будет слушать ответ тестируемого контроллера
			w := httptest.NewRecorder()
			// саем моку ожидаемое поведение
			test.mockExpect(mockService)
			// запускаем тестируемый контроллер, который должен записать свой ответ в
			// инициализированный выше подставной обработчик
			adController.Create(w, req)

			// проверка соответствия кода ответа
			assert.Equal(t, test.code, w.Code)
			// проверка записи котроллера
			assert.Equal(t, test.expectOutput, w.Body.String())
			// проверка, что ожадаемое моком поведение в точности выполнено
			mockService.AssertExpectations(t)
		})
	}
}
func TestGetOneAd(t *testing.T) {
	// добавить инициализацию логгера, чтобы логировать в файл
	testAd := &models.Ad{
		Title:       "title",
		Description: "description",
		Price:       100,
		Photos:      []string{"photo1", "photo2"},
	}

	tests := []struct {
		name         string
		mockExpect   func(adService *mocks.AdServiceInterface)
		id           string
		expectOutput *models.Ad
		err          string
		code         int
	}{
		{
			name: "Валидный кейс",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("GetOne", 5).Return(testAd, nil)
			},
			id:           "5",
			expectOutput: testAd,
			code:         http.StatusOK,
		},
		{
			name:         "Невалидный идентификатор объявления",
			mockExpect:   func(adService *mocks.AdServiceInterface) {},
			id:           "y",
			expectOutput: testAd,
			err:          types.ErrInvalidId.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name: "Любая ошибка из сервиса",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("GetOne", 5).Return(nil, errors.New("some error"))
			},
			id:           "5",
			expectOutput: testAd,
			err:          errors.New("some error").Error() + "\n",
			code:         http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// инициализация контекста
			ctx := context.Background()
			// инициализация мок-сервиса
			mockService := &mocks.AdServiceInterface{}
			// инициализация тестируемого контроллера
			mux := http.NewServeMux()
			adController := InitAdController(ctx, mockService, mux)

			req, _ := http.NewRequest("GET", "/ad/", strings.NewReader(""))
			// создаем подставной обработчик, который будет слушать ответ тестируемого контроллера
			w := httptest.NewRecorder()

			test.mockExpect(mockService)
			// запускаем тестируемый контроллер, который должен записать свой ответ в
			// инициализированный выше подставной обработчик
			adController.GetOne(w, req)

			// Если ожидается ошибка - проверяем ее
			if test.err != "" {
				assert.Equal(t, test.err, w.Body.String())
			} else {
				// Превращаем ожидаемый результат в JSON
				marshalled, _ := json.Marshal(test.expectOutput)
				// затем сверяем его с полученным
				assert.Equal(t, marshalled, w.Body.Bytes())
			}
			// проверка соответствия кода ответа
			assert.Equal(t, test.code, w.Code)
			// проверка, что ожадаемое моком поведение в точности выполнено
			mockService.AssertExpectations(t)
		})
	}
}
func TestGetAllAds(t *testing.T) {
	// добавить инициализацию логгера, чтобы логировать в файл
	testAds := []*models.Ad{
		{
			Title:       "title",
			Description: "description",
			Price:       100,
			Photos:      []string{"photo1", "photo2"},
		},
		{
			Title:       "title321",
			Description: "description321",
			Price:       100312,
			Photos:      []string{"photo13451", "phodsfvto2"},
		},
	}

	tests := []struct {
		name         string
		mockExpect   func(adService *mocks.AdServiceInterface)
		expectOutput []*models.Ad
		queryParams  map[string]string
		err          string
		code         int
	}{
		{
			name: "Валидный кейс, все параметры присутствуют",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("GetAll", "price asc", "date desc", 1).Return(testAds, nil)
			},
			expectOutput: testAds,
			queryParams:  map[string]string{"page": "1", "price": "asc", "date": "desc"},
			code:         http.StatusOK,
		},
		{
			name:         "Невалидный номер страницы",
			mockExpect:   func(adService *mocks.AdServiceInterface) {},
			expectOutput: testAds,
			queryParams:  map[string]string{"page": "j", "price": "asc", "date": "desc"},
			err:          types.ErrInvalidPageNumber.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:         "Невалидный параметр валидации по цене",
			mockExpect:   func(adService *mocks.AdServiceInterface) {},
			expectOutput: testAds,
			queryParams:  map[string]string{"page": "1", "price": "hello,world", "date": "desc"},
			err:          types.ErrInvalidPriceSort.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name:         "Невалидный параметр валидации по дате",
			mockExpect:   func(adService *mocks.AdServiceInterface) {},
			expectOutput: testAds,
			queryParams:  map[string]string{"page": "1", "price": "asc", "date": ""},
			err:          types.ErrInvalidDateSort.Error() + "\n",
			code:         http.StatusBadRequest,
		},
		{
			name: "Любая ошибка в сервисе",
			mockExpect: func(adService *mocks.AdServiceInterface) {
				adService.On("GetAll", "", "", 1).Return(nil, errors.New("some error"))
			},
			queryParams: map[string]string{"page": "1"},
			err:         "some error\n",
			code:        http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			mockService := &mocks.AdServiceInterface{}
			mux := http.NewServeMux()
			adController := InitAdController(ctx, mockService, mux)

			// Создаем тестовый запрос и добавляем туда query параметры
			req, _ := http.NewRequest("GET", "/ad/all/", strings.NewReader(""))
			query := req.URL.Query()
			for k, v := range test.queryParams {
				query.Add(k, v)
			}
			req.URL.RawQuery = query.Encode()

			// создаем обработчик, добавляем моку ожидание и вызываем тестируемую функцию
			w := httptest.NewRecorder()
			test.mockExpect(mockService)
			adController.GetAll(w, req)

			if test.err != "" {
				assert.Equal(t, test.err, w.Body.String())
			} else {
				marshalled, _ := json.Marshal(test.expectOutput)
				assert.Equal(t, marshalled, w.Body.Bytes())
			}
			assert.Equal(t, test.code, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
