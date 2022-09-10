import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  LinearProgress,
  TextField
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import {
  createFile,
  createFolder,
  updateFile,
  updateFolder,
  updateFunction
} from '../../store/foldersSlice'
import { ErrRequiredFolderName } from '../../store/errors'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  dialogActions: {
    justifyContent: 'center'
  }
})

export const EditDialog = () => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const foldersState = useSelector(state => state.folders)

  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [name, setName] = useState('')
  const [nameError, setNameError] = useState('')
  const [fileData, setFileData] = useState(null)

  const onClose = () => {
    dispatch(updateFunction(null))
    setOpen(false)
  }

  const onChange = (e) => {
    e.preventDefault()
    setName(e.target.value)
  }

  const onClick = () => {
    if (!open || foldersState.function === null) {
      return
    }

    if (name === '') {
      setNameError(ErrRequiredFolderName)
      return
    }

    if (foldersState.function.entity !== undefined && name === foldersState.function.entity.name) {
      setNameError('no changes to name')
      return
    }

    setLoading(true)
    if (foldersState.function.type === 'create') {
      if (foldersState.function.folder) {
        dispatch(createFolder({
          name: name,
          parent_folder_id: foldersState.function.id
        }))
        return
      }

      if (fileData === null) {
        setNameError('required file')
        setLoading(false)
        return
      }
      dispatch(createFile({
        file: {
          name: name,
          folder_id: foldersState.function.id
        },
        data: fileData
      }))
      return
    }

    if (foldersState.function.folder) {
      dispatch(updateFolder({
        id: foldersState.function.entity.id,
        name: name
      }))
      return
    }

    dispatch(updateFile({
      id: foldersState.function.entity.id,
      name: name
    }))
  }

  useEffect(() => {
    if (open || foldersState.function === null) {
      return
    }

    if (foldersState.function.type !== 'create' && foldersState.function.type !== 'edit') {
      return
    }

    if (foldersState.function.type === 'edit') {
      setName(foldersState.function.entity.name)
    } else {
      setName('')
    }

    setOpen(true)
    setNameError('')
  }, [foldersState])

  useEffect(() => {
    if (!open || !loading || foldersState.loading) {
      return
    }

    setLoading(false)

    if (foldersState.error) {
      setNameError(foldersState.error)
      return
    }

    onClose()
  })

  const handleUpload = (e) => {
    setFileData(e.target.files[0])
    setName(e.target.files[0].name)
  }

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent>
        {
          loading
            ? <LinearProgress/>
            : ''
        }
        <TextField
          autoFocus
          label='Name'
          type='text'
          error={nameError !== ''}
          helperText={nameError}
          value={name}
          onChange={onChange}
        />
      </DialogContent>
      {
        foldersState.function !== null && !foldersState.function.folder && foldersState.function.type === 'create'
          ? (
            <DialogActions>
              <input type='file' onChange={handleUpload}/>
            </DialogActions>
            )
          : ''
      }
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onClick} disabled={loading}>Save</Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
