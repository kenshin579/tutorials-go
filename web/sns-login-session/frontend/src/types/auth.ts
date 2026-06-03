export interface User {
  id: number;
  email: string;
  name: string;
  avatar_url: string;
  provider: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface AuthResponse {
  tokens: TokenPair;
  user: User;
}
