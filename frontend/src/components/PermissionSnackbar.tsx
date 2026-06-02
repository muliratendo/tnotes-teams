import React, { useEffect } from 'react';
import { ShieldAlert } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface PermissionSnackbarProps {
  isOpen: boolean;
  message: string;
  onClose: () => void;
}

export const PermissionSnackbar: React.FC<PermissionSnackbarProps> = ({ isOpen, message, onClose }) => {
  useEffect(() => {
    if (isOpen) {
      const timer = setTimeout(onClose, 4000);
      return () => clearTimeout(timer);
    }
  }, [isOpen, onClose]);

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0, y: 50, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: 20, scale: 0.95 }}
          transition={{ type: 'spring', damping: 25, stiffness: 350 }}
          className="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-4 py-3.5 bg-[#78350F] border border-[#F59E0B]/20 text-[#FEF3C7] rounded-xl shadow-2xl max-w-sm elevation-4"
        >
          <div className="bg-[#F59E0B]/20 p-2 rounded-lg text-[#F59E0B]">
            <ShieldAlert className="h-5 w-5" />
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-bold text-white leading-tight">Permission Restrained</span>
            <span className="text-xs text-[#FEF3C7] mt-0.5 leading-snug">{message}</span>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
