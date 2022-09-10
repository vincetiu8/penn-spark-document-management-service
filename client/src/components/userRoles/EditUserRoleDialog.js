import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  LinearProgress,
  TextField
} from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import React, { useEffect, useState } from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { ErrRequiredUserRoleName } from '../../store/errors'
import { createUserRole, updateUserRole } from '../../store/userRolesSlice'

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

export const EditUserRoleDialog = ({
  open,
  setOpen,
  create,
  userRoleData
}) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const userRolesState = useSelector(state => state.userRoles)

  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const onSubmit = () => {
    if (name === '') {
      setError(ErrRequiredUserRoleName)
      return
    }

    setLoading(true)
    if (create) {
      dispatch(createUserRole({
        name: name
      }))
      return
    }
    dispatch(updateUserRole({
      id: userRoleData.id,
      name: name
    }))
  }

  useEffect(() => {
    if (!loading || userRolesState.loading) return

    setLoading(false)
    if (userRolesState.error) {
      setError(userRolesState.error)
      return
    }

    setOpen(false)
  })

  useEffect(() => {
    if (!open) return

    setName(create ? '' : userRoleData.name)
  }, [open, create, userRoleData])

  return (
    <Dialog open={open} onClose={() => setOpen(false)}>
      <DialogContent className={classes.dialogContent}>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <TextField
          autoFocus
          label='name'
          type='text'
          value={name}
          onChange={e => setName(e.target.value)}
          helperText={error}
          error={error !== ''}
        />
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onSubmit}>
          {create ? 'Create' : 'Edit'} User Role
        </Button>
        <Button color='secondary' onClick={() => setOpen(false)}>
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}
