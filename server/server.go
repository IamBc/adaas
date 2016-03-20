package main

import (
    "net/http"
    "encoding/json"
    "mime/multipart"
    "os"
    "io"
    "io/ioutil"
    "flag"
    "time"

    "github.com/golang/glog"
    "github.com/gorilla/mux"
)

// curl --form upload=@README.md localhost:9001/api/v1/dataset_file
// curl localhost:9001/api/v1/zl1/list_dataset_files

func main() {
    flag.Parse()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/api/v1/ping", Ping)
    router.HandleFunc("/api/v1/{userId}/dataset_file", UploadDatasetFile)
    router.HandleFunc("/api/v1/{userId}/list_dataset_files", ListDatasetFiles)
    router.HandleFunc("/api/v1/{userId}/compute_request/builtin/{builtinJobId}/{datasetFileId}", BuiltinJob)

    glog.Info(http.ListenAndServe(":9001", router))
}

type apiRequest struct{
    Status  string `json:"Status"`
    Msg     string `json:"Msg"`
    Payload interface{} `json:"Payload"`
}


type datasetFileInfo struct {
  Name     string
  Size	   int64
  ModTime  time.Time
}


func BuiltinJob(w http.ResponseWriter, r *http.Request) {
    //vars := mux.Vars(r)
    //BuiltinMethodsHandler 
}



func Ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`Pong!`))
}



func ListDatasetFiles(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET"{
	WriteResp(w, 405, `Wrong method!`)
	return
    }

    vars := mux.Vars(r)

    path := os.Getenv("UPLOAD_DATASET_FILE_DIR") + vars[`userId`]
    glog.Info("ListDatasetFiles: List path: " + path)

    _, err := os.Stat(path)
    files, err := ioutil.ReadDir(path)
    if err != nil {
	WriteResp(w, 400, `Bad Request!`)
	glog.Error(err)
    }

    var avaliableDatasetFiles []datasetFileInfo
    for _, file := range files {
	newFile := datasetFileInfo{Name: file.Name(), Size: file.Size(), ModTime: file.ModTime()}
	avaliableDatasetFiles = append(avaliableDatasetFiles, newFile)
    }

    bytes, err1 := json.Marshal(&apiRequest{Status: "ok", Msg: `ok`, Payload: &avaliableDatasetFiles})
    if err1 != nil {
	WriteResp(w, 500, `Internal error. Try again later!`)
	glog.Error(err1)
	return
    }

    glog.Info( string(bytes) )
    glog.Info(avaliableDatasetFiles)
    w.Write(bytes)
}


func UploadDatasetFile(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST"{
	WriteResp(w, 405, `Wrong method!`)
	return
    }

    var err error
    defer func() {
	if nil != err {
	    glog.Error("Error string: errStr", err)
	    w.Write([]byte(`{"status":"error", "message":"` + err.Error()  + `"}`))
	}
    }()
    // max bytes in mem at a time  
    const _24K = (1 << 20) * 24
    err = r.ParseMultipartForm(_24K)
    if err != nil {
	WriteResp(w, 500, `Internal error. Try again later!`)
	return
    }
    for _, fheaders := range r.MultipartForm.File {
	for _, hdr := range fheaders {
	    // open uploaded  
	    var infile multipart.File
	    infile, err = hdr.Open()
	    if err != nil {
		WriteResp(w, 500, `Internal error. Try again later!`)
		return
	    }
	    // open destination  
	    var outfile *os.File
	    outfile, err = os.Create(os.Getenv("UPLOAD_DATASET_FILE_DIR") + hdr.Filename)
	    if err != nil {
		WriteResp(w, 500, `Internal error. Try again later!`)
		return
	    }
	    var written int64
	    written, err = io.Copy(outfile, infile)
	    if err != nil || 0 == written {
		WriteResp(w, 500, `Internal error. Try again later!`)
		return
	    }
	    datasetFileApiStruct := struct {
					DatasetFileId string
				    } {
					hdr.Filename,
				    }

	    bytes, err1 := json.Marshal(&apiRequest{Status: "ok", Msg: `ok`, Payload: datasetFileApiStruct})
	    if err1 != nil {
		WriteResp(w, 500, `Internal error. Try again later!`)
		glog.Error(err1)
		return
	    }
	    w.Write(bytes)
	}
    }
}

func WriteResp(w http.ResponseWriter, status int, msg string){
    http.Error(w, msg, status)
}
