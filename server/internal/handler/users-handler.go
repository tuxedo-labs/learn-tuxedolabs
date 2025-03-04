package handler

import (
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/pkg/database"
	"learn-tuxedolabs/pkg/utils"
	"net/http"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	if token == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	var user entity.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	filter := map[string]interface{}{
		"name":       user.Name,
		"first_name": user.FirstName,
		"last_name":  *user.LastName,
		"avatar":     user.Avatar,
		"email":      user.Email,
		"role":       user.Role,
		"verify":     user.Verify,
		"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at": user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.RespondJSON(w, http.StatusOK, filter)
}
