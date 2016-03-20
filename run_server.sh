echo "Starting server"
export UPLOAD_DATASET_FILE_DIR="/tmp/"
cd server
/usr/local/go/bin/go build 
cd ..
./server/server  -logtostderr=true
