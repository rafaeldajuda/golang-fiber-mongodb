# golang-fiber-mongodb
API made in golang with fiber and mongodb


## MongoDB Docker

O comando abaixo irá baixar e rodar a imagem padrão do MongoDB. Neste comando foi definido o usuário, senha a e porta de acesso ao MongoDB. O nome do container será definido como 'my-mongo'.

```command
docker run -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 --name my-mongo mongo
``` 

## Acesso ao Container

Para acessar o container basta rodar o comando abaixo passando o usuário e senha que foi definido anteriormente.

```command
docker exec -it my-mongo mongosh -u admin -p admin
``` 











{"name":"Toby","type":"cachorro","age":4,"owner":"Rafael","castrated":false,"surgery":new Date("2023-10-10")}