package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func PostRequest(path string, body []byte, headers *map[string]string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("path nil.")
	}
	if body == nil {
		return nil, fmt.Errorf("body nil.")
	}

	// body
	bodyBuff := bytes.NewBuffer(body)
	requestReader := io.MultiReader(bodyBuff)
	request, err := http.NewRequest("POST", path, requestReader)
	if err != nil {
		return nil, err
	}

	// header
	request.Header.Add("Content-Type", "text/html")
	request.ContentLength = int64(bodyBuff.Len())
	if headers != nil {
		for k, v := range *headers {
			request.Header.Add(k, v)
		}
	}

	// request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("%v\r\n%v\r\n%s", response.StatusCode, response.Header, respBody)
	}

	return respBody, nil
}

func GetRequest(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("URL nil")
	}

	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%d\r\n%s", response.StatusCode, string(body))
	}

	return body, nil
}

func HttpErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
