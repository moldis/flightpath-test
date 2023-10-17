# Microservice to determine the flight path of a person

## Statement of the problem
Story: There are over 100,000 flights a day, with millions of people and cargo being transferred around the world.
With so many people and different carrier/agency groups, it can be hard to track where a person might be.
In order to determine the flight path of a person, we must sort through all of their flight records.

Goal: To create a simple microservice API that can help us understand and track how a particular person's flight path
may be queried. The API should accept a request that includes a list of flights, which are defined by a source and
destination airport code. These flights may not be listed in order and will need to be sorted to find the total
flight paths starting and ending airports.

## Idea of a possible solution
This is basically a graph problem and we can use the Depth First Search (DFS) algorithm to resolve it. More information about DFS is available here: https://www.geeksforgeeks.org/depth-first-search-or-dfs-for-a-graph/.

In the code, I am using directional search, i.e. all nodes should be reachable within the edges (flight -> flight -> flight without roll over to a different airport). Big credits to the https://pkg.go.dev/github.com/dominikbraun/graph library, since it already supports DFS; I used it as a base, cutting some unnecessary functions, and adapting it to the current task.

Implementation details:
* Disconnected routes are not supported (example `[["IND", "FDF"], ["DAD", "EED"]]`), i.e. it will return back one of the route in the edge.
* Integration-tests not included, since code don't have any external resources and logic embedded to single file.
* Have protection against cycling, i.e `[["IND", "IND"], ["DAD", "EED"]]` will response with error.
* Graph is based on Vertex and Edges, where Edges is the route and Vertex is the node.

## Possible improvements
* Include disconnected routes.
* Weighted routing, where some routes might be faster then others.
* Request and response might be using more complex types.
* Integrated source point, i.e. user can select airport for departure

## Run
```shell
make up 
``` 

Execute test route with curl
```shell
curl --location --request GET 'localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data '[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["GSO", "IND"], ["ATL", "GSO"]]'
```

Wrong routes examples
```shell
curl --location --request GET 'localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data '[["IND", "EWR"],["EWR", "EWR"]]'
```

```shell
{"error":"edge would create a cycle"}
```

## Postman

Collections included.