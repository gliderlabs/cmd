
service:
	sigil -f secrets.yaml.tmpl | kubectl apply -f -
	kubectl apply -f ../kubernetes/dockerbox.yaml
	kubectl apply -f dynamodb.yaml
	kubectl apply -f cmd.yaml
