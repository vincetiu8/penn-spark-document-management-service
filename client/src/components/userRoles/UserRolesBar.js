import { Button, ButtonGroup, Card, CardActions, CardContent, TextField } from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import React from 'react'

const useStyles = makeStyles({
  card: {
    width: '100%',
    justifyContent: 'left',
    display: 'inline-flex'
  },
  cardContent: {
    display: 'flex',
    flexGrow: 1,
    width: 'auto',
    paddingLeft: '1.25rem'
  },
  cardActions: {
    float: 'left'
  },
  button: {
    fontSize: '1.25rem',
    textTransform: 'none'
  }
})

export const UserRolesBar = ({
  searchTerm,
  setSearchTerm,
  openEdit
}) => {
  const classes = useStyles()

  return (
    <div>
      <Card square className={classes.card}>
        <CardContent className={classes.cardContent}>
          <TextField
            label='user role search term'
            value={searchTerm}
            onChange={e => setSearchTerm(e.target.value.toLowerCase())}
          />
        </CardContent>
        <CardActions className={classes.cardActions}>
          <ButtonGroup>
            <Button className={classes.button} onClick={() => openEdit(true)}>
              Create New User Role
            </Button>
          </ButtonGroup>
        </CardActions>
      </Card>
    </div>
  )
}
