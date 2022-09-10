import { Drawer, List, ListItem, ListItemIcon, ListItemText, Toolbar } from '@material-ui/core'
import AssignmentIndIcon from '@material-ui/icons/AssignmentInd'
import SupervisedUserCircleIcon from '@material-ui/icons/SupervisedUserCircle'
import { Dashboard, ExitToApp, Folder } from '@material-ui/icons'
import { Link } from 'react-router-dom'
import React from 'react'
import { logoutUser } from '../../store/authSlice'
import { clearFolder } from '../../store/foldersSlice'
import { clearUsers } from '../../store/usersSlice'
import { useDispatch, useSelector } from 'react-redux'

export const NavDrawer = ({
  open,
  onClose
}) => {
  const dispatch = useDispatch()

  const authState = useSelector(state => state.auth)

  const logout = () => {
    dispatch(logoutUser())
    dispatch(clearFolder())
    dispatch(clearUsers())
    onClose()
  }

  return (
    <div>
      <Drawer
        anchor='left'
        open={open}
        onClose={onClose}
        style={{ zIndex: 1250 }}
      >
        <Toolbar/>
        <List>
          <ListItem button key='Dashboard' component={Link} to='/' onClick={onClose}>
            <ListItemIcon><Dashboard/></ListItemIcon>
            <ListItemText primary='Dashboard'/>
          </ListItem>
          <ListItem button key='Filesystem' component={Link} to='/folders/1' onClick={onClose}>
            <ListItemIcon><Folder/></ListItemIcon>
            <ListItemText primary='Filesystem'/>
          </ListItem>
          {
            authState.isAuthenticated && authState.userData.is_admin
              ? (
                <div>
                  <ListItem button key='Manage Users' component={Link} to='/users' onClick={onClose}>
                    <ListItemIcon><SupervisedUserCircleIcon/></ListItemIcon>
                    <ListItemText primary='Manage Users'/>
                  </ListItem>
                  <ListItem button key='Manage User Roles' component={Link} to='/user-roles'
                            onClick={onClose}>
                    <ListItemIcon><AssignmentIndIcon/></ListItemIcon>
                    <ListItemText primary='Manage User Roles'/>
                  </ListItem>
                </div>)
              : ''
          }
          <ListItem button key='Logout' onClick={logout}>
            <ListItemIcon><ExitToApp/></ListItemIcon>
            <ListItemText primary='Logout'/>
          </ListItem>
        </List>
      </Drawer>
    </div>
  )
}
