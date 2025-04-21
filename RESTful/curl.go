package custom_util

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// 统一的HTTP客户端
var client = &http.Client{}

// 文件上传
func HttpProxyFileUpload(file *multipart.FileHeader, fileKey string, addFields map[string]string,
	addHeaders map[string]string, urlPath string) (body []byte, err error) {
	return HttpProxyFileUploadCustom(file, fileKey, file.Filename, addFields, addHeaders, urlPath)
}

func HttpPostJson(addHeaders map[string]string, data interface{}, urlPath string) (body []byte, err error) {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return
	}
	
	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(bytesData))
	if err != nil {
		return
	}
	
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for headerKey, headerVal := range addHeaders {
		req.Header.Set(headerKey, headerVal)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	return io.ReadAll(resp.Body)
}

func HttpPostJsonDownload(addHeaders map[string]string, data interface{}, urlPath string, filePathName string) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(bytesData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for headerKey, headerVal := range addHeaders {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePathName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func HttpPostJsonPut(addHeaders map[string]string, data interface{}, urlPath string) (body []byte, err error) {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPut, urlPath, bytes.NewReader(bytesData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for headerKey, headerVal := range addHeaders {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func HttpGet(addHeaders map[string]string, data map[string]string, urlPath string) (body []byte, err error) {
	params := url.Values{}
	urlInfo, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	for dataKey, dataVal := range data {
		params.Set(dataKey, dataVal)
	}
	urlInfo.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, urlInfo.String(), nil)
	if err != nil {
		return
	}

	for headerKey, headerVal := range addHeaders {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func HttpProxyFileUploadCustom(file *multipart.FileHeader, fileKey, filename string, addFields map[string]string,
	addHeaders map[string]string, urlPath string) (body []byte, err error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	
	formFile, err := writer.CreateFormFile(fileKey, filename)
	if err != nil {
		log.Fatalf("Upload Create form file failed: %v", err)
	}

	srcFile, err := file.Open()
	if err != nil {
		log.Fatalf("Upload Create form file failed: %v", err)
	}
	defer srcFile.Close()

	if _, err = io.Copy(formFile, srcFile); err != nil {
		log.Fatalf("Write to form file failed: %v", err)
	}

	for fieldKey, fieldVal := range addFields {
		if err = writer.WriteField(fieldKey, fieldVal); err != nil {
			log.Fatalf("WriteField failed: %v", err)
		}
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, urlPath, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for headerKey, headerVal := range addHeaders {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Post failed: %v", err)
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}