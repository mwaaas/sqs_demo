cache_dir=.cache
cache_from := $(shell [ -f $(cache_dir)/index.json ] && echo "--cache-from=type=local,src=$(cache_dir)" || echo )

build:
	docker buildx build $(cache_from) --cache-to=type=local,dest=$(cache_dir) --output=type=docker,name=sqs_demo_sqs_demo .

push_image:
	docker tag sqs_demo_sqs_demo mwaaas/sqs_go_demo
	docker push mwaaas/sqs_go_demo

build_and_push: push_image build

list_queues:
	 aws --endpoint-url=http://localhost:4100 sqs list-queues