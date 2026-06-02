import React from 'react';
import { usePresenceStore } from '../stores/usePresenceStore';
import { MousePointer } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

export const CursorOverlay: React.FC = () => {
  const cursors = usePresenceStore(state => state.cursors);

  return (
    <div className="absolute inset-0 pointer-events-none z-40 overflow-hidden">
      <AnimatePresence>
        {Object.values(cursors).map(cursor => (
          <motion.div
            key={cursor.userId}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ 
              opacity: 1, 
              scale: 1,
              x: cursor.x, 
              y: cursor.y 
            }}
            exit={{ opacity: 0 }}
            transition={{ type: 'spring', damping: 30, stiffness: 280, mass: 0.8 }}
            className="absolute flex items-start gap-1"
          >
            {/* Elegant tinted mouse indicator pointer */}
            <MousePointer className="h-5 w-5 fill-[#6366F1] text-[#6366F1] drop-shadow-md transform -rotate-90" />
            
            {/* Active peer nickname label */}
            <div className="px-2 py-0.5 bg-[#6366F1] text-white text-[10px] font-bold rounded shadow-lg border border-white/10 uppercase tracking-wider">
              {cursor.username}
            </div>
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
};
