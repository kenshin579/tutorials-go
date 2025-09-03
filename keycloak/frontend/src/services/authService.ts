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

  // Authorization Code Flow 로그인 (최소 구현)
  public initiateLogin(): void {
    const authUrl = new URL(`${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/auth`);
    authUrl.searchParams.append('client_id', this.config.clientId);
    authUrl.searchParams.append('redirect_uri', window.location.origin + '/callback');
    authUrl.searchParams.append('response_type', 'code');
    authUrl.searchParams.append('scope', 'openid profile email');
    
    window.location.href = authUrl.toString();
  }



  // Authorization Code를 토큰으로 교환
  public async handleCallback(code: string): Promise<boolean> {
    try {
      // 이미 인증된 상태라면 성공으로 처리
      if (this.isAuthenticated()) {
        console.log('Already authenticated, skipping token exchange');
        return true;
      }

      console.log('Starting token exchange with code:', code.substring(0, 10) + '...');
      
      const tokenUrl = `${this.config.url}/realms/${this.config.realm}/protocol/openid-connect/token`;
      
      const formData = new URLSearchParams();
      formData.append('grant_type', 'authorization_code');
      formData.append('client_id', this.config.clientId);
      formData.append('code', code);
      formData.append('redirect_uri', window.location.origin + '/callback');

      const response = await fetch(tokenUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Token exchange failed:', response.status, errorText);
        
        // Log error for debugging
        
        throw new Error(`Token exchange failed: ${response.status} - ${errorText}`);
      }

      const tokenResponse: TokenResponse = await response.json();
      console.log('Token exchange successful, saving tokens');
      this.saveTokensToStorage(tokenResponse);
      
      // Token exchange successful
      
      return true;
    } catch (error) {
      console.error('Callback handling error:', error);
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
