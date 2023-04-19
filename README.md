# golang-fiber-mongodb

API feita em golang utilizando Fiber e MongoDB.

O objetivo deste projeto é mostrar como fazer uma API Rest em golang utilizando o Fiber para montar o servidor e utilizando o MongoDB para gravar e capturar dados.

## Dependências

|Lib        |Link                               |
|----       |----                               |
|Fiber      |github.com/gofiber/fiber/v2        |
|MongoDB    |github.com/joho/godotenv           |
|godotenv   |go.mongodb.org/mongo-driver/mongo  |


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

## Dados Para Testes

Aqui terá alguns comandos para alimentar o banco para os testes.

* Criar um banco.
```mongodb
use dev
```

* Criar uma coleção. 
```mongodb
db.createCollection("animals")
```

* Inserir dados na coleção.
```mongodb
db.pessoa.insertOne({"name":"Toby","type":"cachorro","age":4,"owner":"Rafael","castrated":false,"surgery":new Date("2023-10-10")})
```

# Iniciando o Servidor e o MongoDB

Antes de iniciar a explicação do código irei monstrar como estão configuradas as variáveis de ambiente. As variáveis foram criadas no arquivo **.env**.

(.env)
```env
PROJECT_NAME="golang-fiber-mongodb"
PROJECT_HOST="0.0.0.0"
PROJECT_PORT="8000"

MONGO_URI="mongodb://admin:admin@localhost:27017/"
MONGO_USER=admin
MONGO_PASSWORD=admin
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DATABASE=dev
MONGO_COLLECTION=animals
```

No arquivo **main.go** é onde iniciamos o servidor e conectamos ao banco de dados. No início do códogo declaramos uma variável do tipo **API**. Ela conterá todos os valores das variáveis de ambiente, funções das rotas e a conexão com o banco de dados.

A função **init()** é responsável por carregar as variáveis de ambiente, conectar ao MongoDB e criar uma instância da entidade **HandlerV1**, que é responsável pelas funções das rotas.

(main.go)
```golang
var api entity.API

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.Name = os.Getenv("PROJECT_NAME")
	api.Host = os.Getenv("PROJECT_HOST")
	api.Port = os.Getenv("PROJECT_PORT")

	api.Ctx = context.TODO()
	api.MongoConfig = database.MongoConfig{
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
	client := api.MongoConfig.CreateConnection(api.Ctx)

	api.HandlerV1.Ctx = api.Ctx
	api.HandlerV1.Database = os.Getenv("MONGO_DATABASE")
	api.HandlerV1.Collection = os.Getenv("MONGO_COLLECTION")
	api.HandlerV1.MongoClient = client
	api.HandlerV1.MongoCollection = client.Database(api.HandlerV1.Database).Collection(api.HandlerV1.Collection)
}
```

Na função **main()** é iniciado o Fiber, que é resposável pelo servidor e rotas da API. Para configurar as rotas é necessário passar a instância do fiber e a da API para a função **Routes**. A instância do fiber é responsável pela criação das rotas e grupos, enquanto a API possui o **HandlerV1** que contém as funções de cada rota.

(main.go)
```golang
func main() {
	defer func() {
		err := api.MongoClient.Disconnect(api.Ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New()

	// Routes
	routes.Routes(app, &api)

	err := app.Listen(api.Host + ":" + api.Port)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Entity

Pacote que contém a struct principal da aplicação, **API**. 

(entity/entity.go)
```golang
type API struct {
	Name        string
	Host        string
	Port        string
	Ctx         context.Context
	MongoConfig database.MongoConfig
	v1.HandlerV1
}
```

## Database

Neste pacote contém a função responsável por conectar ao banco de dados. Nela possui a struct **MongoConfig**, resposável por conter os dados de acesso ao banco e por chamar a função **CreateConnection()**.

(database/mongodb.go)
```golang
type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func (mc *MongoConfig) CreateConnection(ctx context.Context) *mongo.Client {
	// mongodb uri format
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", mc.User, mc.Password, mc.Host, mc.Port)
	if uri == "" {
		log.Fatal("You must set your 'uri' variable.")
	}

	// connection
	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("MongoDB connection error: " + err.Error())
	}

	// check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
```

## Routes

Pacote responsável por criar os grupos e a rotas da API. Possui duas funções, **Routes()** e **routeV1()**.

A função Routes chama todas as funções que criam grupos e rotas, para cada uma dessas funções ela passa a instância do Fiber e da API.

A função routeV1 cria um grupo de rotas, iniciado pelo path `/api/v1/animal`. Nessa função utilizamos o Fiber para criar o grupo e as rotas, e a API para invocar as funções que serão chamadas quando as rotas forem acionadas.

(routes/animal-v1.go)
```golang
func Routes(app *fiber.App, api *entity.API) {
	routeV1(app, api)
}

func routeV1(app *fiber.App, api *entity.API) {
	v1 := app.Group("/api/v1/animal")

	v1.Get("/", api.HandlerV1.GetAnimals)
	v1.Get("/:id", api.HandlerV1.GetAnimal)
	v1.Post("/", api.HandlerV1.PostAnimal)
	v1.Put("/:id", api.HandlerV1.PutAnimal)
	v1.Delete("/:id", api.HandlerV1.DeleteAnimal)
}
```

## Handler/V1

Pacote que contém a lógica das rotas. As principais funções serão chamadas pelo pacote **routes**. Para cada versão da aplicação será criado um pacote novo. No caso desse projeto só possui uma versão, então só terá o pocote **v1** dentro do handler.

Dentro de v1 possui estrutura **HandlerV1**, responsável por chamar as rotas. Essa estrutura recebe o nome do banco de dados e da coleção, a instância do cliente do MongoDB e a instância da coleção do MongoDB. Também recebe um contexto que deve ser usado para todas as chamadas do banco.

As funções das rotas possuem um comportamento simples, todas só podem ser chamadas pelo HandlerV1, que é instânciado no arquivo **main.go** e passado para as rotas. 

A função `GetAnimals` retorna todos os documentos da coleção. A consulta é feita sem nenhum filtro.

A função `GetAnimal` retorna somente um documento da coleção. É preciso passar o ID do documento para que seja feita a filtragem no momemento da consulta. Caso passe um ID inválido será retornado um erro de `Bad Request`. Caso passe um ID válido, mas inexistente na coleção, será retornado um erro de `Not Found`.

A função `PostAnimal` adiciona um novo documento na coleção e retorna o ID do documento criado. Antes de adicionar é feita uma consulta pelo nome do dono do animal e pelo nome animal para checar se o animal já existe na coleção, caso exista será retornado um erro de `Bad Request`. 

A função `PutAnimal` atualiza um documento na coleção a partir de um ID. Caso o ID não sejá válido ou não exista na coleção, será retornado um erro. Caso o ID seja válido a atualização ocorrerá e será retornado o ID do documento como resposta.

A função `DeleteAnimal` remove um documento na coleção a partir do ID passado. Como resposta será retornado o número de documentos removidos na coleção.

(handler/v1/animal-v1.go)
```golang
type HandlerV1 struct {
	Ctx             context.Context
	Database        string
	Collection      string
	MongoCollection *mongo.Collection
	MongoClient     *mongo.Client
}

func (h *HandlerV1) GetAnimals(c *fiber.Ctx) error {
	cursor, err := h.MongoCollection.Find(h.Ctx, bson.M{})
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	var result []bson.M
	for cursor.Next(h.Ctx) {
		raw := cursor.Current
		item := bson.M{}
		err := bson.Unmarshal(raw, &item)
		if err != nil {
			e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
			return c.Status(http.StatusBadRequest).JSON(e)
		}
		result = append(result, item)
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (h *HandlerV1) GetAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	result := bson.M{}
	filter := bson.M{"_id": id}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Decode(result)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (h *HandlerV1) PostAnimal(c *fiber.Ctx) error {
	animal := Animal{}
	err := c.BodyParser(&animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	// animal exist?
	filter := bson.M{"owner": animal.Owner, "name": animal.Name}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Err()
	if err == nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: "this animal already exist"}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	result, err := h.MongoCollection.InsertOne(h.Ctx, animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := MsgOK{ID: result.InsertedID}
	return c.Status(http.StatusOK).JSON(msgOk)
}

func (h *HandlerV1) PutAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	// animal exist?
	filter := bson.M{"_id": id}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Err()
	if err != nil {
		e := MsgError{Code: http.StatusNotFound, Msg: "animal not found"}
		return c.Status(http.StatusNotFound).JSON(e)
	}

	animal := bson.M{}
	err = c.BodyParser(&animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}
	bodyUpdate := bson.M{
		"$set": animal,
	}

	filter = bson.M{"_id": id}
	_, err = h.MongoCollection.UpdateOne(h.Ctx, filter, bodyUpdate)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := MsgOK{ID: id}
	return c.Status(http.StatusOK).JSON(msgOk)
}

func (h *HandlerV1) DeleteAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	filter := bson.M{"_id": id}
	result, err := h.MongoCollection.DeleteOne(h.Ctx, filter)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := struct {
		DeletedCount interface{} `json:"deleted_count"`
	}{DeletedCount: result.DeletedCount}
	return c.Status(http.StatusOK).JSON(msgOk)
}
```

No pocote v1 também existe as estruturas que representam as respostas de sucesso e de erro. Também tem a estrutura `Animal` que representa o documento da coleção do MongoDB e o formato de JSON utilizado nas requisições.

(handler/v1/entity.go)
```golang
type MsgError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type MsgOK struct {
	ID interface{} `json:"_id" bson:"_id"`
}

type Animal struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Owner     string    `json:"owner" bson:"owner"`
	Type      string    `json:"type" bson:"type"`
	Age       int       `json:"age" bson:"age"`
	Castrated bool      `json:"castrated" bson:"castrated"`
	Surgery   time.Time `json:"surgery" bson:"surgery"`
}
```

## Exempos dos Arquivos Completos

(.env)
```.env
PROJECT_NAME="golang-fiber-mongodb"
PROJECT_HOST="0.0.0.0"
PROJECT_PORT="8000"

MONGO_URI="mongodb://admin:admin@localhost:27017/"
MONGO_USER=admin
MONGO_PASSWORD=admin
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DATABASE=dev
MONGO_COLLECTION=animals
```

(main.go)
```golang
package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rafaeldajuda/database"
	"github.com/rafaeldajuda/entity"
	"github.com/rafaeldajuda/routes"
)

var api entity.API

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.Name = os.Getenv("PROJECT_NAME")
	api.Host = os.Getenv("PROJECT_HOST")
	api.Port = os.Getenv("PROJECT_PORT")

	api.Ctx = context.TODO()
	api.MongoConfig = database.MongoConfig{
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
	client := api.MongoConfig.CreateConnection(api.Ctx)

	api.HandlerV1.Ctx = api.Ctx
	api.HandlerV1.Database = os.Getenv("MONGO_DATABASE")
	api.HandlerV1.Collection = os.Getenv("MONGO_COLLECTION")
	api.HandlerV1.MongoClient = client
	api.HandlerV1.MongoCollection = client.Database(api.HandlerV1.Database).Collection(api.HandlerV1.Collection)
}

func main() {
	defer func() {
		err := api.MongoClient.Disconnect(api.Ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New()

	// Routes
	routes.Routes(app, &api)

	err := app.Listen(api.Host + ":" + api.Port)
	if err != nil {
		log.Fatal(err)
	}
}
```

(entity/entity.go)
```golang
package entity

import (
	"context"

	"github.com/rafaeldajuda/database"
	v1 "github.com/rafaeldajuda/handler/v1"
)

type API struct {
	Name        string
	Host        string
	Port        string
	Ctx         context.Context
	MongoConfig database.MongoConfig
	v1.HandlerV1
}
```

(database/mongodb.go)
```golang
package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func (mc *MongoConfig) CreateConnection(ctx context.Context) *mongo.Client {
	// mongodb uri format
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", mc.User, mc.Password, mc.Host, mc.Port)
	if uri == "" {
		log.Fatal("You must set your 'uri' variable.")
	}

	// connection
	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("MongoDB connection error: " + err.Error())
	}

	// check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
```

(routes/animal-v1.go)
```golang
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaeldajuda/entity"
)

func Routes(app *fiber.App, api *entity.API) {
	routeV1(app, api)
}

func routeV1(app *fiber.App, api *entity.API) {
	v1 := app.Group("/api/v1/animal")

	v1.Get("/", api.HandlerV1.GetAnimals)
	v1.Get("/:id", api.HandlerV1.GetAnimal)
	v1.Post("/", api.HandlerV1.PostAnimal)
	v1.Put("/:id", api.HandlerV1.PutAnimal)
	v1.Delete("/:id", api.HandlerV1.DeleteAnimal)
}
```

(handler/v1/animal-v1.go)
```golang
package v1

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerV1 struct {
	Ctx             context.Context
	Database        string
	Collection      string
	MongoCollection *mongo.Collection
	MongoClient     *mongo.Client
}

func (h *HandlerV1) GetAnimals(c *fiber.Ctx) error {
	cursor, err := h.MongoCollection.Find(h.Ctx, bson.M{})
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	var result []bson.M
	for cursor.Next(h.Ctx) {
		raw := cursor.Current
		item := bson.M{}
		err := bson.Unmarshal(raw, &item)
		if err != nil {
			e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
			return c.Status(http.StatusBadRequest).JSON(e)
		}
		result = append(result, item)
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (h *HandlerV1) GetAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	result := bson.M{}
	filter := bson.M{"_id": id}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Decode(result)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (h *HandlerV1) PostAnimal(c *fiber.Ctx) error {
	animal := Animal{}
	err := c.BodyParser(&animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	// animal exist?
	filter := bson.M{"owner": animal.Owner, "name": animal.Name}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Err()
	if err == nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: "this animal already exist"}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	result, err := h.MongoCollection.InsertOne(h.Ctx, animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := MsgOK{ID: result.InsertedID}
	return c.Status(http.StatusOK).JSON(msgOk)
}

func (h *HandlerV1) PutAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	// animal exist?
	filter := bson.M{"_id": id}
	err = h.MongoCollection.FindOne(h.Ctx, filter).Err()
	if err != nil {
		e := MsgError{Code: http.StatusNotFound, Msg: "animal not found"}
		return c.Status(http.StatusNotFound).JSON(e)
	}

	animal := bson.M{}
	err = c.BodyParser(&animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}
	bodyUpdate := bson.M{
		"$set": animal,
	}

	filter = bson.M{"_id": id}
	_, err = h.MongoCollection.UpdateOne(h.Ctx, filter, bodyUpdate)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := MsgOK{ID: id}
	return c.Status(http.StatusOK).JSON(msgOk)
}

func (h *HandlerV1) DeleteAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	filter := bson.M{"_id": id}
	result, err := h.MongoCollection.DeleteOne(h.Ctx, filter)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := struct {
		DeletedCount interface{} `json:"deleted_count"`
	}{DeletedCount: result.DeletedCount}
	return c.Status(http.StatusOK).JSON(msgOk)
}
```

(handler/v1/entity.go)
```golang
package v1

import "time"

type MsgError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type MsgOK struct {
	ID interface{} `json:"_id" bson:"_id"`
}

type Animal struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Owner     string    `json:"owner" bson:"owner"`
	Type      string    `json:"type" bson:"type"`
	Age       int       `json:"age" bson:"age"`
	Castrated bool      `json:"castrated" bson:"castrated"`
	Surgery   time.Time `json:"surgery" bson:"surgery"`
}
```