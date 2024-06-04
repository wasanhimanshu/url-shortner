url shortner app created in golang (backend only)
build the docker file using command
docker build . {IMAGE_NAME}:{TAG} .

application is running on port 8080
you will need a mysql running for the app to work ,you can create it in docker
command- `docker run --rm -it --network=host -e MYSQL_ROOT_PASSWORD="mypassword" mysql:latest`
this will start a mysql container on port 3306

connect to mysql running in docker (do ensure mysql binary is installed on your local pc)
command - `mysql -h 127.0.0.1 -u root -p`

you will need to create a database also with name urlshortner.
command - `CREATE DATABASE urlshortner`

after this you can register,login and create short urls.