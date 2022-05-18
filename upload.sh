set -e

IMAGE_NAME="subscription-service"
SERVER="fliptable-2"
GO_BUILD="build"

BUILD_NAME="$IMAGE_NAME.tar"

# build go build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $GO_BUILD

sudo docker build --platform amd64 . -t $IMAGE_NAME
sudo docker save -o $BUILD_NAME $IMAGE_NAME
sudo docker system prune
sudo chown "$USER" $BUILD_NAME

scp $BUILD_NAME $SERVER:~
rm -rf $BUILD_NAME, $GO_BUILD

ssh $SERVER "
  sudo docker load -i $BUILD_NAME
  rm -rf $BUILD_NAME
  cd ~/containers/$IMAGE_NAME/
  sudo ./run.sh
  sudo docker system prune
  exit
"