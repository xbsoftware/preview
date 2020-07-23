Preview generation service
=============

## How to build

```
// txt, code
go build

// txt, code, images, documents
go build -tags extralibs
```

#### Build with docker

```
docker build . --tag preview
```

#### Dev. build with docker

```
docker build ./ -f Dockefile.build --tag preview:build
docker run --rm -it -p 3201:3201 -v `pwd`:/app registry.webix.io/preview:build /bin/bash
go build -tags extralibs
./preview
```


### Dependencies in extralibs mode

- libreoffice
- libreofficekit-dev
- libvips
- libvips-dev

#### Build with docker

just use the provided Dockerfile

## How to start

```
preview
```

#### configuration

Can be set through config.yml 

```yaml
port: 3201
key: any
uploadlimit: 33000000
text:
  fontdpi: 72
  fontsize: 12
  fontfile: "./fonts/DroidSansMono.ttf"
```

or through ENV parameters ( APP_TEXT_FONTSIZE, etc.)


**for all modes**

- port - http port of the service, default is 3201
- key - secret key
- uploadlimit - max file size allowed

**for code previews**

- fontdpi - screen resolution in Dots Per Inch
- fontfile - filename of the ttf font
- fontsize - font size in points

## How to use

### Preview 

Service expects POST request with a document and respond with preview image

```
POST http://localhost:3201/preview
```

#### Incoming parameters

form-multipart

- width - requested preview width, number in pixels
- height - requested preview height, number in pixels
- file - body of document for which preview need to be generated

#### Response

Response is a png of jpg image of preview

### Convert 

Service expects POST request with a document and respond with preview image

```
POST http://localhost:3201/convert
```

#### Incoming parameters

form-multipart

- name - output file name ( app.pdf )
- type - target type ( pdf, png, etc. )
- file - body of document for which preview need to be generated

#### Response

Response is a png of jpg image of preview

### Authorization

If `key` parameter was set in the configuration, the incoming request must contains matching URL parameter

```
POST https://some.com/preview?key={key}
```

### Supported document types

- Documents ( doc, docx, xls, xlsx, pdf )
- Images ( jpg, png, tiff, svg, gif, webp )
- Code files ( js, json, css, html, yaml, yml, xml, ts, md, ini, java, go, sql, sh  )

## License

MIT
