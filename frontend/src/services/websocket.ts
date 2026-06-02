type MessageHandler = (event: string, payload: Record<string, unknown>) => void;

const BACKOFF_MS = [1000, 2000, 4000, 8000, 16000];

export class BoardWebSocket {
  private ws: WebSocket | null = null;
  private boardId: string | null = null;
  private attempt = 0;
  private closed = false;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private onMessage: MessageHandler;

  constructor(onMessage: MessageHandler) {
    this.onMessage = onMessage;
  }

  connect(boardId: string): void {
    this.boardId = boardId;
    this.closed = false;
    this.open();
  }

  disconnect(): void {
    this.closed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.ws?.close();
    this.ws = null;
  }

  send(event: string, payload: Record<string, unknown>): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ event, payload }));
    }
  }

  private open(): void {
    const token = localStorage.getItem('jwt_token');
    if (!token || !this.boardId) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    this.ws = new WebSocket(`${protocol}//${window.location.host}/ws?token=${token}`);

    this.ws.onopen = () => {
      this.attempt = 0;
      this.send('join_board', { board_id: this.boardId });
    };

    this.ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data as string);
        this.onMessage(msg.event, msg.payload || {});
      } catch {
        /* ignore malformed frames */
      }
    };

    this.ws.onclose = () => {
      this.ws = null;
      if (!this.closed && this.boardId) {
        const delay = BACKOFF_MS[Math.min(this.attempt, BACKOFF_MS.length - 1)];
        this.attempt += 1;
        this.reconnectTimer = setTimeout(() => this.open(), delay);
      }
    };
  }
}
