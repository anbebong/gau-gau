// Logic chuyển hướng sau đăng nhập
const API_BASE_URL = window.location.protocol + '//' + window.location.hostname + ':8082/api';
async function handleLogin(event) {
    event.preventDefault();
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value.trim();
    const errorDiv = document.getElementById('login-error');
    errorDiv.innerText = '';
    try {
        const res = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        if (!res.ok) {
            errorDiv.innerText = 'Sai tài khoản hoặc mật khẩu!';
            return;
        }
        const data = await res.json();
        // Sửa lại để lấy token và role từ data.data
        if (data.success && data.data && data.data.token) {
            localStorage.setItem('jwt_token', data.data.token);
            localStorage.setItem('username', data.data.user.username);
            localStorage.setItem('role', data.data.user.role);
            localStorage.setItem('user_info', JSON.stringify(data.data.user)); // Lưu thông tin user
            // Chuyển hướng theo role
            if (data.data.user.role === 'admin') {
                window.location.href = '../pages/dashboard_admin.html';
            } else {
                window.location.href = '../pages/dashboard_user.html';
            }
        } else {
            errorDiv.innerText = 'Đăng nhập thất bại!';
        }
    } catch (e) {
        errorDiv.innerText = 'Không thể kết nối máy chủ!';
    }
}
document.addEventListener('DOMContentLoaded', function() {
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');
    // Enter để submit, tự động focus
    if (usernameInput) usernameInput.focus();
    [usernameInput, passwordInput].forEach(input => {
        if (input) input.addEventListener('input', function() {
            document.getElementById('login-error').innerText = '';
        });
    });
    const form = document.getElementById('login-form');
    if (form) form.addEventListener('submit', handleLogin);
});
