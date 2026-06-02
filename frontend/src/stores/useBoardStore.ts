import { create } from 'zustand';
import { apiRequest } from '../services/api';
import { db, type LocalBoard, type LocalList, type LocalCard, type LocalSubtask, type LocalComment, type LocalNote, type OfflineMutation } from '../utils/db';

export interface BoardState {
  // Sync Status
  isOnline: boolean;
  isSyncing: boolean;
  syncQueueLength: number;
  
  // Active Data
  activeBoard: LocalBoard | null;
  lists: LocalList[];
  cards: LocalCard[];
  subtasks: Record<string, LocalSubtask[]>; // card_id -> subtasks
  comments: Record<string, LocalComment[]>; // card_id -> comments
  quickNotes: Record<string, LocalNote[]>; // card_id -> notes

  // Actions
  toggleSyncMode: () => void;
  setOnlineStatus: (status: boolean) => void;
  loadBoardFromLocal: (boardId: string) => Promise<void>;
  initializeBoard: (boardData: any) => Promise<void>;

  // Board Mutators
  addCard: (title: string, listId: string, createdBy: string) => Promise<void>;
  moveCard: (cardId: string, sourceListId: string, destListId: string, newPosition: number) => Promise<void>;
  updateCardDetails: (cardId: string, title?: string, description?: string, dueDate?: string, labels?: string[]) => Promise<void>;
  deleteCard: (cardId: string) => Promise<void>;

  // Checklist & Notes Mutators
  addSubtask: (cardId: string, title: string) => Promise<void>;
  toggleSubtask: (subtaskId: string, isCompleted: boolean) => Promise<void>;
  addQuickNote: (cardId: string, content: string) => Promise<void>;
  addComment: (cardId: string, content: string, userId: string, username: string) => Promise<void>;

  // Sync Mechanics
  triggerQueueSync: () => Promise<void>;
}

export const useBoardStore = create<BoardState>((set, get) => {
  // Event listener to check network connectivity dynamically
  window.addEventListener('online', () => get().setOnlineStatus(true));
  window.addEventListener('offline', () => get().setOnlineStatus(false));

  return {
    isOnline: navigator.onLine,
    isSyncing: false,
    syncQueueLength: 0,
    activeBoard: null,
    lists: [],
    cards: [],
    subtasks: {},
    comments: {},
    quickNotes: {},

    toggleSyncMode: () => {
      const nextOnline = !get().isOnline;
      get().setOnlineStatus(nextOnline);
    },

    setOnlineStatus: async (status) => {
      set({ isOnline: status });
      if (status) {
        await get().triggerQueueSync();
      }
    },

    loadBoardFromLocal: async (boardId) => {
      const board = await db.boards.get(boardId);
      if (!board) return;

      const lists = await db.lists.where({ board_id: boardId }).sortBy('position');
      const listIds = lists.map(l => l.id);
      
      const cards = await db.cards.where('list_id').anyOf(listIds).sortBy('position');
      const cardIds = cards.map(c => c.id);

      const allSubtasks = await db.subtasks.where('card_id').anyOf(cardIds).toArray();
      const allComments = await db.comments.where('card_id').anyOf(cardIds).toArray();
      const allNotes = await db.notes.where('card_id').anyOf(cardIds).toArray();

      const subtasksMap: Record<string, LocalSubtask[]> = {};
      const commentsMap: Record<string, LocalComment[]> = {};
      const notesMap: Record<string, LocalNote[]> = {};

      allSubtasks.forEach(s => {
        if (!subtasksMap[s.card_id]) subtasksMap[s.card_id] = [];
        subtasksMap[s.card_id].push(s);
      });

      allComments.forEach(c => {
        if (!commentsMap[c.card_id]) commentsMap[c.card_id] = [];
        commentsMap[c.card_id].push(c);
      });

      allNotes.forEach(n => {
        if (!notesMap[n.card_id]) notesMap[n.card_id] = [];
        notesMap[n.card_id].push(n);
      });

      // Sort checklists by position
      Object.keys(subtasksMap).forEach(key => {
        subtasksMap[key].sort((a, b) => a.position - b.position);
      });

      const queueLen = await db.mutationsQueue.where({ board_id: boardId }).count();

      set({
        activeBoard: board,
        lists,
        cards,
        subtasks: subtasksMap,
        comments: commentsMap,
        quickNotes: notesMap,
        syncQueueLength: queueLen
      });
    },

    initializeBoard: async (boardData) => {
      // Seed full board data fetched from backend into Dexie local database
      await db.transaction('rw', [db.boards, db.lists, db.cards, db.subtasks, db.comments, db.notes], async () => {
        const localBoard: LocalBoard = {
          id: boardData.id,
          name: boardData.name,
          description: boardData.description || '',
          color_theme: boardData.color_theme || '#6366F1',
          is_archived: boardData.is_archived ? 1 : 0,
          workspace_id: boardData.workspace_id,
          created_at: boardData.created_at,
          updated_at: boardData.created_at
        };
        await db.boards.put(localBoard);

        if (boardData.lists) {
          for (const l of boardData.lists) {
            const localList: LocalList = {
              id: l.id,
              board_id: boardData.id,
              title: l.title,
              position: l.position
            };
            await db.lists.put(localList);

            if (l.cards) {
              for (const c of l.cards) {
                const localCard: LocalCard = {
                  id: c.id,
                  list_id: l.id,
                  title: c.title,
                  description: c.description || '',
                  position: c.position,
                  due_date: c.due_date,
                  labels: c.labels || [],
                  progress_percentage: c.progress_percentage || 0,
                  created_by: c.created_by,
                  created_at: c.created_at || new Date().toISOString(),
                  updated_at: c.updated_at || new Date().toISOString()
                };
                await db.cards.put(localCard);

                if (c.subtasks) {
                  for (const s of c.subtasks) {
                    await db.subtasks.put({
                      id: s.id,
                      card_id: c.id,
                      title: s.title,
                      is_completed: s.is_completed ? 1 : 0,
                      position: s.position
                    });
                  }
                }

                if (c.quick_notes) {
                  for (const n of c.quick_notes) {
                    await db.notes.put({
                      id: n.id,
                      card_id: c.id,
                      content: n.content,
                      created_at: n.created_at
                    });
                  }
                }

                if (c.comments) {
                  for (const cm of c.comments) {
                    await db.comments.put({
                      id: cm.id,
                      card_id: c.id,
                      user_id: cm.user_id,
                      username: cm.username || 'Anonymous',
                      avatar_url: cm.avatar_url,
                      content: cm.content,
                      created_at: cm.created_at
                    });
                  }
                }
              }
            }
          }
        }
      });

      await get().loadBoardFromLocal(boardData.id);
    },

    addCard: async (title, listId, createdBy) => {
      const cardId = crypto.randomUUID();
      const boardId = get().activeBoard?.id || '';
      const position = get().cards.filter(c => c.list_id === listId).length;

      const newCard: LocalCard = {
        id: cardId,
        list_id: listId,
        title,
        description: '',
        position,
        labels: [],
        progress_percentage: 0,
        created_by: createdBy,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };

      // Optimistic state write
      await db.cards.put(newCard);
      
      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'card_created',
        entity_type: 'card',
        entity_id: cardId,
        data: { list_id: listId, title, description: '' },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);
      
      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    moveCard: async (cardId, _sourceListId, destListId, newPosition) => {
      const boardId = get().activeBoard?.id || '';

      await db.transaction('rw', db.cards, async () => {
        const activeCard = await db.cards.get(cardId);
        if (!activeCard) return;

        activeCard.list_id = destListId;
        activeCard.position = newPosition;
        activeCard.updated_at = new Date().toISOString();
        await db.cards.put(activeCard);
      });

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'card_moved',
        entity_type: 'card',
        entity_id: cardId,
        data: { new_list_id: destListId, new_position: newPosition },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    updateCardDetails: async (cardId, title, description, dueDate, labels) => {
      const boardId = get().activeBoard?.id || '';
      
      const card = await db.cards.get(cardId);
      if (!card) return;

      if (title !== undefined) card.title = title;
      if (description !== undefined) card.description = description;
      if (dueDate !== undefined) card.due_date = dueDate;
      if (labels !== undefined) card.labels = labels;
      card.updated_at = new Date().toISOString();

      await db.cards.put(card);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'card_updated',
        entity_type: 'card',
        entity_id: cardId,
        data: { title, description, due_date: dueDate, labels },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    deleteCard: async (cardId) => {
      const boardId = get().activeBoard?.id || '';

      await db.cards.delete(cardId);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'card_deleted',
        entity_type: 'card',
        entity_id: cardId,
        data: {},
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    addSubtask: async (cardId, title) => {
      const boardId = get().activeBoard?.id || '';
      const subtaskId = crypto.randomUUID();
      const existing = get().subtasks[cardId] || [];
      const position = existing.length;

      const newSub: LocalSubtask = {
        id: subtaskId,
        card_id: cardId,
        title,
        is_completed: 0,
        position
      };

      await db.subtasks.put(newSub);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'subtask_created',
        entity_type: 'subtask',
        entity_id: subtaskId,
        data: { card_id: cardId, title },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    toggleSubtask: async (subtaskId, isCompleted) => {
      const boardId = get().activeBoard?.id || '';
      const sub = await db.subtasks.get(subtaskId);
      if (!sub) return;

      sub.is_completed = isCompleted ? 1 : 0;
      await db.subtasks.put(sub);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'subtask_toggled',
        entity_type: 'subtask',
        entity_id: subtaskId,
        data: { is_completed: isCompleted },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    addQuickNote: async (cardId, content) => {
      const boardId = get().activeBoard?.id || '';
      const noteId = crypto.randomUUID();

      const newNote: LocalNote = {
        id: noteId,
        card_id: cardId,
        content,
        created_at: new Date().toISOString()
      };

      await db.notes.put(newNote);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'note_created',
        entity_type: 'note',
        entity_id: noteId,
        data: { card_id: cardId, content },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    addComment: async (cardId, content, userId, username) => {
      const boardId = get().activeBoard?.id || '';
      const commentId = crypto.randomUUID();

      const newComment: LocalComment = {
        id: commentId,
        card_id: cardId,
        user_id: userId,
        username,
        avatar_url: `https://api.dicebear.com/7.x/bottts/svg?seed=${username}`,
        content,
        created_at: new Date().toISOString()
      };

      await db.comments.put(newComment);

      const newMutation: OfflineMutation = {
        id: crypto.randomUUID(),
        action: 'comment_added',
        entity_type: 'comment',
        entity_id: cardId,
        data: { content },
        timestamp: new Date().toISOString(),
        board_id: boardId
      };
      await db.mutationsQueue.put(newMutation);

      await get().loadBoardFromLocal(boardId);

      if (get().isOnline) {
        await get().triggerQueueSync();
      }
    },

    triggerQueueSync: async () => {
      const boardId = get().activeBoard?.id;
      if (!boardId || get().isSyncing) return;

      const mutations = await db.mutationsQueue.where({ board_id: boardId }).toArray();
      if (mutations.length === 0) return;

      set({ isSyncing: true });

      try {
        await apiRequest(`/api/boards/${boardId}/sync`, {
          method: 'POST',
          body: JSON.stringify({ operations: mutations }),
        });

        const mutationIds = mutations.map(m => m.id);
        await db.mutationsQueue.bulkDelete(mutationIds);

        const fullBoard = await apiRequest<Record<string, unknown>>(`/api/boards/${boardId}`);
        await get().initializeBoard(fullBoard);
      } catch (err) {
        console.error('Offline batch sync error:', err);
      } finally {
        set({ isSyncing: false });
      }
    }
  };
});
