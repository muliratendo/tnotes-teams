import React, { useEffect, useState } from 'react';
import { useBoardStore } from '../stores/useBoardStore';
import { apiRequest } from '../services/api';
import { 
  Columns, 
  Wifi, 
  WifiOff, 
  LogOut, 
  Briefcase, 
  Plus, 
  User, 
  Loader2 
} from 'lucide-react';

export interface WorkspaceSummary {
  id: string;
  name: string;
  description?: string;
  role?: string;
}

interface AppShellProps {
  children: React.ReactNode;
  onLogout: () => void;
  user: { id: string; username: string; email: string; role?: string };
  onRoleChange?: (role: string) => void;
}

export const AppShell: React.FC<AppShellProps> = ({ children, onLogout, user, onRoleChange }) => {
  const { isOnline, toggleSyncMode, syncQueueLength } = useBoardStore();
  const [boards, setBoards] = useState<{ id: string; name: string }[]>([]);
  const [activeWorkspace, setActiveWorkspace] = useState<WorkspaceSummary | null>(null);
  const [activeBoardId, setActiveBoardId] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [showCreateBoard, setShowCreateBoard] = useState(false);
  const [newBoardName, setNewBoardName] = useState('');

  const initializeBoard = useBoardStore(state => state.initializeBoard);

  useEffect(() => {
    fetchWorkspaces();
  }, []);

  const fetchWorkspaces = async () => {
    try {
      let list = await apiRequest<WorkspaceSummary[]>('/api/workspaces');
      if (list.length === 0) {
        const created = await apiRequest<WorkspaceSummary>('/api/workspaces', {
          method: 'POST',
          body: JSON.stringify({ name: 'My Workspace', description: 'Default team sandbox' }),
        });
        list = [created];
      }
      const first = list[0];
      setActiveWorkspace(first);
      onRoleChange?.(first.role || 'admin');
      fetchBoards(first.id);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchBoards = async (wsId: string) => {
    setLoading(true);
    try {
      const data = await apiRequest<{ id: string; name: string }[]>(`/api/workspaces/${wsId}/boards`);
      setBoards(data);
      if (data.length > 0 && !activeBoardId) {
        handleBoardSelect(data[0].id);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleBoardSelect = async (boardId: string) => {
    setActiveBoardId(boardId);
    try {
      const data = await apiRequest<Record<string, unknown>>(`/api/boards/${boardId}`);
      await initializeBoard(data);
    } catch (err) {
      console.error(err);
    }
  };

  const handleCreateBoard = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newBoardName.trim() || !activeWorkspace) return;

    try {
      const newBoard = await apiRequest<{ id: string; name: string }>(
        `/api/workspaces/${activeWorkspace.id}/boards`,
        {
          method: 'POST',
          body: JSON.stringify({ name: newBoardName, description: 'Team Kanban board' }),
        }
      );
      setBoards([...boards, newBoard]);
      handleBoardSelect(newBoard.id);
      setNewBoardName('');
      setShowCreateBoard(false);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-[#0B0F19] text-[#F3F4F6]">
      {/* Sidebar Navigation */}
      <aside className="w-64 border-r border-[#1F2937] bg-[#111827] flex flex-col justify-between">
        <div>
          {/* Brand Header */}
          <div className="p-6 border-b border-[#1F2937] flex items-center gap-3">
            <Columns className="h-6 w-6 text-[#6366F1]" />
            <span className="text-xl font-bold tracking-tight text-white font-outfit">TNotes Teams</span>
          </div>

          {/* Workspace Switcher */}
          <div className="p-4 border-b border-[#1F2937]">
            <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-[#9CA3AF] mb-2">
              <Briefcase className="h-3.5 w-3.5" />
              <span>Workspace</span>
            </div>
            {activeWorkspace && (
              <div className="px-3 py-2.5 rounded-lg bg-[#1F2937] text-white font-medium flex items-center justify-between border border-[#374151]">
                <span>{activeWorkspace.name}</span>
                <span className="text-[10px] bg-[#6366F1] px-2 py-0.5 rounded-full font-bold uppercase">
                  {user.role || 'Admin'}
                </span>
              </div>
            )}
          </div>

          {/* Boards List */}
          <div className="p-4">
            <div className="flex items-center justify-between text-xs font-semibold uppercase tracking-wider text-[#9CA3AF] mb-3">
              <span>Boards</span>
              <button 
                onClick={() => setShowCreateBoard(true)}
                className="hover:text-white transition-colors"
                title="Create Board"
              >
                <Plus className="h-4 w-4" />
              </button>
            </div>

            {loading ? (
              <div className="flex justify-center p-4">
                <Loader2 className="h-5 w-5 animate-spin text-[#6366F1]" />
              </div>
            ) : (
              <nav className="space-y-1">
                {boards.map(b => (
                  <button
                    key={b.id}
                    onClick={() => handleBoardSelect(b.id)}
                    className={`w-full text-left px-3 py-2 rounded-lg font-medium text-sm transition-all flex items-center justify-between ${
                      activeBoardId === b.id 
                        ? 'bg-[#6366F1] text-white shadow-lg shadow-[#6366F1]/10' 
                        : 'text-[#9CA3AF] hover:bg-[#1F2937] hover:text-white'
                    }`}
                  >
                    <span>{b.name}</span>
                  </button>
                ))}
              </nav>
            )}
          </div>
        </div>

        {/* User Footer Profile & Offline Toggle */}
        <div className="p-4 border-t border-[#1F2937] space-y-4">
          {/* Offline switch */}
          <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-[#1F2937] border border-[#374151]">
            <div className="flex items-center gap-2">
              {isOnline ? (
                <Wifi className="h-4 w-4 text-[#10B981]" />
              ) : (
                <WifiOff className="h-4 w-4 text-[#EF4444]" />
              )}
              <span className="text-xs font-semibold">
                {isOnline ? 'Online Sync' : 'Offline Mode'}
              </span>
            </div>
            <button
              onClick={toggleSyncMode}
              className={`w-10 h-5 rounded-full p-0.5 transition-colors duration-200 focus:outline-none ${
                isOnline ? 'bg-[#10B981]' : 'bg-[#EF4444]'
              }`}
            >
              <div className={`bg-white w-4 h-4 rounded-full shadow-md transform duration-200 ${
                isOnline ? 'translate-x-5' : 'translate-x-0'
              }`} />
            </button>
          </div>

          {/* Sync status length indicator */}
          {syncQueueLength > 0 && (
            <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-[#78350F] border border-[#F59E0B]/20 text-[#FEF3C7] text-xs font-medium">
              <span>Queue Sync</span>
              <span className="bg-[#F59E0B] text-black px-2 py-0.5 rounded-full font-bold font-mono">
                {syncQueueLength}
              </span>
            </div>
          )}

          {/* Logged in User info */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <div className="h-8 w-8 rounded-full bg-[#374151] flex items-center justify-center text-[#6366F1] border border-[#6366F1]/20">
                <User className="h-4 w-4" />
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-semibold text-white leading-tight">@{user.username}</span>
                <span className="text-[10px] text-[#9CA3AF] leading-none">{user.email}</span>
              </div>
            </div>
            <button 
              onClick={onLogout}
              className="p-1.5 rounded-lg hover:bg-[#1F2937] text-[#9CA3AF] hover:text-white transition-colors"
              title="Logout"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        </div>
      </aside>

      {/* Main Board View Container */}
      <main className="flex-1 flex flex-col overflow-hidden">
        {children}
      </main>

      {/* Create Board Modal */}
      {showCreateBoard && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <form onSubmit={handleCreateBoard} className="w-full max-w-md p-6 bg-[#111827] rounded-xl border border-[#1F2937] shadow-2xl space-y-4">
            <h3 className="text-lg font-bold text-white font-outfit">Create New Kanban Board</h3>
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-[#9CA3AF]">Board Name</label>
              <input
                type="text"
                required
                value={newBoardName}
                onChange={e => setNewBoardName(e.target.value)}
                className="w-full px-3.5 py-2 rounded-lg bg-[#1F2937] border border-[#374151] text-white focus:outline-none focus:border-[#6366F1]"
                placeholder="e.g. Sprint Board, Product Backlog"
              />
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <button
                type="button"
                onClick={() => setShowCreateBoard(false)}
                className="px-4 py-2 text-sm font-semibold rounded-lg hover:bg-[#1F2937] text-[#9CA3AF] hover:text-white transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-4 py-2 text-sm font-semibold bg-[#6366F1] hover:bg-[#4F46E5] text-white rounded-lg transition-colors"
              >
                Create
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
};
