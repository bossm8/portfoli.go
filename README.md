# Portfoli.go

<p align="center">
    <img src="./public/img/portfoli.go-yellow.svg" style="width: 200px" />
</p>

The simple and flexible portfolio template written with [Go](https://golang.org) and [Bootstrap](https://getbootstrap.com)
Build your portfolio with simple yaml files!

**See the live example on [bossm8.ch](https://me.bossm8.ch) or [GitHub Pages](https://bossm8.github.io/portfoli.go)** 

This template can be used for either hosting a static webpage with e.g. 
[GitLab](https://docs.gitlab.com/ee/user/project/pages/) / [GitHub](https://pages.github.com) pages. Or if you
like, with a server written in go, this server brings benefits like a contact form to send emails directly
to you via the portfolio page.

## Getting Started

(For more detailed information see below)

### Dynamic

1. Create a directory containing two folders: `configs` and `custom`
2. Add your own images into the directory `custom`
3. Copy the example yaml configurations into the `configs` directory and adjust them (for images use the path `/static/img/custom/<your-image>`)
4. Use the prebuilt docker container to host the portfolio (the command is being run in the directory containing the two folders):
   (make sure the files have the right permissions (container uses 100:82) or run the container as a different user)
   ```bash
   docker run -it --rm \
            --name portfoli.go \
            -p 127.0.0.1:8080:8080 \
            -v ${PWD}/configs:/var/www/portfoli.go/configs:ro \
            -v ${PWD}/content:/var/www/portfoli.go/public/img/custom:ro \
            ghcr.io/bossm8/portfoli-go:latest
   ```
5. This command will mount your configs and custom images in the expected location inside the docker container and start the portfolio with your content.

### Static

1. Do steps 1-2 of the [dynamic](#dynamic) approach
2. Run the following commands to build your static webpage with the prebuilt docker container:
   ```bash
   # create the directory for the static build
   mkdir dist
   # render the html content
   docker run -it --rm \
            --name portfoli.go \
            -v ${PWD}/configs:/var/www/portfoli.go/configs:ro \
            -v ${PWD}/dist:/var/www/portfoli.go/dist \
            -e DIST_PATH=/var/www/portfoli.go/dist \
            -e CONF_PATH=/var/www/portfoli.go/configs \
            --entrypoint portfoli-go-static.sh \
            ghcr.io/bossm8/porfoli-go:latest
    # move your custom images into the static build
    mv ${PWD}/custom dist/static/img/custom
    ```
3. Host the build with any fileserver of your choice (see for example the nginx configuration [below](#local))


## Configuration

The portfolio template expects its content to come from yaml configuration files.
Those files need to be in a single directory and their names must be the same as the 
content type they represent (with the `.yml` extension). They are loaded each time
the page is requested, meaning you can edit them on the fly without having to restart
the server (when not using the static version).

The following content types are currently supported:

* experience
* education
* projects
* certifications
* bio

Each of them might support a different configuration, for possible values and explanaiton see `examples/configs`.

**NOTE** Any *HTML* content in the configurations may 
also contain [go templates](https://pkg.go.dev/text/template), it will be passed through 
the templating engine when loaded. You might want to use the `Assemble` function, which
adds the configured base path to relative paths from `/static`, i.e. you can reference your custom images like this:

```go
{{ "/static/custom/avatar.jpg" | Assemble }}
```

So you do not have to adjust the configurations should the base path ever change.

### Recommendations

I recommend putting your custom content into a subdirectory of `public/img` (e.g. `custom`), and referncing
this directory in the `yaml` configurations when specifying images. This makes it easier for the usage with e.g. Docker. 
You are able to use a relative path starting with `/static` i.e. `/static/img/custom/avatar.jpg` the rendering process will 
make sure that any base path (specified with `-srv.base`) of your server is prepended to this path (e.g. when hosing on GitLab pages).


## Usage

There are different approaches on how to use this template, select the one which might fit you the most.
(`Live` means you use the binary to serve the portfolio, which enables the contact form.)

* [Live (Source)](#build-from-Source)
* [Live (Docker)](#docker)
* [Static (GitLab/GitHub)](#gitlab--github-pages)
* [Static (local)](#local)

### Build from Source 

Get started easily by pulling this repository and running `make run`, this will
start the portfolio with the example configuration (go is required).

The portfoli-go binary can be built with `make build`, then use
the help message to see available options:

```bash
portfoli-go -help
```

### Docker

There exists a pre-build Docker image which you can use to host the portfolio website 
just mount your custom content and configurations to use:
(make sure the files have the right permissions (container uses 100:82) or run the container as a different user)

```bash
docker run -it --rm \
           --name portfoli.go \
           -p 127.0.0.1:8080:8080 \
           -v ${PWD}/configs:/var/www/portfoli.go/configs:ro \
           -v ${PWD}/content:/var/www/portfoli.go/public/img/custom:ro \
           ghcr.io/bossm8/portfoli-go:latest
```

There is also a sample `docker-compose.yml` in the `examples` directory.

### Static Build

As stated earlier, you are able to build the portfolio for being hosted as a static website.
The static build can be used on e.g. [GitLab](https://docs.gitlab.com/ee/user/project/pages/) 
or [GitHub](https://pages.github.com) pages.

It can be built by using the `-dist flag` with the binary or locally with `make dist`, this will output
the content for being served with a static file server in the specified output directory.
However, when using the binary you need to make sure to also copy over the contents of the directory
`public` into the dist path (see for example the `script` in `examples/.gitlab-ci.yml`). 
As described in the [config](#recommendations) section, I recommend putting 
custom images into a subdirectory of `public/img` and specifying the corresponding path in the yaml configs.

#### GitLab / GitHub Pages

For building with GitLab or GitHub you may use the Docker image of portfoli.go (ghcr.io/bossm8/portfoligo:latest) 
as there is everything prepared inside, just have a repository with your custom content and yaml configuration
which you can use to build the static website. There are pipeline configurations examples for bothin the `examples` directory.
Just copy one them over to your repository containing configuration files and custom content and run the pipeline
(some adjustemts might be needed though - e.g. `BASE_PATH`).

#### Local

For running locally, there is an example [nginx](https://nginx.com) configuration which shows how the dist build may be used.
The command below starts this configuration with the nginx [Docker container](https://hub.docker.com/_/nginx) - 
it assumes you have run `make dist` before and that you have copied your custom content into the `dist` directory.

```bash
docker run -it --rm -p 8080:80 \
           --name porfoli.go \
           -v ${PWD}/dist:/usr/share/nginx/html:ro \
           -v ${PWD}/examples/nginx.conf:/etc/nginx/conf.d/default.conf:ro \
           nginx:latest
```

## Development

There is a [devcontainer](https://code.visualstudio.com/docs/devcontainers/containers) 
setup (which can also be used for building the binary). Simply fork (or clone) this repository,
open the devcontainer with [VSCode](https://code.visualstudio.com/) and run `make setup`.

## Authors

* [bossm8](https://github.com/bossm8)

## TODO

* Write tests

## Gopher Artwork

Here are some gopher images I created for this page:

![NotFound](public/img/status/404.svg)
![Error](public/img/status/error.svg)
![Delivered](public/img/status/delivered.svg)
![Undelivered](public/img/status/undelivered.svg)