# Multiplexer - collects information by URLs


### To run server use:
``
make run
``

### Request example:
``
curl -d '["http://jsonplaceholder.typicode.com/posts/1", "http://jsonplaceholder.typicode.com/posts/2", "http://jsonplaceholder.typicode.com/posts/3"]' localhost:8081/multiplexer | jq
``