package main

import (
    "net/http"
    "encoding/json"
    "mime/multipart"
    "os"
    "io"
    "io/ioutil"
    "flag"

    "github.com/golang/glog"
    "github.com/gorilla/mux"
)


func main() {
    flag.Parse()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/api/v1/ping", Ping)
    router.HandleFunc("/api/v1/dataset_file", UploadDatasetFile)
    router.HandleFunc("/api/v1/compute_request/builtin/{builtinJobId}/{datasetFileId}", BuiltinJob)

    glog.Info(http.ListenAndServe(":9001", router))
}


func BuiltinJob(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    bytes, err := ioutil.ReadFile( os.Getenv("UPLOAD_DATASET_FILE_DIR") + vars[`datasetFileId`])
    // File doesn't exist
    // Doesn't have read permissions
    if err != nil {
	glog.Error("Error when oppening file: ", err) //no such file or directory
    }
    y := [][]int{}
    err  = json.Unmarshal(bytes, &y)
    if err != nil {
	//cannot unmarshal bad data
	glog.Error(err)
    }

    if vars[`builtinJobId`] == `get_max_value` {
	for _, ele := range y {
	    for _, ele1 := range ele {
		glog.Info("ELEMENT: ", ele1)
		//w.Write()
	    }
	}

    }
}



func Ping(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, `http://www.google.bg`, http.StatusMovedPermanently)
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
    if nil != err {
	return
    }
    for _, fheaders := range r.MultipartForm.File {
	for _, hdr := range fheaders {
	    // open uploaded  
	    var infile multipart.File
	    infile, err = hdr.Open()
	    if nil != err {
		 return
	    }
	    // open destination  
	    var outfile *os.File
	    outfile, err = os.Create(os.Getenv("UPLOAD_DATASET_FILE_DIR") + hdr.Filename)
	    if nil != err {
		 return
	    }
	    var written int64
	    written, err = io.Copy(outfile, infile)
	    if nil != err || 0 == written {
		 return
	    }
	    w.Write([]byte(`{"status":"ok", "resource_id":"` + hdr.Filename  + `"}`))
	}
    }
}

func WriteResp(w http.ResponseWriter, status int, msg string){
    http.Error(w, msg, status)
}
