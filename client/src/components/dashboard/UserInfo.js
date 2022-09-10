import { Grid, TextField } from '@material-ui/core'
import React from 'react'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  textField: {
    '& .MuiInputBase-root.Mui-disabled': {
      color: 'black'
    }
  }
})

export const UserInfo = ({ userData }) => {
  const classes = useStyles()

  return (
    <Grid container alignItems='center' direction='column' spacing={3}>
      <Grid item>
        <TextField
          type='text'
          label='username'
          disabled
          value={userData.username}
          className={classes.textField}
        />
      </Grid>
      <Grid item>
        <TextField
          type='text'
          label='first name'
          disabled
          value={userData.first_name}
          className={classes.textField}
        />
      </Grid>
      <Grid item>
        <TextField
          type='text'
          label='last name'
          disabled
          value={userData.last_name}
          className={classes.textField}
        />
      </Grid>
    </Grid>
  )
}
