golint | egrep -v  "_(test|string)\.go:.*(don't use underscores in Go names)" && go vet # && go generate
