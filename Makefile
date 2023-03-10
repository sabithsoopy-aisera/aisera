mod:
	go mod tidy && go mod vendor


setup: 
	which ginkgo || go install github.com/onsi/ginkgo/v2/ginkgo


test: setup
	ginkgo \
		--junit-report=report.xml \
		./...

