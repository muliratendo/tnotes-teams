import { create } from 'zustand';

export interface PeerCursor {
  userId: string;
  username: string;
  cardId: string;
  x: number;
  y: number;
}

export interface PeerUser {
  userId: string;
  username: string;
  avatarUrl: string;
}

export interface PresenceState {
  peers: PeerUser[];
  cursors: Record<string, PeerCursor>; // userId -> PeerCursor
  typingUsers: Record<string, string[]>; // cardId -> array of usernames
  
  // Actions
  addPeer: (peer: PeerUser) => void;
  removePeer: (userId: string) => void;
  updateCursor: (cursor: PeerCursor) => void;
  removeCursor: (userId: string) => void;
  setTyping: (cardId: string, username: string, isTyping: boolean) => void;
  clearAll: () => void;
}

export const usePresenceStore = create<PresenceState>((set) => ({
  peers: [],
  cursors: {},
  typingUsers: {},

  addPeer: (peer) => set((state) => {
    // Prevent duplicate entries
    if (state.peers.some(p => p.userId === peer.userId)) return {};
    return { peers: [...state.peers, peer] };
  }),

  removePeer: (userId) => set((state) => ({
    peers: state.peers.filter(p => p.userId !== userId),
    cursors: (() => {
      const nextCursors = { ...state.cursors };
      delete nextCursors[userId];
      return nextCursors;
    })()
  })),

  updateCursor: (cursor) => set((state) => ({
    cursors: {
      ...state.cursors,
      [cursor.userId]: cursor
    }
  })),

  removeCursor: (userId) => set((state) => {
    const nextCursors = { ...state.cursors };
    delete nextCursors[userId];
    return { cursors: nextCursors };
  }),

  setTyping: (cardId, username, isTyping) => set((state) => {
    const active = state.typingUsers[cardId] || [];
    let next: string[];
    
    if (isTyping) {
      if (active.includes(username)) return {};
      next = [...active, username];
    } else {
      next = active.filter(u => u !== username);
    }

    return {
      typingUsers: {
        ...state.typingUsers,
        [cardId]: next
      }
    };
  }),

  clearAll: () => set({ peers: [], cursors: {}, typingUsers: {} })
}));
