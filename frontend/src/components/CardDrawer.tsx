import React, { useState, useEffect } from 'react';
import { useBoardStore } from '../stores/useBoardStore';
import { usePresenceStore } from '../stores/usePresenceStore';
import { 
  X, 
  CheckSquare, 
  MessageSquare, 
  FileText, 
  Calendar, 
  Plus, 
  Trash2,
  AlertCircle
} from 'lucide-react';
import { motion } from 'framer-motion';

interface CardDrawerProps {
  cardId: string;
  onClose: () => void;
  user: { id: string; username: string; email: string; role?: string };
  triggerViewerWarning: (message: string) => void;
}

export const CardDrawer: React.FC<CardDrawerProps> = ({ cardId, onClose, user, triggerViewerWarning }) => {
  const isViewer = user.role === 'viewer';
  
  // Board states
  const cards = useBoardStore(state => state.cards);
  const subtasksMap = useBoardStore(state => state.subtasks);
  const notesMap = useBoardStore(state => state.quickNotes);
  const commentsMap = useBoardStore(state => state.comments);

  // Store mutators
  const updateCardDetails = useBoardStore(state => state.updateCardDetails);
  const deleteCard = useBoardStore(state => state.deleteCard);
  const addSubtask = useBoardStore(state => state.addSubtask);
  const toggleSubtask = useBoardStore(state => state.toggleSubtask);
  const addQuickNote = useBoardStore(state => state.addQuickNote);
  const addComment = useBoardStore(state => state.addComment);

  // Real-time states
  const typingUsers = usePresenceStore(state => state.typingUsers);

  // Local drawer states
  const card = cards.find(c => c.id === cardId);
  const subtasks = subtasksMap[cardId] || [];
  const notes = notesMap[cardId] || [];
  const comments = commentsMap[cardId] || [];
  
  const [description, setDescription] = useState(card?.description || '');
  const [newSubtask, setNewSubtask] = useState('');
  const [newNote, setNewNote] = useState('');
  const [newComment, setNewComment] = useState('');

  useEffect(() => {
    if (card) {
      setDescription(card.description);
    }
  }, [card]);

  if (!card) return null;

  // Calculate Progress Percent Bar
  const totalSubtasks = subtasks.length;
  const completedSubtasks = subtasks.filter(s => s.is_completed === 1).length;
  const progressPercent = totalSubtasks > 0 ? Math.round((completedSubtasks / totalSubtasks) * 100) : 0;

  const handleDescBlur = () => {
    if (isViewer) {
      setDescription(card.description);
      triggerViewerWarning("Viewers cannot update card descriptions.");
      return;
    }
    updateCardDetails(cardId, undefined, description);
  };

  const handleAddSubtaskSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewer) {
      triggerViewerWarning("Viewers cannot add subtasks.");
      return;
    }
    if (!newSubtask.trim()) return;
    addSubtask(cardId, newSubtask);
    setNewSubtask('');
  };

  const handleToggleSub = (subId: string, currentStatus: number) => {
    if (isViewer) {
      triggerViewerWarning("Viewers cannot toggle checklist subtasks.");
      return;
    }
    toggleSubtask(subId, currentStatus === 0);
  };

  const handleAddNoteSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewer) {
      triggerViewerWarning("Viewers cannot append micro-notes.");
      return;
    }
    if (!newNote.trim()) return;
    addQuickNote(cardId, newNote);
    setNewNote('');
  };

  const handleAddCommentSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewer) {
      triggerViewerWarning("Viewers cannot write comments.");
      return;
    }
    if (!newComment.trim()) return;
    addComment(cardId, newComment, user.id, user.username);
    setNewComment('');
  };

  const handleDeleteCard = () => {
    if (isViewer) {
      triggerViewerWarning("Viewers cannot delete Kanban cards.");
      return;
    }
    deleteCard(cardId);
    onClose();
  };

  return (
    <motion.div
      initial={{ x: '100%' }}
      animate={{ x: 0 }}
      exit={{ x: '100%' }}
      transition={{ type: 'spring', damping: 26, stiffness: 220 }}
      className="absolute top-0 right-0 h-full w-full max-w-lg border-l border-[#1F2937] bg-[#111827] shadow-2xl flex flex-col z-30 elevation-5 font-outfit text-white"
    >
      {/* Header Info */}
      <div className="p-6 border-b border-[#1F2937] flex items-center justify-between">
        <div className="flex items-center gap-3">
          <AlertCircle className="h-5 w-5 text-[#6366F1]" />
          <span className="text-lg font-bold">Card Inspector</span>
        </div>
        <div className="flex items-center gap-2">
          <button 
            onClick={handleDeleteCard}
            className="p-2 rounded-lg hover:bg-red-500/10 text-red-400 hover:text-red-500 transition-colors"
            title="Delete Card"
          >
            <Trash2 className="h-4 w-4" />
          </button>
          <button 
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-[#1F2937] text-[#9CA3AF] hover:text-white transition-colors"
          >
            <X className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Drawer scroll content body */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Title */}
        <div className="space-y-1">
          <h2 className="text-2xl font-bold tracking-tight text-white">{card.title}</h2>
          <span className="text-xs text-[#9CA3AF] font-medium flex items-center gap-1.5">
            <Calendar className="h-3.5 w-3.5" />
            Created {new Date(card.created_at).toLocaleDateString()}
          </span>
        </div>

        {/* Description */}
        <div className="space-y-2">
          <span className="text-xs font-semibold uppercase tracking-wider text-[#9CA3AF] flex items-center gap-2">
            <FileText className="h-4 w-4 text-[#6366F1]" />
            Description
          </span>
          <textarea
            value={description}
            onChange={e => setDescription(e.target.value)}
            onBlur={handleDescBlur}
            disabled={isViewer}
            rows={3}
            className="w-full px-3.5 py-2.5 rounded-xl bg-[#1F2937] border border-[#374151] text-white focus:outline-none focus:border-[#6366F1] font-outfit text-sm resize-none disabled:opacity-60"
            placeholder="Add detailed task notes here..."
          />
        </div>

        {/* Progress Bar (Auto computed from checklists) */}
        <div className="space-y-2 p-4 rounded-xl bg-[#1F2937] border border-[#374151]">
          <div className="flex justify-between text-xs font-bold text-[#9CA3AF]">
            <span className="flex items-center gap-2">
              <CheckSquare className="h-4 w-4 text-[#6366F1]" />
              Checklist Progress
            </span>
            <span className="text-[#6366F1]">{progressPercent}%</span>
          </div>
          <div className="w-full bg-[#111827] h-2 rounded-full overflow-hidden border border-[#374151]">
            <div 
              className="bg-gradient-to-r from-[#6366F1] to-[#F59E0B] h-full rounded-full transition-all duration-300"
              style={{ width: `${progressPercent}%` }}
            />
          </div>
        </div>

        {/* Checklist Subtasks */}
        <div className="space-y-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-[#9CA3AF]">
            Subtasks ({completedSubtasks}/{totalSubtasks})
          </span>
          <div className="space-y-2">
            {subtasks.map(s => (
              <div 
                key={s.id}
                onClick={() => handleToggleSub(s.id, s.is_completed)}
                className="flex items-center gap-3 p-3 rounded-lg bg-[#1F2937] hover:bg-[#273549] border border-[#374151] cursor-pointer transition-colors"
              >
                <input
                  type="checkbox"
                  checked={s.is_completed === 1}
                  onChange={() => {}} // handled by div click
                  className="rounded border-[#374151] text-[#6366F1] focus:ring-0 focus:ring-offset-0 pointer-events-none"
                />
                <span className={`text-sm font-medium ${s.is_completed === 1 ? 'line-through text-[#9CA3AF]' : 'text-white'}`}>
                  {s.title}
                </span>
              </div>
            ))}
          </div>

          <form onSubmit={handleAddSubtaskSubmit} className="flex gap-2">
            <input
              type="text"
              value={newSubtask}
              onChange={e => setNewSubtask(e.target.value)}
              disabled={isViewer}
              className="flex-1 px-3.5 py-1.5 rounded-lg bg-[#1F2937] border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
              placeholder="Add subtask checklist item..."
            />
            <button
              type="submit"
              disabled={isViewer}
              className="p-2 rounded-lg bg-[#6366F1] hover:bg-[#4F46E5] text-white transition-colors disabled:opacity-50"
            >
              <Plus className="h-4 w-4" />
            </button>
          </form>
        </div>

        {/* Embedded Quick Notes */}
        <div className="space-y-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-[#9CA3AF] flex items-center gap-2">
            <FileText className="h-4 w-4 text-[#F59E0B]" />
            Quick Notes
          </span>
          <div className="space-y-2">
            {notes.map(n => (
              <div key={n.id} className="p-3 bg-[#FEF3C7]/5 border border-[#F59E0B]/20 rounded-xl text-sm text-[#FEF3C7] leading-relaxed">
                {n.content}
              </div>
            ))}
          </div>
          <form onSubmit={handleAddNoteSubmit} className="flex gap-2">
            <input
              type="text"
              value={newNote}
              onChange={e => setNewNote(e.target.value)}
              disabled={isViewer}
              className="flex-1 px-3.5 py-1.5 rounded-lg bg-[#1F2937] border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
              placeholder="Append spontaneous micro-note..."
            />
            <button
              type="submit"
              disabled={isViewer}
              className="p-2 rounded-lg bg-[#F59E0B] hover:bg-[#D97706] text-black transition-colors disabled:opacity-50"
            >
              <Plus className="h-4 w-4" />
            </button>
          </form>
        </div>

        {/* Comments Feed with Typing Indicators */}
        <div className="space-y-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-[#9CA3AF] flex items-center gap-2">
            <MessageSquare className="h-4 w-4 text-[#6366F1]" />
            Comments Section
          </span>

          {/* Typing indicator */}
          {typingUsers[cardId] && typingUsers[cardId].length > 0 && (
            <span className="text-xs text-[#6366F1] font-semibold animate-pulse block">
              {typingUsers[cardId].join(', ')} typing...
            </span>
          )}

          <div className="space-y-3">
            {comments.map(c => (
              <div key={c.id} className="flex gap-2.5 items-start">
                <div className="h-6 w-6 rounded-full bg-[#1F2937] flex items-center justify-center border border-[#6366F1]/10 text-xs text-[#6366F1] font-bold">
                  {c.username[0].toUpperCase()}
                </div>
                <div className="flex-1 p-3 bg-[#1F2937] border border-[#374151] rounded-xl space-y-1">
                  <div className="flex justify-between items-center">
                    <span className="text-xs font-bold text-white">@{c.username}</span>
                    <span className="text-[10px] text-[#9CA3AF]">
                      {new Date(c.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </span>
                  </div>
                  <p className="text-sm text-[#E5E7EB] leading-relaxed">{c.content}</p>
                </div>
              </div>
            ))}
          </div>

          <form onSubmit={handleAddCommentSubmit} className="flex gap-2 pt-2">
            <input
              type="text"
              value={newComment}
              onChange={e => setNewComment(e.target.value)}
              disabled={isViewer}
              className="flex-1 px-3.5 py-1.5 rounded-lg bg-[#1F2937] border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
              placeholder="Write a comment..."
            />
            <button
              type="submit"
              disabled={isViewer}
              className="px-4 py-1.5 bg-[#6366F1] hover:bg-[#4F46E5] text-white text-xs font-bold rounded-lg transition-colors disabled:opacity-50"
            >
              Post
            </button>
          </form>
        </div>
      </div>
    </motion.div>
  );
};
