golang mongo redis minio vue
this is a project for reserve services
 

for use minio az public you must
download mc.exe (MinIO Client Releases)
and set path in Environment Variables
and use cmd 
.\mc.exe --help
.\mc.exe alias set local http://127.0.0.1:9000 admin password123
.\mc.exe anonymous set public local/media
so you can use this with docker in local
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