Agent, you are a devops engineer in charge of implementing the ability for the data to persist across restarts or if the docker container is recreated.

### Why is the data being lost?
By default, Docker stores a database's files on the container's ephemeral filesystem. If the container is destroyed (`docker-compose down` or simple removal) and recreated, its internal filesystem is also destroyed, taking your database records with it. 

### How to fix it?
If you look inside `api/docker-compose.yml`, you will see a `volumes:` block specifically for the database:

```yaml
    volumes:
      - ./migrations:/docker-entrypoint-initdb.d
      # - postgres_data:/var/lib/postgresql/data
```

The line `- postgres_data:/var/lib/postgresql/data` is currently commented out. 

To persist the database:
1. **Uncomment that line** in the `docker-compose.yml` file. 
2. This maps the directory where Postgres natively saves files (`/var/lib/postgresql/data`) to a dedicated internal host volume managed by Docker named `postgres_data`. 
3. Now, when the container restarts or is rebuilt, Docker re-mounts the `postgres_data` volume, and your database state is perfectly preserved!
