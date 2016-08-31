package common

import (
	"bytes"
	"crypto/md5"
	//"log"
	"io"
	"io/ioutil"
	"strings"

	. "github.com/qiniu/api.v6/conf"
	qiniu_io "github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
)

func UploadFile(dir string, read io.Reader) (string, error) {
	ACCESS_KEY = Webconfig.UploadConfig.QiniuAccessKey
	SECRET_KEY = Webconfig.UploadConfig.QiniuSecretKey

	extra := &qiniu_io.PutExtra{
		// Bucket:   "messagedream",
		MimeType: "",
	}

	var policy = rs.PutPolicy{
		Scope: "joe-data-1",
		//设置Token过期时间
		Expires: 3600,
	}

	filename := strings.Replace(UUID(), "-", "", -1) + ".jpg"
	body, _ := ioutil.ReadAll(read)

	h := md5.New()
	h.Write(body)

	key := "static/upload/" + dir + "/" + filename
	ret := new(qiniu_io.PutRet)

	buf := bytes.NewBuffer(body)

	return key, qiniu_io.Put(nil, ret, policy.Token(nil), key, buf, extra)
}
