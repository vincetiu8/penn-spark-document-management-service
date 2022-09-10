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
import { deleteUser } from '../../store/usersSlice'

export const DeleteUserDialog = ({
  open,
  onClose,
  id
}) => {
  const dispatch = useDispatch()
  const usersState = useSelector(state => state.users)

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const onClick = () => {
    setLoading(true)

    console.log(id)
    dispatch(deleteUser(id))
  }

  useEffect(() => {
    if (!loading || usersState.loading) return

    setLoading(false)
    if (!usersState.error) {
      onClose()
      return
    }

    setError(usersState.error)
  }, [usersState])

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <DialogTitle>
          Are you sure you want to delete this user?
        </DialogTitle>
        {
          error
            ? <DialogContentText>{error}</DialogContentText>
            : ''
        }
      </DialogContent>
      <DialogActions>
        <Button onClick={onClick}>
          Delete User
        </Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
