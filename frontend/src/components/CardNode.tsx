import React from 'react';
import { useBoardStore } from '../stores/useBoardStore';
import { usePresenceStore } from '../stores/usePresenceStore';
import { motion } from 'framer-motion';
import { CheckSquare, Calendar, MessageSquare } from 'lucide-react';
import { type LocalCard } from '../utils/db';

interface CardNodeProps {
  card: LocalCard;
  onClick: () => void;
  onDragStart: (e: React.DragEvent, cardId: string) => void;
  isViewer: boolean;
}

export const CardNode: React.FC<CardNodeProps> = ({ card, onClick, onDragStart, isViewer }) => {
  const subtasks = useBoardStore(state => state.subtasks[card.id] || []);
  const comments = useBoardStore(state => state.comments[card.id] || []);

  const cursors = usePresenceStore(state => state.cursors);

  // Determine if other teammates are actively inspecting this card
  const viewingPeers = Object.values(cursors).filter(c => c.cardId === card.id);

  // Checklist counts
  const totalSubtasks = subtasks.length;
  const completedSubtasks = subtasks.filter(s => s.is_completed === 1).length;

  return (
    <motion.div
      layout
      whileHover={{ scale: 1.015, y: -2 }}
      transition={{ type: 'spring', stiffness: 380, damping: 28 }}
      draggable={!isViewer}
      onDragStartCapture={(e) => onDragStart(e, card.id)}
      onClick={onClick}
      className={`p-4 rounded-xl bg-[#111827]/85 border border-[#1F2937] hover:border-[#6366F1]/50 cursor-pointer shadow-lg transition-colors select-none ${
        isViewer ? 'cursor-default' : 'cursor-grab active:cursor-grabbing'
      }`}
    >
      <div className="space-y-3">
        {/* Title */}
        <h4 className="text-sm font-bold text-white leading-snug tracking-tight">{card.title}</h4>

        {/* Dynamic Progress percent bar */}
        {totalSubtasks > 0 && (
          <div className="space-y-1.5">
            <div className="flex justify-between items-center text-[10px] font-bold text-[#9CA3AF]">
              <span className="flex items-center gap-1">
                <CheckSquare className="h-3 w-3 text-[#6366F1]" />
                Checklist
              </span>
              <span>{completedSubtasks}/{totalSubtasks}</span>
            </div>
            <div className="w-full bg-[#1F2937] h-1.5 rounded-full overflow-hidden border border-[#374151]/20">
              <div 
                className="bg-gradient-to-r from-[#6366F1] to-[#F59E0B] h-full rounded-full transition-all duration-300"
                style={{ width: `${(completedSubtasks / totalSubtasks) * 100}%` }}
              />
            </div>
          </div>
        )}

        {/* Card Metadata & Teammate viewing avatars */}
        <div className="flex justify-between items-center pt-1 border-t border-[#1F2937]/50">
          <div className="flex items-center gap-3 text-[10px] font-semibold text-[#9CA3AF]">
            {card.due_date && (
              <span className="flex items-center gap-1 text-[#F59E0B]">
                <Calendar className="h-3 w-3" />
                {new Date(card.due_date).toLocaleDateString([], { month: 'short', day: 'numeric' })}
              </span>
            )}
            {comments.length > 0 && (
              <span className="flex items-center gap-1">
                <MessageSquare className="h-3 w-3" />
                {comments.length}
              </span>
            )}
          </div>

          {/* Peer viewer highlights (avatars displaying who's hovering/focused) */}
          {viewingPeers.length > 0 && (
            <div className="flex -space-x-1.5 overflow-hidden">
              {viewingPeers.map(peer => (
                <div 
                  key={peer.userId}
                  className="h-5 w-5 rounded-full bg-[#6366F1] border-2 border-[#111827] flex items-center justify-center text-[9px] font-bold text-white shadow-md uppercase"
                  title={`${peer.username} is inspecting`}
                >
                  {peer.username[0]}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
};
