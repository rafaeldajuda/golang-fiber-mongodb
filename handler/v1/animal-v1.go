package v1

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MsgError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type HandlerV1 struct {
	Ctx         context.Context
	Collection  *mongo.Collection
	MongoClient *mongo.Client
}

func (h *HandlerV1) GetAnimals(c *fiber.Ctx) error {
	cursor, err := h.Collection.Find(h.Ctx, bson.M{})
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
	err = h.Collection.FindOne(h.Ctx, filter).Decode(result)
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
	filter := bson.M{"dono": animal.Dono, "nome": animal.Nome}
	err = h.Collection.FindOne(h.Ctx, filter).Err()
	if err == nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: "this animal already exist"}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	result, err := h.Collection.InsertOne(h.Ctx, animal)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := struct {
		ID interface{} `json:"_id"`
	}{ID: result.InsertedID}

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
	err = h.Collection.FindOne(h.Ctx, filter).Err()
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: "animal not found"}
		return c.Status(http.StatusBadRequest).JSON(e)
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
	_, err = h.Collection.UpdateOne(h.Ctx, filter, bodyUpdate)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := struct {
		ID interface{} `json:"_id"`
	}{ID: id}

	return c.Status(http.StatusOK).JSON(msgOk)
}

func (h *HandlerV1) DeleteAnimal(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	filter := bson.M{"_id": id}
	result, err := h.Collection.DeleteOne(h.Ctx, filter)
	if err != nil {
		e := MsgError{Code: http.StatusBadRequest, Msg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(e)
	}

	msgOk := struct {
		DeletedCount interface{} `json:"deleted_count"`
	}{DeletedCount: result.DeletedCount}

	return c.Status(http.StatusOK).JSON(msgOk)
}
