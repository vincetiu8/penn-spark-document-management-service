import { Profile } from './Profile'
import React from 'react'
import { Updates } from './Updates'
import { Grid } from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  grid: {
    padding: '1rem',
    width: 'fit-content'
  }
})

export const Dashboard = () => {
  const classes = useStyles()

  return (
    <Grid container className={classes.grid} spacing={2}>
      <Grid item>
        <Profile/>
      </Grid>
      <Grid item>
        <Updates/>
      </Grid>
    </Grid>
  )
}
