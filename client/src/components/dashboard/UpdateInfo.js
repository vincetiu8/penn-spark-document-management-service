import { Link } from 'react-router-dom'
import { ListItem, ListItemIcon, ListItemText } from '@material-ui/core'
import { Folder, InsertDriveFile } from '@material-ui/icons'
import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles'
import { getUser, selectUserById } from '../../store/usersSlice'

const useStyles = makeStyles({
  listItem: {
    padding: '0.5rem'
  }
})

const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']

export const getTime = (time) => {
  const date = new Date(Date.parse(time))
  const mins = date.getMinutes()
  return days[date.getDay()] + ' ' + date.getHours() + ':' + (mins > 9 ? mins : '0' + mins)
}

export const UpdateInfo = ({ update }) => {
  const dispatch = useDispatch()
  const classes = useStyles()

  const user = useSelector(state => selectUserById(state, update.last_editor_id))
  const usersState = useSelector(state => state.users)

  useEffect(() => {
    if (user || usersState.loading) return

    dispatch(getUser(update.last_editor_id))
  }, [user, usersState])

  return (
    <ListItem
      key={update.id}
      button
      component={Link}
      to={'/folders/' + (update.access_roles !== undefined ? update.id : update.folder_id)}
    >
      <ListItemIcon>
        {update.access_roles !== undefined ? <Folder/> : <InsertDriveFile/>}
      </ListItemIcon>
      <ListItemText className={classes.listItem}>{update.name}</ListItemText>
      <ListItemText className={classes.listItem}>
        Updated By {user ? user.username : ''} On {getTime(update.updated_at)}
      </ListItemText>
    </ListItem>
  )
}
