import { Card, CardActions, CardContent, List, Typography } from '@material-ui/core'
import React, { useEffect } from 'react'
import { getFolder } from '../../store/foldersSlice'
import { useDispatch, useSelector } from 'react-redux'
import { UpdateInfo } from './UpdateInfo'

export const Updates = () => {
  const dispatch = useDispatch()
  const foldersState = useSelector(state => state.folders)

  useEffect(() => {
    if (!foldersState.ids.includes(1)) {
      if (!foldersState.loading && !foldersState.error) {
        dispatch(getFolder(1))
      }
    }
  }, [foldersState])

  const renderContents = () => {
    return foldersState.updates.map(update => <UpdateInfo key={update.id} update={update}/>)
  }

  if (foldersState.updates.length === 0) {
    return ''
  }

  return (
    <Card>
      <CardContent>
        <Typography variant='h5'>Recent Updates</Typography>
      </CardContent>
      <CardActions>
        <List>
          {
            renderContents()
          }
        </List>
      </CardActions>
    </Card>
  )
}
