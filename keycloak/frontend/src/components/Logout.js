import React from 'react'
import { logout } from '../services/keycloak'

export default function Logout() {
  return (
    <button className="btn btn-secondary" onClick={() => logout()}>
      로그아웃
    </button>
  )
}
