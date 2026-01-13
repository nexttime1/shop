# 前端與後端接口同步說明 (基於 shop.md)

本文件詳細記錄了前端接口與後端 `shop.md` 文檔之間的差異、缺失的請求以及多餘的請求。

## 1. 基礎配置 (Request Base)

- **狀態碼變更**: 前端已改為根據後端返回的 `code: 0` 判斷成功。
- **多服務支持**: 根據 `shop.md`，前端 `request.js` 已配置為自動根據路徑前綴選擇正確的服務地址：
  - `/u/`: 用戶服務 (8080)
  - `/g/`: 商品服務 (8081)
  - `/o/`: 訂單服務 (8083)
  - `/up/`: 用戶操作服務 (8084)
  - `/oss/`: 上傳服務 (8084)

---

## 2. 缺失的接口 (在前端模板中有但後端 shop.md 未定義)

以下接口在前端管理後台中常用，但當前後端文檔中尚未定義，目前前端已使用 Mock 或預測路徑處理：

### A. 獲取當前用戶信息 (GetUserInfo)
- **路徑**: `GET /u/v1/user/info`
- **功能**: 登錄成功後，前端需要獲取當前用戶的權限、暱稱和頭像，以渲染導航欄和處理頁面權限。
- **請求參數**: 無 (通過 Header 中的 `x-token` 識別)
- **預期返回 (Data 字段)**:
  ```json
  {
    "roles": ["admin"],
    "name": "管理員",
    "avatar": "https://example.com/avatar.png"
  }
  ```

### B. 退出登錄 (Logout)
- **路徑**: `POST /u/v1/user/logout`
- **功能**: 通知後端銷毀當前 Token 的有效性。
- **請求參數**: 無 (通過 Header 中的 `x-token` 識別)
- **預期返回**: `{ "code": 0, "msg": "成功", "data": {} }`

### C. 批量更新商品狀態 (BatchUpdateProductStatus)
- **路徑**: `PATCH /g/v1/good/batch`
- **功能**: 在商品列表頁面批量上架、下架、設置新品或熱門。
- **請求參數**:
  ```json
  {
    "ids": [1, 2, 3],
    "on_sale": true,
    "is_new": true,
    "is_hot": false
  }
  ```
- **預期返回**: `{ "code": 0, "msg": "成功", "data": {} }`

---

## 3. 前端缺失的功能 (在 shop.md 中有但前端尚未完全對接)

### A. 購物車管理 (src/api/cart.js)
- **接口**:
  - `GET /o/v1/shopcarts`: 獲取某人的購物車。
  - `POST /o/v1/shopcarts`: 加入商品到購物車。
  - `PATCH /o/v1/shopcarts/:id`: 更新數量或選中狀態。
  - `DELETE /o/v1/shopcarts/:id`: 刪除條目。
- **說明**: 已在 `src/api/cart.js` 中新增，需在前端視圖中調用。

### B. 用户收货地址 (src/api/userAddress.js)
- **接口**:
  - `GET /up/v1/address`
  - `POST /up/v1/address`
  - `PUT /up/v1/address/:id`
  - `DELETE /up/v1/address/:id`
- **說明**: 已在 `src/api/userAddress.js` 中新增。

### C. 留言與收藏
- **接口**: `/up/v1/message`, `/up/v1/userfavs`
- **說明**: 目前尚未在 API 文件中單獨列出，如有需求可進一步封裝。

---

## 4. 參數命名規範差異

根據 `shop.md`，請求參數已統一調整為：
- 商品列表使用 `PageInfo` 結構。
- 登錄使用 `mobile`, `password`, `captcha_id`, `answer`。
- 註冊使用 `mobile`, `password`, `code`, `role`。

---

## 5. 其他建議

1. **統一 ID 類型**: 後端大部分 ID 為 `int32`，前端發送時請確保不為字符串。
2. **圖片上傳流程**: 後端流程為先調用 `/oss/v1/token` 獲得七牛雲 token，前端直傳七牛後，由七牛回調後端或前端將 URL 傳給商品創建接口。前端 `src/api/oss.js` 已更新。
