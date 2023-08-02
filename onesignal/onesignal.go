package onesignal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OneSignalConn struct {
	client *http.Client
}

func NewOneSignal(client *http.Client) OneSignalConn {
	return OneSignalConn{
		client: client,
	}
}

func (os OneSignalConn) UpdateUserTag(data, appID, userID string) (error, int) {
	url := fmt.Sprintf("https://onesignal.com/api/v1/apps/%s/users/%s", appID, userID)
	tags := strings.NewReader(fmt.Sprintf("{\"tags\":%s}", data))
	req, err := http.NewRequest("PUT", url, tags)
	if err != nil {
		return err, 400
	}

	req.Header.Add("accept", "text/plain")
	req.Header.Add("Content-Type", "application/json")
	res, err := os.client.Do(req)
	if err != nil {
		return err, res.StatusCode
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err, res.StatusCode
	}

	return nil, res.StatusCode
}
