## Docker installation

TODO/WORK IN PROGRESS

Build from Dockerfile URL:

`$ cd <DOCKER DIR>`   
`$ docker build --no-cache https://raw.githubusercontent.com/stts-se/pronlex/master/docker/Dockerfile -t sttsse/lexserver`   
`$ docker run -v <DOCKERDIR>/lexserver_files:/go/lexserver_files -p 8787:8787 -it sttsse/lexserver sh import_lex`


Build from local Dockerfile:

`$ cd <DOCKER DIR>`   
`$ docker build --no-cache $GOPATH/src/github.com/stts-se/pronlex/docker -t sttsse/lexserver`   
`$ docker run -v <DOCKERDIR>/lexserver_files:/go/lexserver_files -p 8787:8787 -it sttsse/lexserver sh import_lex`


