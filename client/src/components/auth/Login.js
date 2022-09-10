import React, { useState } from 'react'
import {
  Button,
  Card,
  CardActions,
  Grid,
  IconButton,
  InputAdornment,
  TextField
} from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { Redirect } from 'react-router-dom'
import { makeStyles } from '@material-ui/core/styles'
import { loginUser } from '../../store/authSlice'
import {
  ErrIncorrectPassword,
  ErrRequiredUserPassword,
  ErrRequiredUserUsername,
  ErrUserNotFound
} from '../../store/errors'
import { Visibility, VisibilityOff } from '@material-ui/icons'

const useStyles = makeStyles({
  card: {
    paddingTop: 24
  },
  container: {
    flexGrow: 1,
    margin: '0.5rem'
  }
})

const Login = () => {
  const classes = useStyles()
  const dispatch = useDispatch()

  const authState = useSelector((state) => state.auth)

  const [loginDetails, setLoginDetails] = useState({
    username: '',
    password: ''
  })
  const [attemptedLogin, setAttemptedLogin] = useState(false)
  const [showPassword, setShowPassword] = useState(false)

  const onChange = (e) => {
    setLoginDetails({
      ...loginDetails,
      [e.target.name]: e.target.value
    })
  }

  const onClick = (e) => {
    e.preventDefault()
    if (!attemptedLogin) {
      setAttemptedLogin(true)
    }

    if (loginDetails.username !== '' && loginDetails.password !== '' && !authState.loading) {
      dispatch(loginUser(loginDetails))
    }
  }

  const onKeyPress = (e) => {
    if (e.key === 'Enter') {
      onClick(e)
    }
  }

  if (authState.isAuthenticated) {
    return <Redirect to='/'/>
  }

  let usernameError = attemptedLogin && loginDetails.username === '' ? ErrRequiredUserUsername : ''
  let passwordError = attemptedLogin && loginDetails.password === '' ? ErrRequiredUserPassword : ''

  if (authState.error !== '') {
    if (authState.error === ErrUserNotFound || authState.error === ErrRequiredUserUsername) {
      usernameError = authState.error
    } else if (authState.error === ErrRequiredUserPassword || authState.error === ErrIncorrectPassword) {
      passwordError = authState.error
    }
  }

  return (
    <Grid container direction='column' alignItems='center' className={classes.card}>
      <Grid item sm={2}>
        <Card elevation={6}>
          <CardActions>
            <Grid
              container alignItems='center' direction='column' className={classes.container}
              spacing={3}
            >
              <Grid item>
                <TextField
                  name='username'
                  placeholder='username'
                  value={loginDetails.username}
                  onChange={onChange}
                  onKeyPress={onKeyPress}
                  error={usernameError !== ''}
                  helperText={usernameError}
                />
              </Grid>
              <Grid item>
                <TextField
                  name='password'
                  placeholder='password'
                  type={showPassword ? 'text' : 'password'}
                  value={loginDetails.password}
                  onChange={onChange}
                  onKeyPress={onKeyPress}
                  error={passwordError !== ''}
                  helperText={passwordError}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position='end'>
                        <IconButton
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          {showPassword ? <Visibility/> : <VisibilityOff/>}
                        </IconButton>
                      </InputAdornment>
                    )
                  }}
                >
                </TextField>
              </Grid>
              <Grid item>
                <Button onClick={onClick} variant='contained'>
                  Login
                </Button>
              </Grid>
            </Grid>
          </CardActions>
        </Card>
      </Grid>
    </Grid>
  )
}

export default Login
