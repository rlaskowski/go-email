GOBIN := go

BUILDNAME := go-email
CMDPATH := cmd/main.go


run:
	go run cmd/main.go -f /home/rafal/tmp/email-store
all:
	@-${MAKE} clean
	@-${MAKE} build-all 
build-all: 
	@-${MAKE} build-windows
	@-${MAKE} build-linux
	@-${MAKE} build-darwin
	@echo "Finished building all OS"
build-windows:
	@echo "Building to windows..."
	@GOOS=windows ${GOBIN} build -o dist/windows/${BUILDNAME}.exe ${CMDPATH}
build-linux:
	@echo "Building to linux..."
	@GOOS=linux ${GOBIN} build -o dist/linux/${BUILDNAME} ${CMDPATH}
build-darwin:
	@echo "Building to darwin..."
	@GOOS=darwin ${GOBIN} build -o dist/darwin/${BUILDNAME} ${CMDPATH}
clean:
	@echo "Cleaning..."
	@rm -Rf dist
	@echo "All cleaned"
generate-protobuf:	
	@echo "Generating new proto files..."
	@find grpc/protobuf -type f -iname \*.pb.go -delete
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/protobuf/emailservice/*