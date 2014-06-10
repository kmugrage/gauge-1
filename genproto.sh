#!/bin/sh
PATH=$PATH:$GOPATH/bin protoc --go_out=. gauge/spec.proto
PATH=$PATH:$GOPATH/bin protoc --go_out=. gauge/messages.proto
PATH=$PATH:$GOPATH/bin protoc --go_out=. gauge/api.proto
protoc --java_out=gauge-java/src/main/java/ gauge/spec.proto
protoc --java_out=gauge-java/src/main/java/ gauge/messages.proto
protoc --java_out=gauge-java/src/main/java/ gauge/api.proto
#ruby-protoc -o ruby messages.proto 
