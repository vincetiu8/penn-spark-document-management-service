import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  Grid,
  LinearProgress,
  TextField
} from '@material-ui/core'
import Checkbox from '@material-ui/core/Checkbox'
import React, { useEffect, useState } from 'react'
import { makeStyles } from '@material-ui/core/styles'
import {
  ErrRequiredFirstName,
  ErrRequiredLastName,
  ErrRequiredUserPassword,
  ErrRequiredUserUsername
} from '../../store/errors'
import { useDispatch, useSelector } from 'react-redux'
import { createUser, updateUser } from '../../store/usersSlice'

const useStyles = makeStyles({
  dialogActions: {
    justifyContent: 'center'
  },
  grid: {
    paddingTop: '1rem'
  },
  dialogContent: {
    overflow: 'visible'
  },
  dialogContentText: {
    marginTop: '0.5rem',
    marginBottom: '0'
  }
})

export const EditUserDialog = ({
  open,
  onClose,
  create,
  userData
}) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const [newUser, setNewUser] = useState({})

  useEffect(() => {
    if (!open) return

    if (create) {
      setNewUser({
        username: '',
        first_name: '',
        last_name: '',
        password: '',
        is_admin: false
      })
      return
    }

    setNewUser({
      id: userData.id,
      username: userData.username,
      first_name: userData.first_name,
      last_name: userData.last_name,
      is_admin: userData.is_admin
    })
  }, [open])

  const [userError, setUserError] = useState('')

  const onChange = (property, value) => {
    setNewUser({
      ...newUser,
      [property]: value
    })
  }

  const onSubmit = () => {
    if (newUser.username === '') {
      setUserError(ErrRequiredUserUsername)
      return
    }

    if (newUser.first_name === '') {
      setUserError(ErrRequiredFirstName)
      return
    }

    if (newUser.last_name === '') {
      setUserError(ErrRequiredLastName)
      return
    }

    if (newUser.password === '') {
      setUserError(ErrRequiredUserPassword)
      return
    }

    setLoading(true)
    if (create) {
      dispatch(createUser(newUser))
      return
    }

    dispatch(updateUser(newUser))
  }

  const getUserError = property => {
    if (property === 'first_name') {
      if (userError === ErrRequiredFirstName) return userError
      return ''
    }

    if (property === 'last_name') {
      if (userError === ErrRequiredLastName) return userError
      return ''
    }

    if (property === 'password') {
      if (userError === ErrRequiredUserPassword) return userError
      return ''
    }

    if (userError === ErrRequiredFirstName || userError === ErrRequiredLastName || userError === ErrRequiredUserPassword) return ''
    return userError
  }

  const usersState = useSelector(state => state.users)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!loading || usersState.loading) return

    setLoading(false)
    if (!usersState.error) {
      onClose()
      return
    }

    setUserError(usersState.error)
  }, [usersState])

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent className={classes.dialogContent}>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <Grid container direction='column'>
          {
            Object.keys(newUser)
              .map(
                property => (
                  property === 'id'
                    ? ''
                    : (
                      <Grid item key={property}>
                        {
                          property === 'is_admin'
                            ? (
                              <Grid container className={classes.grid} spacing={4}>
                                <Grid item>
                                  <DialogContentText className={classes.dialogContentText}>
                                    admin
                                  </DialogContentText>
                                </Grid>
                                <Grid item>
                                  <Checkbox
                                    checked={newUser.is_admin}
                                    onChange={() => onChange('is_admin', !newUser.is_admin)}
                                  />
                                </Grid>
                              </Grid>)
                            : (
                              <TextField
                                label={property.replace('_', ' ')}
                                type={property === 'password' ? 'password' : 'text'}
                                value={newUser[property]}
                                onChange={e => onChange(property, e.target.value)}
                                helperText={getUserError(property)}
                                error={getUserError(property) !== ''}
                              />)
                        }
                      </Grid>)
                )
              )
          }
        </Grid>
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onSubmit}>
          {create ? 'Create' : 'Edit'} User
        </Button>
        <Button color='secondary' onClick={onClose}>
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}
