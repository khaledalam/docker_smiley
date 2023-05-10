Level 1:

`env_file:` and `environment:` in docker-compose.yml

    environment:
      - NODE_ENV=production

Has high priority than

    env_file:
      - ./Docker/api/api.env # NODE_ENV=test

inside container:
```
process.env.NODE_ENV => `production`
```

____

Level 2:

    environment:
      - NODE_ENV=${NODE_ENV}
