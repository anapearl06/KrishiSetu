import { useState } from 'react'

import loginBg from '../assets/login-bg.png'
import RoleTabs from './RoleTabs'
import TextField from './TextField'
import Icon from './Icon'

function validate(phone, password) {
  const errors = {}

  if (!/^\d{10}$/.test(phone)) {
    errors.phone = 'Enter a valid 10-digit phone number.'
  }

  if (!password) {
    errors.password = 'Password is required.'
  }

  return errors
}

function LoginScreen() {
  const [role, setRole] = useState('farmer')
  const [phone, setPhone] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [attempted, setAttempted] = useState(false)

  const errors = validate(phone, password)
  const handleChangePhone = (event) => {
    setPhone(event.target.value.replace(/\D/g, '').slice(0, 10))
  }

  const handleSubmit = (event) => {
    event.preventDefault()
    setAttempted(true)
  }

  return (
    <div className="ks-login-page">
      <div className="ks-login-bg">
        <img className="ks-login-bg-photo" src={loginBg} alt="" />
        <div className="ks-login-bg-overlay" />
      </div>

      <div className="ks-login-card-wrap">
        <header className="ks-login-header">
          <h1 className="ks-login-brand">KrishiSetu</h1>
          <p className="ks-login-tagline">Agriculture Simplified</p>
        </header>

        <div className="ks-login-card">
          <RoleTabs active={role} onChange={setRole} />

          <div className="ks-login-form-body">
            <h2 className="ks-login-welcome">Welcome Back</h2>

            <form onSubmit={handleSubmit} noValidate>
              <div className="ks-login-fields">
                <TextField
                  id="phone"
                  label="Phone Number"
                  value={phone}
                  onChange={handleChangePhone}
                  placeholder="Enter your 10-digit number"
                  type="tel"
                  inputMode="numeric"
                  autoComplete="tel-national"
                  prefix="+91"
                  error={attempted && Boolean(errors.phone)}
                  errorMessage={errors.phone}
                  trailing={
                    attempted && errors.phone ? (
                      <Icon name="error" tone="error" />
                    ) : undefined
                  }
                />

                <TextField
                  id="password"
                  label="Password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  placeholder="Enter your password"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="current-password"
                  error={attempted && Boolean(errors.password)}
                  errorMessage={errors.password}
                  labelExtra={
                    <a
                      className="ks-link"
                      href="#"
                      onClick={(event) => event.preventDefault()}
                    >
                      Forgot?
                    </a>
                  }
                  trailing={
                    attempted && errors.password ? (
                      <Icon name="error" tone="error" />
                    ) : (
                      <button
                        type="button"
                        className="ks-icon-btn"
                        aria-label={
                          showPassword ? 'Hide password' : 'Show password'
                        }
                        onClick={() => setShowPassword((value) => !value)}
                      >
                        <Icon name={showPassword ? 'visibility_off' : 'visibility'} />
                      </button>
                    )
                  }
                />
              </div>

              <div className="ks-submit-wrap">
                <button type="submit" className="ks-btn ks-btn-primary">
                  Login
                </button>
              </div>
            </form>

            <p className="ks-login-subtext">
              Don&apos;t have an account?{' '}
              <a
                className="ks-link ks-link-strong"
                href="#"
                onClick={(event) => event.preventDefault()}
              >
                Register
              </a>
            </p>
          </div>
        </div>

        <div className="ks-login-support">
          <a
            className="ks-support-link"
            href="#"
            onClick={(event) => event.preventDefault()}
          >
            <Icon name="help" size={20} />
            Need Help? Contact Support
          </a>
        </div>
      </div>
    </div>
  )
}

export default LoginScreen