#!/bin/bash
docker build -t client-pow-image .

docker run --rm -it \
--name client-pow \
--link server-pow:server-pow \
-v $(pwd)/:/usr/src/myapp \
client-pow-image
