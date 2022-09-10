import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  LinearProgress
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { deleteUserRole } from '../../store/userRolesSlice'

export const DeleteUserRoleDialog = ({
  open,
  onClose,
  id
}) => {
  const dispatch = useDispatch()
  const userRolesState = useSelector(state => state.userRoles)

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const onClick = () => {
    setLoading(true)

    dispatch(deleteUserRole(id))
  }

  useEffect(() => {
    if (!loading || userRolesState.loading) return

    setLoading(false)
    if (!userRolesState.error) {
      onClose()
      return
    }

    setError(userRolesState.error)
  }, [userRolesState])

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <DialogTitle>
          Are you sure you want to delete this user role?
        </DialogTitle>
        {
          error
            ? <DialogContentText>{error}</DialogContentText>
            : ''
        }
      </DialogContent>
      <DialogActions>
        <Button onClick={onClick}>
          Delete User Role
        </Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
