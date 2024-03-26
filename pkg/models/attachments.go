package models

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type AttachmentType = uint8

const (
	AttachmentOthers = AttachmentType(iota)
	AttachmentPhoto
	AttachmentVideo
	AttachmentAudio
)

type Attachment struct {
	BaseModel

	FileID      string         `json:"file_id"`
	Filesize    int64          `json:"filesize"`
	Filename    string         `json:"filename"`
	Mimetype    string         `json:"mimetype"`
	Hashcode    string         `json:"hashcode"`
	Type        AttachmentType `json:"type"`
	ExternalUrl string         `json:"external_url"`
	Author      Account        `json:"author"`
	ArticleID   *uint          `json:"article_id"`
	MomentID    *uint          `json:"moment_id"`
	CommentID   *uint          `json:"comment_id"`
	AuthorID    uint           `json:"author_id"`
}

func (v Attachment) GetStoragePath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.FileID)
}

func (v Attachment) GetAccessPath() string {
	return fmt.Sprintf("/api/attachments/o/%s", v.FileID)
}
