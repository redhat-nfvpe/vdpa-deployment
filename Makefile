GOBIN=${PWD}/bin

# Default to build
default: build
all: build


help:
	@echo "Make Targets:"
	@echo " make                - Build all the local sub-projects locally."
	@echo " make clean          - Cleanup all build artifacts."
	@echo " make server         - Build the GO code that handles the gRPC Server."
	@echo " make server-image   - Make the docker image that runs the gRPC Server code."
	@echo " make init           - Build the GO code that scans the vDPA VF Interfaces"
	@echo "                       and writes the associated PCI Addresses to a file."
	@echo " make init-image     - Make the docker image that runs the init code"
	@echo "                       (as init container)."
	@echo ""
	@echo "Other:"
	@echo " glide update --strip-vendor - Recalculate dependancies and update *vendor\* with proper packages."
	@echo ""
#	@echo "Makefile variables (debug):"
#	@echo "   SUDO=$(SUDO) OS_ID=$(OS_ID) OS_VERSION_ID=$(OS_VERSION_ID) PKG=$(PKG) VPPVERSION=$(VPPVERSION) $(VPPDOTVERSION)"
#	@echo "   VPPLIBDIR=$(VPPLIBDIR)"
#	@echo "   VPPINSTALLED=$(VPPINSTALLED) VPPLCLINSTALLED=$(VPPLCLINSTALLED)"
#	@echo ""


test:
	@cd cnivpp/test/memifAddDel && go build -v
	@cd cnivpp/test/vhostUserAddDel && go build -v
	@cd cnivpp/test/ipAddDel && go build -v

server:
	@cd server-image && go build -o ${GOBIN}/vdpa-server -v

client:
	@cd client-image && go build -o ${GOBIN}/vdpa-client -v

init:
	@cd init-image && go build -o ${GOBIN}/vdpa-init -v


server-image:
	@docker build --rm -t nfvpe/vdpa-grpc-server -f ./server-image/Dockerfile .


clean:
	@rm -rf bin/

generate:

lint:

.PHONY: build test install extras clean generate server-image

