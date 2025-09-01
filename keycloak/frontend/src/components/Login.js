import React from 'react'
import { login } from '../services/keycloak'

export default function Login() {
  return (
    <button className="btn" onClick={() => login()}>
      로그인
    </button>
  )
}
