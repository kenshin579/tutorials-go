import Keycloak from 'keycloak-js';

const keycloak = new (Keycloak as any)({
  url: 'http://localhost:8080',
  realm: 'myrealm',
  clientId: 'myclient'
});

export default keycloak;
