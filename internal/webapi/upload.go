package webapi

import "errors"

type UploadToken struct {
	AppId   string
	AppVer  string
	IspType int
	UserKey string
}

func (t *UploadToken) Available() bool {
	return t.UserKey != ""
}

type UploadResultData struct {
	AreaId     IntString `json:"aid"`
	CategoryId IntString `json:"cid"`
	FileId     string    `json:"file_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	PickCode   string    `json:"pick_code"`
	Sha1       string    `json:"sha1"`
}

type UploadInfoResponse struct {
	BasicResponse
	AppId       IntString `json:"app_id"`
	AppVersion  IntString `json:"app_version"`
	UploadLimit int64     `json:"size_limit"`
	IspType     int       `json:"isp_type"`
	UserId      int       `json:"user_id"`
	UserKey     string    `json:"userkey"`
}

type UploadInitResponse struct {
	Request   string `json:"request"`
	ErrorCode int    `json:"statuscode"`
	ErrorMsg  string `json:"statusmsg"`

	Status   int    `json:"status"`
	PickCode string `json:"pickcode"`

	// OSS upload fields
	Bucket   string `json:"bucket"`
	Object   string `json:"object"`
	Callback struct {
		Callback    string `json:"callback"`
		CallbackVar string `json:"callback_var"`
	} `json:"callback"`

	// Useless fields
	FileId   int    `json:"fileid"`
	FileInfo string `json:"fileinfo"`
	Target   string `json:"target"`
}

func (r *UploadInitResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return errors.New(r.ErrorMsg)
}

type UploadOssParams struct {
	Bucket      string
	Object      string
	Callback    string
	CallbackVar string
}

type UploadOssTokenResponse struct {
	StatusCode      string `json:"StatusCode"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
}

func (r *UploadOssTokenResponse) Err() error {
	if r.StatusCode == "200" {
		return nil
	}
	return ErrUnexpected
}