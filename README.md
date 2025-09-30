golang mongo redis minio vue
this is a project for reserve services
you can use this with docker in local
test with 
docker-compose -f docker-compose.yml -f docker-compose.override.yml up --build
or run with
docker-compose up --build
and remove
  server: {
    host: true,
    port: 5173,
    proxy: {
      '/api': 'http://go:8080'
    }
  },
  from vite.config.js