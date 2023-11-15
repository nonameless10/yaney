
#/bin/bash

cd idl && protoc -I. --go_out=plugins=grpc:../.. securekv.proto