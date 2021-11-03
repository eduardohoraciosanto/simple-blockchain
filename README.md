# Blockchain Experiment

This is a Simple Blockchain Experiment coded from scratch

---

## Build & Run

When first cloning the repository, run the following to build the service
```shell
make init
make build
```

`init` will populate a `.env` file with default values.  
`build` will proceed to create the docker image to run the service.  


To run the service use the following:
```shell
make up
```

By default this will run on `https://localhost:18080/`
---

## Endpoints Available

The service has three basic endpoints.  

### GET Health
A Health endpoint under the route `/health` that will return 200 OK and a simple response to show the service is up.  
```
{
    "meta": {
        "version": "local"
    },
    "data": {
        "services": [
            {
                "name": "service",
                "alive": true
            }
        ]
    }
}
```

### GET Fetch Blockchain
The fetch Blockchain endpoint will return the entire chain with all the blocks inside, with the following format:
```
{
    "meta": {
        "version": "local"
    },
    "data": [
        {
            "Index": 0,
            "Timestamp": "2021-11-03 06:11:22.331096801 +0000 UTC m=+0.055173460",
            "Data": {
                "BMP": 0
            },
            "Nuance": 0,
            "Hash": "",
            "PrevHash": ""
        }
    ]
}
```

### POST Insert Block Data 
This endpoint begins an insertion to the Blockchain. It makes sure the hash complies with the Leading Zeroes Rule and the amount of leading zeroes can be configured in the `.env` file.

The body for this request is as follows:
```
{
    "block_data":{
        "bmp": 100
    }
}
```

Where `bmp` is a simple integer.  

The response will be a 202 Accepted once the hash is found. This operation is essentially the mining of the block in this simple logic.
