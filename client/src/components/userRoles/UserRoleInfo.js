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
import { AccessRoleInfo } from './AccessRoleInfo'
import { makeStyles } from '@material-ui/core/styles'
import AssignmentIndIcon from '@material-ui/icons/AssignmentInd'

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

export const UserRoleInfo = ({
  userRoleData,
  openEdit,
  openDelete,
  openAdd
}) => {
  const classes = useStyles()

  const createNewAccessRole = e => {
    e.stopPropagation()
    openAdd(userRoleData)
  }

  const editUserRole = e => {
    e.stopPropagation()
    openEdit(false, userRoleData)
  }

  const deleteUserRole = e => {
    e.stopPropagation()
    openDelete(userRoleData.id)
  }

  return (
    <Accordion square>
      <AccordionSummary>
        <Grid container spacing={2} className={classes.grid}>
          <Grid item>
            <AssignmentIndIcon/>
          </Grid>
          <Grid item>
            <Typography className={classes.typography}>{userRoleData.name}</Typography>
          </Grid>
        </Grid>
        <ButtonGroup className={classes.buttonGroup}>
          <Button className={classes.button} onClick={createNewAccessRole}>
            Create New Access Role
          </Button>
          <Button className={classes.button} onClick={editUserRole}>
            Edit User Role
          </Button>
          {
            userRoleData.id === 1
              ? ''
              : (
                <Button className={classes.button} onClick={deleteUserRole}>
                  Delete User Role
                </Button>)
          }
        </ButtonGroup>
      </AccordionSummary>
      <AccordionActions>
        <Grid container direction='column'>
          {
            userRoleData.access_roles.map(
              accessRole => <AccessRoleInfo key={accessRole.id} accessRoleData={accessRole}/>
            )
          }
        </Grid>
      </AccordionActions>
    </Accordion>
  )
}
