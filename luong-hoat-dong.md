# Luồng Hoạt Động Toàn Diện của Credential Provider

Dưới đây là mô tả chi tiết từng bước về luồng hoạt động của dự án Credential Provider, từ lúc Windows khởi động màn hình đăng nhập cho đến khi người dùng được xác thực. Luồng này là sự phối hợp nhịp nhàng giữa **LogonUI (Hệ thống của Windows)** và các thành phần trong DLL của bạn (`Dll.cpp`, `CSampleProvider.cpp`, `CSampleCredential.cpp`, và `helpers.cpp`).

---

## Giai đoạn 1: Khởi tạo và Hiển thị Giao diện

*Mục tiêu: Windows tải Credential Provider của bạn và hiển thị ô đăng nhập (tile) ban đầu lên màn hình.*

1.  **LogonUI Bắt đầu**: Người dùng khởi động máy, khóa màn hình, hoặc chuyển đổi người dùng. Tiến trình `LogonUI.exe` của Windows được kích hoạt.

2.  **Tải DLL**: LogonUI quét Registry của Windows để tìm các Credential Provider đã được đăng ký. Nó tìm thấy CLSID (mã định danh duy nhất) của bạn và tải file DLL tương ứng vào bộ nhớ.

3.  **Yêu cầu "Nhà máy" (Class Factory)**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `DllGetClassObject` (trong `Dll.cpp`).
    *   **Mục đích**: LogonUI hỏi, "Này DLL, hãy cho tôi một đối tượng 'nhà máy' có khả năng tạo ra Credential Provider của anh."

4.  **Tạo "Nhà cung cấp" (Provider)**:
    *   **Ai gọi**: LogonUI, sử dụng "nhà máy" vừa nhận được.
    *   **Hàm nào được gọi**: `CClassFactory::CreateInstance` (trong `Dll.cpp`).
    *   **Kết quả**: Hàm này gọi `CSample_CreateInstance` (trong `CSampleProvider.cpp`) để tạo ra một đối tượng `CSampleProvider`. Đây chính là "nhà quản lý" chính.

5.  **Thiết lập Bối cảnh (Usage Scenario)**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `CSampleProvider::SetUsageScenario`.
    *   **Mục đích**: Windows thông báo cho Provider biết nó đang được dùng cho mục đích gì (ví dụ: `CPUS_LOGON` - đăng nhập, `CPUS_UNLOCK_WORKSTATION` - mở khóa).
    *   **Hành động**: Dựa vào bối cảnh này, `CSampleProvider` quyết định xem có nên hoạt động hay không. Nếu có, nó sẽ tạo ra một đối tượng `CSampleCredential` (đại diện cho ô đăng nhập) và lưu trữ nó.

6.  **Liệt kê các Ô đăng nhập (Enumeration)**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `CSampleProvider::GetCredentialCount` và `CSampleProvider::GetCredentialAt`.
    *   **Mục đích**: LogonUI hỏi, "Anh có bao nhiêu ô đăng nhập muốn hiển thị?" và "Hãy cho tôi ô đăng nhập ở vị trí số X."
    *   **Kết quả**: `CSampleProvider` trả về số lượng ô (ví dụ: 1) và con trỏ đến đối tượng `CSampleCredential` đã tạo ở bước 5.

7.  **Hiển thị Ô đăng nhập**:
    *   **Ai gọi**: LogonUI, sau khi đã có đối tượng `CSampleCredential`.
    *   **Hàm nào được gọi**: `CSampleCredential::GetUserSid`, `GetStringValue`, `GetBitmapValue`.
    *   **Mục đích**: LogonUI lấy các thông tin cần thiết để hiển thị ô đăng nhập: SID (để biết ô này của ai), tên người dùng, ảnh đại diện.
    *   **Kết quả**: Một ô đăng nhập với ảnh đại diện và tên người dùng được hiển thị trên màn hình, ở trạng thái chưa được chọn.

---

## Giai đoạn 2: Tương tác của Người dùng

*Mục tiêu: Xử lý các hành động của người dùng như nhấp chuột vào ô đăng nhập và nhập mật khẩu.*

8.  **Người dùng Chọn Ô đăng nhập**: Người dùng nhấp chuột vào ô đăng nhập của bạn.

9.  **Thông báo Ô được chọn**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `CSampleCredential::SetSelected(TRUE)`.
    *   **Hành động**: Bên trong hàm này, `CSampleCredential` gọi lại cho LogonUI thông qua `_pCredentialEvents->CredentialsChanged(...)`, báo rằng "Trạng thái của tôi đã thay đổi, hãy vẽ lại giao diện cho tôi."

10. **Cập nhật Giao diện Người dùng**:
    *   **Ai gọi**: LogonUI, sau khi nhận được thông báo thay đổi.
    *   **Hàm nào được gọi**: `CSampleCredential::GetFieldState` (được gọi cho TẤT CẢ các trường UI).
    *   **Mục đích**: LogonUI hỏi, "Trường UI (ô mật khẩu, nút bấm, link...) này nên ở trạng thái nào (hiển thị, ẩn, hay được focus)?"
    *   **Kết quả**: `CSampleCredential` trả về trạng thái mong muốn. Ví dụ, nó sẽ yêu cầu hiển thị ô mật khẩu (`CPFS_DISPLAY_IN_SELECTED_TILE`).

11. **Người dùng Nhập Mật khẩu**: Người dùng gõ từng ký tự vào ô mật khẩu.

12. **Thu thập Dữ liệu**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `CSampleCredential::SetStringValue`.
    *   **Mục đích**: Với mỗi thay đổi trong ô mật khẩu, LogonUI gọi hàm này để cung cấp chuỗi mật khẩu hiện tại.
    *   **Hành động**: `CSampleCredential` nhận chuỗi này và lưu vào biến thành viên của nó (ví dụ: `_szPassword`).

---

## Giai đoạn 3: Quá trình Xác thực

*Mục tiêu: Đóng gói thông tin đăng nhập, gửi cho hệ thống bảo mật của Windows và xử lý kết quả.*

13. **Người dùng Nhấn Đăng nhập**: Người dùng nhấn nút mũi tên (hoặc Enter) để bắt đầu xác thực.

14. **Yêu cầu Đóng gói Dữ liệu (Serialization)**:
    *   **Ai gọi**: LogonUI.
    *   **Hàm nào được gọi**: `CSampleCredential::GetSerialization`. **Đây là hàm cốt lõi của quá trình xác thực.**

15. **Chuẩn bị và Đóng gói Dữ liệu**:
    *   **Bên trong `GetSerialization`**:
        a.  Lấy tên người dùng và mật khẩu đã lưu (`_szUsername`, `_szPassword`).
        b.  Gọi `ProtectIfNecessaryAndCopyPassword` (từ `helpers.cpp`) để **mã hóa mật khẩu**.
        c.  Gọi `KerbInteractiveUnlockLogonInit` và `KerbInteractiveUnlockLogonPack` (từ `helpers.cpp`) để đóng gói tên người dùng và mật khẩu đã mã hóa vào một vùng nhớ đệm (buffer) duy nhất theo định dạng mà LSA yêu cầu.

16. **Gửi Dữ liệu cho LSA**:
    *   `GetSerialization` trả về buffer đã đóng gói cho LogonUI.
    *   LogonUI chuyển tiếp buffer này đến **LSA (Local Security Authority)** - bộ não bảo mật của Windows.

17. **Xác thực Thực sự**:
    *   **Ai thực hiện**: LSA.
    *   **Hành động**: LSA giải mã buffer, lấy ra thông tin đăng nhập và kiểm tra chúng với cơ sở dữ liệu tài khoản của hệ thống (có thể là tài khoản local hoặc Active Directory). **Đây là lúc mật khẩu đúng hay sai được quyết định.**

18. **Xử lý Kết quả**:
    *   LSA trả kết quả (thành công/thất bại) về cho LogonUI.
    *   Nếu thành công, LogonUI tiến hành đăng nhập cho người dùng.
    *   Nếu thất bại, LogonUI hiển thị thông báo lỗi.

19. **Dọn dẹp**:
    *   Sau khi quá trình đăng nhập hoàn tất, Windows sẽ gọi `Unadvise` và các hàm hủy (destructor) của `CSampleCredential` và `CSampleProvider` để giải phóng tài nguyên. Toàn bộ chu trình kết thúc.
