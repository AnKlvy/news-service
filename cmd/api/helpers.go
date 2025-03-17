package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"net/url"

	"news_service.andreyklimov.net/internal/validator"

	"github.com/julienschmidt/httprouter"
)

// Получает параметр "id" из URL текущего запроса, затем преобразует его в целое число
// и возвращает. Если операция не удалась, возвращает 0 и ошибку.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// Определяем тип envelope.
type envelope map[string]any

// Изменяем тип параметра data на envelope вместо any.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Используем http.MaxBytesReader() для ограничения размера тела запроса до 1 МБ.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	// Инициализируем json.Decoder и вызываем метод DisallowUnknownFields(),
	// чтобы запретить неизвестные поля в JSON. Если клиент отправит поле,
	// которое не может быть сопоставлено с целевым объектом, декодер вернёт ошибку.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// Декодируем тело запроса в назначенную структуру.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// Добавляем новую переменную maxBytesError.
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// Если JSON содержит поле, которое не может быть сопоставлено с целевой структурой,
		// Decode() вернёт ошибку в формате "json: unknown field "<name>"". Извлекаем
		// имя поля из ошибки и вставляем его в наше сообщение об ошибке.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// Используем errors.As() для проверки, имеет ли ошибка тип *http.MaxBytesError.
		// Если да, значит, тело запроса превысило лимит в 1 МБ, и возвращаем сообщение об ошибке.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Вызываем Decode() снова, используя пустую анонимную структуру как целевой объект.
	// Если тело запроса содержит только одно JSON-значение, вернётся ошибка io.EOF.
	// Если мы получаем что-то ещё, значит, в теле запроса есть лишние данные.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// Вспомогательная функция readString() возвращает строковое значение из строки запроса
// или указанное значение по умолчанию, если соответствующий ключ не найден.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Извлекаем значение для заданного ключа из строки запроса. Если ключ отсутствует,
	// будет возвращена пустая строка "".
	s := qs.Get(key)
	// Если ключ отсутствует (или значение пустое), возвращаем значение по умолчанию.
	if s == "" {
		return defaultValue
	}
	// В противном случае возвращаем строку.
	return s
}

// Вспомогательная функция readCSV() получает строковое значение из строки запроса,
// затем разбивает его на срез по символу запятой. Если ключ не найден, возвращает
// указанное значение по умолчанию.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Извлекаем значение из строки запроса.
	csv := qs.Get(key)
	// Если ключ отсутствует (или значение пустое), возвращаем значение по умолчанию.
	if csv == "" {
		return defaultValue
	}
	// В противном случае разбираем значение в срез []string и возвращаем его.
	return strings.Split(csv, ",")
}

// Вспомогательная функция readInt() получает строковое значение из строки запроса,
// затем преобразует его в целое число перед возвратом. Если ключ не найден, возвращает
// указанное значение по умолчанию. Если значение нельзя преобразовать в целое число,
// записывает сообщение об ошибке в переданный экземпляр Validator.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Извлекаем значение из строки запроса.
	s := qs.Get(key)
	// Если ключ отсутствует (или значение пустое), возвращаем значение по умолчанию.
	if s == "" {
		return defaultValue
	}
	// Пытаемся преобразовать значение в int. Если не удаётся, добавляем сообщение
	// об ошибке в экземпляр валидатора и возвращаем значение по умолчанию.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	// В противном случае возвращаем преобразованное целое число.
	return i
}
