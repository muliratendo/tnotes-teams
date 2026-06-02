import React, { useEffect, useState } from 'react';
import { apiRequest } from './services/api';
import { Login } from './pages/Login';
import { AppShell } from './components/AppShell';
import { BoardCanvas } from './components/BoardCanvas';
import { PermissionSnackbar } from './components/PermissionSnackbar';
import { Loader2 } from 'lucide-react';

export const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(localStorage.getItem('jwt_token'));
  const [user, setUser] = useState<{ id: string; username: string; email: string; role?: string } | null>(null);
  const [loading, setLoading] = useState(true);

  // Permission warning states
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMsg, setSnackbarMsg] = useState('');

  useEffect(() => {
    if (token) {
      verifyToken();
    } else {
      setLoading(false);
    }
  }, [token]);

  const verifyToken = async () => {
    try {
      const data = await apiRequest<{ id: string; username: string; email: string }>('/api/auth/me');
      setUser({ ...data, role: 'admin' });
    } catch {
      localStorage.removeItem('jwt_token');
      setToken(null);
    } finally {
      setLoading(false);
    }
  };

  const handleAuthSuccess = (newToken: string, authenticatedUser: any) => {
    setToken(newToken);
    setUser(authenticatedUser);
  };

  const handleLogout = () => {
    localStorage.removeItem('jwt_token');
    setToken(null);
    setUser(null);
  };

  const triggerViewerWarning = (message: string) => {
    setSnackbarMsg(message);
    setSnackbarOpen(true);
  };

  if (loading) {
    return (
      <div className="h-screen w-screen bg-[#0B0F19] flex flex-col items-center justify-center text-white space-y-4 font-outfit">
        <Loader2 className="h-10 w-10 text-[#6366F1] animate-spin" />
        <span className="text-xs font-semibold text-[#9CA3AF]">Initializing TNotes Teams Workspace...</span>
      </div>
    );
  }

  if (!token || !user) {
    return <Login onAuthSuccess={handleAuthSuccess} />;
  }

  return (
    <>
      <AppShell
        user={user}
        onLogout={handleLogout}
        onRoleChange={(role) => setUser((u) => (u ? { ...u, role } : u))}
      >
        <BoardCanvas user={user} triggerViewerWarning={triggerViewerWarning} />
      </AppShell>

      {/* Permission alert toasts */}
      <PermissionSnackbar
        isOpen={snackbarOpen}
        message={snackbarMsg}
        onClose={() => setSnackbarOpen(false)}
      />
    </>
  );
};

export default App;
