.PHONY: build build-test test build-integration-test integration-test test-event

build:
	docker build -t pdf2png src

test: build
	docker run --rm -v ${PWD}/test:/tmp pdf2png --test

build-integration-test: build
	docker build -t pdf2png-test test

integration-test: build-integration-test
	docker run --rm -v ${HOME}/.aws:/root/.aws -e AWS_REGION -e AWS_PROFILE -p 9000:8080 pdf2png-test

test-event:
	sed -e s/TEST_BUCKET/${BUCKET}/ -e s/TEST_KEY/${KEY}/ test/event.json | \
		curl -d @- http://localhost:9000/2015-03-31/functions/function/invocations
