# ImageBox
ImageBox is a service that is used to upload files that are uploaded 
to an external Object Storage. You can only upload images in png & jpeg formats, 
no larger than 10 mb.
---

Technologies used:
* [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/)
* [REST](https://ru.wikipedia.org/wiki/REST)
* [JWT Authentication](https://jwt.io/)
* [MinIO](https://min.io/)
---

## Getting started
This is an example of how you may give instructions on setting up your project locally. To get a local copy up and running follow these example steps.

### Installation
1. Clone the repository:
```shell
git clone https://github.com/zenorachi/image-box
```
2. Setup environment variables (create .env file in the project's root):
```dotenv
export DB_HOST=
export DB_PORT=

export DB_USERNAME=
export DB_NAME=
export DB_SSLMODE=
export DB_PASSWORD=

export POSTGRES_PASSWORD=

export MINIO_ROOT_USER=
export MINIO_ROOT_PASSWORD=

export HASH_SALT=
export HASH_SECRET=
```
> **Note:** if you build the project using Docker, setup *DB_HOST=db* (as the container name)
3. Compile and build the project:
```shell
make build
```