GOBIN=${PWD}/bin

SCRATCH?=n
NO_CACHE?=
ifeq ($(SCRATCH),y)
NO_CACHE=--no-cache
endif

# Default to build
all: sriov-dp httpd-init-image httpd-image dpdk-app vdpa-cni

help:
	@echo "Make Targets:"
	@echo " make dpdk-app         - Make the docker image that runs the DPDK l3fwd/l2fwd/testpmd."
	@echo "                         Append SCRATCH=y to build image using '--no-cache'."
	@echo " make httpd-init-image - Make the docker image that runs the Seastar httpd Init code."
	@echo "                         Append SCRATCH=y to build image using '--no-cache'."
	@echo " make httpd-image      - Make the docker image that runs the Seastar httpd."
	@echo "                         Append SCRATCH=y to build image using '--no-cache'."
	@echo " make sriov-dp         - Make the docker image that runs the SR-IOV Device"
	@echo "                         Plugin with vDPA changes integrated. Append SCRATCH=y"
	@echo "                         re-download upstream repo and to build image using '--no-cache'."
	@echo ""
	@echo " make                  - Build all the local sub-projects locally."
	@echo " make clean            - Cleanup all build artifacts."
	@echo " make all              - Build all images for a deployment. Same as:"
	@echo "                           make httpd-init-image; make httpd-image;"
	@echo "                           make sriov-dp;"
	@echo ""
	@echo "Local/Debug (not used in actual deployment):"
	@echo " make httpd-init       - Build the GO code that runs in the Seastar-httpd Init container."
	@echo ""
	@echo "Archive (not used anymore and may be deleted in future):"
	@echo " make scylla-init      - Build the GO code that runs in the Scylla Init container."
	@echo " make scylla-image     - Make the docker image that runs the Scylla Init code."
	@echo " make sriov-cni        - Make the SR-IOV CNI binary with the vDPA changes"
	@echo "                         integrated. Binary needs to copied to proper location"
	@echo "                         once complete (i.e. - /opt/cni/bin/.)."
	@echo ""
	@echo "Other:"
	@echo " glide update --strip-vendor - Recalculate dependancies and update *vendor\*"
	@echo "   with proper packages."
	@echo ""
#	@echo " make vdpa-cni-image   - Build the vDPA CNI in a docker image. When run as a"
#	@echo "                         daemonset, will install the built CNI binary in /opt/cni/bin/."
#	@echo " make sriov-cni-image  - Build the SR-IOV CNI as a docker image. When run as a"
#	@echo "                         daemonset, will install the built CNI binary in /opt/cni/bin/."


#
# Make Binaries
#
httpd-init:
	@cd seastar-httpd/init-container && go build -o ${GOBIN}/httpd-init -v

scylla-init:
	@cd scylla-init-container && go build -o ${GOBIN}/scylla-init -v

export ORG_PATH="github.com/intel"
export REPO_PATH_CNI="${ORG_PATH}/sriov-cni"
export REPO_PATH_DP="${ORG_PATH}/sriov-network-device-plugin"
export GOBIN=${PWD}/bin

# TBD
# sriov-cni: GOPATH=${PWD}/gopath
# sriov-cni:
# ifeq ($(SCRATCH),y)
# 	@rm -rf gopath/src/$(REPO_PATH_CNI)
# endif
# 	@if [ ! -d gopath/src/$(REPO_PATH_CNI) ]; then \
# 		echo ""; \
# 		echo "Making sriov-cni ..."; \
# 		echo "Downloading $(REPO_PATH_CNI)"; \
# 		mkdir -p gopath/src/$(ORG_PATH); \
# 		mkdir -p $(GOBIN); \
# 		pushd gopath/src/ > /dev/null; \
# 		go get $(REPO_PATH_CNI) 2>&1 > /tmp/sriov-cni.log || echo "Can ignore no GO files."; \
# 		popd > /dev/null; \
# 		echo "Patching $(REPO_PATH_CNI)"; \
# 		cp sriov-cni/* gopath/src/$(REPO_PATH_CNI)/.; \
# 		pushd gopath/src/$(REPO_PATH_CNI)/ > /dev/null; \
# 		mkdir -p pkg/vdpa/; \
# 		mv vdpadpdk-client.go pkg/vdpa/.; \
# 		mv vdpa.go pkg/vdpa/.; \
# 		patch -p1 < vdpa_cni_0001.patch; \
# 		echo "Glide Update"; \
# 		glide update --strip-vendor; \
# 		echo "Build CNI binary"; \
# 		make; \
# 		cp build/sriov $(GOBIN)/.; \
# 		echo ""; \
# 		echo "Run \"sudo cp bin/sriov /opt/cni/bin/.\""; \
# 		echo ""; \
# 		popd > /dev/null; \
# 	fi

# Make Docker Images
#
#vdpa-cni-image:
#	@docker build --rm -t vdpa-cni -f ./vdpa-cni/images/Dockerfile .

dpdk-app:
	@echo ""
	@echo "dpdk-app $(NO_CACHE) ..."
	@docker build $(NO_CACHE) --rm -t dpdk-app-centos -f ./dpdk-app-centos/Dockerfile .

httpd-image:
	@echo ""
	@echo "Making httpd-image $(NO_CACHE) ..."
	@docker build $(NO_CACHE) --rm -t seastar-httpd -f ./seastar-httpd/httpd/Dockerfile .

httpd-init-image:
	@echo ""
	@echo "Making httpd-init-image $(NO_CACHE) ..."
	@docker build $(NO_CACHE) --rm -t httpd-init-container -f ./seastar-httpd/init-container/Dockerfile .

scylla-image:
	@echo ""
	@echo "Making scylla-image $(NO_CACHE) ..."
	@docker build $(NO_CACHE) --rm -t scylla-init-container -f ./scylla-init-container/Dockerfile .

#sriov-cni-image: GOPATH=${PWD}/gopath
#sriov-cni-image:
#	@if [ ! -d gopath/src/$(REPO_PATH_CNI) ]; then \
#		echo ""; \
#		echo "Making sriov-cni-image ..."; \
#		echo "Downloading $(REPO_PATH_CNI)"; \
#		mkdir -p gopath/src/$(ORG_PATH); \
#		mkdir -p $(GOBIN); \
#		pushd gopath/src/ > /dev/null; \
#		go get $(REPO_PATH_CNI) 2>&1 > /tmp/sriov-cni.log || echo "Can ignore no GO files."; \
#		popd > /dev/null; \
#		echo "Patching $(REPO_PATH_CNI)"; \
#		cp sriov-cni/* gopath/src/$(REPO_PATH_CNI)/.; \
#		pushd gopath/src/$(REPO_PATH_CNI)/ > /dev/null; \
#		mkdir -p pkg/vdpa/; \
#		mv vdpadpdk-client.go pkg/vdpa/.; \
#		mv vdpa.go pkg/vdpa/.; \
#		patch -p1 < vdpa_cni_0001.patch; \
#		echo "Glide Update"; \
#		glide update --strip-vendor; \
#		echo "Build CNI binary"; \
#		docker build -t sriov-cni -f ./Dockerfile .; \
#		popd > /dev/null; \
#	fi

export ALT_DP_REPO=https://github.com/amorenoz/sriov-network-device-plugin.git
export ALT_DP_REF=vdpaInfoProvider

clean-sriov-dp:
	if [ -d gopath/src/$(REPO_PATH_DP) ]; then \
	    pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
	    make clean >/dev/null; \
	    popd > /dev/null; \
	    rm -fr gopath/src/$(REPO_PATH_DP); \
	fi \

ifeq ($(SCRATCH),y)
sriov-dp: clean-sriov-dp
 endif
sriov-dp: export GOPATH=${PWD}/gopath
sriov-dp:
	if [ ! -d gopath/src/$(REPO_PATH_DP) ]; then \
		echo ""; \
		echo "Making sriov-dp ..."; \
		echo "Downloading $(REPO_PATH_DP)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_DP) 2>&1 > /tmp/sriov-dp.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		if [ -n "$(ALT_DP_REPO)" ];  then \
		    pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
		    git remote add alt $(ALT_DP_REPO) && git fetch alt; \
		    if [ -n "$(ALT_DP_REF)" ]; then  \
			git checkout alt/$(ALT_DP_REF) -b alt; \
		    fi; \
		    popd > /dev/null; \
		fi; \
		pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
		echo "Build binary"; \
		make; \
		echo "Build docker image \"sriov-device-plugin\""; \
		make image; \
		popd > /dev/null; \
	fi

clean: clean-sriov-dp 
	@export GOPATH=${PWD}/gopath && go clean --modcache
	@rm -rf bin/
	@rm -rf gopath/

.PHONY: build clean sriov-dp clean-sriov-dp httpd-init-image httpd-image scylla-init scylla-image sriov-cni sriov-cni-image dpdk-app

