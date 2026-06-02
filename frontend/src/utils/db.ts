import Dexie, { type Table } from 'dexie';

export interface LocalBoard {
  id: string;
  name: string;
  description: string;
  color_theme: string;
  is_archived: number; // 0 or 1
  workspace_id: string;
  created_at: string;
  updated_at: string;
}

export interface LocalList {
  id: string;
  board_id: string;
  title: string;
  position: number;
}

export interface LocalCard {
  id: string;
  list_id: string;
  title: string;
  description: string;
  position: number;
  due_date?: string;
  labels: string[];
  progress_percentage: number;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface LocalSubtask {
  id: string;
  card_id: string;
  title: string;
  is_completed: number; // 0 or 1
  position: number;
}

export interface LocalComment {
  id: string;
  card_id: string;
  user_id: string;
  username: string;
  avatar_url?: string;
  content: string;
  created_at: string;
}

export interface LocalNote {
  id: string;
  card_id: string;
  content: string;
  created_at: string;
}

export interface OfflineMutation {
  id: string; // unique mutation ID (idempotency key)
  action: string; // e.g. "card_moved", "card_created", "subtask_toggled"
  entity_type: string;
  entity_id: string;
  data: any;
  timestamp: string;
  board_id: string;
}

export class TNotesLocalDB extends Dexie {
  boards!: Table<LocalBoard>;
  lists!: Table<LocalList>;
  cards!: Table<LocalCard>;
  subtasks!: Table<LocalSubtask>;
  comments!: Table<LocalComment>;
  notes!: Table<LocalNote>;
  mutationsQueue!: Table<OfflineMutation>;

  constructor() {
    super('TNotesLocalDB');
    this.version(1).stores({
      boards: 'id, workspace_id, is_archived',
      lists: 'id, board_id, position',
      cards: 'id, list_id, position',
      subtasks: 'id, card_id, position',
      comments: 'id, card_id, created_at',
      notes: 'id, card_id, created_at',
      mutationsQueue: 'id, board_id, timestamp'
    });
  }
}

export const db = new TNotesLocalDB();
