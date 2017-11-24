Instructions for expermental use, may not be updated with each release.

For info on recommended installation procedure(s), see the <a href="http://stts-se.github.io/wikispeech/">main doc page for Wikispeech</a>.

## Docker installation

### I. Install Docker CE

1. Install Docker CE for your OS: https://docs.docker.com/engine/installation/   
   * Ubuntu installation: https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/

2. If you're on Linux, make sure you set all permissions and groups as specified in the post-installation instructions: https://docs.docker.com/engine/installation/linux/linux-postinstall/ 


### II. Obtain a Docker image

1. Visit the following URL and decide which release (`<TAGNAME>` below) you want to install   
   https://hub.docker.com/r/sttsse/pronlex/tags/

2. `$ docker pull sttsse/pronlex:<TAGNAME>`


### III. Run the Docker app

To set up and run the lexicon server, you will use the [docker_run.sh](https://raw.githubusercontent.com/stts-se/pronlex/master/docker/docker_run.sh) script. It is a convenience script for calling `docker run` with a few switches. Requires a working `bash` installation on Linux.


1. Import lexicon files (optional)

    `$ bash docker_run.sh -a <APPDIR> -t sttsse/pronlex:<TAGNAME> import_all`

        Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.


3. Run lex server

      `$ bash docker_run.sh -a <APPDIR> -t sttsse/pronlex:<TAGNAME>`


You can also investigate the server environment using `bash`:   

`$ bash docker_run.sh -a <APPDIR> -t sttsse/pronlex:<TAGNAME> bash`
  

###
Server data files and databases are saved in the folder `<APPDIR>`. Please note that this folder will be owned by `root`. If this is a problem, make sure you change the ownership and/or permissions to whatever is best for your environmemnt.
