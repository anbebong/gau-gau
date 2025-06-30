// JS cho dashboard admin (có thể thêm logic động ở đây)

// Bảo vệ dashboard, lấy tổng quan user và thiết bị, render giao diện
function getToken() {
    return localStorage.getItem('jwt_token');
}
function getRole() {
    return localStorage.getItem('role');
}
function protectAdminDashboard() {
    const token = getToken();
    const role = getRole();
    if (!token || role !== 'admin') {
        window.location.href = 'login.html';
    }
}
// Cấu hình base URL API
const API_BASE_URL = 'http://localhost:8082/api';

function apiUrl(path) {
    if (path.startsWith('/')) path = path.slice(1);
    return `${API_BASE_URL}/${path}`;
}
async function fetchWithAuth(path, options = {}) {
    const token = getToken();
    const res = await fetch(apiUrl(path), {
        ...options,
        headers: {
            ...(options.headers || {}),
            'Authorization': 'Bearer ' + token
        }
    });
    if (res.status === 401) {
        // Token hết hạn hoặc không hợp lệ, tự động logout
        localStorage.clear();
        window.location.href = 'login.html';
        return new Response(null, { status: 401 });
    }
    return res;
}
async function loadOverview() {
    // Lấy tổng số user
    let userCount = 0;
    let clientCount = 0;
    let onlineCount = 0;
    let offlineCount = 0;
    try {
        const userRes = await fetchWithAuth('users');
        if (userRes.ok) {
            const userData = await userRes.json();
            userCount = Array.isArray(userData.data) ? userData.data.length : 0;
        }
        const clientRes = await fetchWithAuth('clients');
        if (clientRes.ok) {
            const clientData = await clientRes.json();
            const devices = Array.isArray(clientData.data) ? clientData.data : [];
            clientCount = devices.length;
            onlineCount = devices.filter(d => d.online).length;
            offlineCount = devices.filter(d => !d.online).length;
        }
    } catch (e) {
        // Có thể show lỗi nếu muốn
    }
    document.querySelector('.overview-card:nth-child(1) .count').innerText = clientCount;
    document.querySelector('.overview-card:nth-child(2) .count').innerText = userCount;
    document.getElementById('device-online').innerText = onlineCount + ' Online';
    document.getElementById('device-offline').innerText = offlineCount + ' Offline';
}
function setupLogout() {
    const btn = document.querySelector('.logout-btn');
    if (btn) {
        btn.addEventListener('click', function() {
            localStorage.clear();
            window.location.href = 'login.html';
        });
    }
}
async function loadLogs(agent = '') {
    let url = 'logs/archive';
    if (agent) url += `?agent=${encodeURIComponent(agent)}`;
    try {
        const res = await fetchWithAuth(url);
        if (res.ok) {
            const data = await res.json();
            const logs = Array.isArray(data.data) ? data.data : [];
            renderLogs(logs);
        } else {
            renderLogs([]);
        }
    } catch (e) {
        renderLogs([]);
    }
}
function extractLevelFromMessage(msg) {
    // Tìm [LEVEL] trong message
    const match = msg.match(/\[(.*?)\]/);
    return match ? match[1] : '';
}
function extractEventFromMessage(msg) {
    // Lấy phần sau dấu - (nếu có)
    const idx = msg.indexOf(' - ');
    return idx !== -1 ? msg.slice(idx + 3) : msg;
}
function renderLogs(logs) {
    const tbody = document.querySelector('.log-table tbody');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (!logs.length) {
        tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;color:#888;">Không có dữ liệu log.</td></tr>';
        return;
    }
    for (const log of logs) {
        const level = extractLevelFromMessage(log.message || '');
        const event = extractEventFromMessage(log.message || '');
        tbody.innerHTML += `
        <tr>
            <td>${log.time || ''}</td>
            <td>${log.agent_id || ''}</td>
            <td><span class="${level ? level.toLowerCase() : 'info'}">${level || ''}</span></td>
            <td>${event || ''}</td>
        </tr>`;
    }
}
function setupLogFilter() {
    const select = document.querySelector('.log-filter');
    if (select) {
        select.addEventListener('change', function() {
            const agent = select.value === 'Tất cả Agent' ? '' : select.value;
            loadLogs(agent);
        });
    }
}
async function loadAgentOptions() {
    // Lấy danh sách agent_id từ API clients
    try {
        const res = await fetchWithAuth('clients');
        if (res.ok) {
            const data = await res.json();
            const agents = Array.isArray(data.data) ? data.data.map(c => c.agent_id).filter(Boolean) : [];
            const select = document.querySelector('.log-filter');
            if (select) {
                // Xóa các option cũ (trừ option đầu)
                select.innerHTML = '<option value="">Tất cả Agent</option>';
                // Thêm option cho từng agent_id (không trùng lặp)
                [...new Set(agents)].forEach(agent => {
                    const opt = document.createElement('option');
                    opt.value = agent;
                    opt.textContent = agent;
                    select.appendChild(opt);
                });
            }
        }
    } catch (e) {}
}
function showSection(section) {
    document.querySelectorAll('.sidebar ul li a').forEach((a, idx) => {
        a.classList.toggle('active', a.dataset.section === section);
    });
    document.querySelectorAll('.main-section').forEach(content => {
        content.style.display = content.id === 'section-' + section ? '' : 'none';
    });
    if (section === 'devices') loadDevices();
}
function renderDeviceStatus(online) {
    return online ? '<span class="device-status online">ONLINE</span>' : '<span class="device-status offline">OFFLINE</span>';
}
function renderAssignUser(users, current) {
    let html = `<select class="assign-user">
        <option value="">Chưa gán</option>`;
    users.forEach(u => {
        html += `<option value="${u.username}"${u.username===current?' selected':''}>${u.username}</option>`;
    });
    html += '</select>';
    return html;
}
function showToast(message, type = 'success') {
    let toast = document.createElement('div');
    toast.className = 'custom-toast ' + type;
    toast.innerText = message;
    document.body.appendChild(toast);
    setTimeout(() => { toast.classList.add('show'); }, 10);
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 400);
    }, 2000);
}
async function assignUserToDevice(agent_id, username) {
    try {
        const res = await fetchWithAuth(`clients/assign-agentid`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ agent_id, username })
        });
        if (res.ok) {
            showToast('Gán user thành công!', 'success');
            loadDevices();
        } else {
            showToast('Gán user thất bại!', 'error');
        }
    } catch (e) {
        showToast('Lỗi kết nối khi gán user!', 'error');
    }
}
async function loadDevices() {
    const tbody = document.getElementById('device-tbody');
    if (!tbody) return;
    tbody.innerHTML = '<tr><td colspan="9" style="text-align:center;color:#888;">Đang tải dữ liệu...</td></tr>';
    try {
        const res = await fetchWithAuth('clients');
        const userRes = await fetchWithAuth('users');
        let users = [];
        if (userRes.ok) {
            const userData = await userRes.json();
            users = Array.isArray(userData.data) ? userData.data : [];
        }
        if (res.ok) {
            const data = await res.json();
            const devices = Array.isArray(data.data) ? data.data : [];
            if (!devices.length) {
                tbody.innerHTML = '<tr><td colspan="9" style="text-align:center;color:#888;">Không có thiết bị nào.</td></tr>';
                return;
            }
            tbody.innerHTML = '';
            devices.forEach(device => {
                const info = device.device_info || {};
                tbody.innerHTML += `
                <tr>
                    <td>${device.agent_id||''}</td>
                    <td>${info.hostName||''}</td>
                    <td>${info.ipAddress||''}</td>
                    <td>${info.hardwareID||''}</td>
                    <td>${renderAssignUser(users, device.username||'')}</td>
                    <td>${renderDeviceStatus(device.online)}</td>
                    <td>${device.last_seen||''}</td>
                    <td>
                        <button class="btn-otp" data-agent="${device.agent_id}">OTP</button>
                        <button class="btn-delete" data-agent="${device.agent_id}">Xóa</button>
                    </td>
                </tr>`;
            });
            // Gán sự kiện cho select assign-user
            tbody.querySelectorAll('.assign-user').forEach((select, idx) => {
                select.addEventListener('change', function() {
                    const tr = select.closest('tr');
                    const agent_id = tr ? tr.children[0].textContent : null;
                    const username = select.value;
                    if (agent_id !== null) {
                        assignUserToDevice(agent_id, username);
                    }
                });
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="9" style="text-align:center;color:#e74c3c;">Không thể tải thiết bị.</td></tr>';
        }
    } catch (e) {
        tbody.innerHTML = '<tr><td colspan="9" style="text-align:center;color:#e74c3c;">Lỗi kết nối server.</td></tr>';
    }
}
document.addEventListener('DOMContentLoaded', function() {
    protectAdminDashboard();
    loadOverview();
    setupLogout();
    loadAgentOptions();
    loadLogs();
    setupLogFilter();
    // Gắn data-section cho sidebar
    const sidebarLinks = document.querySelectorAll('.sidebar ul li a');
    const sectionMap = ['overview','devices','users','alerts','reports','settings'];
    sidebarLinks.forEach((a, idx) => {
        a.dataset.section = sectionMap[idx];
        a.addEventListener('click', function(e) {
            e.preventDefault();
            showSection(a.dataset.section);
        });
    });
    // Hiển thị mặc định section overview
    showSection('overview');
});
