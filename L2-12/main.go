package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

/*
L2.12 «HTTP-сервер»
Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.

Требования
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.

Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.

Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.

# Реализовать middleware для логирования запросов

Методы API:

POST /create_event

POST /update_event

POST /delete_event

GET /events_for_day

GET /events_for_week

GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09). В GET методах параметры передаются через queryString, в POST — через тело запроса.

В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."} в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:

Реализовать все методы.

Бизнес логика НЕ должна зависеть от кода HTTP сервера.

В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500.

Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

// Event представляет собой структуру события
type Event struct {
	UserID    int       `json:"user_id"`
	EventDate time.Time `json:"event_date"`
	Note      string    `json:"note"`
}

var events = []Event{}

// Функция для сериализации события в JSON
func serializeEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}

// Функция для десериализации JSON в событие
func deserializeEvent(data []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

//Создадим middleware для логирования запросов.

// LoggingMiddleware логирует входящие HTTP-запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем запрос
		logRequest(r)
		next.ServeHTTP(w, r)
	})
}

func logRequest(r *http.Request) {
	// Здесь можно использовать более сложное логирование, например, logrus или zap
	println(r.Method, r.URL.String())
}

//Валидация параметров

func validateEventParams(userIDStr string, dateStr string) (int, time.Time, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, time.Time{}, errors.New("invalid user_id")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, errors.New("invalid date format")
	}
	return userID, date, nil
}

//HTTP-обработчики

// Создание события
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	userIDStr := r.FormValue("user_id")
	dateStr := r.FormValue("date")
	note := r.FormValue("note")

	userID, eventDate, err := validateEventParams(userIDStr, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := Event{UserID: userID, EventDate: eventDate, Note: note}
	events = append(events, event)

	response := map[string]string{"result": "event created"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Обновление события
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для обновления события
}

// Удаление события
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для удаления события
}

// Получение событий за день
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для получения событий за день
}

// Получение событий за неделю
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для получения событий за неделю
}

// Получение событий за месяц
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	// Логика для получения событий за месяц
}

//Основная функция и роутер

func main() {
	http.Handle("/", LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)
	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)

	port := ":8080" // Укажите ваш порт
	println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
