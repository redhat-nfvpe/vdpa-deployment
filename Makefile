GOBIN=${PWD}/bin

# Default to build
default: server client
local: server client
all: server-image vdpa-image sriov-dp sriov-cni

help:
	@echo "Make Targets:"
	@echo " make                - Build all the local sub-projects locally."
	@echo " make clean          - Cleanup all build artifacts."
	@echo " make server-image   - Make the docker image that runs the gRPC Server code."
	@echo " make vdpa-image     - Make the docker image that runs the DPDK vDPA sample"
	@echo "                       APP. Manages the socketfiles for host."
	@echo " make sriov-dp       - Make the docker image that runs the SR-IOV Device"
	@echo "                       Plugin with vDPA changes integrated."
	@echo " make sriov-cni      - Make the SR-IOV CNI binary with the vDPA changes"
	@echo "                       integrated. Binary needs to copied to proper location"
	@echo "                       once complete (i.e. - /opt/cni/bin/.)."
	@echo " make all            - Build all images for a deployment. Same as:"
	@echo "                         make server-image; make vdpa-image; make sriov-dp; make sriov-cni"
	@echo ""
	@echo "Local/Debug (not used in actual deployment):"
	@echo " make server         - Build the GO code that handles the gRPC Server."
	@echo " make client         - Build the GO code that handles the gRPC Client (test code)."
	@echo " make client-image   - Make the docker image that runs the gRPC Client test code."
	@echo " make local          - Build the GO code locally, same as:"
	@echo "                         make server; make client"
	@echo ""
	@echo "Other:"
	@echo " glide update --strip-vendor - Recalculate dependancies and update *vendor\*"
	@echo "   with proper packages."
	@echo ""


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



export ORG_PATH="github.com/intel"
export REPO_PATH_CNI="${ORG_PATH}/sriov-cni"
export REPO_PATH_DP="${ORG_PATH}/sriov-network-device-plugin"
export GOBIN=${PWD}/bin

sriov-cni: GOPATH=${PWD}/gopath
sriov-cni:
	@if [ ! -d gopath/src/$(REPO_PATH_CNI) ]; then \
		echo "Downloading $(REPO_PATH_CNI)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		mkdir -p $(GOBIN); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_CNI) 2>&1 > /tmp/sriov-cni.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		echo "Patching $(REPO_PATH_CNI)"; \
		cp sriov-cni/* gopath/src/$(REPO_PATH_CNI)/.; \
		pushd gopath/src/$(REPO_PATH_CNI)/ > /dev/null; \
		mkdir -p pkg/vdpa/; \
		mv vdpadpdk-client.go pkg/vdpa/.; \
		mv vdpa.go pkg/vdpa/.; \
		patch -p1 < vdpa_cni_0001.patch; \
		echo "Glide Update"; \
		glide update --strip-vendor; \
		echo "Build CNI binary"; \
		make; \
		cp build/sriov $(GOBIN)/.; \
		echo "Run \"sudo cp bin/sriov /opt/cni/bin/.\""; \
		popd > /dev/null; \
	fi

sriov-dp: GOPATH=${PWD}/gopath
sriov-dp:
	@if [ ! -d gopath/src/$(REPO_PATH_DP) ]; then \
		echo "Downloading $(REPO_PATH_DP)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_DP) 2>&1 > /tmp/sriov-dp.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		echo "Patching $(REPO_PATH_DP)"; \
		cp sriov-dp/* gopath/src/$(REPO_PATH_DP)/.; \
		pushd gopath/src/$(REPO_PATH_DP)/ > /dev/null; \
		patch -p1 < vdpa_dp_0001.patch; \
		echo "Build binary"; \
		make; \
		echo "Build docker image \"nfvpe/sriov-device-plugin\""; \
		make image; \
		popd > /dev/null; \
	fi

clean:
	@rm -rf bin/
	@rm -rf gopath/

.PHONY: build clean server client server-image client-image sriov-cni sriov-dp

