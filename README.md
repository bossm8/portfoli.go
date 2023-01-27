# Portfoli.go

<p align="center">
    <img src="./public/img/portfoli.go-yellow.svg" style="width: 200px" />
</p>

The simple and flexible portfolio template written witgh [Go](https://golang.org) and [Bootstrap](https://getbootstrap.com)


## Static Build

```bash
docker run -it --rm -p 8080:80 \
           --name porfoli.go \
           -v ${PWD}/dist:/usr/share/nginx/html \
           -v ${PWD}/nginx.conf:/etc/nginx/conf.d/default.conf \
           nginx:latest
```