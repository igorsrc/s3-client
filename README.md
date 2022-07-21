# App description
Filestorage client based on s3 interface

# Config example
```
app:
  port: "8080"
  max-upload-size: 100
s3:
  host: "127.0.1.1"
  port: "9000"
  access: "access_key"
  secret: "secret_key"
  bucket: "dev"
  region: "eu-west-1"
```