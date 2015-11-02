# Nibbleshell

Nibbleshell is a proxy server for processing images on the fly. It allows you to dynamically crop/resize/flip images hosted on S3, a local filesystem or an http source via query parameters. It supports creating “families” of images which can read from distinct image sources and enable different configuration values for image processing and retrieval.

Current version: `0.2.0`

## Concepts

### Sources

Sources are repositories from which an “original” image can be loaded. They return an image given a path. Currently, sources for downloading images from S3, a local filesystem and http are included.

### Routes

Routes bind URL rules (regular expressions) with a source and a processor. Nimbleshell supports setting up an arbitrary number of routes and sources.

When Nibbleshell receives a request, it determines the matching route, retrieves the image from its source, and processes the image according to the request processing options.

## Usage

To start the server pass a configuration file path as an argument.

```bash
$ ./bin/nibbleshell config.json
```

This will start the server on port 8080, and service requests whose path begins with /users/ or /blog/, e.g.:

    http://localhost:8080/users/joe/default.jpg?w=100&h=100
    http://localhost:8080/blog/posts/announcement.jpg?w=600&h=200

The image_host named group in the route pattern match (e.g., `^/users(?P<image_path>/.*)$`) gets extracted as the request path for the source. In this instance, the file “joe/default.jpg” is requested from the “my-company-profile-photos” S3 bucket. The processor resizes the image to a width and height of 100.

### Server

The `server` configuration block accepts the following settings:

##### port

The port to run the server on.

##### read_timeout

The timeout in seconds for reading the initial data from the connection.

##### write_timeout

The timeout in seconds for writing the image data backto the connection.

### Sources

The `sources` block is a mapping of source names to source configuration values.
Values from a source named `default` will be inherited by all other sources.

##### type

The type of image source. Currently `s3` or `filesystem`.

##### s3_access_key

For the S3 source type, the access key to read from S3.

##### s3_secret_key

For the S3 source type, the secret key to read from S3.

##### s3_bucket

For the S3 source type, the bucket to request images from.

##### directory

For the Filesystem source type, the local directory to request images from. Required.
For the S3 source type, `directory` corresponds to an optional base directory in the S3 bucket.

### Routes

The `routes` block is a mapping of route patterns to route configuration values.

The route pattern is a regular expression with a captured group for `image_path`.
The subexpression match is the path that is requested from the image source.

##### name

The name to use for the route. This is currently used in logging and StatsD key
names.

##### source

The name of the source to use for the route.

##### cache_control

The Cache-Control response header to set. If left empty or unspecified, `no-transform,public,max-age=86400,s-maxage=2592000` will be set.

### Health Checks

You can check the server health at `/health`. If the server is up and running, the HTTP client will receive a response with status code
`200`.
