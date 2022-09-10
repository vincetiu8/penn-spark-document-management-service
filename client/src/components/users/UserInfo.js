import {
  Accordion,
  AccordionActions,
  AccordionSummary,
  Button,
  ButtonGroup,
  Grid,
  Typography
} from '@material-ui/core'
import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { UserRoleInfo } from './UserRoleInfo'
import { AccountCircle, SupervisedUserCircle } from '@material-ui/icons'
import { useSelector } from 'react-redux'

const useStyles = makeStyles({
  grid: {
    display: 'flex',
    flexGrow: 1,
    width: 'auto',
    paddingLeft: '1.25rem',
    paddingTop: '0.5rem'
  },
  buttonGroup: {
    float: 'left'
  },
  button: {
    fontSize: '1.25rem',
    textTransform: 'none'
  },
  typography: {
    fontSize: '1.25rem'
  }
})

export const UserInfo = ({
  userData,
  openAdd,
  openEdit,
  openPassword,
  openDelete
}) => {
  const classes = useStyles()

  const authState = useSelector(state => state.auth)

  const addUserRole = e => {
    e.stopPropagation()
    openAdd(userData)
  }

  const editUser = e => {
    e.stopPropagation()
    openEdit(false, userData)
  }

  const changePassword = e => {
    e.stopPropagation()
    openPassword(false, userData)
  }

  const deleteUser = e => {
    e.stopPropagation()
    openDelete(userData)
  }

  return (
    <Accordion square>
      <AccordionSummary>
        <Grid container spacing={2} className={classes.grid}>
          <Grid item>
            {
              userData.is_admin
                ? <SupervisedUserCircle/>
                : <AccountCircle/>
            }
          </Grid>
          <Grid item>
            <Typography
              className={classes.typography}>{userData.username} ({userData.first_name} {userData.last_name})</Typography>
          </Grid>
        </Grid>
        <ButtonGroup className={classes.buttonGroup}>
          <Button className={classes.button} onClick={addUserRole}>
            Add User Role
          </Button>
          <Button className={classes.button} onClick={editUser}>
            Edit User
          </Button>
          <Button className={classes.button} onClick={changePassword}>
            Change Password
          </Button>
          {
            userData.id === 1 || userData.id === authState.userData.id
              ? ''
              : (
                <Button className={classes.button} onClick={deleteUser}>
                  Delete User
                </Button>)
          }

        </ButtonGroup>
      </AccordionSummary>
      <AccordionActions>
        <Grid container direction='column'>
          {
            userData.user_roles.map(
              userRole => <UserRoleInfo
                key={userRole.id}
                userRoleData={userRole}
                userData={userData}
              />
            )
          }
        </Grid>
      </AccordionActions>
    </Accordion>
  )
}
