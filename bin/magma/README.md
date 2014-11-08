

Exit the service with:

    $ curl -XPOST -H "Content-Type: application/json" \
        -d '{"method": "Magma.Exit", "params": [], "id": 1}' \
        localhost:8080/api

Create a ZIP archive with encrypted files in it with:

    $ curl -XPOST -H "Content-Type: application/json" \
        -d '{"method": "Magma.CreateArchive", \
        "params": [{ \
            "password": "mypasswd", "outputname": "/tmp/archive.zip", \
            "files": ["/tmp/file1.jpg", "/tmp/file2.pdf"] \
        }], "id": 1}' localhost:8080/api
