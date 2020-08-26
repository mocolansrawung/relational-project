package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evermos/boilerplate-go/mocks"
	"github.com/evermos/boilerplate-go/internal/dto"

	"github.com/golang/mock/gomock"
)

func TestExampleHandler_Example(t *testing.T) {
	tests := []struct {
		name          string
		configureMock func(*mocks.MockExampleServiceContract)
		want          string
	}{
		{
			name: "health check is good",
			configureMock: func(contract *mocks.MockExampleServiceContract) {
				contract.EXPECT().Get().Return(dto.Example{
					Status: "good",
				}, nil)
			},
			want: "{\"message\":\"success\",\"data\":{\"status\":\"good\"}}\n",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exampleService := mocks.NewMockExampleServiceContract(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ExampleHandler{
				ExampleService: exampleService,
			}

			tt.configureMock(exampleService)

			ts := httptest.NewServer(http.HandlerFunc(h.Example))
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			if string(body) != tt.want {
				t.Errorf("Example() got = %v, want %v", string(body), tt.want)
			}
		})
	}
}
