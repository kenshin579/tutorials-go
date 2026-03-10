import type { ReactNode } from 'react';
import { usePermission } from '../auth/usePermission';

interface PermissionGateProps {
  children: ReactNode;
  permission: string;
  ownerId?: number;
}

export default function PermissionGate({ children, permission, ownerId }: PermissionGateProps) {
  const { hasPermission, hasRole, isOwner } = usePermission();

  if (!hasPermission(permission)) {
    return null;
  }

  // If ownerId check is needed
  if (ownerId !== undefined) {
    if (hasRole('admin')) return <>{children}</>;
    if (!isOwner(ownerId)) return null;
  }

  return <>{children}</>;
}
