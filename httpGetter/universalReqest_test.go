package httpGetter

import (
	"github.com/go-resty/resty/v2"
	"strings"
	"testing"
)

// HttpbinStruct структура
type HttpbinStruct struct {
	Args  map[string]string `json:"args"`
	Data  string            `json:"data"`
	Files struct {
	} `json:"files"`
	Form struct {
		Bod1 string `json:"bod1"`
	} `json:"form"`
	Headers map[string]string `json:"headers"`

	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

func TestUniversalRequest_UniversalRequest(t *testing.T) {
	t.Parallel()
	type fields struct {
		Client *resty.Client
		Method
		URI        string
		RespStatus int
		Body       map[string]interface{}
		Headers    []map[string]string
		Params     []map[string]string
	}

	type args struct {
		resultStruct interface{}
	}

	// Параметры запроса
	params := []map[string]string{{"param1": "1"}, {"param2": "2"}}

	tests := []struct {
		name          string
		args          args
		fields        fields
		wantReqDetail string
		bodyContain   string // Проверка существования в ответе текста
		wantErr       bool
	}{
		{
			name: "Без структуры",
			args: args{
				resultStruct: "",
			},
			fields: fields{
				Client:     resty.New(),
				Method:     MethodPost,
				URI:        "https://httpbin.org/post",
				RespStatus: 200,
				Body:       nil,
				Headers:    nil,
				Params:     params,
			},
			wantReqDetail: "Method:POST", // ищем слово в информации о запросе
			bodyContain:   "args",
			wantErr:       false,
		},
		{
			name: "Со структурой",
			args: args{
				resultStruct: HttpbinStruct{},
			},
			fields: fields{
				Client:     resty.New(),
				Method:     MethodPost,
				URI:        "https://httpbin.org/post",
				RespStatus: 200,
				Body:       nil,
				Headers:    nil,
				Params:     params,
			},
			wantReqDetail: "Method:POST", // ищем слово в информации о запросе
			bodyContain:   "args",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := UniversalRequest{
				Client:     tt.fields.Client,
				Method:     tt.fields.Method,
				URI:        tt.fields.URI,
				RespStatus: tt.fields.RespStatus,
				Body:       tt.fields.Body,
				Headers:    tt.fields.Headers,
				Params:     tt.fields.Params,
			}
			var reqDetail string
			var err error
			switch tt.args.resultStruct.(type) {
			case string: // Передача троки
				result := tt.args.resultStruct.(string)
				reqDetail, err = ur.UniversalRequest(&result)

				if !strings.Contains(result, tt.bodyContain) {
					t.Errorf("Неожиданное тело ответа. Не нашли в боди %v. Боди:%v", tt.bodyContain, result)
				}

			case HttpbinStruct: // Передача структуры
				result := tt.args.resultStruct.(HttpbinStruct)
				reqDetail, err = ur.UniversalRequest(&result)

				if !strings.Contains(result.URL, "https://httpbin.org/post") {
					t.Errorf("Не нашли в стуктуре URL https://httpbin.org/post. Полученная структура %+v", result)
				}
				if result.Args["param1"] != "1" {
					t.Errorf("Не нашли параметр param1 в запросе\n %+v", result)
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Вернулась ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(reqDetail, tt.wantReqDetail) {
				t.Errorf("В информации о заросе не нашли %v, вернулось %v", tt.wantReqDetail, reqDetail)
			}

			t.Logf(reqDetail)

		})
	}
}
