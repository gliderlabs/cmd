
service:
	sigil -f secrets.yaml.tmpl | kubectl apply -f -
	kubectl apply -f dockerbox.yaml
	kubectl apply -f dynamodb.yaml
	kubectl apply -f cmd.yaml
