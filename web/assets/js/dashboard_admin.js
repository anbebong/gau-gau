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
const API_BASE_URL = 'http://192.168.15.12:8082/api';

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
// ========== LOGS WITH PAGINATION ========== //
let currentLogPage = 1;
let logPageSize = 10;
let logTotal = 0;
let currentLogAgent = '';

async function loadLogsPaged(page = 1, pageSize = 5, agent = '') {
    // Nếu agent không truyền vào, dùng agent đang filter
    if (typeof agent === 'undefined') agent = currentLogAgent;
    else currentLogAgent = agent;
    let url = `logs/paged?page=${page}&pageSize=${pageSize}`;
    if (agent) url += `&agent=${encodeURIComponent(agent)}`;
    try {
        const res = await fetchWithAuth(url);
        if (res.ok) {
            const data = await res.json();
            renderLogsPaged(data.logs || []);
            renderLogPagination(page, pageSize, data.total || 0);
            logTotal = data.total || 0;
            currentLogPage = page;
        } else {
            renderLogsPaged([]);
            renderLogPagination(1, pageSize, 0);
        }
    } catch (e) {
        renderLogsPaged([]);
        renderLogPagination(1, pageSize, 0);
    }
}

function renderLogsPaged(logs) {
    const tbody = document.querySelector('.log-table tbody');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (!logs.length) {
        tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;color:#888;">Không có dữ liệu log.</td></tr>';
        return;
    }
    for (const log of logs) {
        tbody.innerHTML += `
        <tr>
            <td>${log.time || ''}</td>
            <td>${log.agent_id || ''}</td>
            <td></td>
            <td>${log.message || ''}</td>
        </tr>`;
    }
}

function renderLogPagination(page, pageSize, total) {
    const container = document.querySelector('.log-pagination');
    if (!container) return;
    const totalPages = Math.ceil(total / pageSize);
    if (totalPages <= 1) {
        container.innerHTML = '';
        return;
    }
    let html = `<button class="log-page-btn" data-page="${page-1}" ${page<=1?'disabled':''}>&laquo;</button>`;
    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || Math.abs(i-page)<=1) {
            html += `<button class="log-page-btn${i===page?' active':''}" data-page="${i}">${i}</button>`;
        } else if (i === page-2 || i === page+2) {
            html += `<span style="padding:0 4px;">...</span>`;
        }
    }
    html += `<button class="log-page-btn" data-page="${page+1}" ${page>=totalPages?'disabled':''}>&raquo;</button>`;
    html += `<span style="margin-left:12px;color:#888;">Tổng: ${total}</span>`;
    container.innerHTML = html;

    // Gán sự kiện chuyển trang, luôn truyền agent đang filter
    container.querySelectorAll('.log-page-btn').forEach(btn => {
        btn.onclick = function() {
            const p = parseInt(btn.getAttribute('data-page'));
            if (p >= 1 && p <= totalPages && p !== page) {
                loadLogsPaged(p, pageSize, currentLogAgent);
            }
        };
    });
}

// Sửa setupLogFilter để gọi loadLogsPaged
function setupLogFilter() {
    const select = document.querySelector('.log-filter');
    if (select) {
        select.addEventListener('change', function() {
            const agent = select.value === 'Tất cả Agent' ? '' : select.value;
            loadLogsPaged(1, logPageSize, agent);
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
    if (section === 'users') loadUsers();
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
function showOtpBox(otp, expireIn, agentId) {
    // Xóa box cũ nếu có
    document.querySelectorAll('.otp-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'otp-modal';
    modal.innerHTML = `
        <div class="otp-modal-content">
            <div class="otp-title">Mã OTP Thiết Bị <span class="otp-agent-id">(${agentId})</span></div>
            <div class="otp-code">${otp}</div>
            <div class="otp-expire">Hết hạn sau <span>${expireIn}</span> giây</div>
            <button class="otp-close">Đóng</button>
        </div>
    `;
    document.body.appendChild(modal);
    // Đếm ngược thời gian hết hạn
    let time = expireIn;
    const expireSpan = modal.querySelector('.otp-expire span');
    const timer = setInterval(() => {
        time--;
        if (time <= 0) {
            clearInterval(timer);
            modal.remove();
        } else {
            expireSpan.textContent = time;
        }
    }, 1000);
    // Đóng khi bấm nút
    modal.querySelector('.otp-close').onclick = () => {
        clearInterval(timer);
        modal.remove();
    };
    // Tự động ẩn sau 10s nếu chưa hết hạn
    setTimeout(() => { if (document.body.contains(modal)) modal.remove(); }, 10000);
}
function showDeleteConfirm(agentId, onConfirm) {
    // Xóa modal cũ nếu có
    document.querySelectorAll('.delete-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'delete-modal';
    modal.innerHTML = `
        <div class="delete-modal-content">
            <div class="delete-title">Xác nhận xóa thiết bị</div>
            <div class="delete-desc">Bạn có chắc chắn muốn xóa thiết bị với AGENT ID: <b>${agentId}</b>?</div>
            <div class="delete-actions">
                <button class="delete-cancel">Hủy</button>
                <button class="delete-confirm">Xóa</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    modal.querySelector('.delete-cancel').onclick = () => modal.remove();
    modal.querySelector('.delete-confirm').onclick = () => {
        modal.remove();
        onConfirm();
    };
}
// Thêm hàm render bảng user
function renderUserTable(users) {
    const tbody = document.getElementById('user-tbody');
    if (!tbody) return;
    if (!users.length) {
        tbody.innerHTML = '<tr><td colspan="6" style="text-align:center;color:#888;">Không có user nào.</td></tr>';
        return;
    }
    tbody.innerHTML = '';
    users.forEach(user => {
        tbody.innerHTML += `
        <tr>
            <td>${user.username || ''}</td>
            <td>${user.full_name || ''}</td>
            <td>${user.email || ''}</td>
            <td>${user.role || ''}</td>
            <td>${user.created_at ? user.created_at.split('T')[0] : ''}</td>
            <td>
                <button class="btn-edit" data-uid="${user.id}">Chỉnh sửa</button>
                <button class="btn-password" data-uid="${user.id}">Đổi MK</button>
                <button class="btn-delete-user" data-uid="${user.id}"${user.role==='admin'?' disabled':''}>Xóa</button>
            </td>
        </tr>`;
    });
}

// ========== USER MANAGEMENT ========== //
async function loadUsers() {
    const tbody = document.getElementById('user-tbody');
    if (!tbody) return;
    tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:#888;">Đang tải dữ liệu...</td></tr>';
    try {
        const res = await fetchWithAuth('users');
        if (res.ok) {
            const data = await res.json();
            const users = Array.isArray(data.data) ? data.data : [];
            if (!users.length) {
                tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:#888;">Không có user nào.</td></tr>';
                return;
            }
            tbody.innerHTML = '';
            users.forEach(user => {
                tbody.innerHTML += `
                <tr data-username="${user.username}">
                    <td>${user.username}</td>
                    <td>${user.full_name || ''}</td>
                    <td>${user.email || ''}</td>
                    <td>${user.role}</td>
                    <td>${user.created_at ? new Date(user.created_at).toLocaleString() : ''}</td>
                    <td>
                        <button class="btn-edit">Chỉnh sửa</button>
                        <button class="btn-password">Đổi MK</button>
                        <button class="btn-delete-user" ${user.role==='admin'?'disabled':''}>Xóa</button>
                    </td>
                </tr>`;
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:#e74c3c;">Không thể tải user.</td></tr>';
        }
    } catch (e) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:#e74c3c;">Lỗi kết nối server.</td></tr>';
    }
}

// Modal chỉnh sửa user
function showEditUserModal(user, onSave) {
    document.querySelectorAll('.edit-user-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'edit-user-modal';
    modal.innerHTML = `
        <div class="edit-user-modal-content">
            <div class="edit-title">Chỉnh sửa user: <b>${user.username}</b></div>
            <div class="edit-form">
                <label>Họ tên</label>
                <input type="text" class="edit-fullname" value="${user.full_name||''}">
                <label>Email</label>
                <input type="email" class="edit-email" value="${user.email||''}">
                <label>Role</label>
                <select class="edit-role">
                    <option value="user" ${user.role==='user'?'selected':''}>user</option>
                    <option value="admin" ${user.role==='admin'?'selected':''}>admin</option>
                </select>
            </div>
            <div class="edit-actions">
                <button class="edit-cancel">Hủy</button>
                <button class="edit-save">Lưu</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    modal.querySelector('.edit-cancel').onclick = () => modal.remove();
    modal.querySelector('.edit-save').onclick = () => {
        const full_name = modal.querySelector('.edit-fullname').value.trim();
        const email = modal.querySelector('.edit-email').value.trim();
        const role = modal.querySelector('.edit-role').value;
        onSave({ full_name, email, role });
        modal.remove();
    };
}

// Modal đổi mật khẩu user
function showChangePasswordModal(username, onSave) {
    document.querySelectorAll('.change-password-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'change-password-modal';
    modal.innerHTML = `
        <div class="change-password-modal-content">
            <div class="change-title">Đổi mật khẩu cho <b>${username}</b></div>
            <div class="change-form">
                <label>Mật khẩu mới</label>
                <input type="password" class="change-password" autocomplete="new-password">
            </div>
            <div class="change-actions">
                <button class="change-cancel">Hủy</button>
                <button class="change-save">Lưu</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    modal.querySelector('.change-cancel').onclick = () => modal.remove();
    modal.querySelector('.change-save').onclick = () => {
        const password = modal.querySelector('.change-password').value.trim();
        if (!password) {
            showToast('Vui lòng nhập mật khẩu mới!', 'error');
            return;
        }
        onSave(password);
        modal.remove();
    };
}

// Modal xác nhận xóa user
function showDeleteUserConfirm(username, onConfirm) {
    document.querySelectorAll('.delete-user-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'delete-user-modal';
    modal.innerHTML = `
        <div class="delete-user-modal-content">
            <div class="delete-title">Xác nhận xóa user</div>
            <div class="delete-desc">Bạn có chắc chắn muốn xóa user <b>${username}</b>?</div>
            <div class="delete-actions">
                <button class="delete-cancel">Hủy</button>
                <button class="delete-confirm">Xóa</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    modal.querySelector('.delete-cancel').onclick = () => modal.remove();
    modal.querySelector('.delete-confirm').onclick = () => {
        modal.remove();
        onConfirm();
    };
}

// Modal tạo mới user
function showAddUserModal(onSave) {
    document.querySelectorAll('.add-user-modal').forEach(e => e.remove());
    const modal = document.createElement('div');
    modal.className = 'add-user-modal';
    modal.innerHTML = `
        <div class="edit-user-modal-content">
            <div class="edit-title">Tạo mới user</div>
            <div class="edit-form">
                <label>Tài khoản (username)</label>
                <input type="text" class="add-username" autocomplete="off" required>
                <label>Họ tên</label>
                <input type="text" class="add-fullname" autocomplete="off" required>
                <label>Email</label>
                <input type="email" class="add-email" autocomplete="off">
                <label>Quyền</label>
                <select class="add-role">
                    <option value="user">user</option>
                    <option value="admin">admin</option>
                </select>
                <label>Mật khẩu</label>
                <input type="password" class="add-password" autocomplete="new-password" required>
            </div>
            <div class="edit-actions">
                <button class="edit-cancel">Hủy</button>
                <button class="edit-save">Tạo mới</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    modal.querySelector('.edit-cancel').onclick = () => modal.remove();
    modal.querySelector('.edit-save').onclick = () => {
        const username = modal.querySelector('.add-username').value.trim();
        const full_name = modal.querySelector('.add-fullname').value.trim();
        const email = modal.querySelector('.add-email').value.trim();
        const role = modal.querySelector('.add-role').value;
        const password = modal.querySelector('.add-password').value.trim();
        if (!username || !full_name || !password) {
            showToast('Vui lòng nhập đủ thông tin bắt buộc!', 'error');
            return;
        }
        onSave({ username, full_name, email, role, password });
        modal.remove();
    };
}

// Sự kiện thao tác user
function setupUserActions() {
    const tbody = document.getElementById('user-tbody');
    if (!tbody) return;
    tbody.addEventListener('click', async function(e) {
        const target = e.target;
        const tr = target.closest('tr');
        if (!tr) return;
        const username = tr.getAttribute('data-username');
        if (target.classList.contains('btn-edit')) {
            // Lấy dữ liệu user hiện tại
            const tds = tr.children;
            const user = {
                username,
                full_name: tds[1].textContent,
                email: tds[2].textContent,
                role: tds[3].textContent
            };
            showEditUserModal(user, async (data) => {
                try {
                    const res = await fetchWithAuth(`users/update`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username, ...data })
                    });
                    if (res.ok) {
                        showToast('Cập nhật user thành công!', 'success');
                        loadUsers();
                    } else {
                        showToast('Cập nhật user thất bại!', 'error');
                    }
                } catch (e) {
                    showToast('Lỗi kết nối khi cập nhật user!', 'error');
                }
            });
        } else if (target.classList.contains('btn-password')) {
            showChangePasswordModal(username, async (password) => {
                try {
                    const res = await fetchWithAuth(`users/change-password`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username, password })
                    });
                    if (res.ok) {
                        showToast('Đổi mật khẩu thành công!', 'success');
                    } else {
                        showToast('Đổi mật khẩu thất bại!', 'error');
                    }
                } catch (e) {
                    showToast('Lỗi kết nối khi đổi mật khẩu!', 'error');
                }
            });
        } else if (target.classList.contains('btn-delete-user')) {
            showDeleteUserConfirm(username, async () => {
                try {
                    const res = await fetchWithAuth(`users/delete`, {
                        method: 'DELETE',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username })
                    });
                    if (res.ok) {
                        showToast('Đã xóa user.', 'success');
                        loadUsers();
                    } else {
                        showToast('Xóa user thất bại!', 'error');
                    }
                } catch (e) {
                    showToast('Lỗi kết nối khi xóa user!', 'error');
                }
            });
        }
    });
}

document.addEventListener('DOMContentLoaded', function() {
    protectAdminDashboard();
    loadOverview();
    setupLogout();
    loadAgentOptions();
    loadLogsPaged(1, logPageSize);
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

    // Gán sự kiện cho các nút OTP và Xóa trong bảng thiết bị
    const deviceTbody = document.getElementById('device-tbody');
    if (deviceTbody) {
        deviceTbody.addEventListener('click', async function(e) {
            const target = e.target;
            const tr = target.closest('tr');
            if (!tr) return;
            const agent_id = tr.children[0]?.textContent;
            // Lấy client_id từ thuộc tính data-clientid nếu cần
            let client_id = tr.getAttribute('data-clientid');
            if (!client_id && tr.children[0]) client_id = tr.children[0].getAttribute('data-clientid');
            if (target.classList.contains('btn-otp') && agent_id) {
                // Gọi API lấy OTP
                try {
                    const res = await fetchWithAuth(`clients/${encodeURIComponent(agent_id)}/otp`);
                    if (res.ok) {
                        const data = await res.json();
                        if (data && data.data && data.data.otp) {
                            showOtpBox(data.data.otp, data.data.expire_in || 0, agent_id);
                        } else {
                            showToast('Lấy OTP thành công!', 'success');
                        }
                    } else {
                        showToast('Lấy OTP thất bại!', 'error');
                    }
                } catch (e) {
                    showToast('Lỗi kết nối khi lấy OTP!', 'error');
                }
            } else if (target.classList.contains('btn-delete')) {
                // Gọi API xóa thiết bị
                if (!agent_id) {
                    showToast('Không tìm thấy agent_id để xóa!', 'error');
                    return;
                }
                showDeleteConfirm(agent_id, async () => {
                    try {
                        const res = await fetchWithAuth('clients/delete-agentid', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ agent_id })
                        });
                        if (res.ok) {
                            showToast('Đã xóa thiết bị.', 'success');
                            loadDevices();
                        } else {
                            showToast('Xóa thiết bị thất bại!', 'error');
                        }
                    } catch (e) {
                        showToast('Lỗi kết nối khi xóa thiết bị!', 'error');
                    }
                });
            }
        });
    }
    // Khi chuyển tab users thì load lại bảng user
    const showSectionOrigin = showSection;
    window.showSection = function(section) {
        showSectionOrigin(section);
        if (section === 'users') {
            loadUsers();
        }
    };
    loadUsers();
    setupUserActions();
    // Sự kiện nút tạo mới user
    const btnAddUser = document.getElementById('btn-add-user');
    if (btnAddUser) {
        btnAddUser.onclick = function() {
            showAddUserModal(async (data) => {
                try {
                    const res = await fetchWithAuth('users/create', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(data)
                    });
                    if (res.ok) {
                        showToast('Tạo user thành công!', 'success');
                        loadUsers();
                    } else {
                        showToast('Tạo user thất bại!', 'error');
                    }
                } catch (e) {
                    showToast('Lỗi kết nối khi tạo user!', 'error');
                }
            });
        };
    }
});
