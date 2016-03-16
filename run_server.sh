echo "Starting server"
export UPLOAD_DATASET_FILE_DIR="/tmp/"
/usr/local/go/bin/go run server/server.go  -logtostderr=true
