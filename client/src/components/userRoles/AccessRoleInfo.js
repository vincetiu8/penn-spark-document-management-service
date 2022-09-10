import React, { useEffect, useState } from 'react'
import {
  Button,
  ButtonGroup,
  Card,
  CardActions,
  CardContent,
  Grid,
  Typography
} from '@material-ui/core'
import { Folder } from '@material-ui/icons'
import { makeStyles } from '@material-ui/core/styles'
import { useDispatch, useSelector } from 'react-redux'
import { getFolder } from '../../store/foldersSlice'
import { removeAccessRole } from '../../store/userRolesSlice'

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

export const accessLevels = ['Unset', 'None', 'Viewer', 'Uploader', 'Publisher']

export const AccessRoleInfo = ({ accessRoleData }) => {
  const classes = useStyles()
  const dispatch = useDispatch()

  const foldersState = useSelector(state => state.folders)

  const [path, setPath] = useState('')

  useEffect(() => {
    if (path !== '') return

    if (foldersState.loading) return

    let id = accessRoleData.folder_id
    let count = 0
    let currentPath = ''
    while (id !== 0) {
      if (!foldersState.ids.includes(id)) {
        dispatch(getFolder(accessRoleData.folder_id))
        return
      }

      if (count < 3) {
        currentPath = (foldersState.entities[id].parent_folder_id === 0 ? '' : ' / ') + foldersState.entities[id].name + currentPath
      }
      id = foldersState.entities[id].parent_folder_id
      count++
    }

    if (count >= 3) {
      currentPath = foldersState.entities[id].name + ' / ...' + currentPath
    }
    setPath(currentPath)
  }, [foldersState])

  const onRemove = () => {
    dispatch(removeAccessRole(accessRoleData))
  }

  return (
    <Grid item>
      <Card variant='outlined' square className={classes.card}>
        <CardContent
          className={classes.cardActionArea}
        >
          <Grid container spacing={2} className={classes.grid}>
            <Grid item>
              <Folder/>
            </Grid>
            <Grid item>
              <Typography className={classes.typography}>
                {path} ({accessLevels[accessRoleData.access_level]})
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
        {
          accessRoleData.id === 1
            ? ''
            : (
              <CardActions className={classes.cardActions}>
              <ButtonGroup>
                <Button className={classes.button} onClick={onRemove}>
                  Remove Access Role
                </Button>
              </ButtonGroup>
            </CardActions>)
        }

      </Card>
    </Grid>
  )
}
