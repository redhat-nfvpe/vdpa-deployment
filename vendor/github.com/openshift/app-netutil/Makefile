go_sample:
	./hack/build.sh
c_sample:
	./hack/build-c.sh
dpdk_app:
	./hack/build-dpdkapp.sh
testpod:
	./hack/build-testpod.sh
# 'make' and 'make image' are left around for legacy, but not
# documented anywhere.
default:
	./hack/build.sh
image:
	./hack/build-testpod.sh
clean:
	rm -rf gopath/
	rm -rf bin/
