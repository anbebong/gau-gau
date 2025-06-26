package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8082/api"

func main() {
	fmt.Println("--- API Test Tool ---")
	var testLog []map[string]interface{} // slice lưu log từng API

	// 1. Login
	loginReq := map[string]string{"username": "admin", "password": "1"}
	loginResp := postJSON("/login", loginReq)
	fmt.Println("[POST] /login\nInput:", loginReq, "\nOutput:", loginResp, "\n")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/login",
		"method": "POST",
		"input":  loginReq,
		"output": loginResp,
	})

	// Lấy token từ loginResp["data"]["token"]
	token := ""
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
	}
	if token == "" {
		fmt.Println("Login failed, cannot test further APIs.")
		os.Exit(1)
	}
	header := map[string]string{"Authorization": fmt.Sprintf("Bearer %v", token)}

	// 2. List users (admin only)
	usersResp := getWithHeader("/users", header)
	fmt.Println("[GET] /users\nInput: (token)\nOutput:", usersResp, "\n")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/users",
		"method": "GET",
		"input":  "(token)",
		"output": usersResp,
	})

	// Lấy username đầu tiên (không phải admin) để test assign
	username := ""
	if users, ok := usersResp["data"].([]interface{}); ok {
		for _, u := range users {
			if userMap, ok := u.(map[string]interface{}); ok {
				if uname, _ := userMap["username"].(string); uname != "admin" {
					username = uname
					break
				}
			}
		}
	}
	if username == "" {
		fmt.Println("Không tìm thấy user nào để test assign-user!")
		os.Exit(1)
	}

	// 3. List clients
	clientsResp := getWithHeader("/clients", header)
	fmt.Println("[GET] /clients\nInput: (token)\nOutput:", clientsResp, "\n")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/clients",
		"method": "GET",
		"input":  "(token)",
		"output": clientsResp,
	})

	// Lấy agent_id đầu tiên để test
	agentID := ""
	if data, ok := clientsResp["data"].(map[string]interface{}); ok {
		if clients, ok := data["clients"].([]interface{}); ok {
			for _, c := range clients {
				if cMap, ok := c.(map[string]interface{}); ok {
					agentID, _ = cMap["agent_id"].(string)
					if agentID != "" {
						break
					}
				}
			}
		}
	}
	if agentID == "" {
		agentID = "001" // fallback nếu không có client nào
	}

	// 4. Assign user to client (admin only, dùng username)
	assignReq := map[string]string{"agent_id": agentID, "username": username}
	assignResp := postJSONWithHeader("/clients/assign-user", assignReq, header)
	fmt.Println("[POST] /clients/assign-user\nInput:", assignReq, "\nOutput:", assignResp, "\n")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/clients/assign-user",
		"method": "POST",
		"input":  assignReq,
		"output": assignResp,
	})

	// 5. Get logs archive (admin only)
	logsResp := getWithHeader("/logs/archive", header)
	fmt.Println("==============================")
	fmt.Println("[GET] /logs/archive")
	fmt.Println("Input: (token)")
	var logsOut interface{}
	if logs, ok := logsResp["data"].([]interface{}); ok {
		logsOut = logs
		b, _ := json.MarshalIndent(logs, "", "  ")
		fmt.Println("Output:")
		fmt.Println(string(b))
	} else {
		logsOut = logsResp
		b, _ := json.MarshalIndent(logsResp, "", "  ")
		fmt.Println("Output:")
		fmt.Println(string(b))
	}
	fmt.Println("==============================")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/logs/archive",
		"method": "GET",
		"input":  "(token)",
		"output": logsOut,
	})

	// --- FLOW TEST USER MỚI (KHÔNG LẶP LẠI CÁC API ĐÃ TEST VỚI ADMIN) ---
	fmt.Println("\n==============================")
	fmt.Println("[USER FLOW] Tạo user, login, gán thiết bị, lấy log thiết bị")
	// 1. Kiểm tra user test đã tồn tại chưa
	userExisted := false
	if users, ok := usersResp["data"].([]interface{}); ok {
		for _, u := range users {
			if userMap, ok := u.(map[string]interface{}); ok {
				if uname, _ := userMap["username"].(string); uname == "testuser1" {
					userExisted = true
					break
				}
			}
		}
	}
	// 2. Tạo user mới nếu chưa tồn tại
	if !userExisted {
		newUserReq := map[string]string{
			"username":  "testuser1",
			"password":  "testpass1",
			"full_name": "Test User 1",
			"email":     "testuser1@example.com",
		}
		// Gọi đúng endpoint tạo user: POST /users/create
		createUserResp := postJSONWithHeader("/users/create", newUserReq, header)
		fmt.Println("[POST] /users/create\nInput:", newUserReq, "\nOutput:", createUserResp, "\n")
		testLog = append(testLog, map[string]interface{}{
			"api":    "/users/create",
			"method": "POST",
			"input":  newUserReq,
			"output": createUserResp,
		})
	}
	// 3. Login bằng user test
	loginUserReq := map[string]string{"username": "testuser1", "password": "testpass1"}
	loginUserResp := postJSON("/login", loginUserReq)
	fmt.Println("[POST] /login (user)\nInput:", loginUserReq, "\nOutput:", loginUserResp, "\n")
	testLog = append(testLog, map[string]interface{}{
		"api":    "/login",
		"method": "POST",
		"input":  loginUserReq,
		"output": loginUserResp,
	})
	userToken := ""
	if data, ok := loginUserResp["data"].(map[string]interface{}); ok {
		userToken, _ = data["token"].(string)
	}
	if userToken == "" {
		fmt.Println("Login user test thất bại, dừng test user flow!")
		return
	}
	userHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %v", userToken)}
	// 4. Gán thiết bị cho user test (nếu có client)
	if agentID != "" {
		assignUserReq := map[string]string{"agent_id": agentID, "username": "testuser1"}
		assignUserResp := postJSONWithHeader("/clients/assign-user", assignUserReq, header)
		fmt.Println("[POST] /clients/assign-user (testuser1)\nInput:", assignUserReq, "\nOutput:", assignUserResp, "\n")
		testLog = append(testLog, map[string]interface{}{
			"api":    "/clients/assign-user",
			"method": "POST",
			"input":  assignUserReq,
			"output": assignUserResp,
		})
	}
	// 5. Lấy log thiết bị theo user test
	myLogResp := getWithHeader("/logs/my-device", userHeader)
	fmt.Println("[GET] /logs/my-device (testuser1)\nInput: (user token)")
	var myLogOut interface{}
	if logs, ok := myLogResp["data"].([]interface{}); ok {
		myLogOut = logs
		b, _ := json.MarshalIndent(logs, "", "  ")
		fmt.Println("Output:")
		fmt.Println(string(b))
	} else {
		myLogOut = myLogResp
		b, _ := json.MarshalIndent(myLogResp, "", "  ")
		fmt.Println("Output:")
		fmt.Println(string(b))
	}
	testLog = append(testLog, map[string]interface{}{
		"api":    "/logs/my-device",
		"method": "GET",
		"input":  "(user token)",
		"output": myLogOut,
	})

	// 6. Lấy OTP cho agent_id (nếu có client)
	if agentID != "" {
		otpResp := getWithHeader("/clients/"+agentID+"/otp", header)
		fmt.Println("[GET] /clients/"+agentID+"/otp\nInput: (token)\nOutput:", otpResp, "\n")
		testLog = append(testLog, map[string]interface{}{
			"api":    "/clients/" + agentID + "/otp",
			"method": "GET",
			"input":  "(token)",
			"output": otpResp,
		})
	}

	// Ghi toàn bộ log test ra file
	f, err := os.Create("test_output.json")
	if err == nil {
		json.NewEncoder(f).Encode(testLog)
		f.Close()
	}
}

func postJSON(path string, data map[string]string) map[string]interface{} {
	b, _ := json.Marshal(data)
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var out map[string]interface{}
	json.Unmarshal(body, &out)
	return out
}

func postJSONWithHeader(path string, data map[string]string, header map[string]string) map[string]interface{} {
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var out map[string]interface{}
	json.Unmarshal(body, &out)
	return out
}

func getWithHeader(path string, header map[string]string) map[string]interface{} {
	req, _ := http.NewRequest("GET", baseURL+path, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var out map[string]interface{}
	json.Unmarshal(body, &out)
	return out
}
