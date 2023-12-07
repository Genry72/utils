package httpGetter

import (
	"fmt"
	"net/http"
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
	type args struct {
		Method       string
		resultStruct interface{}
		RespStatus   int
		URI          string
		Body         interface{}
	}
	ur := UniversalRequest{
		Client: NewUsClient(5, 0),
	}

	// Параметры запроса
	params := []map[string]string{{"param1": "1"}, {"param2": "2"}}
	headers := []map[string]string{{"Myhead1": "1"}, {"Myhead2": "2"}}
	bodMap := map[string]interface{}{"bodvalue": 1}
	resulStruct := HttpbinStruct{}
	tests := []struct {
		name string
		args
		wantErr bool
	}{
		//Post
		{
			name: "Без структуры",
			args: args{
				Method:       http.MethodPost,
				resultStruct: nil,
				RespStatus:   200,
				URI:          "https://httpbin.org/post",
				Body:         "bodvalue",
			},
			wantErr: false,
		},
		{
			name: "Со структурой",
			args: args{
				Method:       http.MethodPost,
				resultStruct: &resulStruct,
				RespStatus:   200,
				URI:          "https://httpbin.org/post",
				Body:         bodMap,
			},
			wantErr: false,
		},
		//Get
		{
			name: "Get Со структурой",
			args: args{
				Method:       http.MethodGet,
				resultStruct: &resulStruct,
				RespStatus:   200,
				URI:          "https://httpbin.org/get",
				Body:         nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := NewRequestParams(tt.args.URI, tt.args.Method, tt.args.RespStatus, headers, params, tt.args.Body)

			body, resp, err := ur.GetRresponse(params, tt.args.resultStruct)
			if (err != nil) != tt.wantErr {
				t.Errorf("Вернулась ошибка = %v, wantErr %v, %+v body:%s url: %s", err, tt.wantErr, tt.args, body, resp.Request.Method)
				return
			}
			bodyontains := "ttps://httpbin.org" // ищем в боди
			if !strings.Contains(body, bodyontains) {
				t.Errorf("Неожиданное тело ответа. Не нашли в боди %v. Боди:%v", bodyontains, body)
			}

			if !strings.Contains(resp.Request.URL, "https://httpbin.org") {
				t.Errorf("Не нашли в стуктуре URL https://httpbin.org/post. Полученная структура %+v", resp.Request)
			}

			if tt.args.resultStruct != nil {
				if tt.args.resultStruct.(*HttpbinStruct).Args["param1"] != "1" {
					t.Errorf("Не нашли параметр param1 в запросе\n %+v", tt.args.resultStruct.(HttpbinStruct))
				}
				if tt.args.resultStruct.(*HttpbinStruct).Headers["Myhead1"] != "1" {
					t.Errorf("Не нашли хедер Myhead1 в запросе\n %+v", tt.args.resultStruct.(HttpbinStruct))
				}
			}

			// todo добавить проверку reqdetail, хеедер, параметры
			if tt.args.Method == http.MethodPost {
				//Для метода post проверяем корректную отправку тела запроса
				if !strings.Contains(body, "bodvalue") {
					t.Errorf("Не нашли в теле запроса bodvalue\n %+v", body)
				}
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

	params := []map[string]string{{"param1": "1"}, {"param2": "z"}}
	headers := []map[string]string{{"Myhead1": "1"}, {"Myhead2": "2"}}
	bodMap := map[string]interface{}{"bodvalue": 1}
	ur := UniversalRequest{
		Client: NewUsClient(5, 0),
	}

	reqParams := NewRequestParams("https://httpbin.org/post", http.MethodPost, 200, headers, params, bodMap)
	result := HttpbinStruct{}
	_, _, err := ur.GetRresponse(reqParams, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.URL)
	// Output: https://httpbin.org/post?param1=1&param2=z
}
