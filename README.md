# Run script for create and start docker container
```sh
chmod +x ./docker_start.sh && ./docker_start.sh
```

# Run with command
## Create docker container
```sh
docker build -t client-pow-image .
```
## Run docker container
```sh
docker run --rm -it \
--name client-pow \
--link server-pow:server-pow \
-v $(pwd)/:/usr/src/myapp \
client-pow-image
```
