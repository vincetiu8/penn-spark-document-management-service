import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  FormControl,
  FormHelperText,
  InputLabel,
  LinearProgress,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Select
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { getFolder } from '../../store/foldersSlice'
import { makeStyles } from '@material-ui/core/styles'
import { Folder } from '@material-ui/icons'
import { addAccessRole } from '../../store/userRolesSlice'
import { ErrRequiredAccessLevel } from '../../store/errors'
import { accessLevels } from './AccessRoleInfo'

const useStyles = makeStyles({
  dialog: {
    width: '100%',
    height: 'auto'
  },
  dialogActions: {
    justifyContent: 'left'
  },
  form: {
    justifyContent: 'center'
  },
  button: {
    textTransform: 'none',
    '&:disabled': { color: 'black' }
  },
  formControl: {
    minWidth: 150
  }
})

export const AddAccessRoleDialog = ({
  open,
  onClose,
  userRoleData
}) => {
  const classes = useStyles()
  const dispatch = useDispatch()
  const foldersState = useSelector(state => state.folders)

  const [loading, setLoading] = useState(false)
  const [folder, setFolder] = useState(null)
  const [error, setError] = useState('')

  const loadFolder = id => {
    if (foldersState.ids.includes(id)) {
      setFolder(foldersState.entities[id])
      return
    }

    setFolder({
      id: id,
      name: '',
      child_folders: []
    })
    setLoading(true)
    dispatch(getFolder(id))
  }

  useEffect(() => {
    if (!open) return

    loadFolder(1)
  }, [open])

  useEffect(() => {
    if (!loading || foldersState.loading) return

    setLoading(false)
    if (foldersState.error) {
      setError(foldersState.error)
      return
    }

    if (folder.name === '') {
      setFolder(foldersState.entities[folder.id])
      return
    }

    onClose()
  }, [loading, foldersState])

  const [accessLevel, setAccessLevel] = useState('')

  const onSubmit = () => {
    if (accessLevel === '') {
      setError(ErrRequiredAccessLevel)
    }
    setLoading(true)
    dispatch(addAccessRole({
      user_role_id: userRoleData.id,
      folder_id: folder.id,
      access_level: accessLevel
    }))
  }

  return (
    <Dialog open={open} onClose={onClose}>
      {
        loading
          ? <LinearProgress/>
          : ''
      }
      {
        folder
          ? <div>
            <DialogActions className={classes.dialogActions}>
              {
                folder.name !== '' && folder.parent_folder_id !== 0
                  ? <Button
                    className={classes.button}
                    onClick={() => loadFolder(folder.parent_folder_id)}
                  >{foldersState.entities[folder.parent_folder_id].name}</Button>
                  : ''
              }
              {
                folder.parent_folder_id !== 0
                  ? <Button disabled className={classes.button}>/</Button>
                  : ''
              }
              <Button className={classes.button} disabled>{folder.name}</Button>
            </DialogActions>
            <DialogContent>
              <Divider/>
            </DialogContent>
            <DialogActions className={classes.dialogActions}>
              <List>
                {
                  folder.child_folders.map(folder => (
                    <ListItem button key={folder.id} onClick={() => loadFolder(folder.id)}>
                      <ListItemIcon><Folder/></ListItemIcon>
                      <ListItemText primary={folder.name}/>
                    </ListItem>
                  ))
                }
              </List>
            </DialogActions>
          </div>
          : ''
      }
      <DialogActions className={classes.form}>
        <FormControl className={classes.formControl}>
          <InputLabel>Access Level</InputLabel>
          <Select
            value={accessLevel}
            onChange={e => setAccessLevel(e.target.value)}
            error={error !== ''}
          >
            {
              accessLevels.map((level, index) =>
                <MenuItem key={index} value={index}>{level}</MenuItem>)
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
        <Button onClick={onSubmit}>
          Add Access Role
        </Button>
        <Button onClick={onClose} color='secondary'>Cancel</Button>
      </DialogActions>
    </Dialog>
  )
}
