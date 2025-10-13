package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/middleware"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type RequestBody struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"success":false,"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.Page < 1 {
		body.Page = 1
	}
	if body.Limit <= 0 || body.Limit > 50 {
		body.Limit = 10
	}
	skip := (body.Page - 1) * body.Limit
	postColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("posts")
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(body.Limit))
	cursor, err := postColl.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	var posts []models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		http.Error(w, `{"success":false,"error":"error reading data"}`, http.StatusInternalServerError)
		return
	}
	total, _ := postColl.CountDocuments(context.Background(), bson.M{})
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    posts,
		"page":    body.Page,
		"limit":   body.Limit,
		"total":   total,
	})
}
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	userID := user.ID
	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max
		http.Error(w, `{"success":false,"error":"cannot parse form"}`, http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" {
		http.Error(w, `{"success":false,"error":"title is required"}`, http.StatusBadRequest)
		return
	}
	form := r.MultipartForm
	files := form.File["media"]
	var mediaList []models.PostMedia
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		defer file.Close()
		fileName := fmt.Sprintf("%s_%d_%s", userID.Hex(), time.Now().UnixNano(), header.Filename)
		contentType := header.Header.Get("Content-Type")
		if _, err := service.MinioUpload(fileName, file, header.Size, contentType); err != nil {
			fmt.Println("Upload error:", err)
			continue
		}
		mediaList = append(mediaList, models.PostMedia{
			URL:      fmt.Sprintf("%s/%s", os.Getenv("MINIO_PUBLIC_URL"), fileName),
			Type:     contentType,
			Filename: fileName,
		})
	}
	post := models.Post{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		Media:     mediaList,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	postColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("posts")
	newPost, err := postColl.InsertOne(config.Ctx, post)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newPost)
}
