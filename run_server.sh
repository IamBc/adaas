
#Env setup
echo 'Setting the envornment...'
export UPLOAD_DATASET_FILE_DIR="/tmp/"
export WEBUI_DIR="./web-client" 

#Build
echo 'Building the project...'
cd server
/usr/local/go/bin/go build 

#Run
echo 'Running the project...'
cd ..
./server/server  -logtostderr=true
