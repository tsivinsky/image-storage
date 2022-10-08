# image-storage

REST API for uploading and serving images

## Usage

### POST /api/images/upload

allows upload new images

#### Request body

```
image={file}
```

### GET /:filename

serves uploaded earlier images
