import { Button, ButtonGroup, Card, CardActions } from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import React from 'react'

const useStyles = makeStyles({
  card: {
    width: '100%',
    verticalAlign: 'middle',
    justifyContent: 'center',
    display: 'inline-flex'
  },
  button: {
    fontSize: '1.25rem',
    textTransform: 'none'
  }
})

export const CreateUserBar = ({ openEditDialog }) => {
  const classes = useStyles()

  return (
    <Card square className={classes.card}>
      <CardActions>
        <ButtonGroup>
          <Button className={classes.button} onClick={() => openEditDialog(true)}>
            Create New User
          </Button>
        </ButtonGroup>
      </CardActions>
    </Card>
  )
}
