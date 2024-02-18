package binance

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

func (as *apiService) CoinStartFutureUserDataStream() (*Stream, error) {
	params := make(map[string]string)

	res, err := as.request("POST", "dapi/v1/listenKey", params, true, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from userDataStream.post")
	}
	defer res.Body.Close()

	log.Println(string(textRes))
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	var s Stream
	if err := json.Unmarshal(textRes, &s); err != nil {
		return nil, errors.Wrap(err, "stream unmarshal failed")
	}
	return &s, nil
}
func (as *apiService) CoinKeepAliveFutureUserDataStream(s *Stream) error {
	params := make(map[string]string)
	params["listenKey"] = s.ListenKey

	res, err := as.request("PUT", "dapi/v1/listenKey", params, true, false)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from userDataStream.put")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}
	return nil
}
func (as *apiService) CoinCloseFutureUserDataStream(s *Stream) error {
	params := make(map[string]string)
	params["listenKey"] = s.ListenKey

	res, err := as.request("DELETE", "dapi/v1/listenKey", params, true, false)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from userDataStream.delete")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}
	return nil
}
