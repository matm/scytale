

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

Extract all encrypted files to plain text files from the archive:

    $ curl -XPOST -H "Content-Type: application/json" \
        -d '{"method": "Magma.ExtractAll", \
        "params": [{ \
            "password": "mypasswd", "archive": "/tmp/archive.zip", \
            "outputdir": "/tmp" \
        }], "id": 1}' localhost:8080/api

Extract one file at position 2 (the third one, 0-indexed):

    $ curl -XPOST -H "Content-Type: application/json" \
        -d '{"method": "Magma.ExtractAt", \
        "params": [{ \
            "password": "mypasswd", "archive": "/tmp/archive.zip", \
            "outputdir": "/tmp", "pos": 2 \
        }], "id": 1}' localhost:8080/api
