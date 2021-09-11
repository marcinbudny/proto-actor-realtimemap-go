# To run these commands 
# 1. Install protoc compiler (e.g using chocolatey.org if on Windows)
#
# choco install protoc 
#
# 2. Install plugins
# 
# go install github.com/gogo/protobuf/protoc-gen-gogoslick@latest
# go install github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograinv2@latest

protoc -I=. -I=$GOPATH/src --gogoslick_out=. protos.proto
protoc -I=. -I=$GOPATH/src --gograinv2_out=. protos.proto