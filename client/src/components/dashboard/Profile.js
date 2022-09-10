import { Button, Card, CardActions, CardContent, Typography } from '@material-ui/core'
import React, { useState } from 'react'
import { useSelector } from 'react-redux'
import { ChangePasswordDialog } from '../users/ChangePasswordDialog'
import { UserInfo } from './UserInfo'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  card: {
    width: 'fit-content'
  },
  cardActions: {
    justifyContent: 'center'
  }
})

export const Profile = () => {
  const classes = useStyles()

  const authState = useSelector(state => state.auth)

  const [open, setOpen] = useState(false)
  const onClose = () => {
    setOpen(false)
  }
  return (
    <div>
      <ChangePasswordDialog open={open} onClose={onClose} id={authState.userData.id}/>
      <Card elevation={6} className={classes.card}>
        <CardContent>
          <Typography align='center' variant='h5'>
            Profile
          </Typography>
          <UserInfo userData={authState.userData}/>
        </CardContent>
        <CardActions className={classes.cardActions}>
          <Button onClick={() => setOpen(true)}>Change Password</Button>
        </CardActions>
      </Card>
    </div>
  )
}
