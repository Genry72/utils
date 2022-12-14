package httpGetter

import (
	"fmt"
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
		Method
		URI     string
		Body    interface{}
		Headers []map[string]string
		Params  []map[string]string
	}

	type args struct {
		resultStruct interface{}
	}

	ur := NewUniversalRequest(5, 0)
	ur.RespStatus = 200

	// Параметры запроса
	params := []map[string]string{{"param1": "1"}, {"param2": "2"}}
	headers := []map[string]string{{"Myhead1": "1"}, {"Myhead2": "2"}}
	bodMap := map[string]interface{}{"bodvalue": 1}
	tests := []struct {
		name          string
		args          args
		fields        fields
		wantReqDetail string
		respContain   string // Проверка существования в ответе текста
		wantErr       bool
	}{
		//Post
		{
			name: "Без структуры",
			args: args{
				resultStruct: "",
			},
			fields: fields{
				Method:  MethodPost,
				URI:     "https://httpbin.org/post",
				Body:    "bodvalue",
				Headers: headers,
				Params:  params,
			},
			wantReqDetail: "Method:", // ищем слово в информации о запросе
			respContain:   "args",
			wantErr:       false,
		},
		{
			name: "Со структурой",
			args: args{
				resultStruct: HttpbinStruct{},
			},
			fields: fields{
				Method:  MethodPost,
				URI:     "https://httpbin.org/post",
				Body:    bodMap,
				Headers: headers,
				Params:  params,
			},
			wantReqDetail: "Method:", // ищем слово в информации о запросе
			respContain:   "args",
			wantErr:       false,
		},
		//Get
		{
			name: "Get Со структурой",
			args: args{
				resultStruct: HttpbinStruct{},
			},
			fields: fields{
				Method:  MethodGet,
				URI:     "https://httpbin.org/get",
				Body:    nil,
				Headers: headers,
				Params:  params,
			},
			wantReqDetail: "Method:", // ищем слово в информации о запросе
			respContain:   "args",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur.Method = tt.fields.Method
			ur.URI = tt.fields.URI
			ur.Body = tt.fields.Body
			ur.Headers = tt.fields.Headers
			ur.Params = tt.fields.Params

			result := fmt.Sprintf("%+v", tt.args.resultStruct)
			reqDetail, err := ur.UniversalRequest(&result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Вернулась ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			switch tt.args.resultStruct.(type) {
			case string: // Передача cтроки
				result := tt.args.resultStruct.(string)
				reqDetail, err = ur.UniversalRequest(&result)

				if !strings.Contains(result, tt.respContain) {
					t.Errorf("Неожиданное тело ответа. Не нашли в боди %v. Боди:%v", tt.respContain, result)
				}

			case HttpbinStruct: // Передача структуры
				result := tt.args.resultStruct.(HttpbinStruct)
				reqDetail, err = ur.UniversalRequest(&result)

				if !strings.Contains(result.URL, "https://httpbin.org") {
					t.Errorf("Не нашли в стуктуре URL https://httpbin.org/post. Полученная структура %+v", result)
				}
				if result.Args["param1"] != "1" {
					t.Errorf("Не нашли параметр param1 в запросе\n %+v", result)
				}

				if result.Headers["Myhead1"] != "1" {
					t.Errorf("Не нашли хедер Myhead1 в запросе\n %+v", result)
				}
			}

			switch tt.fields.Method {
			case MethodPost:
				//Для метода post проверяем корректную отправку тела запроса
				if !strings.Contains(result, "bodvalue") {
					t.Errorf("Не нашли в теле запроса bodvalue\n %+v", result)
				}
			}

			if !strings.Contains(fmt.Sprintf("%+v", reqDetail), tt.wantReqDetail) {
				t.Errorf("В информации о заросе не нашли %v, вернулось %v", tt.wantReqDetail, reqDetail)
			}

			if !strings.Contains(fmt.Sprint(reqDetail.URL), "https://httpbin.org") {
				t.Errorf("В информации о запросе не вернулся URL\n %+v", reqDetail)
			}

		})
	}
}

func ExampleUniversalRequest_UniversalRequest() {
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
	ur := NewUniversalRequest(5, 0)
	params := []map[string]string{{"param1": "1"}, {"param2": "2"}}
	headers := []map[string]string{{"Myhead1": "1"}, {"Myhead2": "2"}}
	bodMap := map[string]interface{}{"bodvalue": 1}
	ur.Method = MethodPost
	ur.URI = "https://httpbin.org/post"
	ur.RespStatus = 200
	ur.Body = bodMap
	ur.Headers = headers
	ur.Params = params
	result := HttpbinStruct{}
	_, err := ur.UniversalRequest(&result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.URL)
	// Output: https://httpbin.org/post?param1=1&param2=2
}
