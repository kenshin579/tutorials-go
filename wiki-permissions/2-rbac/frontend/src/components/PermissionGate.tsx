import type { ReactNode } from 'react';
import { useAuth } from '../auth/AuthContext';

interface Props {
  permission: string;
  children: ReactNode;
}

// PermissionGate는 사용자가 want permission을 가졌을 때만 children을 렌더링한다.
// 사용 예: <PermissionGate permission="pages:edit"><button>편집</button></PermissionGate>
//
// 1편(ACL)에는 이 컴포넌트가 없었다 — 편집 버튼은 모두에게 보였고 서버가 403으로 거부했다.
// RBAC에서는 login 응답에 사용자 permissions가 함께 오므로 사전 게이팅이 가능해졌다.
export default function PermissionGate({ permission, children }: Props) {
  const { user } = useAuth();
  if (!user?.permissions.includes(permission)) return null;
  return <>{children}</>;
}
