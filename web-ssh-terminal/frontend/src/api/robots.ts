export interface Robot {
  id: string;
  name: string;
  host: string;
  port: number;
  username: string;
  authType: 'password' | 'privateKey';
  description?: string;
  isOnline: boolean;
}

export async function fetchRobots(): Promise<Robot[]> {
  const res = await fetch('/api/robots');
  if (!res.ok) {
    throw new Error(`Failed to fetch robots: ${res.statusText}`);
  }
  return res.json();
}
