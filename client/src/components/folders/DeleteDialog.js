import React, { useEffect, useState } from 'react'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  LinearProgress
} from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { deleteFile, deleteFolder, updateFunction } from '../../store/foldersSlice'

export const DeleteDialog = () => {
  const dispatch = useDispatch()

  const foldersState = useSelector(state => state.folders)

  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const onClose = () => {
    dispatch(updateFunction(null))
    setOpen(false)
  }

  const onClick = () => {
    setLoading(true)
    if (foldersState.function.folder) {
      dispatch(deleteFolder(foldersState.function.ids))
      return
    }
    dispatch(deleteFile(foldersState.function.ids))
  }

  useEffect(() => {
    if (open || foldersState.function === null || foldersState.function.type !== 'delete') {
      return
    }

    setOpen(true)
    setError('')
  }, [foldersState])

  useEffect(() => {
    if (!open || !loading || foldersState.loading) {
      return
    }

    setLoading(false)
    if (foldersState.error) {
      setError(foldersState.error)
      return
    }

    onClose()
  }, [foldersState])

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <DialogTitle>
          Are you sure you want to delete this {foldersState.function === null || foldersState.function.folder ? 'folder' : 'file'}?
        </DialogTitle>
        {
          error
            ? <DialogContentText>{error}</DialogContentText>
            : ''
        }
      </DialogContent>
      <DialogActions>
        <Button onClick={onClick}>
          Delete {foldersState.function === null || foldersState.function.folder ? 'Folder' : 'File'}
        </Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
