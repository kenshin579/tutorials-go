interface KeycloakConfig {
  url: string;
  realm: string;
  clientId: string;
}

interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

interface UserInfo {
  sub: string;
  name: string;
  email: string;
  preferred_username: string;
}

class AuthService {
  private config: KeycloakConfig;
  private accessToken: string | null = null;
  private refreshToken: string | null = null;
  private tokenExpiry: number | null = null;

  constructor(config: KeycloakConfig) {
    this.config = config;
    this.loadTokensFromStorage();
  }

  private loadTokensFromStorage() {
    this.accessToken = localStorage.getItem('access_token');
    this.refreshToken = localStorage.getItem('refresh_token');
    const expiry = localStorage.getItem('token_expiry');
    this.tokenExpiry = expiry ? parseInt(expiry) : null;
  }

  private saveTokensToStorage(tokenResponse: TokenResponse) {
    const expiryTime = Date.now() + (tokenResponse.expires_in * 1000);
    
    this.accessToken = tokenResponse.access_token;
    this.refreshToken = tokenResponse.refresh_token;
    this.tokenExpiry = expiryTime;

    localStorage.setItem('access_token', tokenResponse.access_token);
    localStorage.setItem('refresh_token', tokenResponse.refresh_token);
    localStorage.setItem('token_expiry', expiryTime.toString());
  }

  private clearTokensFromStorage() {
    this.accessToken = null;
    this.refreshToken = null;
    this.tokenExpiry = null;

    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('token_expiry');
  }

  public isAuthenticated(): boolean {
    if (!this.accessToken || !this.tokenExpiry) {
      return false;
    }
    
    // Check if token is expired (with 30 second buffer)
    return Date.now() < (this.tokenExpiry - 30000);
  }

  public getAccessToken(): string | null {
    if (this.isAuthenticated()) {
      return this.accessToken;
    }
    return null;
  }

  public async login(username: string, password: string): Promise<boolean> {
    try {
      const tokenUrl = `${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/token`;
      
      const formData = new URLSearchParams();
      formData.append('grant_type', 'password');
      formData.append('client_id', this.config.clientId);
      formData.append('username', username);
      formData.append('password', password);

      const response = await fetch(tokenUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`Login failed: ${response.status}`);
      }

      const tokenResponse: TokenResponse = await response.json();
      this.saveTokensToStorage(tokenResponse);
      
      return true;
    } catch (error) {
      console.error('Login error:', error);
      return false;
    }
  }

  public async refreshAccessToken(): Promise<boolean> {
    if (!this.refreshToken) {
      return false;
    }

    try {
      const tokenUrl = `${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/token`;
      
      const formData = new URLSearchParams();
      formData.append('grant_type', 'refresh_token');
      formData.append('client_id', this.config.clientId);
      formData.append('refresh_token', this.refreshToken);

      const response = await fetch(tokenUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData,
      });

      if (!response.ok) {
        this.clearTokensFromStorage();
        return false;
      }

      const tokenResponse: TokenResponse = await response.json();
      this.saveTokensToStorage(tokenResponse);
      
      return true;
    } catch (error) {
      console.error('Token refresh error:', error);
      this.clearTokensFromStorage();
      return false;
    }
  }

  public async getUserInfo(): Promise<UserInfo | null> {
    const token = this.getAccessToken();
    if (!token) {
      return null;
    }

    try {
      const userInfoUrl = `${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/userinfo`;
      
      const response = await fetch(userInfoUrl, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          // Try to refresh token
          const refreshed = await this.refreshAccessToken();
          if (refreshed) {
            return this.getUserInfo(); // Retry with new token
          }
        }
        throw new Error(`Failed to get user info: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Get user info error:', error);
      return null;
    }
  }

  public async logout(): Promise<void> {
    if (this.refreshToken) {
      try {
        const logoutUrl = `${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/logout`;
        
        const formData = new URLSearchParams();
        formData.append('client_id', this.config.clientId);
        formData.append('refresh_token', this.refreshToken);

        await fetch(logoutUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: formData,
        });
      } catch (error) {
        console.error('Logout error:', error);
      }
    }

    this.clearTokensFromStorage();
  }
}

// Create singleton instance
const authService = new AuthService({
  url: 'http://localhost:8080',
  realm: 'myrealm',
  clientId: 'myclient'
});

export default authService;
export type { UserInfo };
