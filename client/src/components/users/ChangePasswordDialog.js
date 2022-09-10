import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  Grid,
  LinearProgress,
  TextField
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { ErrRequiredUserPassword } from '../../store/errors'
import { makeStyles } from '@material-ui/core/styles'
import { updateUser } from '../../store/usersSlice'

const useStyles = makeStyles({
  dialogActions: {
    justifyContent: 'center'
  }
})

export const ChangePasswordDialog = ({
  open,
  onClose,
  id
}) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const usersState = useSelector(state => state.users)

  const [loading, setLoading] = useState(false)
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [passwordError, setPasswordError] = useState('')
  const [confirmPasswordError, setConfirmPasswordError] = useState('')

  const onClick = (e) => {
    e.preventDefault()
    if (password === '') {
      setPasswordError(ErrRequiredUserPassword)

      return
    }
    setPasswordError('')

    if (password !== confirmPassword) {
      setConfirmPasswordError('passwords do not match')
      return
    }
    setConfirmPasswordError('')

    setLoading(true)
    dispatch(updateUser({
      id: id,
      is_admin: true,
      password: password
    }))
  }

  useEffect(() => {
    if (!loading || usersState.loading) return

    setLoading(false)
    if (usersState.error) {
      setPasswordError(usersState.error)
      return
    }

    onClose()
  }, [usersState, loading])

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <Grid container direction='column' spacing={2}>
          <Grid item>
            <TextField
              autoFocus
              label='new password'
              type='password'
              error={passwordError !== ''}
              helperText={passwordError}
              value={password}
              onChange={e => setPassword(e.target.value)}
            />
          </Grid>
          <Grid item>
            <TextField
              label='confirm password'
              type='password'
              error={confirmPasswordError !== ''}
              helperText={confirmPasswordError}
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
            />
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onClick} disabled={loading}>Save</Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
