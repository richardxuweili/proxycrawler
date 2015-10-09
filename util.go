package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func DefaultCheck(resp *http.Response) error {
	if resp.StatusCode/400 > 0 {
		return fmt.Errorf("wrong status code:%d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}
