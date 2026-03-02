import { useState, useEffect } from 'react';
import { getCurrentUser, signIn, signOut as amplifySignOut, fetchAuthSession, AuthUser } from 'aws-amplify/auth';

function useAuth() {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [hasProfile, setHasProfile] = useState<boolean | null>(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const getSessionToken = async (): Promise<string | null> => {
    try {
      const session = await fetchAuthSession();
      return session.tokens?.accessToken?.toString() ?? session.tokens?.idToken?.toString() ?? null;
    } catch {
      return null;
    }
  };

  const handleSignOut = async () => {
    try {
      await amplifySignOut();
    } finally {
      setUser(null);
      setIsAuthenticated(false);
      setHasProfile(null);
    }
  };

  const checkAuth = async () => {
    try {
      const currentUser = await getCurrentUser();
      const token = await getSessionToken();
      if (!token) {
        throw new Error('No valid auth token');
      }

      setUser(currentUser);
      setIsAuthenticated(true);
      await checkProfile(currentUser.userId, token);
    } catch {
      await handleSignOut();
    } finally {
      setLoading(false);
    }
  };

  const checkProfile = async (userId: string, token: string) => {
    try {
      const response = await fetch(`/users?user_id=${encodeURIComponent(userId)}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.ok) {
        setHasProfile(true);
      } else {
        if (response.status === 401) {
          await handleSignOut();
          return;
        }
        setHasProfile(false);
      }
    } catch {
      setHasProfile(false);
    }
  }

  const getToken = async (): Promise<string | null> => {
    return getSessionToken();
  };

  return { user, isAuthenticated, hasProfile, loading, signIn, signOut: handleSignOut, getToken, checkAuth };
}
export default useAuth;