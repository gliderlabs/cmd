NAME := comlab
INSTALL_DIR ?= /usr/local/bin

install:
	glide install
	go build -o $(NAME) ./cmd/comlab
	install -m 755 $(NAME) $(INSTALL_DIR)/$(NAME)

test-env:
	docker build -t comlab-env -f dev/setup/Dockerfile .
	docker rmi comlab-env
