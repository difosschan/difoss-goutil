package http_util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path/filepath"
)

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

var HttpContentType = map[string]string{
	".avi":  "video/avi",
	".mp3":  "audio/mp3",
	".mp4":  "video/mp4",
	".wmv":  "video/x-ms-wmv",
	".asf":  "video/x-ms-asf",
	".rm":   "application/vnd.rn-realmedia",
	".rmvb": "application/vnd.rn-realmedia-vbr",
	".mov":  "video/quicktime",
	".m4v":  "video/mp4",
	".flv":  "video/x-flv",
	".jpg":  "image/jpeg",
	".png":  "image/png",
}

func HttpContentTypeByExt(filePath string) string {
	ext := filepath.Ext(filePath)
	if len(ext) > 0 {
		if contentType, ok := HttpContentType[ext]; ok {
			return contentType
		}
	}
	return ""
}

func HttpResponseBodyToStruct(response *http.Response, resultPtr interface{}) error {
	body := &bytes.Buffer{}
	_, e := body.ReadFrom(response.Body)
	if e != nil {
		return e
	}
	defer response.Body.Close()
	if e = json.Unmarshal(body.Bytes(), resultPtr); e != nil {
		return e
	}
	return nil
}
