PATH=$PATH:$GOPATH/bin protoc --go_out=gauge/ gauge/messages.proto
PATH=$PATH:$GOPATH/bin protoc --go_out=gauge/ gauge/api.proto
protoc --java_out=gauge-java/src/main/java/ gauge/messages.proto
protoc --java_out=gauge-java/src/main/java/ gauge/api.proto
#ruby-protoc -o ruby messages.proto 
