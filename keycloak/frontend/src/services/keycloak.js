import Keycloak from 'keycloak-js'

// 환경 변수로 설정할 수 있도록 기본값 제공
const kcConfig = {
  url: import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:8080',
  realm: import.meta.env.VITE_KEYCLOAK_REALM || 'myrealm',
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'myclient',
}

let keycloak

export function getKeycloak() {
  return keycloak
}

export async function initKeycloak() {
  if (!keycloak) {
    keycloak = new Keycloak(kcConfig)
  }

  // silent check를 위한 iframe 설정 (옵션)
  const authenticated = await keycloak.init({
    onLoad: 'check-sso',
    checkLoginIframe: true,
    pkceMethod: 'S256',
    silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
  })

  // 토큰 자동 갱신
  if (authenticated) {
    scheduleTokenRefresh()
  }

  return authenticated
}

function scheduleTokenRefresh() {
  if (!keycloak) return
  const refresh = async () => {
    try {
      await keycloak.updateToken(30) // 30초 이내 만료 시 갱신
    } catch (e) {
      console.warn('Token refresh failed, logging out', e)
      keycloak.logout({ redirectUri: window.location.origin })
    }
  }
  // 10초마다 갱신 시도 (가벼운 주기)
  const interval = setInterval(refresh, 10 * 1000)
  // 페이지 이탈 시 정리
  window.addEventListener('beforeunload', () => clearInterval(interval))
}

export function login() {
  if (!keycloak) return
  return keycloak.login({ redirectUri: window.location.origin })
}

export function logout() {
  if (!keycloak) return
  return keycloak.logout({ redirectUri: window.location.origin })
}

export function getToken() {
  return keycloak?.token || null
}
