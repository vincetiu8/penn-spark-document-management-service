import { Button, ButtonGroup, Card, CardActions, Typography } from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles'
import { Link } from 'react-router-dom'
import { getFolder, updateFunction } from '../../store/foldersSlice'
import React, { useEffect, useState } from 'react'
import { getTime } from '../dashboard/UpdateInfo'
import { getUser } from '../../store/usersSlice'

const useStyles = makeStyles({
  card: {
    width: '100%',
    verticalAlign: 'middle',
    display: 'inline-flex'
  },
  button: {
    fontSize: '1.25rem',
    textTransform: 'none'
  },
  cardActions: {
    float: 'left'
  },
  cardContent: {
    display: 'flex',
    flexGrow: 1,
    width: 'auto',
    paddingLeft: '1.25rem'
  },
  cardWithParentContent: {
    display: 'flex',
    flexGrow: 1,
    width: 'auto',
    paddingLeft: 0
  },
  parentFolderActions: {
    paddingRight: 0
  }
})

export const ParentFoldersBar = ({
  id,
  showUser
}) => {
  const classes = useStyles()
  const dispatch = useDispatch()

  const foldersState = useSelector(state => state.folders)
  const [parentFolders, setParentFolders] = useState([])

  const onCreateFolder = () => {
    dispatch(updateFunction({
      type: 'create',
      folder: true,
      id: parentFolders[0].id
    }))
  }

  const onEditFolder = () => {
    dispatch(updateFunction({
      type: 'edit',
      folder: true,
      entity: parentFolders[0].id
    }))
  }

  const onDeleteFolder = () => {
    dispatch(updateFunction({
      type: 'delete',
      folder: true,
      ids: {
        id: parentFolders[0].id,
        parentFolderID: parentFolders[0].id.parent_folder_id
      }
    }))
  }

  const onCreateFile = () => {
    dispatch(updateFunction({
      type: 'create',
      folder: false,
      id: parentFolders[0].id
    }))
  }

  useEffect(() => {
    if (foldersState.loading) return

    let currentId = id
    const folders = []
    while (currentId !== 0) {
      if (!foldersState.ids.includes(currentId)) {
        dispatch(getFolder(currentId))
        return
      }

      folders.push(foldersState.entities[currentId])
      currentId = foldersState.entities[currentId].parent_folder_id
    }

    setParentFolders(folders)
  }, [id, foldersState])

  const generateFolders = () => {
    const parentFolderComponents = []
    for (let i = 0; i < (parentFolders.length > 3 ? 2 : 3); i++) {
      if (i >= parentFolders.length) break

      if (i !== 0) {
        parentFolderComponents.push(
          <Typography variant='h6'>
            /
          </Typography>
        )
      }

      parentFolderComponents.push(
        <Button
          className={classes.button}
          component={Link}
          to={'/folders/' + parentFolders[i].id}
          disabled={i === 0}
          style={{ color: 'black' }}
        >
          {parentFolders[i].name}
        </Button>
      )
    }

    if (parentFolders.length > 3) {
      parentFolderComponents.push(
        <Typography variant='h6'>
          / ... /
        </Typography>
      )
      parentFolderComponents.push(
        <Button
          className={classes.button}
          component={Link}
          to={'/folders/' + parentFolders[parentFolders.length - 1].id}
        >
          {parentFolders[parentFolders.length - 1].name}
        </Button>
      )
    }

    return parentFolderComponents.reverse()
  }

  const [lastEditor, setLastEditor] = useState(null)
  const usersState = useSelector(state => state.users)
  useEffect(() => {
    console.log(parentFolders[0])
    if (lastEditor || usersState.loading || !parentFolders[0]) return

    if (usersState.ids.includes(parentFolders[0].last_editor_id)) {
      setLastEditor(usersState.entities[parentFolders[0].last_editor_id])
      return
    }
    dispatch(getUser(parentFolders[0].last_editor_id))
  }, [lastEditor, parentFolders, usersState])

  if (foldersState.loading || parentFolders.length === 0) {
    return ''
  }

  return (
    <Card square elevation={6} className={classes.card}>
      <CardActions className={classes.cardContent}>
        {generateFolders()}
      </CardActions>
      <CardActions className={classes.cardActions}>
        <Button className={classes.button} onClick={() => showUser(lastEditor)}>
          Last Edited
          By {lastEditor ? lastEditor.username : ''} On {getTime(parentFolders[0].updated_at)}
        </Button>
      </CardActions>
      <CardActions className={classes.cardActions}>
        <ButtonGroup>
          <Button className={classes.button} onClick={onCreateFolder}>
            Create Child Folder
          </Button>
          <Button className={classes.button} onClick={onCreateFile}>
            Upload File
          </Button>
          <Button className={classes.button} onClick={onEditFolder}>Edit Folder</Button>
          {
            id !== 1
              ? <Button className={classes.button} onClick={onDeleteFolder}>
                Delete Folder
              </Button>
              : ''
          }
        </ButtonGroup>
      </CardActions>
    </Card>
  )
}
