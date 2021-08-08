
SERVER := localhost:61535

.PHONY:app-test1
app-test1:
	@curl -I --header 'X-ECHO-CODE: 404' $(SERVER)
	@curl -I $(SERVER)/?echo_code=404

.PHONY: app-test2
app-test2:
	@curl --header 'X-ECHO-BODY: amazing' $(SERVER)
	@curl $(SERVER)/?echo_body=amazing

.PHONY: app-test3
app-test3:
	@curl --header 'X-ECHO-TIME: 5000' $(SERVER)
	@curl $(SERVER)/?echo_time=5000

