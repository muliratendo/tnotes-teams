import React, { useState, useEffect, useRef } from 'react';
import { apiRequest, downloadExport } from '../services/api';
import { BoardWebSocket } from '../services/websocket';
import { useBoardStore } from '../stores/useBoardStore';
import { usePresenceStore } from '../stores/usePresenceStore';
import { ListColumn } from './ListColumn';
import { CardDrawer } from './CardDrawer';
import { CursorOverlay } from './CursorOverlay';
import { 
  Plus, 
  Download, 
  Table, 
  FolderLock, 
  Sparkles,
  RefreshCw
} from 'lucide-react';

interface BoardCanvasProps {
  user: { id: string; username: string; email: string; role?: string };
  triggerViewerWarning: (message: string) => void;
}

export const BoardCanvas: React.FC<BoardCanvasProps> = ({ user, triggerViewerWarning }) => {
  const isViewer = user.role === 'viewer';
  
  // Board stores
  const activeBoard = useBoardStore(state => state.activeBoard);
  const lists = useBoardStore(state => state.lists);
  const cards = useBoardStore(state => state.cards);
  
  // Mutators
  const addCard = useBoardStore(state => state.addCard);
  const moveCard = useBoardStore(state => state.moveCard);
  const isOnline = useBoardStore(state => state.isOnline);
  const isSyncing = useBoardStore(state => state.isSyncing);

  // States
  const [showAddList, setShowAddList] = useState(false);
  const [listTitle, setListTitle] = useState('');
  const [selectedCardId, setSelectedCardId] = useState<string | null>(null);
  
  const wsClient = useRef<BoardWebSocket | null>(null);
  const canvasRef = useRef<HTMLDivElement>(null);
  const lastMouseUpdate = useRef<number>(0);

  useEffect(() => {
    if (!activeBoard || !isOnline) {
      wsClient.current?.disconnect();
      wsClient.current = null;
      return;
    }

    const client = new BoardWebSocket((event, payload) => {
      const presence = usePresenceStore.getState();
      switch (event) {
        case 'user_present':
          presence.addPeer({
            userId: String(payload.user_id),
            username: String(payload.username),
            avatarUrl:
              typeof payload.avatar_url === 'string' && payload.avatar_url.length > 0
                ? payload.avatar_url
                : `https://api.dicebear.com/7.x/bottts/svg?seed=${String(payload.username)}`,
          });
          break;
        case 'user_absent':
          presence.removePeer(String(payload.user_id));
          break;
        case 'cursor_broadcast':
          presence.updateCursor({
            userId: String(payload.user_id),
            username: String(payload.username),
            cardId: String(payload.card_id || ''),
            x: Number(payload.x),
            y: Number(payload.y),
          });
          break;
        case 'typing_broadcast':
          presence.setTyping(
            String(payload.card_id),
            String(payload.username),
            Boolean(payload.is_typing)
          );
          break;
        case 'board_mutated':
          useBoardStore.getState().loadBoardFromLocal(activeBoard.id);
          break;
      }
    });

    client.connect(activeBoard.id);
    wsClient.current = client;

    return () => {
      client.disconnect();
      usePresenceStore.getState().clearAll();
    };
  }, [activeBoard, isOnline]);

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!activeBoard) return;
    const now = Date.now();
    if (now - lastMouseUpdate.current < 100) return;
    lastMouseUpdate.current = now;

    if (canvasRef.current) {
      const rect = canvasRef.current.getBoundingClientRect();
      wsClient.current?.send('cursor_update', {
        board_id: activeBoard.id,
        card_id: selectedCardId || '',
        x: e.clientX - rect.left,
        y: e.clientY - rect.top,
      });
    }
  };

  const handleAddListSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewer) {
      triggerViewerWarning("Viewers cannot create columns.");
      return;
    }
    if (!listTitle.trim() || !activeBoard) return;

    try {
      await apiRequest(`/api/boards/${activeBoard.id}/lists`, {
        method: 'POST',
        body: JSON.stringify({ title: listTitle }),
      });
      setListTitle('');
      setShowAddList(false);
      const fullBoard = await apiRequest<Record<string, unknown>>(`/api/boards/${activeBoard.id}`);
      await useBoardStore.getState().initializeBoard(fullBoard);
    } catch (err) {
      console.error(err);
    }
  };

  // HTML5 Drag & Drop handlers
  const handleDragStart = (e: React.DragEvent, cardId: string) => {
    if (isViewer) return;
    e.dataTransfer.setData('text/plain', cardId);
  };

  const handleCardDrop = async (e: React.DragEvent, destListId: string) => {
    if (isViewer) return;
    const cardId = e.dataTransfer.getData('text/plain');
    const sourceCard = cards.find(c => c.id === cardId);
    if (!sourceCard || sourceCard.list_id === destListId) return;

    const destListCards = cards.filter(c => c.list_id === destListId);
    moveCard(cardId, sourceCard.list_id, destListId, destListCards.length);
  };

  const handleExportJSON = async () => {
    if (!activeBoard) return;
    try {
      await downloadExport(`/api/boards/${activeBoard.id}/export`, 'board-export.json');
    } catch (err) {
      console.error(err);
    }
  };

  const handleExportCSV = async () => {
    if (!activeBoard) return;
    try {
      await downloadExport(`/api/boards/${activeBoard.id}/export/csv`, 'board-export.csv');
    } catch (err) {
      console.error(err);
    }
  };

  if (!activeBoard) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center bg-[#0B0F19] text-[#9CA3AF] p-6 space-y-4 font-outfit">
        <FolderLock className="h-16 w-16 text-[#374151]" />
        <h2 className="text-xl font-bold text-white tracking-tight">Select or Create a Kanban Board</h2>
        <p className="text-sm text-[#9CA3AF] max-w-sm text-center">
          Open the sidebar workspace navigator to inspect boards or add a new sprint timeline.
        </p>
      </div>
    );
  }

  return (
    <div 
      ref={canvasRef}
      onMouseMove={handleMouseMove}
      className="flex-1 flex flex-col bg-[#0B0F19] relative overflow-hidden"
    >
      {/* Board Utility Header */}
      <header className="px-6 py-4 border-b border-[#1F2937]/50 flex justify-between items-center bg-[#111827]/10 backdrop-blur-sm z-10 font-outfit">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <h2 className="text-xl font-bold tracking-tight text-white">{activeBoard.name}</h2>
            {isSyncing && (
              <RefreshCw className="h-4 w-4 text-[#6366F1] animate-spin" />
            )}
          </div>
          <p className="text-xs text-[#9CA3AF] font-medium max-w-md truncate">
            {activeBoard.description || 'Collaborative project board'}
          </p>
        </div>

        {/* Action triggers */}
        <div className="flex items-center gap-3">
          {/* Simulated engine indicator */}
          <span className="text-[10px] bg-[#6366F1]/10 text-[#6366F1] px-2.5 py-1 rounded-full font-bold uppercase flex items-center gap-1.5 border border-[#6366F1]/20">
            <Sparkles className="h-3 w-3 animate-pulse" />
            Co-Presence Active
          </span>
          
          <button
            onClick={handleExportJSON}
            className="px-3.5 py-2 text-xs font-bold bg-[#1F2937] hover:bg-[#374151] border border-[#374151] rounded-lg transition-colors flex items-center gap-1.5"
            title="Export JSON backup"
          >
            <Download className="h-3.5 w-3.5" />
            <span>JSON</span>
          </button>
          
          <button
            onClick={handleExportCSV}
            className="px-3.5 py-2 text-xs font-bold bg-[#1F2937] hover:bg-[#374151] border border-[#374151] rounded-lg transition-colors flex items-center gap-1.5"
            title="Export CSV cards list"
          >
            <Table className="h-3.5 w-3.5" />
            <span>CSV</span>
          </button>
        </div>
      </header>

      {/* Horizontal Kanban Canvas */}
      <div className="flex-1 flex gap-5 overflow-x-auto p-6 items-start relative select-none">
        {lists.map(l => (
          <ListColumn
            key={l.id}
            list={l}
            cards={cards.filter(c => c.list_id === l.id)}
            onCardClick={setSelectedCardId}
            onCardDrop={handleCardDrop}
            onCardDragStart={handleDragStart}
            onAddCard={(title, listId) => addCard(title, listId, user.id)}
            isViewer={isViewer}
            triggerViewerWarning={triggerViewerWarning}
          />
        ))}

        {/* Add column controls */}
        {showAddList ? (
          <form onSubmit={handleAddListSubmit} className="w-72 p-4 rounded-2xl bg-[#111827] border border-[#1F2937] space-y-3 font-outfit shadow-xl">
            <input
              type="text"
              required
              autoFocus
              value={listTitle}
              onChange={e => setListTitle(e.target.value)}
              placeholder="e.g. In Progress, Dev Review"
              className="w-full px-3.5 py-2 rounded-lg bg-[#1F2937] border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
            />
            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setShowAddList(false)}
                className="px-3 py-1.5 text-xs font-semibold rounded-lg hover:bg-[#1F2937] text-[#9CA3AF] hover:text-white"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-3.5 py-1.5 bg-[#6366F1] hover:bg-[#4F46E5] text-white text-xs font-bold rounded-lg transition-colors"
              >
                Add Column
              </button>
            </div>
          </form>
        ) : (
          <button
            onClick={() => {
              if (isViewer) {
                triggerViewerWarning("Viewers cannot add board columns.");
                return;
              }
              setShowAddList(true);
            }}
            className="w-72 flex items-center justify-center gap-2 py-4 rounded-2xl border border-dashed border-[#1F2937] hover:border-[#6366F1]/40 bg-[#111827]/10 hover:bg-[#111827]/30 text-[#9CA3AF] hover:text-[#6366F1] text-sm font-bold transition-all shrink-0 font-outfit"
          >
            <Plus className="h-4 w-4" />
            <span>Create Column</span>
          </button>
        )}
      </div>

      {/* Cursors co-presence mapping */}
      <CursorOverlay />

      {/* Details drawer inspect panel */}
      {selectedCardId && (
        <CardDrawer
          cardId={selectedCardId}
          onClose={() => setSelectedCardId(null)}
          user={user}
          triggerViewerWarning={triggerViewerWarning}
        />
      )}
    </div>
  );
};
