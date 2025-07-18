// --- Cấu hình base URL và hàm fetchWithAuth ---
// const API_BASE_URL = 'http://192.168.15.12:8082/api';
// const API_BASE_URL = window.location.origin + '/api';
const API_BASE_URL = window.location.protocol + '//' + window.location.hostname + ':8082/api';

function apiUrl(path) {
    if (path.startsWith('/')) path = path.slice(1);
    return `${API_BASE_URL}/${path}`;
}
function getToken() {
    return localStorage.getItem('jwt_token');
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
        localStorage.clear();
        window.location.href = 'login.html';
        return new Response(null, { status: 401 });
    }
    return res;
}

// --- Lấy thông tin thiết bị được gán cho user ---
async function loadUserDevicesAndLogs() {
    const deviceTableBody = document.querySelector('.device-table tbody');
    const agentFilter = document.getElementById('agent-filter');
    // Tìm dòng template (ẩn) nếu có, hoặc tạo dòng mẫu
    let templateRow = deviceTableBody.querySelector('tr[data-template]');
    if (!templateRow) {
        // Nếu chưa có, tạo dòng mẫu ẩn (chỉ dùng lần đầu)
        templateRow = document.createElement('tr');
        templateRow.setAttribute('data-template', '');
        templateRow.style.display = 'none';
        templateRow.innerHTML = `
            <td class="agent-id"></td>
            <td class="device-info"></td>
            <td class="device-status"></td>
            <td class="action-cell"><button class="otp-btn">Lấy OTP</button></td>
        `;
        deviceTableBody.appendChild(templateRow);
    }
    // Xóa tất cả các dòng không phải template
    Array.from(deviceTableBody.querySelectorAll('tr:not([data-template])')).forEach(tr => tr.remove());

    try {
        const res = await fetchWithAuth('/clients/my');
        if (!res.ok) throw new Error();
        const data = await res.json();
        const devices = data.data || [];
        if (devices.length === 0) {
            const noDevicesRow = document.createElement('tr');
            noDevicesRow.innerHTML = `<td colspan="4" class="no-devices">Không có thiết bị nào được gán.</td>`;
            deviceTableBody.appendChild(noDevicesRow);
            if (agentFilter) agentFilter.innerHTML = '';
            return;
        }
        // Đổ dữ liệu vào bảng bằng cách clone dòng template
        devices.forEach(d => {
            const row = templateRow.cloneNode(true);
            row.removeAttribute('data-template');
            row.style.display = '';
            row.querySelector('.agent-id').innerHTML = `<span class="device-id-chip">${d.agent_id}</span>`;
            row.querySelector('.device-info').innerHTML =
                `<span class="device-info-chip"><b>Tên:</b> ${d.device_info.hostName}</span>` +
                `<span class="device-info-chip"><b>IP:</b> ${d.device_info.ipAddress}</span>` +
                `<span class="device-info-chip"><b>MAC:</b> ${d.device_info.macAddress}</span>` +
                `<span class="device-info-chip"><b>Hardware ID:</b> ${d.device_info.hardwareID}</span>`;
            row.querySelector('.device-status').innerHTML =
                `<span class="device-status-chip ${d.online ? 'online' : 'offline'}">${d.online ? 'Online' : 'Offline'}</span>` +
                `<span class="device-info-chip"><div style="margin-top:4px;"><b>Last seen:</b> ${d.last_seen}</div></span>`;
            // Gán sự kiện cho nút OTP nếu cần
            // row.querySelector('.otp-btn').onclick = ...
            deviceTableBody.appendChild(row);
        });
        // Đổ agent filter
        if (agentFilter) {
            agentFilter.innerHTML = devices.map(d => `<option value="${d.agent_id}">${d.agent_id} - ${d.device_info.hostName}</option>`).join('');
            agentFilter.onchange = function () {
                loadUserLogs(this.value);
            };
        }
        // Lấy log của thiết bị đầu tiên (truyền agent)
        loadUserLogs(devices[0].agent_id);
        if (agentFilter) agentFilter.value = devices[0].agent_id;
    } catch (error) {
        console.error('Error loading devices:', error);
    }
}

// --- Lấy log thiết bị của user (API mới, cần agent) ---
async function loadUserLogs(agent, page = 1, pageSize = 8) {
    const tbody = document.querySelector('.log-table tbody');
    const pagination = document.getElementById('log-pagination');
    const logTotal = document.getElementById('log-total');
    if (!agent) {
        tbody.innerHTML = `<tr><td colspan="5" style="text-align:center;color:#888;">Không có dữ liệu log thiết bị.</td></tr>`;
        if (pagination) pagination.innerHTML = '';
        if (logTotal) logTotal.textContent = '';
        return;
    }
    try {
        const res = await fetchWithAuth(`/logs/my-device-paged?agent=${encodeURIComponent(agent)}&page=${page}&pageSize=${pageSize}`);
        if (!res.ok) throw new Error();
        const data = await res.json();
        const logs = data.logs || [];
        const total = data.total || 0;
        tbody.innerHTML = '';
        // Hiển thị tổng số log
        if (logTotal) {
            logTotal.textContent = `Tổng số: ${total} log`;
        }
        if (!logs.length) {
            tbody.innerHTML = `<tr><td colspan="5" style="text-align:center;color:#888;">Không có dữ liệu log thiết bị.</td></tr>`;
            if (pagination) pagination.innerHTML = '';
            if (logTotal) logTotal.textContent = '';
            return;
        }
        for (const log of logs) {
            tbody.innerHTML += `
                <tr>
                    <td>${log.time || ''}</td>
                    <td>${log.agent_id || ''}</td>
                    <td>${log.level || ''}</td>
                    <td>${log.message || ''}</td>
                </tr>
            `;
        }
        // Phân trang
        if (pagination) {
            const totalPages = Math.ceil(total / pageSize);
            let html = '';
            if (totalPages > 1) {
                // <<
                html += `<button class="log-page-btn" ${page === 1 ? 'disabled' : ''} data-page="1">&laquo;</button>`;

                // Trang 1
                html += `<button class="log-page-btn${page === 1 ? ' active' : ''}" data-page="1">1</button>`;

                // Dấu ...
                if (page > 3) {
                    html += `<span style="margin:0 2px;">...</span>`;
                }

                // Các trang lân cận (trừ 1 và cuối)
                for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) {
                    html += `<button class="log-page-btn${i === page ? ' active' : ''}" data-page="${i}">${i}</button>`;
                }

                // Dấu ...
                if (page < totalPages - 2) {
                    html += `<span style="margin:0 2px;">...</span>`;
                }

                // Trang cuối (nếu > 1)
                if (totalPages > 1) {
                    html += `<button class="log-page-btn${page === totalPages ? ' active' : ''}" data-page="${totalPages}">${totalPages}</button>`;
                }

                // >>
                html += `<button class="log-page-btn" ${page === totalPages ? 'disabled' : ''} data-page="${totalPages}">&raquo;</button>`;
            }
            pagination.innerHTML = html;
            Array.from(pagination.querySelectorAll('button[data-page]')).forEach(btn => {
                btn.onclick = function () {
                    const newPage = parseInt(this.getAttribute('data-page'));
                    if (!isNaN(newPage) && newPage !== page) {
                        loadUserLogs(agent, newPage, pageSize);
                    }
                };
            });
        }
    } catch (e) {
        tbody.innerHTML = `<tr><td colspan=\"5\" style=\"text-align:center;color:#888;\">Không thể tải dữ liệu log thiết bị.</td></tr>`;
        if (pagination) pagination.innerHTML = '';
        if (logTotal) logTotal.textContent = '';
    }
}

document.addEventListener('DOMContentLoaded', function () {
    const user = JSON.parse(localStorage.getItem('user_info') || '{}');
    // document.getElementById('username').value = user.username || '';
    document.getElementById('email').value = user.email || '';
    document.getElementById('fullname').value = user.full_name || '';
    loadUserDevicesAndLogs();

    // Gán sự kiện cho nút reload mới
    const reloadBtn = document.getElementById('reload-btn');
    if (reloadBtn) {
        reloadBtn.onclick = function () {
            window.location.reload();
        };
    }

    // Xử lý modal đổi mật khẩu
    const changePasswordBtn = document.getElementById('change-password-btn');
    const changePasswordModal = document.getElementById('change-password-modal');
    const closeModalBtn = document.getElementById('close-modal-btn');
    if (changePasswordBtn && changePasswordModal) {
        changePasswordBtn.onclick = function () {
            changePasswordModal.style.display = 'flex';
        };
    }
    if (closeModalBtn && changePasswordModal) {
        closeModalBtn.onclick = function () {
            changePasswordModal.style.display = 'none';
        };
    }
    // Đóng modal khi click ra ngoài
    if (changePasswordModal) {
        changePasswordModal.addEventListener('click', function (e) {
            if (e.target === changePasswordModal) {
                changePasswordModal.style.display = 'none';
            }
        });
    }

    // Xử lý cập nhật thông tin user
    const updateBtn = document.querySelector('.profile-form button[type="button"]');
    if (updateBtn) {
        updateBtn.onclick = async function () {
            const fullname = document.getElementById('fullname').value.trim();
            const email = document.getElementById('email').value.trim();
            const user = JSON.parse(localStorage.getItem('user_info') || '{}');
            if (!fullname || !email) {
                alert('Vui lòng nhập đầy đủ họ tên và email.');
                return;
            }
            try {
                const res = await fetchWithAuth('/users/update-info', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        username: user.username,
                        full_name: fullname,
                        email: email
                    })
                });
                if (!res.ok) throw new Error();
                alert('Cập nhật thông tin thành công!');
                // Cập nhật lại localStorage
                user.full_name = fullname;
                user.email = email;
                localStorage.setItem('user_info', JSON.stringify(user));
            } catch (e) {
                alert('Cập nhật thông tin thất bại!');
            }
        };
    }

    // Xử lý đổi mật khẩu
    const changePasswordForm = document.getElementById('change-password-form');
    if (changePasswordForm) {
        changePasswordForm.onsubmit = async function (e) {
            e.preventDefault();
            const oldPass = document.getElementById('old-password').value;
            const newPass = document.getElementById('new-password').value;
            const confirmPass = document.getElementById('confirm-password').value;
            const user = JSON.parse(localStorage.getItem('user_info') || '{}');
            if (!oldPass || !newPass || !confirmPass) {
                alert('Vui lòng nhập đầy đủ các trường!');
                return;
            }
            if (newPass !== confirmPass) {
                alert('Mật khẩu mới và xác nhận không khớp!');
                return;
            }
            try {
                const res = await fetchWithAuth('/users/change-password', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        username: user.username,
                        old_password: oldPass,
                        new_password: newPass
                    })
                });
                if (!res.ok) throw new Error();
                alert('Đổi mật khẩu thành công!');
                document.getElementById('change-password-modal').style.display = 'none';
                changePasswordForm.reset();
            } catch (e) {
                alert('Đổi mật khẩu thất bại! Kiểm tra lại mật khẩu hiện tại.');
            }
        };
    }
});