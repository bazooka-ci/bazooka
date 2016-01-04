default: docker

docker: gox
	docker build -t bazooka/worker -f Dockerfile .

gox:
	gox -osarch="linux/amd64" -output="worker"

run: docker
	-docker rm -f bzk_worker
	docker run -d --name bzk_worker \
		-e BZK_API_URL=http://bzk_server:3000 \
		-e BZK_SYSLOG_URL=tcp://192.168.99.100:3001 \
		-e BZK_QUEUE_URL=bzk_queue:11300 \
		-e BZK_HOME=$(BZK_HOME) \
		-e BZK_SCM_KEYFILE=/Users/jawher/.ssh/id_rsa \
		-e BZK_SLOTS=2 \
		-e BZK_DOCKERSOCK=$(BZK_DOCKERSOCK) \
		-e BZK_NETWORK=bzk_net \
		-v $(BZK_HOME):/bazooka \
		-v $(BZK_DOCKERSOCK):$(BZK_DOCKERSOCK) \
		--net=bzk_net \
		bazooka/worker

queue:
	docker run -d --name bzk_queue -p 11300:11300 --net=bzk_net jawher/beanstalkd

