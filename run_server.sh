echo "Starting server"

#Env setup
export UPLOAD_DATASET_FILE_DIR="/tmp/"
export WEBUI_DIR="" 

#Build
cd server
/usr/local/go/bin/go build 

#Run
cd ..
./server/server  -logtostderr=true
