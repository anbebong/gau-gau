# Phân Tích Chi Tiết File `CSampleCredential.cpp`

## 1. Tổng Quan

File `CSampleCredential.cpp` là một trong những thành phần cốt lõi của dự án Credential Provider. Nó định nghĩa lớp `CSampleCredential`, và mỗi đối tượng của lớp này đại diện cho một **"ô đăng nhập" (credential tile)** cụ thể trên màn hình đăng nhập của Windows.

Nếu `CSampleProvider` là "nhà quản lý" quyết định có bao nhiêu ô đăng nhập và cho những ai, thì `CSampleCredential` chính là **bản thân ô đăng nhập đó**. Nó chịu trách nhiệm cho mọi thứ liên quan đến một ô cụ thể:
*   Hiển thị các trường nhập liệu (như ô mật khẩu, checkbox, link).
*   Thu thập thông tin khi người dùng tương tác với các trường đó.
*   Đóng gói thông tin đã thu thập để gửi cho Windows khi người dùng nhấn đăng nhập.

## 2. Các Giao Diện (Interfaces) Được Triển Khai

Lớp `CSampleCredential` triển khai các giao diện COM quan trọng để có thể giao tiếp với hệ thống LogonUI của Windows:

*   **`ICredentialProviderCredential`**: Giao diện cơ bản và bắt buộc. Nó định nghĩa các phương thức để quản lý trạng thái, giao diện, và quá trình đóng gói thông tin xác thực của một ô đăng nhập.
*   **`ICredentialProviderCredential2`**: Một phiên bản mở rộng của giao diện trên. Nó bổ sung phương thức `GetUserSid`, một phương thức cực kỳ quan trọng để liên kết ô đăng nhập này với một tài khoản người dùng Windows cụ thể thông qua **SID (Security Identifier)** của họ.

## 3. Phân Tích Các Phương Thức Quan Trọng

Dưới đây là giải thích chi tiết về các phương thức chính trong `CSampleCredential.cpp`.

### Phương thức Khởi tạo và Hủy
*   **`CSampleCredential::CSampleCredential()` (Constructor)**: Được gọi khi `CSampleProvider` tạo ra một ô đăng nhập mới. Nhiệm vụ chính là khởi tạo các biến thành viên về trạng thái ban đầu (ví dụ: đặt các con trỏ là `NULL`, các chuỗi là rỗng).
*   **`CSampleCredential::~CSampleCredential()` (Destructor)**: Được gọi khi ô đăng nhập bị hủy. Nhiệm vụ là giải phóng tất cả các tài nguyên đã cấp phát, như bộ nhớ cho các chuỗi, các đối tượng COM, v.v., để tránh rò rỉ bộ nhớ.

### Quản lý Giao tiếp và Trạng thái
*   **`Advise()`**: Windows gọi phương thức này để cung cấp cho ô đăng nhập một "con trỏ gọi lại" (`ICredentialProviderCredentialEvents`). Ô đăng nhập sẽ sử dụng con trỏ này để thông báo cho LogonUI về bất kỳ thay đổi nào (ví dụ: "dữ liệu đã thay đổi, hãy cập nhật giao diện").
*   **`Unadvise()`**: Được gọi khi giao tiếp kết thúc. Ô đăng nhập phải giải phóng con trỏ gọi lại đã nhận trong `Advise()`.
*   **`SetSelected(BOOL* pbAutoLogon)`**: Được gọi khi người dùng nhấp vào hoặc rời khỏi ô đăng nhập này. Nếu ô được chọn (`isSelected` là `true`), đây là lúc để hiển thị các trường nhập liệu.
*   **`GetUserSid(PWSTR* ppszSid)`**: (Từ `ICredentialProviderCredential2`) Windows gọi phương thức này để hỏi "Ô đăng nhập này thuộc về người dùng nào?". Phương thức phải trả về chuỗi SID của người dùng tương ứng. Đây là cách Windows biết rằng khi đăng nhập thành công, nó sẽ đăng nhập vào tài khoản đó.

### Quản lý Giao diện Người dùng (UI)
*   **`GetFieldState(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs, CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis)`**: Một phương thức rất quan trọng, được gọi liên tục. Windows hỏi "Trường giao diện có ID là `dwFieldID` nên ở trạng thái nào?". Bạn trả lời bằng cách thay đổi `pcpfs` (ví dụ: `CPFS_DISPLAY_IN_SELECTED_TILE` - hiển thị khi ô được chọn) và `pcpfis` (ví dụ: `CPFIS_FOCUSED` - trường này nên được focus).
*   **`GetStringValue(DWORD dwFieldID, PWSTR* ppsz)`**: Khi Windows cần hiển thị một chuỗi văn bản cho một trường (ví dụ: tiêu đề của một link), nó sẽ gọi hàm này. Bạn phải trả về chuỗi tương ứng với `dwFieldID`.
*   **`GetBitmapValue(DWORD dwFieldID, HBITMAP* phbmp)`**: Tương tự như `GetStringValue` nhưng dành cho hình ảnh (bitmap). Ví dụ, để cung cấp ảnh đại diện cho người dùng.

### Thu Thập Dữ Liệu Từ Người Dùng
*   **`SetStringValue(DWORD dwFieldID, PCWSTR psz)`**: Được gọi khi người dùng **gõ chữ** vào một trường văn bản (`CPFT_EDIT_TEXT` hoặc `CPFT_PASSWORD_TEXT`). `dwFieldID` cho biết người dùng đang gõ vào ô nào, và `psz` là nội dung họ đã gõ. Đây là nơi bạn **lưu trữ** mật khẩu hoặc tên người dùng vào các biến thành viên (ví dụ: `_szPassword`).
*   **`SetCheckboxValue(DWORD dwFieldID, BOOL bChecked)`**: Được gọi khi người dùng **tick hoặc bỏ tick** một checkbox. `bChecked` cho biết trạng thái mới của checkbox.

### Lõi Logic: Đóng Gói và Xác Thực
*   **`GetSerialization(...)`**: **Đây là phương thức quan trọng nhất về mặt bảo mật.** Nó được gọi khi người dùng nhấn nút "Đăng nhập".
    1.  **Mục đích:** Đóng gói tất cả thông tin đã thu thập (tên người dùng, mật khẩu, v.v.) thành một định dạng nhị phân (binary) mà hệ thống bảo mật của Windows (LSA) có thể hiểu.
    2.  **Cách hoạt động (trong dự án mẫu này):**
        *   Nó lấy tên người dùng (`_szUsername`) và mật khẩu (`_szPassword`) đã được lưu từ các lần gọi `SetStringValue` trước đó.
        *   Nó gọi hàm API của Windows là **`KerbPackAuthenticationBuffer`**.
        *   Hàm này thực hiện công việc "pass-through" (chuyển tiếp). Nó chỉ đơn thuần **đóng gói** tên người dùng và mật khẩu, chứ **không hề kiểm tra** chúng có đúng hay không.
        *   Kết quả đã đóng gói được trả về cho Windows.
    3.  **Điểm can thiệp:** Đây chính là nơi bạn sẽ thêm logic xác thực của riêng mình (ví dụ: kiểm tra mã OTP, gọi API web, v.v.) **trước khi** gọi `KerbPackAuthenticationBuffer`.

## 4. Luồng Hoạt Động Điển Hình

1.  Windows khởi động, `CSampleProvider` được tải.
2.  `CSampleProvider` tạo một đối tượng `CSampleCredential` cho một người dùng cụ thể.
3.  Windows gọi `GetUserSid()` để biết ô này thuộc về ai.
4.  Người dùng nhấp vào ô. Windows gọi `SetSelected(true)`.
5.  Bên trong `SetSelected`, ô đăng nhập yêu cầu Windows cập nhật giao diện.
6.  Windows gọi `GetFieldState()` cho từng trường để biết nên hiển thị/ẩn/focus trường nào.
7.  Người dùng gõ mật khẩu. Windows liên tục gọi `SetStringValue()`. Code của bạn lưu lại mật khẩu vào biến `_szPassword`.
8.  Người dùng nhấn nút đăng nhập. Windows gọi `GetSerialization()`.
9.  `GetSerialization()` gọi `KerbPackAuthenticationBuffer` để đóng gói tên người dùng và mật khẩu đã lưu.
10. Gói dữ liệu được trả về cho Windows. Windows tự mình thực hiện việc xác thực cuối cùng.

---

## Phân tích `CSampleProvider.cpp`

`CSampleProvider.cpp` đóng vai trò là "nhà quản lý" hoặc "nhà cung cấp" các credential. Nó không xử lý việc xác thực cụ thể của một người dùng, mà quản lý việc *khi nào* và *làm thế nào* các ô đăng nhập (credential tiles) được hiển thị trên màn hình đăng nhập.

### Vai trò chính:

1.  **Đăng ký và Khởi tạo:**
    *   Khi Logon UI hoặc Credential UI khởi động, nó sẽ tạo một instance của `CSampleProvider`.
    *   `CSampleProvider` chịu trách nhiệm khởi tạo các tài nguyên cần thiết.

2.  **Quản lý các loại Credential (Credential Types):**
    *   Nó định nghĩa các loại phương thức xác thực mà nó hỗ trợ. Trong ví dụ này, nó chỉ hỗ trợ một loại là xác thực bằng mật khẩu (`CPT_PASSWORD`).
    *   Hàm `SetUsageScenario` được gọi bởi hệ thống để cho provider biết nó đang được dùng trong bối cảnh nào (đăng nhập, mở khóa, thay đổi mật khẩu...). Dựa vào đây, provider sẽ quyết định có hiển thị tile đăng nhập hay không. Ví dụ, nếu kịch bản là `CPUS_UNLOCK_WORKSTATION` (mở khóa máy), nó sẽ cho phép hiển thị.

3.  **Tạo và Quản lý Credential (Enumerate Credentials):**
    *   Hàm `GetCredentialCount` được gọi để hỏi xem provider có bao nhiêu credential muốn hiển thị.
    *   Hàm `GetCredentialAt` được gọi lặp đi lặp lại để lấy từng đối tượng `ICredentialProviderCredential` (chính là các instance của `CSampleCredential` mà chúng ta đã phân tích).
    *   Về cơ bản, nó tạo ra một danh sách các "ô đăng nhập" để hệ thống hiển thị. Trong trường- hợp này, nó chỉ tạo một ô duy nhất.

### Luồng hoạt động chính:

1.  **Khởi tạo Provider:**
    *   Hệ thống tạo một đối tượng `CSampleProvider`.
    *   Constructor `CSampleProvider::CSampleProvider` được gọi, khởi tạo các biến thành viên (ví dụ: `_cRef` - reference count).

2.  **Thiết lập Kịch bản Sử dụng (SetUsageScenario):**
    *   Hệ thống gọi `CSampleProvider::SetUsageScenario`.
    *   Provider kiểm tra `cpus` và `dwFlags`. Nếu là kịch bản hợp lệ (ví dụ: đăng nhập hoặc mở khóa) và người dùng hiện tại có trong danh sách được phép, nó sẽ tạo một instance của `CSampleCredential` (ô đăng nhập) và lưu lại.
    *   Nếu kịch bản không hợp lệ, nó sẽ giải phóng `CSampleCredential` đã có và không hiển thị gì cả.

3.  **Liệt kê Credential (Enumeration):**
    *   Hệ thống gọi `CSampleProvider::GetCredentialCount` để lấy số lượng ô đăng nhập. Provider sẽ trả về 1 nếu `_pCredential` (instance của `CSampleCredential`) đã được tạo, ngược lại trả về 0.
    *   Hệ thống gọi `CSampleProvider::GetCredentialAt` để lấy đối tượng credential tại một chỉ số cụ thể. Provider sẽ trả về con trỏ đến `_pCredential` nếu chỉ số là 0.

4.  **Giải phóng:**
    *   Khi không cần thiết nữa, destructor `CSampleProvider::~CSampleProvider` được gọi, đảm bảo rằng đối tượng `_pCredential` được giải phóng an toàn.

### Tóm tắt:

-   `CSampleProvider.cpp` là **entry point** từ phía hệ thống để quản lý các lựa chọn đăng nhập của bạn.
-   Nó quyết định **có hiển thị** ô đăng nhập hay không dựa trên kịch bản (đăng nhập, mở khóa...).
-   Nó **tạo ra** các instance của `CSampleCredential` để đại diện cho từng ô đăng nhập sẽ được hiển thị trên màn hình.
-   Nó không liên quan trực tiếp đến việc lấy mật khẩu hay xử lý giao diện người dùng của ô đăng nhập. Đó là nhiệm vụ của `CSampleCredential.cpp`.

---

### Phân tích `helpers.cpp` và `helpers.h`

Các file này chứa các hàm tiện ích (utility functions) đóng vai trò hỗ trợ cho các tác vụ phức tạp và lặp đi lặp lại trong dự án. Việc tách các hàm này ra giúp mã nguồn ở các file chính trở nên sạch sẽ và tập trung vào logic chính hơn.

#### Vai trò chính:

1.  **Quản lý bộ nhớ và cấu trúc:**
    *   `FieldDescriptorCoAllocCopy`, `FieldDescriptorCopy`: Sao chép các cấu trúc `CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR`. Điều này cần thiết vì Windows yêu cầu các cấu trúc này phải được cấp phát bằng một loại bộ nhớ đặc biệt (`CoTaskMemAlloc`).
    *   `UnicodeStringInitWithString`: Tạo một chuỗi `UNICODE_STRING` (một định dạng chuỗi đặc biệt của Windows) từ một chuỗi thông thường.

2.  **Đóng gói dữ liệu xác thực (Serialization):**
    *   `KerbInteractiveUnlockLogonInit` và `KerbInteractiveUnlockLogonPack`: Đây là cặp hàm cực kỳ quan trọng, là trái tim của việc "đóng gói" thông tin đăng nhập để gửi cho hệ thống.
        *   `...Init`: Khởi tạo một cấu trúc `KERB_INTERACTIVE_UNLOCK_LOGON` với các con trỏ đến chuỗi domain, username, và password.
        *   `...Pack`: Nhận cấu trúc đã khởi tạo và "đóng gói" nó vào một vùng bộ nhớ duy nhất, liền mạch. Đây là định dạng mà LSA (Local Security Authority) của Windows mong đợi. Nó sẽ thay thế các con trỏ bằng các giá trị offset (vị trí tương đối) bên trong vùng nhớ đó.

3.  **Bảo mật và Tương tác với LSA:**
    *   `RetrieveNegotiateAuthPackage`: Lấy về "mã định danh" của gói xác thực Kerberos/Negotiate. Mã này cần thiết khi giao tiếp với LSA.
    *   `ProtectIfNecessaryAndCopyPassword`: Một hàm bảo mật quan trọng. Nó sử dụng API `CredProtectW` của Windows để **mã hóa mật khẩu** trước khi gửi đi trong các kịch bản đăng nhập hoặc mở khóa. Nó đủ thông minh để không mã hóa trong môi trường CredUI (nơi không cần thiết) hoặc nếu mật khẩu đã được mã hóa từ trước.

4.  **Tiện ích chuỗi:**
    *   `SplitDomainAndUsername`: Một hàm tiện ích đơn giản để tách chuỗi `DOMAIN\Username` thành hai phần riêng biệt.

### Phân tích `Dll.cpp`

File này chứa mã nguồn "tiêu chuẩn" (boilerplate) cho một thư viện COM (Component Object Model). Vai trò của nó là cho phép hệ điều hành Windows (cụ thể là LogonUI) có thể tìm, tải, và tạo các đối tượng từ DLL của chúng ta.

#### Các thành phần chính:

1.  **Quản lý Vòng đời DLL:**
    *   `DllMain`: Điểm vào (entry point) của DLL. Được hệ điều hành gọi khi DLL được nạp vào hoặc giải phóng khỏi một tiến trình.
    *   `DllAddRef`, `DllRelease`, `DllCanUnloadNow`: Các hàm này quản lý một bộ đếm tham chiếu toàn cục. Chúng đảm bảo rằng DLL sẽ không bị gỡ khỏi bộ nhớ khi vẫn còn đối tượng đang sử dụng nó.

2.  **Nhà máy Lớp (Class Factory):**
    *   `DllGetClassObject`: Một hàm xuất chuẩn của COM. Khi LogonUI muốn tạo một đối tượng `CSampleProvider`, nó sẽ gọi hàm này trước tiên để yêu cầu một "nhà máy" có khả năng tạo ra đối tượng đó.
    *   `CClassFactory`: Là lớp triển khai "nhà máy" này.
    *   `CClassFactory::CreateInstance`: Phương thức quan trọng nhất của nhà máy. Khi được gọi, nó sẽ thực sự tạo ra một đối tượng mới của lớp `CSampleProvider` (bằng cách gọi hàm `CSample_CreateInstance` được định nghĩa trong `CSampleProvider.cpp`) và trả về cho LogonUI.

#### Tóm tắt:

-   `Dll.cpp` là **cửa ngõ** để thế giới bên ngoài (Windows) giao tiếp với code của bạn. Nó không chứa logic nghiệp vụ về xác thực, mà chỉ làm nhiệm vụ của một thư viện COM tiêu chuẩn.
-   Nó cho phép hệ thống tạo ra `CSampleProvider`, và từ đó, toàn bộ luồng hoạt động của Credential Provider được bắt đầu.