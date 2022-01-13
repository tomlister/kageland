package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tomlister/kageland/util"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

var db *util.DB

type ShaderPostRequest struct {
	Name       string `json:"name" validate:"required"`
	FragShader string `json:"frag_shader" validate:"required"`
	Image1     string `json:"image_1" validate:"required"`
	Image2     string `json:"image_2" validate:"required"`
	Image3     string `json:"image_3" validate:"required"`
	Image4     string `json:"image_4" validate:"required"`
}

type ShaderDocument struct {
	ID         string `bson:"id"`
	Name       string `bson:"name"`
	AuthorID   string `bson:"author_id"`
	Anon       bool   `bson:"anon"`
	FragShader string `bson:"frag_shader"`
	Image1     string `bson:"image_1"`
	Image2     string `bson:"image_2"`
	Image3     string `bson:"image_3"`
	Image4     string `bson:"image_4"`
	Likes      int    `bson:"likes"`
	Views      int    `bson:"views"`
}

type ShaderPostResponse struct {
	ID string `json:"id"`
}

type ShaderGetResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Anon       bool   `json:"anon"`
	FragShader string `json:"frag_shader"`
	Image1     string `json:"image_1"`
	Image2     string `json:"image_2"`
	Image3     string `json:"image_3"`
	Image4     string `json:"image_4"`
	Likes      int    `json:"likes"`
	Views      int    `json:"views"`
}

func ShaderLikeHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	if qs.Get("id") == "" {
		fmt.Println("wtf")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := qs.Get("id")

	if db == nil {
		db = &util.DB{}
		db.Connect()
	}

	prod := db.Client.Database("prod")
	shaders := prod.Collection("shaders")
	err := shaders.FindOneAndUpdate(db.Ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"likes": 1}}).Err()
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
}

func ShaderUnlikeHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	if qs.Get("id") == "" {
		fmt.Println("wtf")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := qs.Get("id")

	if db == nil {
		db = &util.DB{}
		db.Connect()
	}

	prod := db.Client.Database("prod")
	shaders := prod.Collection("shaders")
	err := shaders.FindOneAndUpdate(db.Ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"likes": -1}}).Err()
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
}

func ShaderPostHandler(w http.ResponseWriter, r *http.Request) {
	reqDecoder := json.NewDecoder(r.Body)
	reqDecoder.DisallowUnknownFields()
	reqData := ShaderPostRequest{}
	err := reqDecoder.Decode(&reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validate := validator.New()
	err = validate.Struct(reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if db == nil {
		db = &util.DB{}
		db.Connect()
	}

	prod := db.Client.Database("prod")
	shaders := prod.Collection("shaders")
	id, err := gonanoid.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	document := ShaderDocument{
		ID:         id,
		Name:       reqData.Name,
		FragShader: reqData.FragShader,
		AuthorID:   "",
		Anon:       true,
		Image1:     reqData.Image1,
		Image2:     reqData.Image2,
		Image3:     reqData.Image3,
		Image4:     reqData.Image4,
		Likes:      0,
		Views:      0,
	}
	_, err = shaders.InsertOne(db.Ctx, document)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resp := ShaderPostResponse{
		ID: document.ID,
	}
	json.NewEncoder(w).Encode(resp)
}

func ShaderGetHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	if qs.Get("id") == "" {
		fmt.Println("wtf")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := qs.Get("id")
	if db == nil {
		db = &util.DB{}
		db.Connect()
	}
	prod := db.Client.Database("prod")
	shaders := prod.Collection("shaders")
	var shader ShaderDocument
	err := shaders.FindOneAndUpdate(db.Ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"views": 1}}).Decode(&shader)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resp := ShaderGetResponse{
		ID:         shader.ID,
		Name:       shader.Name,
		Anon:       shader.Anon,
		FragShader: shader.FragShader,
		Image1:     shader.Image1,
		Image2:     shader.Image2,
		Image3:     shader.Image3,
		Image4:     shader.Image4,
		Likes:      shader.Likes,
		Views:      shader.Views,
	}
	json.NewEncoder(w).Encode(resp)
}
