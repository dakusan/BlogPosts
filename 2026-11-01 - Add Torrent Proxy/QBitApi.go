package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

//qbitApiCall performs one qBittorrent WebUI API request and handles all the boilerplate every call needs. A non-200 status is treated as an error.
//POSTs are sent as multipart/form-data. A postBody key starting with “*” means its value is a file path to open+attach.
//The body is built fresh inside makeReq (opening the file anew, if any) rather than once beforehand, since qbitDo can call makeReq a second time to retry after a mid-session relogin - and a reader/file already drained by the first attempt would send the retry with an empty body instead of erroring loudly
//client is nil for every call except the login bootstrap: those go through qbitDo (reusing the shared session, with its login-and-retry-once-on-403 behavior); login instead executes directly on the given client, since qbitDo’s session is exactly what login is establishing and can’t be used to establish itself
//urlPath is relative to “/api/v2/torrents/” (or, if it starts with “/”, to “/api/v2” instead)
//maskPostFields lists postBody keys whose value is shown as “***” in the debug log instead of in full (e.g. “password”).
func qbitApiCall(client *http.Client, method, urlPath string, postBody map[string]string, maskPostFields []string) (string, error) {
	fullPath := "/api/v2/torrents/" + urlPath
	if strings.HasPrefix(urlPath, "/") {
		fullPath = "/api/v2" + urlPath
	}

	makeReq := func() (*http.Request, error) {
		var bodyReader io.Reader
		var contentType string
		if postBody != nil {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			for key, value := range postBody {
				if !strings.HasPrefix(key, "*") {
					mw.WriteField(key, value)
					continue
				}
				file, err := os.Open(value)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				part, err := mw.CreateFormFile(strings.TrimPrefix(key, "*"), filepath.Base(value))
				if err != nil {
					return nil, err
				}
				if _, err := io.Copy(part, file); err != nil {
					return nil, err
				}
			}
			if err := mw.Close(); err != nil {
				return nil, err
			}
			bodyReader, contentType = &buf, mw.FormDataContentType()
		}

		req, err := http.NewRequest(method, cfg.QbitURL+fullPath, bodyReader)
		if err != nil {
			return nil, err
		}
		if postBody != nil {
			req.Header.Set("Content-Type", contentType)
		}
		return req, nil
	}

	var resp *http.Response
	var err error
	if client != nil {
		var req *http.Request
		if req, err = makeReq(); err == nil {
			resp, err = client.Do(req)
		}
		if err != nil {
			qbitMarkDisconnected(err)
			return "", err
		}
	} else if resp, err = qbitDo(makeReq); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if Debug {
		logBodyStr := "[]"
		if len(postBody) > 0 {
			displayBody := make(map[string]string, len(postBody))
			for key, value := range postBody {
				switch {
				case strings.HasPrefix(key, "*"):
					value = "<file: " + filepath.Base(value) + ">"
				case slices.Contains(maskPostFields, key):
					value = "***"
				}
				displayBody[strings.TrimPrefix(key, "*")] = value
			}
			if pretty, err := json.MarshalIndent(displayBody, "", "\t"); err == nil {
				logBodyStr = string(pretty)
			}
		}
		resultStr := string(body)
		var pretty bytes.Buffer
		if json.Indent(&pretty, body, "", "\t") == nil {
			resultStr = pretty.String()
		}
		log.Printf("Debug: %s %s %s: %s", method, urlPath, logBodyStr, resultStr)
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%s %s rejected (status %d): %s", method, urlPath, resp.StatusCode, string(body))
		if client != nil {
			qbitMarkDisconnected(err)
		}
		return string(body), err
	}
	return string(body), nil
}