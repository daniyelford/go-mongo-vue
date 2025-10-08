package handlers

import (
	"encoding/json"
	"fmt"
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/middleware"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserBusinesses(w http.ResponseWriter, r *http.Request) {
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
	bizColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("businesses")
	cursor, err := bizColl.Find(r.Context(), bson.M{"user_id": user.ID})
	if err != nil {
		http.Error(w, "cannot fetch businesses", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	var businesses []models.Business
	if err := cursor.All(r.Context(), &businesses); err != nil {
		http.Error(w, "cannot parse businesses", http.StatusInternalServerError)
		return
	}
	for i, biz := range businesses {
		for j, m := range biz.Media {
			url, err := service.MinioGetURL(m.URL, 15*time.Minute)
			if err != nil {
				fmt.Println("MinIO get url error:", err)
				continue
			}
			businesses[i].Media[j].URL = url
		}
	}
	json.NewEncoder(w).Encode(businesses)
}
func CreateBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"success":false,"error":"cannot parse form"}`, http.StatusBadRequest)
		return
	}
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user); err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	name := r.FormValue("name")
	category := r.FormValue("category")
	description := r.FormValue("description")
	files := r.MultipartForm.File["media[]"]
	var mediaList []models.Media
	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer src.Close()
		fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		contentType := fileHeader.Header.Get("Content-Type")
		_, err = service.MinioUpload(fileName, src, fileHeader.Size, contentType)
		if err != nil {
			fmt.Println("upload error:", err)
			continue
		}
		fileURL, _ := service.MinioGetURL(fileName, 24*time.Hour)
		fileType := "other"
		if contentType != "" {
			if contentType[:5] == "image" {
				fileType = "image"
			} else if contentType[:5] == "video" {
				fileType = "video"
			}
		}
		mediaList = append(mediaList, models.Media{
			URL:  fileURL,
			Type: fileType,
		})
	}
	bizColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("businesses")
	business := models.Business{
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		Name:        name,
		Category:    category,
		Description: description,
		Status:      "active",
		Media:       mediaList,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := bizColl.InsertOne(r.Context(), business)
	if err != nil {
		http.Error(w, "cannot create business", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(business)
}
func UpdateBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// --- گرفتن ID از URL ---
	idStr := r.PathValue("id")
	businessID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// --- احراز هویت ---
	mobileCtx := r.Context().Value("mobile")
	if mobileCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	mobile := mobileCtx.(string)

	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := userColl.FindOne(r.Context(), bson.M{"mobile": mobile}).Decode(&user); err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// --- گرفتن بیزینس ---
	bizColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("businesses")
	var biz models.Business
	if err := bizColl.FindOne(r.Context(), bson.M{"_id": businessID, "user_id": user.ID}).Decode(&biz); err != nil {
		http.Error(w, "business not found", http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// --- ویرایش فیلدها ---
	update := bson.M{
		"updated_at": time.Now(),
	}
	if name := r.FormValue("name"); name != "" {
		update["name"] = name
	}
	if cat := r.FormValue("category"); cat != "" {
		update["category"] = cat
	}
	if desc := r.FormValue("description"); desc != "" {
		update["description"] = desc
	}

	// --- حذف مدیاهای انتخابی ---
	deleteMedia := r.MultipartForm.Value["delete_media[]"]
	var newMedia []models.Media
	for _, m := range biz.Media {
		toDelete := false
		for _, d := range deleteMedia {
			if m.URL == d {
				toDelete = true
				break
			}
		}
		if !toDelete {
			newMedia = append(newMedia, m)
		}
	}

	// --- آپلود فایل‌های جدید ---
	files := r.MultipartForm.File["media[]"]
	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		contentType := fileHeader.Header.Get("Content-Type")

		_, err = service.MinioUpload(fileName, src, fileHeader.Size, contentType)
		if err != nil {
			fmt.Println("upload error:", err)
			continue
		}

		fileURL, _ := service.MinioGetURL(fileName, 24*time.Hour)
		fileType := "other"
		if len(contentType) >= 5 {
			if contentType[:5] == "image" {
				fileType = "image"
			} else if contentType[:5] == "video" {
				fileType = "video"
			}
		}

		newMedia = append(newMedia, models.Media{
			URL:  fileURL,
			Type: fileType,
		})
	}

	update["media"] = newMedia

	// --- ذخیره در دیتابیس ---
	_, err = bizColl.UpdateOne(r.Context(),
		bson.M{"_id": businessID, "user_id": user.ID},
		bson.M{"$set": update},
	)
	if err != nil {
		http.Error(w, "cannot update business", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bson.M{"message": "updated successfully"})
}
