package httpGetter

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const (
	ErrorWrongResultCode = "неожиданный код ответа"
	ErrorMethodNotDefine = "метод определен не корректно"
)

// RequestParams Параметры запроса
type RequestParams struct {
	Method     string
	URI        string
	RespStatus int
	Headers    []map[string]string
	Params     []map[string]string
	Body       interface{}
}

type UniversalRequest struct {
	Client resty.Client
}

// todo few
func NewRequestParams(url string, method string, respStatus int, headers []map[string]string, params []map[string]string, body interface{}) *RequestParams {
	return &RequestParams{
		Method:     method,
		URI:        url,
		RespStatus: respStatus,
		Headers:    headers,
		Params:     params,
		Body:       body,
	}
}

// GetRresponse выполняет любой запрос. Возвращает распарсенный ответ на основе переданной структуры.
// Если парсить не нужно, передайте nil
func (ur UniversalRequest) GetRresponse(params *RequestParams, resultStruct interface{}) (body string, resp *resty.Response, err error) {
	// Подготавливаем запрос
	req, err := prepare(ur, params)
	if err != nil {
		return "", nil, err
	}
	// Выполняем
	switch params.Method {
	case http.MethodGet:
		resp, err = req.Get(params.URI)
	case http.MethodPost:
		resp, err = req.Post(params.URI)
	default:
		return "", resp, fmt.Errorf(ErrorMethodNotDefine)
	}

	if err != nil {
		return "", resp, err
	}

	body = string(resp.Body())

	// Обрабатываем
	err = getResult(resp, resultStruct, params.RespStatus)
	if err != nil {
		return body, resp, err
	}
	return body, resp, nil
}

func NewUsClient(timeout time.Duration, retryCount int) resty.Client {
	client := *resty.New()
	// Добавляем дефолтный таймаут
	if timeout != 0 {
		client.SetTimeout(timeout * time.Second)
	}

	//Количество повторов запроса, в случае неудачи
	client.SetRetryCount(retryCount)
	return client
}

// Добавляем боди, заголовки и тд
func prepare(ur UniversalRequest, params *RequestParams) (*resty.Request, error) {
	if params.URI == "" {
		return nil, fmt.Errorf("не задан URI")
	}

	if params.RespStatus == 0 {
		return nil, fmt.Errorf("не задан код ожидаемого ответа")
	}

	req := ur.Client.R()

	// Добавляем тело запроса
	if params.Body != nil {
		req.SetBody(params.Body)
	}

	// Добавляем заголовки
	for _, v := range params.Headers {
		req.SetHeaders(v)
	}
	// Добавляем параметры запроса
	for _, v := range params.Params {
		req.SetQueryParams(v)
	}
	return req, nil
}

// Обрабатываем результат ответа: получаем боди, и парсим, если нужно
func getResult(resp *resty.Response, resultStruct interface{}, wantCode int) error {

	if resp.StatusCode() != wantCode {
		return fmt.Errorf(ErrorWrongResultCode)
	}

	if resultStruct == nil {
		return nil
	}

	err := json.Unmarshal(resp.Body(), resultStruct)
	if err != nil {
		return fmt.Errorf("не удалось распарсить тело ответа: %w %s", err, string(resp.Body()))
	}
	return nil
}
