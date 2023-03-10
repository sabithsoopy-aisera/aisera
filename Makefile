mod:
	go mod tidy && go mod vendor


setup: 
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo


test: setup
	ginkgo -v \
		--cover \
		--fail-on-pending \
		--covermode=set \
		--output-dir=test-results \
		--coverprofile=coverage.out \
		--junit-report=report.xml \
		./...

