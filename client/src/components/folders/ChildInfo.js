import {
  Button,
  ButtonGroup,
  Card,
  CardActionArea,
  CardActions,
  Grid,
  Typography
} from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import { Folder, InsertDriveFile, VisibilityOff } from '@material-ui/icons'
import { Link } from 'react-router-dom'
import React, { useEffect } from 'react'
import { getFile, updateFile, updateFunction } from '../../store/foldersSlice'
import { useDispatch, useSelector } from 'react-redux'
import { getUser, selectUserById } from '../../store/usersSlice'
import { getTime } from '../dashboard/UpdateInfo'

const useStyles = makeStyles({
  card: {
    verticalAlign: 'middle',
    display: 'inline-flex',
    width: '100%'
  },
  cardActionArea: {
    paddingLeft: '2ch',
    display: 'flex',
    flexGrow: 1,
    width: 'auto'
  },
  cardActions: {
    float: 'left'
  },
  button: {
    textTransform: 'none',
    fontSize: '1.25rem',
    width: '100%',
    whiteSpace: 'nowrap'
  },
  typography: {
    fontSize: '1.25rem'
  }
})

export const ChildInfo = ({
  folder,
  data,
  showUser
}) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const foldersState = useSelector(state => state.folders)

  const onGetFile = () => {
    dispatch(getFile(data))
  }

  const onPublish = () => {
    if (folder || foldersState.entities[data.folder_id].access_level !== 4) return

    dispatch(updateFile({
      id: data.id,
      is_published: !data.is_published
    }))
  }

  const onEdit = () => {
    if (folder) {
      dispatch(updateFunction({
        type: 'edit',
        folder: true,
        entity: data
      }))
      return
    }
    dispatch(updateFunction({
      type: 'edit',
      folder: false,
      entity: data
    }))
  }

  const onDelete = () => {
    if (folder) {
      dispatch(updateFunction({
        type: 'delete',
        folder: true,
        ids: {
          id: data.id,
          parentFolderID: data.parent_folder_id
        }
      }))
      return
    }
    dispatch(updateFunction({
      type: 'delete',
      folder: false,
      ids: {
        id: data.id,
        folderID: data.folder_id
      }
    }))
  }

  const lastEditor = useSelector(state => selectUserById(state, data.last_editor_id))
  const usersState = useSelector(state => state.users)
  useEffect(() => {
    if (lastEditor || usersState.loading) return

    dispatch(getUser(data.last_editor_id))
  }, [lastEditor, data])

  return (
    <Card variant='outlined' square className={classes.card}>
      <CardActionArea
        className={classes.cardActionArea}
        component={folder ? Link : undefined}
        to={folder ? '/folders/' + data.id : undefined}
        onClick={!folder ? onGetFile : undefined}
      >
        <Grid container spacing={2} className={classes.grid}>
          <Grid item>
            {folder ? <Folder/> : data.is_published ? <InsertDriveFile/> : <VisibilityOff/>}
          </Grid>
          <Grid item>
            <Typography className={classes.typography}>
              {data.name}
            </Typography>
          </Grid>
        </Grid>
      </CardActionArea>
      <CardActions className={classes.cardActions}>
        <Button className={classes.button} onClick={() => showUser(lastEditor)}>
          Last Edited By {lastEditor ? lastEditor.username : ''} On {getTime(data.updated_at)}
        </Button>
      </CardActions>
      <CardActions className={classes.cardActions}>
        <ButtonGroup>
          {
            !folder && foldersState.entities[data.folder_id].access_level === 4
              ? (
                <Button className={classes.button} onClick={onPublish}>
                  {data.is_published ? 'Unp' : 'P'}ublish File
                </Button>)
              : ''
          }
          <Button className={classes.button} onClick={onEdit}>
            Edit {folder ? 'Folder' : 'File'}
          </Button>
          <Button className={classes.button} onClick={onDelete}>
            Delete {folder ? 'Folder' : 'File'}
          </Button>
        </ButtonGroup>
      </CardActions>
    </Card>
  )
}
