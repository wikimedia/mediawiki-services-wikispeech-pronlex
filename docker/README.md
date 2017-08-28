## Docker installation

The `<DOCKERTAG>` variable represents the name of the Docker installation image.

### I. Install Docker CE

1. Install Docker CE for your OS: https://docs.docker.com/engine/installation/   
   * Ubuntu installation: https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/

2. Make sure you set all permissions and groups as specified in the installation instructions above. Log out and log in again.


### II. Obtain a Docker image

Obtain a Docker image using one of the following methods

* Build from a Dockerfile URL:

   `$ docker build https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t <DOCKERTAG>`   

   `<DOCKERTAG>` should normally be set to `stts-lexserver-local`.

* Build from local Dockerfile:

   `$ docker build $(go env GOPATH)/src/github.com/stts-se/pronlex -t <DOCKERTAG>`

   `<DOCKERTAG>` should normally be set to `stts-lexserver-local`.

* Download from Docker (not guaranteed to be updated on a regular basis):
 
  1. Visit the following URL and decide which version you want to install    
     https://hub.docker.com/r/sttsse/wikispeech/tags
 
  2. `$ docker pull sttsse/wikispeech:<TAGNAME>`

  Make a note that your `<DOCKERTAG>` for this version is `sttsse/wikispeech:<TAGNAME>`.
   	

Insert the `--no-cache` switch after the `build` tag if you encounter caching issues (e.g., with updated git repos, etc).


### III. Run the Docker app

Do set up and run the lexicon server, you need the [docker_run.sh](https://raw.githubusercontent.com/stts-se/pronlex/master/docker/docker_run.sh) script. It is a convenience script for calling `docker run` with a few switches.

1. Setup the lexicon server 

    `$ sh docker_run.sh -a <APPDIR> -t <DOCKERTAG> setup`

    Set up the server's required file structure in the specified `<APPDIR>`
  

2. Import lexicon files (optional)

    `$ sh docker_run.sh -a <APPDIR> -t <DOCKERTAG> import_all`

    Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.


3. Run lex server

    `$ sh docker_run.sh -a <APPDIR> -t <DOCKERTAG>`


You can also investigate the server environment using `bash`:   

`$ sh docker_run.sh -a <APPDIR> -t <DOCKERTAG> bash`
  

###
Server data files and databases are saved in the folder `<APPDIR>`. Please note that this folder will be owned by `root`. If this is a problem, make sure you change the ownership and/or permissions to whatever is best for your environmemnt.


<!-- to pass on system user to the Docker environment:
<!---   $ docker build --build-arg USER=$USER https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t stts-lexserver-local	 --->

<!---   $ docker build --build-arg USER=$USER $GOPATH/src/github.com/stts-se/pronlex -t stts-lexserver-local --->


