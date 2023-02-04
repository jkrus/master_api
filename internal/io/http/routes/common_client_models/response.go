package common_client_models

// RestResponse общий для разных роутов ответ если что-то пошло не так
type RestResponse struct {
	Message string // Ответ для отображения клиенту
	Details string // Детали ответа
}
