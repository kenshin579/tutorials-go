import { useAuth } from './AuthContext';

export function usePermission() {
  const { user } = useAuth();

  const hasPermission = (permission: string): boolean => {
    return user?.permissions?.includes(permission) ?? false;
  };

  const hasRole = (role: string): boolean => {
    return user?.roles?.includes(role) ?? false;
  };

  const isOwner = (ownerId: number): boolean => {
    return user?.id === ownerId;
  };

  return { hasPermission, hasRole, isOwner };
}
