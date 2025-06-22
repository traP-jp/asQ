export interface TypewriterMessage {
  id: string
  message: string // 受信した完全なメッセージ（バッファ）
  displayedMessage: string // 現在表示中のメッセージ
  currentIndex: number // 現在表示中の文字インデックス
}

export class TypewriterController {
  private typewriterInterval: number | null = null
  private readonly typingSpeed: number

  constructor(typingSpeed: number = 30) {
    this.typingSpeed = typingSpeed
  }

  // メッセージ更新時の処理（新しい文字が受信されたとき）
  updateMessage(messages: TypewriterMessage[], messageIndex: number): void {
    const message = messages[messageIndex]
    if (!message) return

    // まだタイピング中でない場合は開始
    if (!this.typewriterInterval) {
      this.startTyping(messages, messageIndex)
    }
  }

  // タイプライター効果を開始（プライベートメソッド）
  private startTyping(messages: TypewriterMessage[], messageIndex: number): void {
    if (this.typewriterInterval) {
      clearInterval(this.typewriterInterval)
    }

    this.typewriterInterval = setInterval(() => {
      const message = messages[messageIndex]
      if (!message) {
        this.stopTyping()
        return
      }

      if (message.currentIndex < message.message.length) {
        message.displayedMessage = message.message.substring(0, message.currentIndex + 1)
        message.currentIndex++
      } else {
        // タイピング完了
        this.stopTyping()
      }
    }, this.typingSpeed)
  }

  // タイプライター効果を停止（プライベートメソッド）
  private stopTyping(): void {
    if (this.typewriterInterval) {
      clearInterval(this.typewriterInterval)
      this.typewriterInterval = null
    }
  }

  // クリーンアップ（コンポーネントの破棄時に呼び出し）
  cleanup(): void {
    this.stopTyping()
  }
}

// ファクトリー関数
export function createTypewriter(typingSpeed: number = 30): TypewriterController {
  return new TypewriterController(typingSpeed)
}
