Emotion recognition app
=======================

Design
------

Architecture
------------

Persistence layer
-----------------

Start PostgreSQL container
```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=emokognition -e POSTGRES_USER=postgres -d postgres
```

Deployment
----------
