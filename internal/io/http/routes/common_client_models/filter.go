package common_client_models

// FilterResult результат выборки по фильтру
type FilterResult struct {
	Total uint        // Полное количество объектов в выборке
	Items interface{} // Объекты выборки данных
}
