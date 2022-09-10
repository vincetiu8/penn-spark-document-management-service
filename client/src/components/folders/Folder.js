import { Card, CardContent, Grid, LinearProgress, Typography } from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { useParams } from 'react-router-dom'
import { getFolder, selectFolderById } from '../../store/foldersSlice'
import { ErrUserForbidden } from '../../store/errors'
import React, { useEffect, useState } from 'react'
import { ParentFoldersBar } from './ParentFoldersBar'
import { ChildInfo } from './ChildInfo'
import { makeStyles } from '@material-ui/core/styles'
import { EditDialog } from './EditDialog'
import { DeleteDialog } from './DeleteDialog'
import { UserInfoDialog } from './UserInfoDialog'

const useStyles = makeStyles({
  grid: {
    padding: '2rem'
  }
})

export const Folder = () => {
  const classes = useStyles()

  let { id } = useParams()
  id = parseInt(id)

  const foldersState = useSelector(state => state.folders)
  const folder = useSelector(state => selectFolderById(state, id))

  const dispatch = useDispatch()

  useEffect(() => {
    if (!foldersState.ids.includes(id) && !foldersState.loading && !foldersState.error) {
      dispatch(getFolder(id))
    }
  }, [id, foldersState])

  const getErrorMessage = () => {
    if (foldersState.error === ErrUserForbidden) {
      return 'User forbidden'
    }
    if (!foldersState.loading && folder === undefined) {
      return 'Could not find folder, please check your connection...'
    }
    return ''
  }

  const [userInfoOpen, setUserInfoOpen] = useState(false)
  const [userData, setUserData] = useState({
    username: '',
    first_name: '',
    last_name: ''
  })
  const onShowUser = data => {
    setUserData(data)
    setUserInfoOpen(true)
  }

  const renderContents = () => {
    if (folder.child_folders.length === 0 && folder.files.length === 0 && !foldersState.loading) {
      return (
        <Grid container direction='column' alignItems='center' className={classes.grid}>
          <Grid item>
            <Card>
              <CardContent>
                <Typography variant='h3'>
                  No files or child folders present.
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )
    }

    return renderChildFoldersAndFiles()
  }

  const compareChildren = (a, b) => {
    if (b.name > a.name) {
      return -1
    }
    if (b.name < a.name) {
      return 1
    }
    return 0
  }

  const renderChildFoldersAndFiles = () => {
    return (
      <div>
        {
          folder.child_folders.slice()
            .sort(compareChildren)
            .map(folder =>
              <ChildInfo
                key={folder.id}
                data={folder}
                folder
                showUser={onShowUser}
              />
            )
        }
        {
          folder.files.slice()
            .sort(compareChildren)
            .map(file =>
              <ChildInfo
                key={file.id}
                data={file}
                showUser={onShowUser}
              />
            )
        }
      </div>
    )
  }

  const errorMessage = getErrorMessage()
  if (errorMessage !== '') {
    return (
      <Card>
        <CardContent>
          <Grid container justify='center'>
            <Grid item>
              <Typography variant='h3'>
                {errorMessage}
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>
    )
  }

  return (
    <div>
      <EditDialog/>
      <DeleteDialog/>
      <UserInfoDialog open={userInfoOpen} onClose={() => setUserInfoOpen(false)} userData={userData}/>
      <ParentFoldersBar id={id} showUser={onShowUser}/>
      {
        foldersState.loading
          ? <LinearProgress/>
          : renderContents()
      }
    </div>
  )
}
