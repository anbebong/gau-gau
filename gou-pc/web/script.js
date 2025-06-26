// API Configuration
const API_CONFIG = {
  baseURL: "http://localhost:8082", // Thay đổi theo Go server
  endpoints: {
    login: "/api/login",
    logout: "/api/logout",
    clients: "/api/clients",
    deleteClient: "/clients/delete",
    assignUser: "/api/clients/assign-user",
    otp: "/api/otp",
    createUser: "/api/users/create",
    changePassword: "/api/users/change-password",
    updateUser: "/api/users/update-info", // Updated endpoint
    // Thêm endpoint mock cho việc lấy danh sách người dùng
    getUsers: "/api/users",
    logs: "/api/logs/archive",
    myDeviceLogs: "/api/logs/my-device", // <--- Thêm endpoint mới
  },
}

async function apiCall(endpoint, options = {}) {
  const url = `${API_CONFIG.baseURL}${endpoint}`
  const defaultOptions = {
    headers: {
      "Content-Type": "application/json",
      Authorization: getAuthToken() ? `Bearer ${getAuthToken()}` : undefined,
    },
  }
  try {
    const response = await fetch(url, { ...defaultOptions, ...options })
    if (response.status === 401) {
      localStorage.removeItem("authToken")
      localStorage.removeItem("currentUser")
      window.location.reload()
      return
    }
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return await response.json()
  } catch (error) {
    console.error(`API call to ${endpoint} failed:`, error)
    return getMockData(endpoint, currentUser)
  }
}

function getAuthToken() {
  return localStorage.getItem("authToken") || ""
}

// Global variables to store data
let loadedClients = []
let loadedUsers = [] // New global variable for users

// Mock Data Functions (for development/testing)
function getMockData(endpoint, currentUser) {
  const mockData = {
    "/clients": [
      {
        agentID: "AG-001-X7Y9",
        username: "admin",
        isOnline: true,
        lastSeen: "2025-06-24T09:00:00Z",
        hardwareInfo: {
          hostID: "HW-7F8A9B2C",
          ipAddress: "192.168.1.101",
          hostName: "PC-Main-Office", // Added hostName
        },
      },
      {
        agentID: "AG-002-M4N8",
        username: "user01",
        isOnline: true,
        lastSeen: "2025-06-24T08:45:00Z",
        hardwareInfo: {
          hostID: "HW-3D5E6F1A",
          ipAddress: "192.168.1.102",
          hostName: "Laptop-User01", // Added hostName
        },
      },
      {
        agentID: "AG-003-P7Q2",
        username: null,
        isOnline: false,
        lastSeen: "2025-06-24T08:30:00Z",
        hardwareInfo: {
          hostID: "HW-9B1C2D4E",
          ipAddress: "192.168.1.103",
          hostName: "Server-Dev", // Added hostName
        },
      },
    ],
    "/api/otp": [
      {
        clientID: "72300b6b-3e6e-421b-84d3-24b7d4df4a08",
        agentID: "AG-001-X7Y9", // Changed to match mock client agentID
        otp: "839601",
        expiresTime: 27, // Using expiresTime as per user's output
      },
      {
        clientID: "PC-002",
        agentID: "AG-002-M4N8",
        otp: "789012",
        expiresTime: 45,
      },
    ],
    // Dữ liệu mock cho danh sách người dùng
    "/api/users": [
      {
        username: "admin",
        email: "admin@example.com",
        name: "Administrator",
        role: "admin",
        device: "AG-001-X7Y9", // Changed to match mock client agentID
        status: "active",
        lastLogin: "2024-01-15 14:25:33",
      },
      {
        username: "user01",
        email: "user01@example.com",
        name: "Nguyễn Văn A",
        role: "user",
        device: "AG-002-M4N8", // Changed to match mock client agentID
        status: "active",
        lastLogin: "2024-01-15 13:45:12",
      },
      {
        username: "user02",
        email: "user02@example.com",
        name: "Trần Thị B",
        role: "user",
        device: null,
        status: "active",
        lastLogin: "2024-01-15 13:00:00",
      },
    ],
    // Dữ liệu mock cho log hệ thống
    "/api/logs/archive": [
      {
        agentID: "AG-001-X7Y9",
        level: "WARNING",
        message: "- Memory allocation failed",
        timestamp: "2025/06/24 10:41:05",
      },
      {
        agentID: "AG-001-X7Y9",
        level: "ERROR",
        message: "- Memory allocation failed",
        timestamp: "2025/06/24 10:41:06",
      },
      {
        agentID: "AG-001-X7Y9",
        level: "DEBUG",
        message: "- Disk space is running low",
        timestamp: "2025/06/24 10:41:10",
      },
      { agentID: "AG-002-M4N8", level: "INFO", message: "- System started normally", timestamp: "2025/06/24 10:42:00" },
      {
        agentID: "AG-003-P7Q2",
        level: "CRITICAL",
        message: "- Hard drive failure detected",
        timestamp: "2025/06/24 10:43:15",
      },
    ],
    // Mock response for /api/users/update-info
    "/api/users/update-info": {
      status: "updated",
    },
    // Dữ liệu mock cho log hệ thống của người dùng
    "/api/logs/my-device": (currentUser) => {
      if (!currentUser || !currentUser.device) {
        return [] // Trả về mảng rỗng nếu không có người dùng hoặc không có thiết bị được gán
      }
      // Lọc log dựa trên agentID của thiết bị được gán cho người dùng hiện tại
      const assignedAgentId = currentUser.device
      return mockData["/api/logs/archive"].filter((log) => log.agentID === assignedAgentId)
    },
  }

  // Nếu endpoint là một hàm (như /api/logs/my-device), gọi nó với currentUser
  if (typeof mockData[endpoint] === "function") {
    return mockData[endpoint](currentUser)
  }
  return mockData[endpoint] || {}
}

// User data simulation
const users = {
  admin: {
    username: "admin",
    password: "admin123",
    role: "admin",
    fullName: "Administrator",
    email: "admin@example.com",
    device: "AG-001-X7Y9", // Thêm thông tin thiết bị cho admin mock
  },
  user: {
    username: "user",
    password: "user123",
    role: "user",
    fullName: "Nguyễn Văn A",
    email: "user@example.com",
    device: "AG-002-M4N8", // Thêm thông tin thiết bị cho user mock
  },
}

// Current session
let currentUser = null

// DOM Elements
const loginPage = document.getElementById("loginPage")
const adminDashboard = document.getElementById("adminDashboard")
const userDashboard = document.getElementById("userDashboard")
const loginForm = document.getElementById("loginForm")
const loginError = document.getElementById("loginError")

// Initialize the application
document.addEventListener("DOMContentLoaded", () => {
  console.log("DOM Content Loaded. Initializing application.")
  // Check if user is already logged in
  checkSession()

  // Setup event listeners
  setupEventListeners()
})

// Check for existing session
function checkSession() {
  console.log("Checking for existing session...")
  const savedUser = localStorage.getItem("currentUser")
  if (savedUser) {
    currentUser = JSON.parse(savedUser)
    console.log("Session found for user:", currentUser.username)
    showDashboard(currentUser.role)
  } else {
    console.log("No session found. Showing login page.")
    showPage("loginPage")
  }
}

// Setup all event listeners
function setupEventListeners() {
  console.log("Setting up event listeners...")
  // Login form submission
  loginForm.addEventListener("submit", handleLogin)

  // Logout buttons
  document.getElementById("adminLogout").addEventListener("click", handleLogout)
  document.getElementById("userLogout").addEventListener("click", handleLogout)

  // Navigation links
  setupNavigation()
}

// Handle logout
function handleLogout() {
  console.log("Handling logout...")
  currentUser = null
  localStorage.removeItem("authToken") // Clear auth token
  localStorage.removeItem("currentUser")
  showPage("loginPage")
  showNotification("Đã đăng xuất!", "info")
}

// Show appropriate dashboard based on role
function showDashboard(role) {
  console.log("Showing dashboard for role:", role)
  if (role === "admin") {
    showPage("adminDashboard")
    document.getElementById("adminUsername").textContent = currentUser.fullName
    setupAdminNavigation() // This calls setupAdminDashboard which loads data
  } else if (role === "user") {
    showPage("userDashboard")
    document.getElementById("userUsername").textContent = currentUser.fullName
    setupUserNavigation()
    setupUserDashboard() // New function to load user-specific data
  }
}

// Show specific page
function showPage(pageId) {
  console.log("Activating page:", pageId)
  // Hide all pages
  document.querySelectorAll(".page").forEach((page) => {
    page.classList.remove("active")
  })

  // Show selected page
  document.getElementById(pageId).classList.add("active")
}

// Setup navigation for admin dashboard
function setupAdminNavigation() {
  console.log("Setting up Admin Navigation...")
  const navLinks = document.querySelectorAll("#adminDashboard .nav-link")
  const contentSections = document.querySelectorAll("#adminDashboard .content-section")

  navLinks.forEach((link) => {
    link.addEventListener("click", function (e) {
      e.preventDefault()

      const targetSection = this.getAttribute("data-section")
      console.log("Admin Nav clicked:", targetSection)

      // Remove active class from all links and sections
      navLinks.forEach((l) => l.classList.remove("active"))
      contentSections.forEach((s) => s.classList.remove("active"))

      // Add active class to clicked link and corresponding section
      this.classList.add("active")
      document.getElementById(targetSection).classList.add("active")
    })
  })

  // Setup dashboard specific features
  setupAdminDashboard()
}

// Setup navigation for user dashboard
function setupUserNavigation() {
  console.log("Setting up User Navigation...")
  const navLinks = document.querySelectorAll("#userDashboard .nav-link")
  const contentSections = document.querySelectorAll("#userDashboard .content-section")

  navLinks.forEach((link) => {
    link.addEventListener("click", function (e) {
      e.preventDefault()

      const targetSection = this.getAttribute("data-section")
      console.log("User Nav clicked:", targetSection)

      // Remove active class from all links and sections
      navLinks.forEach((l) => l.classList.remove("active"))
      contentSections.forEach((s) => s.classList.remove("active"))

      // Add active class to clicked link and corresponding section
      this.classList.add("active")
      document.getElementById(targetSection).classList.add("active")
    })
  })
}

// Setup general navigation
function setupNavigation() {
  // This function can be extended for additional navigation features
  console.log("General navigation setup complete")
}

// Utility function to format numbers
function formatNumber(num) {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
}

// Utility function to format currency
function formatCurrency(amount) {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
  }).format(amount)
}

// // Simulate real-time updates (optional)
// function startRealTimeUpdates() {
//   console.log("Starting real-time updates simulation...")
//   setInterval(() => {
//     if (currentUser && currentUser.role === "admin") {
//       // Update stats with random values (simulation)
//       updateStats()
//     }
//   }, 30000) // Update every 30 seconds
// }

// Update statistics (simulation)
function updateStats() {
  console.log("Updating stats (simulation)...")
  const statNumbers = document.querySelectorAll(".stat-number")
  statNumbers.forEach((stat) => {
    if (stat.textContent.includes("₫")) {
      // Update currency values
      const currentValue = Number.parseInt(stat.textContent.replace(/[₫,]/g, ""))
      const newValue = currentValue + Math.floor(Math.random() * 1000000)
    } else if (!isNaN(Number.parseInt(stat.textContent.replace(/,/g, "")))) {
      // Update numeric values
      const currentValue = Number.parseInt(stat.textContent.replace(/,/g, ""))
      const newValue = currentValue + Math.floor(Math.random() * 10)
      stat.textContent = formatNumber(newValue)
    }
  })
}

// Start real-time updates when page loads
// startRealTimeUpdates();

// Additional features can be added here:
// - Form validation
// - AJAX requests simulation
// - Data export functionality
// - Advanced filtering and search
// - Notifications system
// - Theme switching
// - Multi-language support

// Example: Simple notification system
function showNotification(message, type = "info") {
  console.log(`Notification (${type}): ${message}`)
  const notification = document.createElement("div")
  notification.className = `notification ${type}`
  notification.textContent = message

  // Style the notification
  notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem 1.5rem;
        background-color: ${type === "success" ? "#28a745" : type === "error" ? "#dc3545" : "#17a2b8"};
        color: white;
        border-radius: 5px;
        box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        z-index: 1000;
        animation: slideIn 0.3s ease;
    `

  document.body.appendChild(notification)

  // Remove notification after 3 seconds
  setTimeout(() => {
    notification.style.animation = "slideOut 0.3s ease"
    setTimeout(() => {
      document.body.removeChild(notification)
    }, 300)
  }, 3000)
}

// Add CSS for notification animations
const style = document.createElement("style")
style.textContent = `
  @keyframes slideIn {
      from { transform: translateX(100%); opacity: 0; }
      to { transform: translateX(0); opacity: 1; }
  }
  
  @keyframes slideOut {
      from { transform: translateX(0); opacity: 1; }
      to { transform: translateX(100%); opacity: 0; }
  }

  .countdown-text {
    font-size: 1.5rem; /* Larger font size */
    font-weight: bold; /* Make it bold */
    color: #667eea; /* Default color */
    transition: color 0.5s ease-in-out; /* Smooth color transition */
  }

  .countdown-text.warning {
    color: #d69e2e; /* Warning color */
  }

  .countdown-text.critical {
    color: #e53e3e; /* Critical color */
  }
`
document.head.appendChild(style)

// Example usage of notification (can be called from anywhere)
// showNotification('Đăng nhập thành công!', 'success');
// showNotification('Có lỗi xảy ra!', 'error');
// showNotification('Thông tin đã được cập nhật.', 'info');

// Thêm vào cuối file script.js

// Device management functions
function editDevice(agentId) {
  console.log("Editing device:", agentId)
  const modal = document.getElementById("deviceModal")
  const targetDevice = loadedClients.find((device) => device.agentID === agentId)

  if (targetDevice) {
    // Populate modal fields
    document.getElementById("deviceId").value = targetDevice.hardwareInfo.hostName || "N/A" // Now displays Device Name
    document.getElementById("deviceIp").value = targetDevice.hardwareInfo.ipAddress
    document.getElementById("deviceHardwareId").value = targetDevice.hardwareInfo.hostID
    document.getElementById("deviceAgentId").value = targetDevice.agentID
    document.getElementById("deviceUser").value = targetDevice.username || ""

    document.getElementById("modalTitle").textContent = "Chỉnh sửa thiết bị"
    document.getElementById("deviceId").readOnly = true // Device Name is display only
    document.getElementById("deviceAgentId").readOnly = true // Agent ID is display only

    modal.style.display = "block"
  } else {
    showNotification("Không tìm thấy thiết bị để chỉnh sửa!", "error")
    console.error("Device not found for editing:", agentId)
  }
}

function closeDeviceModal() {
  console.log("Closing device modal.")
  document.getElementById("deviceModal").style.display = "none"
}

// Setup user assignment change handlers
// function setupUserAssignments() {
//   const userSelects = document.querySelectorAll(".user-assignment")
//   userSelects.forEach((select) => {
//     select.addEventListener("change", function () {
//       const deviceId = this.getAttribute("data-device")
//       const userId = this.value
//       const userName = userId || "Chưa gán"

//       showNotification(`Thiết bị ${deviceId} đã được gán cho ${userName}`, "success")

//       // Log the assignment change
//       addAssignmentLog(deviceId, userId)
//     })
//   })
// }

// Add assignment log entry
function addAssignmentLog(deviceId, userId) {
  console.log(`Adding assignment log for device ${deviceId} to user ${userId}`)
  const systemLogsTable = document.getElementById("systemLogs")
  if (systemLogsTable) {
    const newRow = document.createElement("tr")
    const currentTime = new Date().toLocaleString("vi-VN")
    const userName = userId || "Không có"

    newRow.innerHTML = `
    <td>${currentTime}</td>
    <td>${deviceId}</td>
    <td><span class="log-level info">INFO</span></td>
    <td>User Assignment</td>
    <td>Thiết bị được gán cho người dùng: ${userName}</td>
`

    // Insert at the beginning
    systemLogsTable.insertBefore(newRow, systemLogsTable.firstChild)

    // Remove last row if too many
    if (systemLogsTable.children.length > 20) {
      systemLogsTable.removeChild(systemLogsTable.lastChild)
    }
  }
}

// Mock functions to resolve the "undeclared variables" errors
function setupDeviceFilter() {
  console.log("setupDeviceFilter called")
}

// Removed setupRefreshLogs as per user request
// function setupRefreshLogs() {
//   console.log("setupRefreshLogs called")
// }

// Enhanced setup function

// Load Dashboard Statistics
async function loadDashboardStats() {
  console.log("Loading dashboard stats...")
  try {
    const response = await apiCall(API_CONFIG.endpoints.clients)
    let clients = []
    if (Array.isArray(response.data)) {
      clients = response.data
    } else if (Array.isArray(response.clients)) {
      clients = response.clients
    } else if (Array.isArray(response)) {
      clients = response
    } else {
      clients = []
    }
    // Fallback mock nếu không có dữ liệu
    if (!clients.length && typeof getMockData === 'function') {
      const mockClients = getMockData(API_CONFIG.endpoints.clients, currentUser)
      if (Array.isArray(mockClients)) clients = mockClients
    }
    clients = clients.filter(c => c && c.agentID && c.hardwareInfo)
    const stats = calculateStatsFromClients(clients)
    updateStatsDisplay(stats)
    console.log("Dashboard stats loaded successfully.")
  } catch (error) {
    console.error("Failed to load stats:", error)
    showNotification("Lỗi khi tải thống kê dashboard!", "error")
  }
}

function updateStatsDisplay(stats) {
  console.log("Updating stats display with:", stats)
  if (document.getElementById("totalPCs")) document.getElementById("totalPCs").textContent = stats.totalPCs || 0
  if (document.querySelector(".online")) document.querySelector(".online").textContent = `${stats.onlinePCs || 0} Online`
  if (document.querySelector(".offline")) document.querySelector(".offline").textContent = `${stats.offlinePCs || 0} Offline`
  if (document.getElementById("totalUsers")) document.getElementById("totalUsers").textContent = stats.totalUsers || 0
  if (document.querySelector(".active-users")) document.querySelector(".active-users").textContent = `${stats.activeUsers || 0} Đang hoạt động`
  if (document.getElementById("totalAlerts")) document.getElementById("totalAlerts").textContent = stats.totalAlerts || 0
  if (document.querySelector(".critical")) document.querySelector(".critical").textContent = `${stats.criticalAlerts || 0} Nghiêm trọng`
  if (document.querySelector(".warning")) document.querySelector(".warning").textContent = `${stats.warningAlerts || 0} Cảnh báo`
}

// Load Devices
async function loadDevices() {
  console.log("Loading devices...")
  try {
    const response = await apiCall(API_CONFIG.endpoints.clients)
    let clients = []
    // Map đúng với API thực tế bạn gửi mẫu
    if (Array.isArray(response.data)) {
      clients = response.data.map(item => ({
        agentID: item.agent_id || item.agentID || item.id || '',
        username: item.user_id || item.username || '',
        isOnline: false, // Nếu backend không trả thì mặc định offline
        lastSeen: null,
        hardwareInfo: {
          hostID: item.device_info?.hardwareID || item.device_info?.hostID || '',
          ipAddress: item.device_info?.ipAddress || '',
          hostName: item.device_info?.hostName || '',
        },
      }))
    }
    // Fallback mock nếu không có dữ liệu
    if (!clients.length && typeof getMockData === 'function') {
      const mockClients = getMockData(API_CONFIG.endpoints.clients, currentUser)
      if (Array.isArray(mockClients)) clients = mockClients
    }
    clients = clients.filter(c => c && c.agentID && c.hardwareInfo)
    loadedClients = clients
    updateDevicesTable(clients)
    console.log("Devices loaded successfully:", clients.length, "devices.")
  } catch (error) {
    console.error("Failed to load devices:", error)
    showNotification("Lỗi khi tải danh sách thiết bị!", "error")
  }
}

function updateDevicesTable(clients) {
  console.log("Updating devices table...")
  const tbody = document.getElementById("devicesTable")
  tbody.innerHTML = ""

  clients.forEach((client) => {
    const row = document.createElement("tr")
    let userOptions = '<option value="">Chưa gán</option>'
    loadedUsers.forEach((user) => {
      const selected = client.username === user.username ? "selected" : ""
      userOptions += `<option value="${user.username}" ${selected}>${user.username}</option>`
    })

    row.innerHTML = `
      <td>${client.agentID}</td>
      <td>${client.hardwareInfo.hostName || "N/A"}</td>
      <td>${client.hardwareInfo.ipAddress}</td>
      <td>${client.hardwareInfo.hostID}</td>
      <td>
        <select class="user-assignment" data-agent-id="${client.agentID}">
          ${userOptions}
        </select>
      </td>
      <td><span class="status ${client.isOnline ? "online" : "offline"}">${client.isOnline ? "Online" : "Offline"}</span></td>
      <td>
        <button class="btn btn-sm btn-danger" onclick="deleteDevice('${client.agentID}')">Xóa</button>
        <button class="btn btn-sm btn-success" onclick="getOTP('${client.agentID}')">OTP</button>
      </td>
    `
    tbody.appendChild(row)
  })

  // Re-setup user assignment handlers
  setupUserAssignments()
  console.log("Devices table updated.")
}

// Load System Logs
async function loadSystemLogs(agentId = "all") {
  console.log("Loading system logs for agent:", agentId)
  try {
    let url
    if (currentUser && currentUser.role !== "admin") {
      // Nếu là user thường, chỉ lấy log thiết bị của mình
      url = API_CONFIG.endpoints.myDeviceLogs
    } else {
      url = API_CONFIG.endpoints.logs
      if (agentId && agentId !== "all") {
        url += `?agent=${agentId}`
      }
    }
    const response = await apiCall(url)
    const logs = response.data || response
    updateLogsTable(logs)
    console.log("System logs loaded successfully:", logs.length, "logs.")
  } catch (error) {
    console.error("Failed to load logs:", error)
    showNotification("Lỗi khi tải log hệ thống!", "error")
  }
}

function updateLogsTable(logs) {
  console.log("Updating logs table...")
  const tbody = document.getElementById("systemLogs")
  tbody.innerHTML = ""

  logs.forEach((log) => {
    const row = document.createElement("tr")
    // Map API data to table columns
    const deviceId = log.agent_id || "N/A"
    const level = log.level ? log.level.toLowerCase() : "info"
    const event = log.message || "Log Entry" // Use message as event if available
    const details = log.details || "" // Details might be empty or not provided
    const timestamp = (log.time && String(log.time).trim() !== "") ? log.time : new Date().toLocaleString("vi-VN")

    row.innerHTML = `
     <td>${timestamp}</td>
     <td>${deviceId}</td>
     <td><span class="log-level ${level}">${level.toUpperCase()}</span></td>
     <td>${event}</td>
     <td>${details}</td>
`
    tbody.appendChild(row)
  })
  console.log("Logs table updated.")
}

// Loading States
function showLoadingState() {
  console.log("Showing loading state...")
  const loadingHTML = `
    <div class="loading-overlay">
      <div class="loading-spinner"></div>
      <p>Đang tải dữ liệu...</p>
    </div>
  `
  document.body.insertAdjacentHTML("beforeend", loadingHTML)
}

function hideLoadingState() {
  console.log("Hiding loading state...")
  const loadingOverlay = document.querySelector(".loading-overlay")
  if (loadingOverlay) {
    loadingOverlay.remove()
  }
}

// API-based Device Management
async function saveDevice() {
  console.log("Saving device...")
  const isEdit = document.getElementById("deviceId").readOnly // Check if deviceId is readonly (meaning it's an edit)

  const deviceName = document.getElementById("deviceId").value
  const ipAddress = document.getElementById("deviceIp").value
  const hardwareId = document.getElementById("deviceHardwareId").value
  const agentId = document.getElementById("deviceAgentId").value // This will be empty for new devices
  const assignedUser = document.getElementById("deviceUser").value

  const deviceData = {
    hardwareInfo: {
      hostName: deviceName,
      ipAddress: ipAddress,
      hostID: hardwareId,
    },
    username: assignedUser,
  }

  let endpoint = API_CONFIG.endpoints.clients
  let method = "POST"

  if (isEdit) {
    // For editing, agentID is already set and used in the URL
    endpoint = `${API_CONFIG.baseURL}${API_CONFIG.endpoints.clients}/${agentId}`
    method = "PUT"
    // Ensure agentID is also in the body for consistency if API expects it
    deviceData.agentID = agentId
  } else {
    // For adding, agentID might be generated by backend or needs to be provided.
    if (!agentId) {
      showNotification("Vui lòng nhập Agent ID cho thiết bị mới!", "error")
      return
    }
    deviceData.agentID = agentId
  }

  try {
    const response = await apiCall(endpoint, {
      method: method,
      body: JSON.stringify(deviceData),
    })

    if (response.status === "success") {
      // Assuming API returns status: "success"
      showNotification("Thông tin thiết bị đã được cập nhật!", "success")
      closeDeviceModal()
      await loadDevices()
      console.log("Device saved successfully.")
    } else {
      showNotification("Lỗi khi lưu thiết bị!", "error")
      console.error("API response error when saving device:", response)
    }
  } catch (error) {
    showNotification("Lỗi khi lưu thiết bị!", "error")
    console.error("Failed to save device:", error)
  }
}

async function handleLogin(e) {
  e.preventDefault()
  console.log("Attempting login...")
  const username = document.getElementById("username").value
  const password = document.getElementById("password").value
  loginError.textContent = ""
  try {
    const response = await apiCall(API_CONFIG.endpoints.login, {
      method: "POST",
      body: JSON.stringify({ username, password }),
    })
    if (response.success && response.data && response.data.token && response.data.user) {
      currentUser = {
        id: response.data.user.id,
        username: response.data.user.username,
        role: response.data.user.role,
        email: response.data.user.email,
        fullName: response.data.user.name,
        device: response.data.user.device,
      }
      localStorage.setItem("authToken", response.data.token)
      localStorage.setItem("currentUser", JSON.stringify(currentUser))
      showDashboard(currentUser.role)
      loginForm.reset()
      showNotification("Đăng nhập thành công!", "success")
      console.log("Login successful for:", username)
    } else {
      loginError.textContent = response.message || "Đăng nhập thất bại!"
      console.warn("Login failed for:", username, "Message:", response.message)
    }
  } catch (error) {
    // Fallback to local authentication for development
    console.error("Login API call failed, attempting local fallback:", error)
    if (users[username] && users[username].password === password) {
      currentUser = {
        username: users[username].username,
        role: users[username].role,
        fullName: users[username].fullName,
        email: users[username].email,
        device: users[username].device,
      }
      localStorage.setItem("currentUser", JSON.stringify(currentUser))
      showDashboard(currentUser.role)
      loginForm.reset()
      showNotification("Đăng nhập thành công (local fallback)!", "success")
      console.log("Login successful via local fallback for:", username)
    } else {
      loginError.textContent = "Tên đăng nhập hoặc mật khẩu không đúng!"
      console.warn("Login failed for:", username, "Invalid credentials (local fallback).")
    }
  }
}

function updateSystemStats() {
  console.log("updateSystemStats called (placeholder)")
}

function addNewLogEntry() {
  console.log("addNewLogEntry called (placeholder)")
}

// Load Users
async function loadUsers() {
  console.log("Loading users...")
  try {
    const response = await apiCall(API_CONFIG.endpoints.getUsers)
    const users = response.data || response // Chuẩn hóa lấy từ response.data nếu có
    loadedUsers = users
    updateUsersTable(users)
    console.log("Users loaded successfully:", users.length, "users.")
  } catch (error) {
    console.error("Failed to load users:", error)
    showNotification("Lỗi khi tải danh sách người dùng!", "error")
  }
}

function updateUsersTable(users) {
  console.log("Updating users table...")
  const tbody = document.getElementById("usersTable")
  tbody.innerHTML = ""

  users.forEach((user) => {
    const row = document.createElement("tr")
    row.innerHTML = `
      <td>${user.username}</td>
      <td>${user.email}</td>
      <td>${user.full_name || "N/A"}</td>
      <td><span class="badge ${user.role}">${user.role.charAt(0).toUpperCase() + user.role.slice(1)}</span></td>
      <td>${user.device || "N/A"}</td>
      <td>
        <button class="btn btn-sm btn-fixed" onclick="editUser('${user.username}')">Sửa</button>
        <button class="btn btn-sm btn-change-password" onclick="changeUserPassword('${user.username}')">Đổi pass</button>
      </td>
    `
    tbody.appendChild(row)
  })
  console.log("Users table updated.")
}

function calculateStatsFromClients(clients) {
  console.log("Calculating stats from clients...")
  const totalPCs = clients.length
  const onlinePCs = clients.filter((client) => client.isOnline).length
  const offlinePCs = totalPCs - onlinePCs
  const assignedUsers = new Set(clients.map((c) => c.username).filter(Boolean)).size // Count unique assigned users

  return {
    totalPCs,
    onlinePCs,
    offlinePCs,
    totalUsers: assignedUsers,
    activeUsers: onlinePCs, // Giả sử user active = PC online
    totalAlerts: 0, // Set to 0 as alerts section is "coming soon"
    criticalAlerts: 0, // Set to 0
    warningAlerts: 0, // Set to 0
    // Removed systemPerformance
  }
}

// Cập nhật hàm setupUserAssignments để sử dụng assign-user API
function setupUserAssignments() {
  console.log("Setting up user assignment handlers...")
  const userSelects = document.querySelectorAll(".user-assignment")
  userSelects.forEach((select) => {
    select.addEventListener("change", async function () {
      const agentId = this.getAttribute("data-agent-id")
      const username = this.value
      console.log(`Attempting to assign device ${agentId} to user ${username}`)
      try {
        // Gửi đúng format backend yêu cầu: agent_id, username
        const response = await apiCall(API_CONFIG.endpoints.assignUser, {
          method: "POST",
          body: JSON.stringify({
            agent_id: agentId,
            username: username,
          }),
        })
        // Accept both status: "success" and success: true as success
        if (response.status === "success" || response.success === true) {
          showNotification(response.message, "success")
          await loadDevices()
          console.log(`Device ${agentId} assigned to ${username} successfully.`)
        } else {
          showNotification("Lỗi khi gán user!", "error")
          console.error("API response error when assigning user:", response)
        }
      } catch (error) {
        showNotification("Lỗi khi gán user!", "error")
        console.error("Failed to assign user!", error)
      }
    })
  })
}

// Thêm hàm deleteDevice sử dụng Go API
async function deleteDevice(agentId) {
  console.log("Attempting to delete device:", agentId)
  if (confirm(`Bạn có chắc chắn muốn xóa thiết bị ${agentId}?`)) {
    try {
      const response = await fetch(`${API_CONFIG.baseURL}${API_CONFIG.endpoints.deleteClient}?id=${agentId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
        },
      })

      if (response.ok) {
        const message = await response.text()
        showNotification(message, "success")
        await loadDevices() // Reload devices table
        console.log(`Device ${agentId} deleted successfully.`)
      } else {
        showNotification("Lỗi khi xóa thiết bị!", "error")
        console.error("API response error when deleting device:", response.status, await response.text())
      }
    } catch (error) {
      showNotification("Lỗi khi xóa thiết bị!", "error")
      console.error("Failed to delete device!", error)
    }
  }
}

// Thêm hàm getOTP
async function getOTP(agentId) {
  console.log("Requesting OTP for device:", agentId)
  try {
    const response = await apiCall(`${API_CONFIG.endpoints.otp}?id=${agentId}`)
    const otpData = response.data || response
    if (otpData && otpData.otp) {
      const expiresIn = otpData.expire_in !== undefined ? otpData.expire_in : otpData.expiresTime
      showOTPModal(otpData.otp, otpData.agentID, expiresIn)
      console.log(`OTP received for ${agentId}: ${otpData.otp}`)
    } else {
      showNotification("Không thể lấy OTP cho thiết bị này!", "error")
      console.error("OTP response invalid or missing OTP:", response)
    }
  } catch (error) {
    showNotification("Lỗi khi lấy OTP!", "error")
    console.error("Failed to get OTP!", error)
  }
}

function showOTPModal(otp, agentID, expiresInSeconds) {
  console.log(`Showing OTP modal for ${agentID} with OTP ${otp}, expires in ${expiresInSeconds}s`)
  const modal = document.createElement("div")
  modal.className = "modal"
  modal.style.display = "block"
  modal.innerHTML = `
    <div class="modal-content">
      <div class="modal-header">
        <h3>OTP cho thiết bị ${agentID}</h3>
        <span class="close" onclick="this.closest('.modal').remove()">&times;</span>
      </div>
      <div class="modal-body">
        <div style="text-align: center; padding: 2rem;">
          <h2 style="font-size: 3rem; color: #667eea; margin-bottom: 1rem;">${otp}</h2>
          <p>Mã OTP sẽ hết hạn sau <span id="countdown" class="countdown-text">${expiresInSeconds}</span> giây</p>
          <button class="btn btn-primary" onclick="copyOTP('${otp}')">Sao chép OTP</button>
        </div>
      </div>
    </div>
  `

  document.body.appendChild(modal)

  // Countdown timer
  let timeLeft = expiresInSeconds
  const countdownElement = modal.querySelector("#countdown") // Get element from within the modal

  const countdownInterval = setInterval(() => {
    timeLeft--
    if (countdownElement) {
      countdownElement.textContent = timeLeft

      // Add visual cues based on remaining time
      if (timeLeft <= 5) {
        countdownElement.classList.add("critical")
        countdownElement.classList.remove("warning")
      } else if (timeLeft <= 15) {
        countdownElement.classList.add("warning")
        countdownElement.classList.remove("critical")
      } else {
        countdownElement.classList.remove("warning", "critical")
      }
    }

    if (timeLeft <= 0) {
      clearInterval(countdownInterval)
      modal.remove()
      showNotification("OTP đã hết hạn!", "warning")
      console.log(`OTP for ${agentID} expired.`)
    }
  }, 1000)
}

function copyOTP(otp) {
  console.log("Copying OTP:", otp)
  navigator.clipboard
    .writeText(otp)
    .then(() => {
      showNotification("Đã sao chép OTP!", "success")
    })
    .catch(() => {
      showNotification("Không thể sao chép OTP!", "error")
    })
}

// Thêm hàm tạo user mới
async function createUser(userData) {
  console.log("Creating user:", userData.username)
  try {
    const response = await apiCall(API_CONFIG.endpoints.createUser, {
      method: "POST",
      body: JSON.stringify(userData),
    })

    if (response.status === "success") {
      showNotification(response.message, "success")
      console.log("User created successfully:", userData.username)
      return true
    } else {
      showNotification("Lỗi khi tạo user!", "error")
      console.error("API response error when creating user:", response)
      return false
    }
  } catch (error) {
    showNotification("Lỗi khi tạo user!", "error")
    console.error("Failed to create user!", error)
    return false
  }
}

// Thêm hàm đổi mật khẩu
async function changePassword(passwordData) {
  console.log("Changing password for user:", passwordData.username)
  try {
    const response = await apiCall(API_CONFIG.endpoints.changePassword, {
      method: "POST",
      body: JSON.stringify(passwordData),
    })

    if (response.status === "success") {
      showNotification(response.message, "success")
      console.log("Password changed successfully for:", passwordData.username)
      return true
    } else {
      showNotification("Lỗi khi đổi mật khẩu!", "error")
      console.error("API response error when changing password:", response)
      return false
    }
  } catch (error) {
    showNotification("Lỗi khi đổi mật khẩu!", "error")
    console.error("Failed to change password!", error)
    return false
  }
}

// Thêm hàm cập nhật user
async function updateUser(userData) {
  const { username, email, name } = userData // Destructure only allowed fields for this API
  console.log("Attempting to update user:", username, "with data:", { username, email, name })

  // Access control check: User can only update their own info
  if (currentUser.role === "user" && currentUser.username !== username) {
    showNotification("Bạn chỉ có thể cập nhật thông tin của chính mình!", "error")
    console.warn("Access denied: User tried to update another user's profile.")
    return false
  }

  try {
    const response = await apiCall(API_CONFIG.endpoints.updateUser, {
      method: "POST", // User specified POST
      body: JSON.stringify({
        username: username, // Identifier, not changeable
        name: name,
        email: email,
      }),
    })

    if (response && response.status === "updated") {
      // Check for the new response format
      showNotification("Thông tin người dùng đã được cập nhật!", "success")
      console.log("User updated successfully:", username)
      return true
    } else {
      showNotification(response.message || "Lỗi khi cập nhật user!", "error")
      console.error("API response error when updating user:", response)
      return false
    }
  } catch (error) {
    showNotification("Lỗi khi cập nhật user!", "error")
    console.error("Failed to update user!", error)
    return false
  }
}

// Thêm vào cuối file các hàm quản lý user

// User Management Functions
function editUser(username) {
  console.log("Editing user (admin side):", username)
  const modal = document.getElementById("userModal")
  const title = document.getElementById("userModalTitle")
  const passwordGroup = document.getElementById("passwordGroup")

  title.textContent = "Chỉnh sửa người dùng"
  passwordGroup.style.display = "none" // Ẩn password khi edit

  // Tìm user data từ loadedUsers
  const targetUser = loadedUsers.find((user) => user.username === username)
  if (targetUser) {
    document.getElementById("userUsername").value = username
    document.getElementById("userUsername").readOnly = true
    document.getElementById("userEmail").value = targetUser.email
    document.getElementById("userName").value = targetUser.name
    // Phone field was removed from HTML, so no need to populate it here.

    const roleText = targetUser.role.toLowerCase()
    document.getElementById("userRole").value = roleText.includes("admin") ? "admin" : "user"
  } else {
    console.error("User not found in loadedUsers for editing:", username)
    showNotification("Không tìm thấy người dùng để chỉnh sửa!", "error")
  }

  modal.style.display = "block"
}

function addUser() {
  console.log("Adding new user (admin side)...")
  const modal = document.getElementById("userModal")
  const title = document.getElementById("userModalTitle")
  const passwordGroup = document.getElementById("passwordGroup")

  title.textContent = "Thêm người dùng mới"
  passwordGroup.style.display = "block" // Hiện password khi add

  // Clear form
  document.getElementById("userForm").reset()
  document.getElementById("userUsername").readOnly = false

  modal.style.display = "block"
}

function closeUserModal() {
  console.log("Closing user modal.")
  document.getElementById("userModal").style.display = "none"
}

async function saveUser() {
  console.log("Saving user (admin side)...")
  const username = document.getElementById("userUsername").value
  const password = document.getElementById("userPassword").value
  const email = document.getElementById("userEmail").value
  const name = document.getElementById("userName").value
  const role = document.getElementById("userRole").value

  if (!username || !email || !name || !role) {
    showNotification("Vui lòng điền đầy đủ thông tin!", "error")
    console.warn("Validation failed: Missing user information.")
    return
  }

  const isEdit = document.getElementById("userUsername").readOnly

  if (isEdit) {
    // Update user (admin editing existing user)
    console.log("Admin editing existing user:", username)
    const success = await updateUser({
      username: username,
      email: email,
      name: name,
      // Do NOT send password, role, or phone to the /update-info API
    })

    if (success) {
      closeUserModal()
      await loadUsers() // Reload users table
    }
  } else {
    // Create new user (existing logic)
    console.log("Admin creating new user:", username)
    if (!password) {
      showNotification("Vui lòng nhập mật khẩu!", "error")
      console.warn("Validation failed: Password missing for new user.")
      return
    }

    const success = await createUser({
      username,
      password,
      email,
      name,
      role,
    })

    if (success) {
      closeUserModal()
      await loadUsers() // Reload users table
    }
  }
}

function changeUserPassword(username) {
  console.log("Opening change password modal for user:", username)
  const modal = document.getElementById("passwordModal")
  document.getElementById("changePasswordUsername").value = username
  document.getElementById("passwordForm").reset()
  document.getElementById("changePasswordUsername").value = username
  modal.style.display = "block"
}

function closePasswordModal() {
  console.log("Closing password modal.")
  document.getElementById("passwordModal").style.display = "none"
}

async function savePassword() {
  console.log("Saving new password...")
  const username = document.getElementById("changePasswordUsername").value
  const oldPassword = document.getElementById("oldPassword").value
  const newPassword = document.getElementById("newPassword").value
  const confirmPassword = document.getElementById("confirmPassword").value

  if (!oldPassword || !newPassword || !confirmPassword) {
    showNotification("Vui lòng điền đầy đủ thông tin!", "error")
    console.warn("Validation failed: Missing password fields.")
    return
  }

  if (newPassword !== confirmPassword) {
    showNotification("Mật khẩu mới không khớp!", "error")
    console.warn("Validation failed: New passwords do not match.")
    return
  }

  const success = await changePassword({
    username,
    oldPassword,
    newPassword,
  })

  if (success) {
    closePasswordModal()
  }
}

// Reload functions for each section
async function reloadOverview() {
  console.log("Reloading Overview section...")
  showLoadingState()
  await Promise.all([loadDashboardStats(), loadSystemLogs()])
  hideLoadingState()
  showNotification("Đã tải lại tổng quan!", "info")
}

async function reloadDevices() {
  console.log("Reloading Devices section...")
  showLoadingState()
  await loadDevices()
  hideLoadingState()
  showNotification("Đã tải lại danh sách thiết bị!", "info")
}

async function reloadUsers() {
  console.log("Reloading Users section...")
  showLoadingState()
  await loadUsers()
  hideLoadingState()
  showNotification("Đã tải lại danh sách người dùng!", "info")
}

// Removed reloadAlerts as per user request
async function reloadAlerts() {
  console.log("Reloading Alerts section (placeholder)...")
  showNotification("Tính năng Cảnh báo đang được phát triển!", "info")
}

// Removed reloadReports as per user request
async function reloadReports() {
  console.log("Reloading Reports section (placeholder)...")
  showNotification("Tính năng Báo cáo đang được phát triển!", "info")
}

async function reloadSettings() {
  console.log("Reloading Settings section (placeholder)...")
  // For now, just a notification as there's no dynamic data loading for settings
  showNotification("Đã tải lại cài đặt!", "info")
}

// Renamed and updated reload function for the combined user overview
async function reloadUserOverview() {
  console.log("Reloading User Overview section...")
  showLoadingState()
  await Promise.all([loadDevices(), loadUsers()]) // Ensure data is fresh
  await loadUserOverviewData() // Load all combined data
  hideLoadingState()
  showNotification("Đã tải lại tổng quan cá nhân!", "info")
}

// New function to populate the device filter dropdown
function populateDeviceFilterDropdown() {
  console.log("Populating device filter dropdown...")
  const deviceFilterSelect = document.getElementById("deviceFilter")
  if (!deviceFilterSelect) {
    console.warn("Device filter select element not found.")
    return
  }

  // Clear existing options
  deviceFilterSelect.innerHTML = '<option value="all">Tất cả Agent</option>' // Default option

  // Add options for each unique agentID from loadedClients
  const uniqueAgentIds = new Set()
  loadedClients.forEach((client) => {
    if (client.agentID) {
      uniqueAgentIds.add(client.agentID)
    }
  })

  Array.from(uniqueAgentIds)
    .sort()
    .forEach((agentId) => {
      const option = document.createElement("option")
      option.value = agentId
      option.textContent = agentId
      deviceFilterSelect.appendChild(option)
    })
  console.log("Device filter dropdown populated.")
}

// Cập nhật hàm setupAdminDashboard để thêm event listeners
async function setupAdminDashboard() {
  console.log("Setting up Admin Dashboard...")
  try {
    // Show loading state
    showLoadingState()

    // Load dashboard data
    await Promise.all([
      loadDashboardStats(),
      loadUsers(), // Load users first so they are available for device assignment
      loadDevices(), // This populates loadedClients
      // Removed loadAlerts() as per user request
    ])

    // After devices are loaded, populate the filter dropdown
    populateDeviceFilterDropdown()
    // Now load system logs with the default filter (which is "all")
    await loadSystemLogs() // Call after dropdown is populated

    // Hide loading state
    hideLoadingState()
    console.log("Admin Dashboard data loaded.")

    // Setup event handlers
    setupDeviceFilter()
    // Removed setupRefreshLogs as per user request
    setupUserAssignments()
    // setupRealTimeUpdates(); // This function is commented out in the provided script.js

    // Setup user management
    const addUserBtn = document.getElementById("addUserBtn")
    if (addUserBtn) {
      addUserBtn.addEventListener("click", addUser)
      console.log("Add User button event listener attached.")
    }
    // Setup add device button
    const addDeviceBtn = document.getElementById("addDeviceBtn")
    if (addDeviceBtn) {
      addDeviceBtn.addEventListener("click", () => {
        console.log("Add Device button clicked.")
        document.getElementById("deviceId").value = ""
        document.getElementById("deviceIp").value = ""
        document.getElementById("deviceHardwareId").value = ""
        document.getElementById("deviceAgentId").value = ""
        document.getElementById("deviceUser").value = ""
        document.getElementById("modalTitle").textContent = "Thêm thiết bị mới"
        document.getElementById("deviceId").readOnly = false
        document.getElementById("deviceAgentId").readOnly = false
        document.getElementById("deviceModal").style.display = "block"
      })
      console.log("Add Device button event listener attached.")
    }

    // Attach reload button listeners for Admin Dashboard
    document.getElementById("reloadOverviewBtn")?.addEventListener("click", () => {
      console.log("Reload Overview button clicked!")
      reloadOverview()
    })
    document.getElementById("reloadDevicesBtn")?.addEventListener("click", () => {
      console.log("Reload Devices button clicked!")
      reloadDevices()
    })
    document.getElementById("reloadUsersBtn")?.addEventListener("click", () => {
      console.log("Reload Users button clicked!")
      reloadUsers()
    })
    document.getElementById("reloadAlertsBtn")?.addEventListener("click", () => {
      console.log("Reload Alerts button clicked!")
      reloadAlerts()
    })
    document.getElementById("reloadReportsBtn")?.addEventListener("click", () => {
      console.log("Reload Reports button clicked!")
      reloadReports()
    })
    document.getElementById("reloadSettingsBtn")?.addEventListener("click", () => {
      console.log("Reload Settings button clicked!")
      reloadSettings()
    })
    console.log("Admin Dashboard reload button listeners attached.")

    // Attach reload button listeners for User Dashboard (if applicable)
    document.getElementById("reloadUserOverviewBtn")?.addEventListener("click", () => {
      console.log("Reload User Overview button clicked!")
      reloadUserOverview()
    })
    console.log("User Dashboard reload button listeners attached (if elements exist).")

    // Attach event listener for device filter dropdown
    document.getElementById("deviceFilter")?.addEventListener("change", function () {
      console.log("Device filter changed to:", this.value)
      loadSystemLogs(this.value)
    })
    console.log("Device filter dropdown event listener attached.")

    // Close modal when clicking outside
    window.addEventListener("click", (event) => {
      const deviceModal = document.getElementById("deviceModal")
      const userModal = document.getElementById("userModal")
      const passwordModal = document.getElementById("passwordModal")

      if (event.target === deviceModal) {
        closeDeviceModal()
      }
      if (event.target === userModal) {
        closeUserModal()
      }
      if (event.target === passwordModal) {
        closePasswordModal()
      }
    })
    console.log("Global modal close listener attached.")
  } catch (error) {
    console.error("Failed to load admin dashboard during setup:", error)
    showNotification("Không thể tải dữ liệu dashboard admin", "error")
  }
}

// New function for user dashboard specific data loading
async function setupUserDashboard() {
  console.log("Setting up User Dashboard...")
  try {
    showLoadingState()
    // Ensure clients and users are loaded for user dashboard
    await Promise.all([loadDevices(), loadUsers()])
    await loadUserOverviewData() // Populate user profile and device info, and logs
    hideLoadingState()
    console.log("User Dashboard data loaded.")

    // Attach event listener for user profile update button
    document.getElementById("updateUserProfileBtn")?.addEventListener("click", saveUserProfile)
    console.log("User Profile update button event listener attached.")
  } catch (error) {
    console.error("Failed to load user dashboard during setup:", error)
    showNotification("Không thể tải dữ liệu dashboard người dùng", "error")
  }
}

// Function to populate user profile data and device details, and logs
async function loadUserOverviewData() {
  console.log("Loading user overview data (profile, device, logs)...")
  const userProfileUsername = document.getElementById("userProfileUsername")
  const userProfileEmail = document.getElementById("userProfileEmail")
  const userProfileName = document.getElementById("userProfileName")

  const userDeviceInfoContent = document.getElementById("userDeviceInfoContent")
  const userDeviceInfoEmpty = document.getElementById("userDeviceInfoEmpty")
  const userGetOtpBtn = document.getElementById("userGetOtpBtn")
  const userDeviceLogsTableBody = document.getElementById("userDeviceLogsTable")
  const userDeviceLogsEmptyState = document.getElementById("userDeviceLogsEmpty")

  // Clear previous data and hide empty states
  userProfileUsername.value = ""
  userProfileEmail.value = ""
  userProfileName.value = ""
  userDeviceInfoContent.style.display = "none"
  userDeviceInfoEmpty.style.display = "none"
  userDeviceLogsTableBody.innerHTML = ""
  userDeviceLogsEmptyState.style.display = "none"

  if (currentUser) {
    // Populate user profile data
    userProfileUsername.value = currentUser.username || "N/A"
    userProfileEmail.value = currentUser.email || "N/A"
    userProfileName.value = currentUser.fullName || "N/A"
    console.log("User profile data populated for:", currentUser.username)

    // Populate user assigned device details and logs
    const assignedDevice = loadedClients.find((client) => client.username === currentUser.username)

    if (assignedDevice) {
      document.getElementById("userDeviceInfoAgentId").value = assignedDevice.agentID || "N/A"
      document.getElementById("userDeviceInfoName").value = assignedDevice.hardwareInfo.hostName || "N/A"
      document.getElementById("userDeviceInfoIp").value = assignedDevice.hardwareInfo.ipAddress || "N/A"
      document.getElementById("userDeviceInfoHardwareId").value = assignedDevice.hardwareInfo.hostID || "N/A"
      document.getElementById("userDeviceInfoStatus").value = assignedDevice.isOnline ? "Online" : "Offline"

      userGetOtpBtn.onclick = () => getOTP(assignedDevice.agentID)

      userDeviceInfoContent.style.display = "block" // Show device info

      // Tải và hiển thị log trực tiếp
      try {
        const logs = await apiCall(API_CONFIG.endpoints.myDeviceLogs) // Gọi API log của người dùng
        if (logs && logs.length > 0) {
          logs.forEach((log) => {
            const row = document.createElement("tr")
            const deviceId = log.agentID || "N/A"
            const level = log.level ? log.level.toLowerCase() : "info"
            const event = log.message || "Log Entry"
            const details = log.details || ""
            const timestamp = log.timestamp || new Date().toLocaleString("vi-VN")

            row.innerHTML = `
                                <td>${timestamp}</td>
                                <td>${deviceId}</td>
                                <td><span class="log-level ${level}">${level.toUpperCase()}</span></td>
                                <td>${event}</td>
                                <td>${details}</td>
                            `
            userDeviceLogsTableBody.appendChild(row)
          })
          console.log(`Loaded ${logs.length} logs for device ${assignedDevice.agentID}.`)
        } else {
          userDeviceLogsEmptyState.style.display = "block" // Hiển thị empty state nếu không có log
          console.log(`No logs found for device ${assignedDevice.agentID}.`)
        }
      } catch (error) {
        console.error("Failed to load user device logs:", error)
        userDeviceLogsEmptyState.style.display = "block"
      }
    } else {
      userDeviceInfoEmpty.style.display = "block" // Show empty state if no device assigned
      console.log("No device assigned to user:", currentUser.username)
    }
  } else {
    console.warn("No current user found. Cannot load user overview data.")
  }
}

// Function to handle saving user profile data
async function saveUserProfile() {
  console.log("Saving user profile...")
  const username = document.getElementById("userProfileUsername").value
  const email = document.getElementById("userProfileEmail").value
  const name = document.getElementById("userProfileName").value

  // Validate input
  if (!username || !email || !name) {
    showNotification("Vui lòng điền đầy đủ thông tin!", "error")
    console.warn("Validation failed: Missing user profile information.")
    return
  }

  // Call the updateUser function
  const success = await updateUser({
    username: username,
    email: email,
    name: name,
  })

  if (success) {
    // Update currentUser object in localStorage
    currentUser.email = email
    currentUser.fullName = name
    localStorage.setItem("currentUser", JSON.stringify(currentUser))
    showNotification("Thông tin cá nhân đã được cập nhật!", "success")
    console.log("User profile updated successfully.")
  } else {
    showNotification("Lỗi khi cập nhật thông tin cá nhân!", "error")
    console.error("Failed to update user profile.")
  }
}
