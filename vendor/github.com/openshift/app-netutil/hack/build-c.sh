set -e

ORG_PATH="github.com/openshift"
REPO_PATH="${ORG_PATH}/app-netutil"

if [ ! -h gopath/src/${REPO_PATH} ]; then
        mkdir -p gopath/src/${ORG_PATH}
        ln -s ../../../.. gopath/src/${REPO_PATH} || exit 255 
fi

export GOBIN=${PWD}/bin
export GOPATH=${PWD}/gopath
export CGO_ENABLED=1

#go install "$@" ${REPO_PATH}/samples/go_app
go build -o ${GOBIN}/libnetutil_api.so -buildmode=c-shared -v ${REPO_PATH}/c_api
gcc -I${GOBIN} -L${GOBIN} -Wall -o ${GOBIN}/c_sample ${GOPATH}/src/${REPO_PATH}/samples/c_app/app_sample.c -lnetutil_api
