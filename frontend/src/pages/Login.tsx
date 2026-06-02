import React, { useState } from 'react';
import { Columns, Mail, Lock, User, Loader2 } from 'lucide-react';

interface LoginProps {
  onAuthSuccess: (token: string, user: any) => void;
}

export const Login: React.FC<LoginProps> = ({ onAuthSuccess }) => {
  const [isRegister, setIsRegister] = useState(false);
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const endpoint = isRegister ? '/api/auth/register' : '/api/auth/login';
    const payload = isRegister ? { email, username, password } : { email, password };

    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      const json = await res.json();
      if (!res.ok) {
        setError(json.error || 'Authentication failed');
        return;
      }

      // Store JWT token locally
      localStorage.setItem('jwt_token', json.data.token);
      onAuthSuccess(json.data.token, json.data.user);
    } catch (err) {
      setError('Connection to backend failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="h-screen w-screen bg-[#0B0F19] flex items-center justify-center font-outfit text-white p-4 relative overflow-hidden">
      {/* Decorative gradient glowing spheres */}
      <div className="absolute top-1/4 left-1/4 h-[350px] w-[350px] bg-[#6366F1]/10 rounded-full blur-[80px] pointer-events-none" />
      <div className="absolute bottom-1/4 right-1/4 h-[350px] w-[350px] bg-[#F59E0B]/5 rounded-full blur-[80px] pointer-events-none" />

      {/* Main card */}
      <div className="w-full max-w-md p-8 bg-[#111827]/70 backdrop-blur-xl border border-[#1F2937] rounded-3xl shadow-2xl space-y-6 relative z-10">
        
        {/* Brand header */}
        <div className="flex flex-col items-center text-center space-y-2">
          <div className="h-12 w-12 bg-[#6366F1]/10 rounded-2xl flex items-center justify-center border border-[#6366F1]/20 text-[#6366F1] shadow-inner mb-2">
            <Columns className="h-6 w-6" />
          </div>
          <h2 className="text-2xl font-bold tracking-tight text-white font-outfit">
            {isRegister ? 'Register Team Account' : 'TNotes Teams Login'}
          </h2>
          <p className="text-xs text-[#9CA3AF] max-w-xs leading-relaxed">
            {isRegister 
              ? 'Join your team and coordinate sprint tasks collaborative with co-presence.' 
              : 'Sign in to access your Kanban boards, quick notes, and sync offline mutations.'}
          </p>
        </div>

        {error && (
          <div className="p-3.5 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl font-medium text-center">
            {error}
          </div>
        )}

        {/* Auth form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-[#9CA3AF]">Email Address</label>
            <div className="relative">
              <Mail className="absolute left-3.5 top-3.5 h-4.5 w-4.5 text-[#9CA3AF]" />
              <input
                type="email"
                required
                value={email}
                onChange={e => setEmail(e.target.value)}
                className="w-full pl-11 pr-4 py-3 rounded-xl bg-[#1F2937]/50 border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
                placeholder="you@workplace.com"
              />
            </div>
          </div>

          {isRegister && (
            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-[#9CA3AF]">Username</label>
              <div className="relative">
                <User className="absolute left-3.5 top-3.5 h-4.5 w-4.5 text-[#9CA3AF]" />
                <input
                  type="text"
                  required
                  value={username}
                  onChange={e => setUsername(e.target.value)}
                  className="w-full pl-11 pr-4 py-3 rounded-xl bg-[#1F2937]/50 border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
                  placeholder="e.g. sarah, alex"
                />
              </div>
            </div>
          )}

          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-[#9CA3AF]">Password</label>
            <div className="relative">
              <Lock className="absolute left-3.5 top-3.5 h-4.5 w-4.5 text-[#9CA3AF]" />
              <input
                type="password"
                required
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="w-full pl-11 pr-4 py-3 rounded-xl bg-[#1F2937]/50 border border-[#374151] text-sm text-white focus:outline-none focus:border-[#6366F1]"
                placeholder="••••••••"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3.5 bg-[#6366F1] hover:bg-[#4F46E5] text-white text-sm font-bold rounded-xl transition-all shadow-lg shadow-[#6366F1]/10 flex items-center justify-center gap-2 disabled:opacity-50"
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                <span>Processing...</span>
              </>
            ) : (
              <span>{isRegister ? 'Register' : 'Login'}</span>
            )}
          </button>
        </form>

        {/* Toggle login / register */}
        <div className="text-center pt-2">
          <button
            onClick={() => {
              setIsRegister(!isRegister);
              setError('');
            }}
            className="text-xs font-semibold text-[#9CA3AF] hover:text-[#6366F1] transition-colors"
          >
            {isRegister ? 'Already have an account? Sign in' : "Don't have an account? Register here"}
          </button>
        </div>

      </div>
    </div>
  );
};
