package mozapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type mozAPIImpl struct {
	accessID  string
	secretKey string
}

// New MozAPI
func New(accessID string, secretKey string) MozAPI {
	return &mozAPIImpl{
		accessID:  accessID,
		secretKey: secretKey,
	}
}

func (m *mozAPIImpl) MetricsForURL(query string, columns int64) (*URLMetrics, error) {

	// Set your expires times for several minutes into the future.
	// An expires time excessively far in the future will not be honored by the Mozscape API.
	expires := time.Now().Add(ExpireTimeInSeconds * time.Second).Unix()

	// Build Signature
	signature := buildURLSignature(m.accessID, m.secretKey, expires)

	// Build URL
	mozURL := "http://lsapi.seomoz.com/linkscape/url-metrics"
	requestURL := fmt.Sprintf("%v/%v?Cols=%v&AccessID=%v&Expires=%v&Signature=%v",
		mozURL, url.QueryEscape(query), columns, m.accessID, expires, signature)

	// Send request
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Retrieve response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	metrics := &URLMetrics{}
	err = json.Unmarshal(body, metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (m *mozAPIImpl) MetricsForURLBatch(urls []string, columns int64) ([]*URLMetrics, error) {

	// Set your expires times for several minutes into the future.
	// An expires time excessively far in the future will not be honored by the Mozscape API.
	expires := time.Now().Add(ExpireTimeInSeconds * time.Second).Unix()

	// Build Signature
	signature := buildURLSignature(m.accessID, m.secretKey, expires)

	// Build URL
	mozURL := "http://lsapi.seomoz.com/linkscape/url-metrics/"
	requestURL := fmt.Sprintf("%v?Cols=%v&AccessID=%v&Expires=%v&Signature=%v",
		mozURL, columns, m.accessID, expires, signature)

	// Prepare URL list
	b, err := json.Marshal(urls)
	if err != nil {
		return nil, err
	}

	// Build request
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Retrieve response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var metrics []*URLMetrics

	if err := json.Unmarshal(body, &metrics); err != nil {

		var response map[string]string

		if errResp := json.Unmarshal(body, &response); errResp != nil {
			fmt.Printf("Response: %v", string(body))
			return nil, errResp
		}

		if status, ok := response["status"]; ok {
			if status == "429" {
				return nil, ErrTooManyRequests
			}
			return nil, fmt.Errorf("Status: %v, Message: %v", status, response["message"])
		}

		return nil, err
	}

	return metrics, nil
}

func buildURLSignature(accessID string, secretKey string, expires int64) string {

	// Prepare credentials
	stringToSign := fmt.Sprintf("%v\n%v", accessID, expires)

	// Get the "raw" or binary output of the hmac hash.
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(stringToSign))
	binarySignature := mac.Sum(nil)

	// Base64-encode it and then url-encode that.
	urlSafeSignature := base64.StdEncoding.EncodeToString(binarySignature)

	return url.QueryEscape(urlSafeSignature)
}
