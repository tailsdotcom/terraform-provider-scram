terraform-provider-scram: *.go */*.go
	go build .

install: terraform-provider-scram
	mkdir -p ~/.terraform.d/plugins/tails.com/tailsdotcom/scram/0.1.0/linux_amd64
	cp $+ ~/.terraform.d/plugins/tails.com/tailsdotcom/scram/0.1.0/linux_amd64
