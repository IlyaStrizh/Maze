CC= ~/go/bin/fyne
NAME= Maze

ifeq ($(shell uname -s), Linux)
OS= linux

else ifeq ($(shell uname -s), Darwin)
OS= darwin

endif

all: style test build

deps:
	@command -v $(CC) >/dev/null 2>&1 || { echo >&2 "Fyne is not installed. Installing..."; go install fyne.io/fyne/v2/cmd/fyne@latest; }

install: deps
	@cd cmd && $(CC) package -os $(OS) -name $(NAME) -icon ../images/maze.png && \
	ln -s $(PWD)/cmd/$(NAME).app ~/Desktop/$(NAME).app

dist: deps
	@cd cmd && $(CC) package -os $(OS) -name $(NAME) -icon ../images/maze.png && \
	cd $(PWD)/cmd && tar -czf ~/Desktop/$(NAME).tar.gz $(NAME).app

uninstall: clean
	rm -rf $(PWD)/cmd/$(NAME).app
	rm -f ~/Desktop/$(NAME).app ~/Desktop/$(NAME).app.tar.gz

build:
	@cd $(PWD)/cmd && go build -o $(NAME)

tests:
	@cd $(PWD)/internal/model && go test -v -cover -coverprofile=coverage.out ./...
	@mv $(PWD)/internal/model/*.txt .
	@mv $(PWD)/internal/model/*.png .

gcov_report: tests
	@cd $(PWD)/internal/model && go tool cover -html=coverage.out

dvi:
	@texi2html $(PWD)/info/info.texi
	@open info.html

style:
	@gofmt -d .

clean:
	@rm -f ./cmd/$(NAME) ./internal/model/*.out
	@rm -f *.txt *.png *.html