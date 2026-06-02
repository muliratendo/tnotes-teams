import React, { useState } from 'react';
import { CardNode } from './CardNode';
import { Plus, X } from 'lucide-react';
import { type LocalList, type LocalCard } from '../utils/db';

interface ListColumnProps {
  list: LocalList;
  cards: LocalCard[];
  onCardClick: (cardId: string) => void;
  onCardDrop: (e: React.DragEvent, destListId: string) => void;
  onCardDragStart: (e: React.DragEvent, cardId: string) => void;
  onAddCard: (title: string, listId: string) => void;
  isViewer: boolean;
  triggerViewerWarning: (message: string) => void;
}

export const ListColumn: React.FC<ListColumnProps> = ({
  list,
  cards,
  onCardClick,
  onCardDrop,
  onCardDragStart,
  onAddCard,
  isViewer,
  triggerViewerWarning
}) => {
  const [showAddCard, setShowAddCard] = useState(false);
  const [cardTitle, setCardTitle] = useState('');
  const [isOver, setIsOver] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    if (isViewer) return;
    setIsOver(true);
  };

  const handleDragLeave = () => {
    setIsOver(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    if (isViewer) return;
    setIsOver(false);
    onCardDrop(e, list.id);
  };

  const handleAddCardSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewer) {
      triggerViewerWarning("Viewers cannot create card nodes.");
      return;
    }
    if (!cardTitle.trim()) return;
    onAddCard(cardTitle, list.id);
    setCardTitle('');
    setShowAddCard(false);
  };

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={`w-72 flex flex-col max-h-[82vh] rounded-2xl bg-[#111827]/40 border transition-all duration-200 overflow-hidden ${
        isOver 
          ? 'border-[#6366F1] bg-[#111827]/60 shadow-lg shadow-[#6366F1]/5' 
          : 'border-[#1F2937]'
      }`}
    >
      {/* Column Title Header */}
      <div className="p-4 border-b border-[#1F2937]/50 flex justify-between items-center bg-[#111827]/30">
        <h3 className="font-bold text-white tracking-tight">{list.title}</h3>
        <span className="bg-[#1F2937] text-[#9CA3AF] px-2 py-0.5 rounded-full text-xs font-bold font-mono border border-[#374151]/30">
          {cards.length}
        </span>
      </div>

      {/* Cards stack body scroll area */}
      <div className="flex-1 overflow-y-auto p-3 space-y-3">
        {cards.map(c => (
          <CardNode
            key={c.id}
            card={c}
            onClick={() => onCardClick(c.id)}
            onDragStart={onCardDragStart}
            isViewer={isViewer}
          />
        ))}
      </div>

      {/* Footer controls: Add Card */}
      <div className="p-3 bg-[#111827]/20 border-t border-[#1F2937]/35">
        {showAddCard ? (
          <form onSubmit={handleAddCardSubmit} className="space-y-2">
            <input
              type="text"
              required
              autoFocus
              value={cardTitle}
              onChange={e => setCardTitle(e.target.value)}
              placeholder="What needs to be done?"
              className="w-full px-3 py-1.5 rounded-lg bg-[#1F2937] border border-[#374151] text-xs text-white focus:outline-none focus:border-[#6366F1]"
            />
            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setShowAddCard(false)}
                className="p-1 rounded-lg hover:bg-[#1F2937] text-[#9CA3AF] hover:text-white transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
              <button
                type="submit"
                className="px-3 py-1 rounded-lg bg-[#6366F1] hover:bg-[#4F46E5] text-white text-xs font-bold transition-colors"
              >
                Add Card
              </button>
            </div>
          </form>
        ) : (
          <button
            onClick={() => {
              if (isViewer) {
                triggerViewerWarning("Viewers cannot add card nodes.");
                return;
              }
              setShowAddCard(true);
            }}
            className="w-full flex items-center justify-center gap-1.5 py-1.5 rounded-xl border border-dashed border-[#374151] hover:border-[#6366F1]/50 text-[#9CA3AF] hover:text-[#6366F1] text-xs font-bold transition-all"
          >
            <Plus className="h-4 w-4" />
            <span>Create Card</span>
          </button>
        )}
      </div>
    </div>
  );
};
