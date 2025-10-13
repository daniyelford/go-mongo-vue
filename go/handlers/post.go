package handlers

import (
	"encoding/json"
	"fmt"
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/middleware"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"success":false,"error":"user not found"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	userID := user.ID
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

	cursor, err := postColl.Find(r.Context(), bson.M{}, opts)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	var posts []models.Post
	if err := cursor.All(r.Context(), &posts); err != nil {
		http.Error(w, `{"success":false,"error":"error reading data"}`, http.StatusInternalServerError)
		return
	}
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		publicEndpoint = os.Getenv("MINIO_ENDPOINT")
	}
	bucket := os.Getenv("MINIO_BUCKET")
	for i := range posts {
		for j := range posts[i].Media {
			if posts[i].Media[j].Filename != "" {
				posts[i].Media[j].URL = fmt.Sprintf(
					"http://%s/%s/%s",
					publicEndpoint,
					bucket,
					url.PathEscape(posts[i].Media[j].Filename),
				)
			}
		}
	}
	type PostWithSelf struct {
		models.Post
		Self bool `json:"self"`
	}
	postsWithSelf := make([]PostWithSelf, len(posts))
	for i, p := range posts {
		postsWithSelf[i] = PostWithSelf{
			Post: p,
			Self: p.UserID == userID,
		}
	}
	total, _ := postColl.CountDocuments(r.Context(), bson.M{})
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    postsWithSelf,
		"page":    body.Page,
		"limit":   body.Limit,
		"total":   total,
	})
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseMultipartForm(50 << 20); err != nil {
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
	files := form.File["media[]"]
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
	_, err = postColl.InsertOne(config.Ctx, post)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"success":false,"error":"user not found"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	userID := user.ID
	type RequestBody struct {
		ID string `json:"id"`
	}
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"success":false,"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	postID, err := primitive.ObjectIDFromHex(body.ID)
	if err != nil {
		http.Error(w, `{"success":false,"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}
	postColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("posts")
	res, err := postColl.DeleteOne(r.Context(), bson.M{
		"_id":     postID,
		"user_id": userID,
	})
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, `{"success":false,"error":"post not found or not authorized"}`, http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
func EditPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user); err != nil {
		http.Error(w, `{"success":false,"error":"user not found"}`, http.StatusUnauthorized)
		return
	}
	userID := user.ID
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, `{"success":false,"error":"cannot parse form"}`, http.StatusBadRequest)
		return
	}
	postIDStr := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	oldMediaJSON := r.FormValue("oldMedia")
	if postIDStr == "" || title == "" || content == "" {
		http.Error(w, `{"success":false,"error":"missing fields"}`, http.StatusBadRequest)
		return
	}
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, `{"success":false,"error":"invalid post id"}`, http.StatusBadRequest)
		return
	}
	var oldMedia []models.PostMedia
	if oldMediaJSON != "" {
		if err := json.Unmarshal([]byte(oldMediaJSON), &oldMedia); err != nil {
			oldMedia = []models.PostMedia{}
		}
	}
	form := r.MultipartForm
	files := form.File["media[]"]
	var newMedia []models.PostMedia
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
		newMedia = append(newMedia, models.PostMedia{
			URL:      fmt.Sprintf("%s/%s", os.Getenv("MINIO_PUBLIC_URL"), fileName),
			Type:     contentType,
			Filename: fileName,
		})
	}
	postColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("posts")
	var existingPost models.Post
	if err := postColl.FindOne(r.Context(), bson.M{"_id": postID, "user_id": userID}).Decode(&existingPost); err != nil {
		http.Error(w, `{"success":false,"error":"post not found or permission denied"}`, http.StatusNotFound)
		return
	}
	for _, old := range existingPost.Media {
		keep := false
		for _, m := range oldMedia {
			if m.Filename == old.Filename {
				keep = true
				break
			}
		}
		if !keep {
			if err := service.MinioRemove(old.Filename); err != nil {
				fmt.Println("Failed to delete from MinIO:", old.Filename)
			}
		}
	}
	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"content":    content,
			"media":      append(oldMedia, newMedia...),
			"updated_at": time.Now(),
		},
	}
	if _, err := postColl.UpdateOne(r.Context(), bson.M{"_id": postID, "user_id": userID}, update); err != nil {
		http.Error(w, `{"success":false,"error":"failed to update post"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
