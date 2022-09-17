package httpGetter

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Method string

const (
	MethodGet  Method = "GET"
	MethodPost Method = "POST"
)

type UniversalRequest struct {
	Client     *resty.Client
	Method     Method
	URI        string
	RespStatus int
	Body       interface{}
	Headers    []map[string]string
	Params     []map[string]string
}

// UniversalRequest выполняет любой запрос. Возвращает расперсенный ответ на основе переданной структуры.
// Если вместо структуры передали string, то вернется тело ответа
func (ur UniversalRequest) UniversalRequest(resultStruct interface{}) (req resty.Request, err error) {

	if ur.URI == "" {
		return req, fmt.Errorf("не задан URI")
	}

	if ur.Method == "" {
		return req, fmt.Errorf("не задан метод")
	}

	if ur.RespStatus == 0 {
		return req, fmt.Errorf("не задан код ожидаемого ответа")
	}

	if resultStruct == nil {
		return req, fmt.Errorf("не передана структура. Передайте поинтер на переменную с пустой строкой, если не нужно парсить тело ответа")
	}

	req = *ur.Client.R()

	// Добавляем тело запроса
	if ur.Body != nil {
		req.SetBody(ur.Body)
	}

	// Добавляем заголовки
	for _, v := range ur.Headers {
		req.SetHeaders(v)
	}
	// Добавляем параметры запроса
	for _, v := range ur.Params {
		req.SetQueryParams(v)
	}

	var resp *resty.Response
	switch ur.Method {
	case MethodPost:
		resp, err = req.Post(ur.URI)
		if err != nil {
			return req, err
		}
	case MethodGet:
		resp, err = req.Get(ur.URI)
		if err != nil {
			return req, err
		}
	default:
		return req, fmt.Errorf("указан некорректный метод %v", ur.Method)
	}

	err = ur.checkStatus(resp)
	if err != nil {
		return req, err
	}
	switch resultStruct.(type) {
	case *string: // Возвращаем тело ответа без парсинга
		*resultStruct.(*string) = string(resp.Body())
		return req, err
	default:
		err = json.Unmarshal(resp.Body(), resultStruct)

		if err != nil {
			return req, fmt.Errorf("не удалось распарсить тело ответа: %w %s", err, string(resp.Body()))
		}
		return req, nil
	}
}

func (ur UniversalRequest) checkStatus(resp *resty.Response) error {
	if resp.StatusCode() != ur.RespStatus {
		return fmt.Errorf("%s %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

func NewUniversalRequest(timeout time.Duration, retryCount int) UniversalRequest {
	client := resty.New()
	// Добавляем дефолтный таймаут
	if timeout != 0 {
		client.SetTimeout(timeout * time.Second)
	}

	//Количество повторов запроса, в случае неудачи
	client.SetRetryCount(retryCount)

	return UniversalRequest{
		Client:     client,
		URI:        "",
		RespStatus: 0,
		Body:       nil,
		Headers:    nil,
	}
}
