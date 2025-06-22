import axios from 'axios'

// 開発環境でのみX-Forwarded-Userヘッダーを付与するaxiosインスタンス
export const api = axios.create()

// 開発環境でのリクエストインターセプター
if (import.meta.env.DEV) {
  api.interceptors.request.use((config) => {
    // 仮のユーザーID（実際の実装では適切に取得）
    config.headers['X-Forwarded-User'] = 'mumumu'
    return config
  })
}

export default api
