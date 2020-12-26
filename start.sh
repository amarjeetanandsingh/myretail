docker run -p 8529:8529 -e ARANGO_ROOT_PASSWORD=root arangodb/arangodb:3.7.5
docker build -t myretail .
docker run --env ARANGO_HOST='127.0.0.1'  --env ARANGO_PORT='8529'  --env ARANGO_USERNAME='root'  --env ARANGO_PASSWORD='root' -p 5229:5229 myretail
