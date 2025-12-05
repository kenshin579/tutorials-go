package send

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	storeHttp "github.com/kenshin579/tutorials-go/go-unit-test/httpmock/store/http"

	"github.com/kenshin579/tutorials-go/go-unit-test/httpmock/myresty"

	"github.com/stretchr/testify/assert"

	"github.com/go-resty/resty/v2"

	"github.com/jarcoal/httpmock"
)

const (
	Url = "https://mock.google.com"
)

func teardown() {
	httpmock.DeactivateAndReset()
}

func Test_Httpmock_MyResty(t *testing.T) {
	client := myresty.New()
	sendStore := storeHttp.NewHttpSendStore(client)

	httpmock.ActivateNonDefault(client.GetClient())

	defer teardown()

	httpmock.RegisterResponder("GET", Url,
		httpmock.NewStringResponder(200, `{
				"success": true,
				"message": "success"
			}`))

	resp, err := sendStore.Send(Url)
	fmt.Println("err", err)
	fmt.Println("resp", string(resp))
	assert.Contains(t, string(resp), "message")
}

func Test_Httpmock_Resty(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	defer teardown()

	httpmock.RegisterResponder("GET", Url,
		httpmock.NewStringResponder(200, `{
				"success": true,
				"message": "success"
			}`))

	resp, _ := client.R().Get(Url)
	fmt.Println(resp.String())
	assert.Contains(t, resp.String(), "message")

}

func Test_Httpmock_Mock_Response를_반환한다(t *testing.T) {
	httpmock.Activate()
	httpmock.RegisterResponder("GET", Url,
		httpmock.NewStringResponder(200, `
			{
				"success": true,
				"message": "success"
			}`),
	)

	defer teardown()

	resp, err := HTTPGetMethod(Url)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(resp))
}

func HTTPGetMethod(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}
