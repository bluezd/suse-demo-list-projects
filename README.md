README
--

## building 

### docker build

```
TAG=`date '+%Y-%m-%d-%H-%M-%S'`
docker build -t bluezd/suse-demo-list-projects:$TAG .
```

### docker run

```
docker run -d bluezd/suse-demo-list-projects:$TAG

```

## Demo 

### delete entry

`curl -k -s -XDELETE http://localhost:8001/projects/3`

### add entry

```
curl -X POST http://localhost:8001/projects/3 \
-H "Content-Type: application/json" \
--data-binary @- << EOF
{
    "id": "3",
    "name": "RKE",
    "repository": "https://github.com/rancher/rke",
    "twitter": "",
    "website": "https://rancher.com/products/rke",
    "description": "Kubernetes simplified. Runs in Docker containers."
}
EOF
```
