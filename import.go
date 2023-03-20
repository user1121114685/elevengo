package elevengo

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/webapi"
)

type ErrImportNeedCheck struct {
	// The sign key your should set to ImportTicket
	SignKey string
	// The sign range in format of "<start>-<end>" in bytes.
	// You can directly use it in ImportCreateTicket. 
	SignRange string
}

func (e *ErrImportNeedCheck) Error() string {
	return "import requires sign check"
}

// ImportTicket container reqiured fields to import(aka. quickly upload) a file
// to your 115 cloud storage.
type ImportTicket struct {
	// File base name
	FileName  string
	// File size in bytes
	FileSize  int64
	// File SHA-1 hash, in upper-case HEX format
	FileSha1  string
	// Sign key from 115 server.
	SignKey string
	// SHA-1 hash value of a segment of the file, in upper-case HEX format
	SignValue string
}

// Import imports(aka. fast-upload) a file to your 115 cloud storage.
// Please check example code for the detailed usage.
func (a *Agent) Import(dirId string, ticket *ImportTicket) (err error) {
	if err = a.uploadInitHelper(); err != nil {
		return
	}
	target := fmt.Sprintf("U_1_%s", dirId)
	initData := &webapi.UploadInitData{
		FileId: ticket.FileSha1,
		FileName: ticket.FileName,
		FileSize: ticket.FileSize,
		Target: target,
		Signature: a.uh.CalculateSignature(ticket.FileSha1, target),
		SignKey: ticket.SignKey,
		SignValue: ticket.SignValue,
	}
	exist, checkRange := false, ""
	if exist, checkRange, err = a.uploadInitInternal(initData, nil); err == nil {
		if checkRange != "" {
			err = &ErrImportNeedCheck{
				SignKey: initData.SignKey,
				SignRange: checkRange,
			}
		} else if !exist {
			err = webapi.ErrNotExist
		}
	}
	return
}

// ImportCreateTicket is a helper function to create an ImportTicket of a file, 
// that you can share to others to import this file to their cloud storage.
// You should also send pickcode together with ticket.
func (a *Agent) ImportCreateTicket(fileId string, ticket *ImportTicket) (pickcode string, err error) {
	file := &File{}
	if err = a.FileGet(fileId, file); err == nil {
		pickcode = file.PickCode
		if ticket != nil {
			ticket.FileName = file.Name
			ticket.FileSize = file.Size
			ticket.FileSha1 = file.Sha1
		}
	}
	return
}

// ImportCalculateSignValue calculates sign value of a file on cloud storage.
// Please check example code for the detailed usage.
func (a *Agent) ImportCalculateSignValue(pickcode string, signRange string) (value string, err error) {
	// Parse range text at first
	var start, end int64
	if _, err = fmt.Sscanf(signRange, "%d-%d", &start, &end); err != nil {
		return
	}
	// Get download URL
	ticket := &DownloadTicket{}
	if err = a.DownloadCreateTicket(pickcode, ticket); err != nil {
		return 
	}
	// Get range content
	var body io.ReadCloser
	if body, err = a.GetRange(ticket.Url, Range{start, end}); err != nil {
		return
	}
	defer util.QuietlyClose(body)
	h := sha1.New()
	if _, err = io.Copy(h, body); err == nil {
		value = hash.ToHexUpper(h)
	}
	return
}
