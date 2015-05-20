package dynago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/underarmour/dynago/schema"
)

/*
This interface defines how all the various queries manage their internal execution logic.

Executor is primarily provided so that testing and mocking can be done on
the API level, not just the transport level.
*/
type Executor interface {
	GetItem(*GetItem) (*GetItemResult, error)
	PutItem(*PutItem) (*PutItemResult, error)
	Query(*Query) (*QueryResult, error)
	UpdateItem(*UpdateItem) (*UpdateItemResult, error)
	CreateTable(*schema.CreateRequest) (*schema.CreateResponse, error)
	BatchWriteItem(*BatchWrite) (*BatchWriteResult, error)
}

type awsExecutor struct {
	endpoint string
	caller   http.Client
	aws      awsInfo
}

// Create an AWS executor with a specified endpoint and AWS parameters.
func NewAwsExecutor(endpoint, region, accessKey, secretKey string) Executor {
	return &awsExecutor{
		endpoint: endpoint,
		aws: awsInfo{
			Region:    region,
			AccessKey: accessKey,
			SecretKey: secretKey,
			Service:   "dynamodb",
		},
	}
}

func (e *awsExecutor) makeRequest(target string, document interface{}) ([]byte, error) {
	buf, err := json.Marshal(document)
	// log.Printf("Request Body: \n%s\n\n", buf)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", e.endpoint, body)
	log.Printf("--------------- Making request ---------------\n\n %s \n\n ---------------------------------------------", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-amz-target", dynamoTargetPrefix+target)
	req.Header.Add("content-type", "application/x-amz-json-1.0")
	req.Header.Set("Host", req.URL.Host)
	e.aws.signRequest(req, buf)
	response, err := e.caller.Do(req)
	fmt.Printf("%+v\n", req)
	if err != nil {
		return nil, err
	}
	respBody, err := responseBytes(response)
	if response.StatusCode != http.StatusOK {
		e := &Error{
			Response:     response,
			ResponseBody: respBody,
		}
		dest := &inputError{}
		if err = json.Unmarshal(respBody, dest); err == nil {
			e.parse(dest)
		} else {
			e.Message = err.Error()
		}
		err = e
	}
	return respBody, err
}

func (e *awsExecutor) makeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error) {
	body, err := e.makeRequest(method, document)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, dest)
	return
}

func responseBytes(response *http.Response) (buf []byte, err error) {
	if response.ContentLength > 0 {
		buf = make([]byte, response.ContentLength)
		_, err = response.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		err = response.Body.Close()
	}
	return
}
