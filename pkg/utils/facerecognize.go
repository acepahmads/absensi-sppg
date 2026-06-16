package utils

import (
	"bytes"
	"absensi-sppg/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

type AttendancePayload struct {
	File string `json:"file"`
}

/*
========================
RESPONSE STRUCT
========================
*/
type FaceResponse struct {
	Result []struct {
		Subjects []struct {
			Subject    string  `json:"subject"`
			Similarity float64 `json:"similarity"`
		} `json:"subjects"`
	} `json:"result,omitempty"`

	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type FaceResponse1 struct {
	Result []Face `json:"result"`
}
type Face struct {
	Box      map[string]interface{} `json:"box"`
	Subjects []Subject1             `json:"subjects"`
}
type Subject1 struct {
	Subject    string  `json:"subject"`
	Similarity float64 `json:"similarity"`
}

// {
//     "image_id": "e68bfbf1-0dd1-41a4-bd46-b5058c4df96f",
//     "subject": "acep"
// }

// {
//     "message": "No face is found in the given image",
//     "code": 28
// }

func FaceRecognizeAttendance(photoBase64 string) (string, float64, error) {
	// =========================
	// CLEAN BASE64 PREFIX
	// =========================
	if strings.Contains(photoBase64, ",") {
		photoBase64 = strings.SplitN(photoBase64, ",", 2)[1]
	}

	payload := AttendancePayload{
		File: photoBase64,
	}

	bodyReq, err := json.Marshal(payload)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest(
		"POST",
		"http://192.168.1.147:8000/api/v1/recognition/recognize",
		bytes.NewBuffer(bodyReq),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "23e225bf-8f28-4493-a870-39019954fdae")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var faceResp FaceResponse
	if err := json.Unmarshal(bodyResp, &faceResp); err != nil {
		return "", 0, err
	}

	// =========================
	// HANDLE ERROR RESPONSE
	// =========================
	if faceResp.Code != 0 {
		return "", 0, errors.New(faceResp.Message)
	}

	// =========================
	// HANDLE NO FACE
	// =========================
	if len(faceResp.Result) == 0 {
		return "", 0, errors.New("no face detected")
	}

	// =========================
	// PICK BEST SIMILARITY
	// =========================
	var bestName string
	var bestScore float64

	for _, r := range faceResp.Result {
		for _, s := range r.Subjects {
			if s.Similarity > bestScore {
				bestScore = s.Similarity
				bestName = s.Subject
			}
		}
	}

	if bestName == "" {
		return "", 0, errors.New("no valid face match")
	}

	// =========================
	// THRESHOLD (OPSIONAL)
	// =========================
	if bestScore < 0.8 {
		return "", bestScore, errors.New("face similarity too low")
	}

	return bestName, bestScore, nil
}

func RegisterFace(subject string, photoBase64 string) (string, float64, error) {
	if strings.Contains(photoBase64, ",") {
		photoBase64 = strings.SplitN(photoBase64, ",", 2)[1]
	}
	fmt.Println("Register Face", photoBase64)

	payload := AttendancePayload{
		File: photoBase64,
	}

	bodyReq, err := json.Marshal(payload)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest(
		"POST",
		"http://192.168.1.147:8000/api/v1/recognition/faces",
		bytes.NewBuffer(bodyReq),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "23e225bf-8f28-4493-a870-39019954fdae")

	//add param subject
	query := req.URL.Query()
	query.Add("subject", subject)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var faceResp FaceResponse
	if err := json.Unmarshal(bodyResp, &faceResp); err != nil {
		return "", 0, err
	}

	// =========================
	// HANDLE ERROR RESPONSE
	// =========================
	if faceResp.Code != 0 {
		return "", 0, errors.New(faceResp.Message)
	}
	return "OK", 1.0, nil
}

func FaceVerify(c *gin.Context, photoBase64 string, db *sqlx.DB) (string, float64, error) {
	// =========================
	// CLEAN BASE64 PREFIX
	// =========================
	if strings.Contains(photoBase64, ",") {
		photoBase64 = strings.SplitN(photoBase64, ",", 2)[1]
	}

	payload := AttendancePayload{
		File: photoBase64,
	}

	bodyReq, err := json.Marshal(payload)
	if err != nil {
		return "", 0, err
	}
	// fmt.Println("bodyreq", photoBase64)

	req, err := http.NewRequest(
		"POST",
		"http://192.168.1.147:8000/api/v1/recognition/recognize",
		bytes.NewBuffer(bodyReq),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "23e225bf-8f28-4493-a870-39019954fdae")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var result FaceResponse1
	if err := json.Unmarshal(bodyResp, &result); err != nil {
		return "", 0, err
	}

	// fmt.Println("result", result)

	// faces = result.get("result", [])
	if len(result.Result) == 0 {
		c.JSON(200, gin.H{
			"authorized": false,
			"message":    "no_face",
		})
		return "", 0, errors.New("no face detected")
	}

	face := result.Result[0]
	box := face.Box
	subjects := face.Subjects

	// fmt.Println("subject", subjects)

	// if not subjects
	if len(subjects) == 0 {
		c.JSON(200, gin.H{
			"authorized": false,
			"message":    "not_recognized",
			"box":        box,
		})
		return "", 0, errors.New("face not recognized")
	}

	top := subjects[0]
	similarity := top.Similarity
	subject := top.Subject

	const THRESHOLD = 0.99 // sesuaikan

	if similarity < THRESHOLD {
		c.JSON(200, gin.H{
			"authorized": false,
			"message":    "low_confidence",
			"similarity": similarity,
			"box":        box,
		})
		return "", 0, errors.New("low confidence")
	}

	token, jabatan, err := getToken(subject, db)
	if err != nil {
		return "Gagal mendapatkan token", 0, err
	}
	fmt.Println("Token:", token, "Jabatan:", jabatan)
	c.JSON(200, gin.H{
		"authorized": true,
		"subject":    subject,
		"similarity": similarity,
		"box":        box,
		"role":       "",
		"department": jabatan,
		"token":      token,
	})

	return subject, similarity, nil
}

func getToken(subject string, db *sqlx.DB) (string, string, error) {
	var user model.UserAccount

	query := `
		SELECT 
			ua.*,
			COALESCE(uk.nama_mesin_absen, '') AS nama_karyawan, 
			uk.jabatan
		FROM user_accounts ua
		LEFT JOIN user_karyawan uk ON uk.id = ua.id_karyawan
		WHERE ua.name = ?
		LIMIT 1
	`
	err := db.Get(&user, query, subject)
	if err != nil {
		fmt.Println("Error finding user by name:", err)
		return "Nama tidak ditemukan", "", err
	}
	fmt.Println("User:", user)
	s, err := CheckStatus(user.Email, db)
	if err != nil {
		fmt.Println("Error checking status:", err)
		return s, user.Jabatan.String, err
	}
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}

	return token, user.Jabatan.String, nil
}

func CheckStatus(email string, db *sqlx.DB) (string, error) {
	var user model.UserAccount
	err := db.Get(&user, "SELECT * FROM user_accounts WHERE email = ? AND status = 1", email)
	if err != nil {
		// fmt.Println("Error finding user by email:", err)
		return "Email tidak aktif", err
	}
	return email, nil
}
