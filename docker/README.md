## Docker installation

TODO/WORK IN PROGRESS

1. Build

  1. from Dockerfile URL:

   `$ cd <DOCKER DIR>`   
   `$ docker build --no-cache https://raw.githubusercontent.com/stts-se/pronlex/master/docker/Dockerfile -t sttsse/lexserver`   

  2. from local Dockerfile:

   `$ cd <DOCKER DIR>`   
   `$ docker build --no-cache $GOPATH/src/github.com/stts-se/pronlex/docker -t sttsse/lexserver`   

3. Setup server

`$ docker run -v <DOCKERDIR>/lexserver_files:/go/lexserver_files -p 8787:8787 -it sttsse/lexserver sh setup`

4. Import lexicon files (optional)

`$ docker run -v <DOCKERDIR>/lexserver_files:/go/lexserver_files -p 8787:8787 -it sttsse/lexserver sh import_lex`


