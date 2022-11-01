package network

import (
	"encoding/json"
	"give-me-setu/main/storage"
	"give-me-setu/util"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
)

const (
	OK             int = 0
	OTHER_ERROR    int = 99
	FILE_NOT_IMAGE     = 100
	LACK_QUERY
	LIB_NOTFOUND
	FILE_ALREADY_EXIST
	CLIENT_ERROR // 服务器的 client 访问外部时出现的错误
	JSON_ERROR
)

type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SetuServer struct {
	lib storage.ImgLib
}

func (s *SetuServer) Run(port string) error {
	return http.ListenAndServe(":"+port, nil)
}

func NewServer(rootLibDir string) *SetuServer {
	s := &SetuServer{
		lib: *storage.GetLib(rootLibDir),
	}
	ws := new(restful.WebService)
	ws.Route(ws.GET("/ping").To(ping))
	ws.Route(ws.PUT("/library/{*}").To(s.libraryAdd))
	restful.Add(ws)
	return s
}

func ping(r1 *restful.Request, r2 *restful.Response) {
	io.WriteString(r2, "hello")
}

// TODO: 编写测试
func (s *SetuServer) libraryAdd(r1 *restful.Request, r2 *restful.Response) {
	resp := BaseResp{
		Code:    0,
		Message: "ok",
	}
	defer writeJson(r2, &resp)

	libName := strings.TrimPrefix(r1.Request.URL.Path, "/library/")

	// 文件前缀
	name := ""
	// 文件后缀
	extName := ""

	lib, failName := s.lib.Go(libName)
	if failName != nil {
		switch len(failName) {
		case 0:
			resp.Code = OTHER_ERROR
			resp.Message = "Must set a [name] after [lib]"
			return
		case 1:
			name = failName[0]
			extName = path.Ext(name)
		default:
			resp.Code = LIB_NOTFOUND
			resp.Message = "Can not find lib: " + failName[0]
			log.Println(resp.Message)
			return
		}
	}

	// 网络图片的 url（如果有）
	imgUrl := ""
	// 网络请求的content-type
	// 第一次赋值用于判断当前请求体的格式，如果不是 json 则作为文件接收
	// 如果是 json 则取出里面的 url 用于下载网络图片
	contentType := r1.HeaderParameter("Content-Type")
	if contentType == "application/json" {
		data := new(LibAddReq)
		body, err := io.ReadAll(r1.Request.Body)
		if err != nil {
			resp.Code = OTHER_ERROR
			resp.Message = "Can not read request body"
			return
		}
		err = json.Unmarshal(body, data)
		if err != nil {
			resp.Code = JSON_ERROR
			resp.Message = "Can not parse json data: " + err.Error()
			return
		}
		imgUrl = data.Url
	}
	// 接收的图片
	var img io.Reader
	if imgUrl != "" {
		webImg, err := http.DefaultClient.Get(imgUrl)
		if err != nil {
			resp.Code = CLIENT_ERROR
			resp.Message = "Failed to fetch image from [" + imgUrl + "]: " + err.Error()
			log.Println(resp.Message)
			return
		}
		img = webImg.Body
		contentType = webImg.Header.Get("Content-Type")
		defer webImg.Body.Close()
	} else {
		img = r1.Request.Body
	}

	// 写入图片到缓存文件
	file, err := os.CreateTemp("", "give-me-setu")
	if err != nil {
		resp.Code = OTHER_ERROR
		resp.Message = "Failed to create temp: " + err.Error()
		log.Printf(resp.Message)
		return
	}
	defer file.Close()
	defer os.Remove(path.Join(os.TempDir(), file.Name()))
	_, err = io.Copy(file, img)
	if err != nil {
		log.Println(err)
		resp.Code = OTHER_ERROR
		resp.Message = "Failed to copy image from url"
		return
	}

	// 获取文件后缀
	if extName == "" {
		_, ext, _ := strings.Cut(contentType, "/")
		extName = "." + ext
	} else {
		name = strings.TrimSuffix(name, extName)
	}

	if !util.IsMIMEType(extName, "image") {
		resp.Code = FILE_NOT_IMAGE
		resp.Message = "Maybe file is not a image, check [name] or Content-Type from request and [url]"
		return
	}

	if _, ok := lib.Setus[name+extName]; ok {
		resp.Code = FILE_ALREADY_EXIST
		resp.Message = "File " + name + " is exist"
		return
	}

	// 添加图像到库
	err = lib.Add(file.Name(), name+extName)
	if err != nil {
		resp.Code = OTHER_ERROR
		resp.Message = "Failed to add image into library"
	}
}

func writeJson(w http.ResponseWriter, obj any) {
	j, err := json.Marshal(obj)
	if err != nil {
		log.Panicf("Write Json error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
