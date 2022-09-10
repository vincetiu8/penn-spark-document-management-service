import { Button, Dialog, DialogActions, DialogContent } from '@material-ui/core'
import { UserInfo } from '../dashboard/UserInfo'
import React from 'react'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  dialogContent: {
    overflow: 'clip'
  }
})

export const UserInfoDialog = ({
  open,
  onClose,
  userData
}) => {
  const classes = useStyles()

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogContent className={classes.dialogContent}>
        <UserInfo userData={userData}/>
      </DialogContent>
      <DialogActions>
        <Button color='secondary' onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  )
}
