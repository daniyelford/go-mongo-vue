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
docker-compose up --build