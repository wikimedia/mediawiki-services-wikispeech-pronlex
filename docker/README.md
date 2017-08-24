## Docker installation

WORK IN PROGRESS

1. Build the Docker image

    1. from Dockerfile URL:

        `$ docker build https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t sttsse/lexserver-local`   

    2. from local Dockerfile:

        `$ docker build $GOPATH/src/github.com/stts-se/pronlex -t sttsse/lexserver-local`

    Insert the `--no-cache` switch after the `build` tag if you encounter caching issues (updated git repos, etc).


2. Run the docker app


   1. Setup the server 

      `$ run_docker.sh setup`


   2. Import lexicon files (optional)

      `$ run_docker.sh import_lex`


   3. Run lex server

      `$ run_docker.sh`


   You can also investigate the server environment using `bash`:   
   `$ run_docker.sh bash`
  

