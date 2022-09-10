import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  FormControl,
  FormHelperText,
  InputLabel,
  LinearProgress,
  MenuItem,
  Select
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { getUserRoles } from '../../store/userRolesSlice'
import { addUserRole } from '../../store/usersSlice'
import { makeStyles } from '@material-ui/core/styles'
import { ErrRequiredUserRoleName } from '../../store/errors'

const useStyles = makeStyles({
  formControl: {
    minWidth: 150
  },
  dialogActions: {
    justifyContent: 'center'
  }
})

export const AddUserRoleDialog = ({
  open,
  onClose,
  userData
}) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const usersState = useSelector(state => state.users)
  const userRolesState = useSelector(state => state.userRoles)

  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (userRolesState.loadedData) return

    if (userRolesState.loading) return

    if (loading) {
      setLoading(false)
      return
    }

    setLoading(true)
    dispatch(getUserRoles(null))
  }, [userRolesState, loading])

  const [userRole, setUserRole] = useState('')
  const [error, setError] = useState('')

  const onClick = () => {
    if (userRole === null) {
      setError(ErrRequiredUserRoleName)
      return
    }

    if (userData.user_roles.filter(role => role.id === userRole).length > 0) {
      setError('user already has user role')
      return
    }

    setLoading(true)
    dispatch(addUserRole({
      userID: userData.id,
      userRole: { id: userRole }
    }))
  }

  useEffect(() => {
    if (!loading || usersState.loading) return

    setLoading(false)
    if (!usersState.error) {
      onClose()
      return
    }

    setError(usersState.error)
  })

  useEffect(() => {
    setError('')
  }, [open])

  return (
    <Dialog open={open} onClose={onClose}>
      {
        loading
          ? <DialogContent>
            <LinearProgress/>
          </DialogContent>
          : ''
      }
      <DialogActions className={classes.dialogActions}>
        <FormControl className={classes.formControl}>
          <InputLabel>User Role</InputLabel>
          <Select
            value={userRole}
            onChange={e => setUserRole(e.target.value)}
            error={error !== ''}
          >
            {
              Object.values(userRolesState.entities)
                .map(userRole =>
                  <MenuItem key={userRole.id} value={userRole.id}>{userRole.name}</MenuItem>
                )
            }
          </Select>
          {
            error
              ? <FormHelperText>{error}</FormHelperText>
              : ''
          }
        </FormControl>
      </DialogActions>
      <DialogActions>
        <Button onClick={onClick}>
          Add User Role
        </Button>
        <Button color='secondary' onClick={onClose}>
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}
