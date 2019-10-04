GOBIN=${PWD}/bin

# Default to build
default: server client
all: server client


help:
	@echo "Make Targets:"
	@echo " make                - Build all the local sub-projects locally."
	@echo " make clean          - Cleanup all build artifacts."
	@echo " make server         - Build the GO code that handles the gRPC Server."
	@echo " make server-image   - Make the docker image that runs the gRPC Server code."
	@echo " make client         - Build the GO code that handles the gRPC Client (test code)."
	@echo " make client-image   - Make the docker image that runs the gRPC Client test code."
	@echo ""
	@echo "Other:"
	@echo " glide update --strip-vendor - Recalculate dependancies and update *vendor\*"
	@echo "   with proper packages."
	@echo ""


test:
	@cd cnivpp/test/memifAddDel && go build -v
	@cd cnivpp/test/vhostUserAddDel && go build -v
	@cd cnivpp/test/ipAddDel && go build -v

server:
	@cd server-image && go build -o ${GOBIN}/vdpa-server -v

client:
	@cd client-image && go build -o ${GOBIN}/vdpa-client -v


server-image:
	@docker build --rm -t nfvpe/vdpa-grpc-server -f ./server-image/Dockerfile .

client-image:
	@docker build --rm -t nfvpe/vdpa-grpc-client -f ./client-image/Dockerfile .

vdpa-image:
	@docker build --rm -t nfvpe/vdpa-daemonset -f ./vdpa-dpdk-image/Dockerfile .


clean:
	@rm -rf bin/

generate:

lint:

.PHONY: build test install extras clean generate server-image client-image

